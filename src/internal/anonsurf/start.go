package anonsurf

import (
	"fmt"
	"strings"

	"gitlab.com/parrotsec/packages/anonsurf/internal/config"
	"gitlab.com/parrotsec/packages/anonsurf/internal/devices"
	"gitlab.com/parrotsec/packages/anonsurf/internal/firewall"
	"gitlab.com/parrotsec/packages/anonsurf/internal/ipv6"
	"gitlab.com/parrotsec/packages/anonsurf/internal/mac"
	"gitlab.com/parrotsec/packages/anonsurf/internal/netutils"
	"gitlab.com/parrotsec/packages/anonsurf/internal/proxy"
	"gitlab.com/parrotsec/packages/anonsurf/internal/systemd"
	"gitlab.com/parrotsec/packages/anonsurf/internal/timezone"
	"gitlab.com/parrotsec/packages/anonsurf/internal/tor"
	"gitlab.com/parrotsec/packages/anonsurf/internal/vpn"
)

func (c Controller) Start() error {
	state := startTransition{}
	c.reportStatus("Checking Tor service state")
	active, err := systemd.NewTorController().IsActiveOrActivating()
	if err != nil {
		return fmt.Errorf("check Tor service: %w", err)
	}
	state.torWasActive, err = captureTorState(active)
	if err != nil {
		return err
	}

	// This flag means "rollback is armed", not "generation succeeded",
	// so it is set before the call.
	//
	// failStart() can then restore the old file or remove the generated one;
	// if generation failed before touching disk, Restore safely does nothing.
	state.torrcGenerated = true
	c.reportStatus("Generating Tor configuration")
	if err := c.generateTorrc(); err != nil {
		return c.failStart(fmt.Errorf("generate torrc: %w", err), state)
	}

	if exists, err := fileExists(tor.TorrcPath); err != nil {
		return c.failStart(fmt.Errorf("check torrc: %w", err), state)
	} else if !exists {
		return c.failStart(ErrTorrcNotFound, state)
	}
	routing, err := c.readTorRoutingConfig()
	if err != nil {
		return c.failStart(err, state)
	}
	cfg, cfgErr := config.Load(config.SystemPath)
	if cfgErr != nil {
		cfg = config.Default()
	}
	if err := c.startSurf(routing, &state, cfg); err != nil {
		return c.failStart(err, state)
	}

	state.torStartAttempted = true
	c.reportStatus("Starting Tor")
	if err := systemd.NewTorController().Start(); err != nil {
		return c.failStart(fmt.Errorf("start Tor service: %w", err), state)
	}
	if !routing.PreserveSystemDNS {
		c.reportStatus("Configuring Tor DNS")
		if err := c.configureDNS(&state); err != nil {
			return c.failStart(err, state)
		}
	}
	if cfg.UseProxy && !cfg.ProxyBeforeTor {
		if err := c.applyProxy(cfg, &state); err != nil {
			return c.failStart(err, state)
		}
	}
	if err := c.markActive(); err != nil {
		return c.failStart(fmt.Errorf("mark daemon active: %w", err), state)
	}
	c.reportStatus("Traffic is routed through Tor")
	return nil
}

func (c Controller) startSurf(routing tor.RoutingConfig, state *startTransition, cfg config.Config) error {
	if cfg.ChangeMac {
		c.reportStatus("Randomizing MAC address")
		iface := cfg.NetworkInterface
		if iface == "" {
			var err error
			iface, err = mac.DefaultInterface()
			if err != nil {
				return fmt.Errorf("detect interface for MAC change: %w", err)
			}
		}
		original, err := mac.GetCurrent(iface)
		if err != nil {
			return fmt.Errorf("read current MAC: %w", err)
		}
		if err := mac.Randomize(iface); err != nil {
			return err
		}
		state.macTouched = true
		state.macIface = iface
		state.macOriginal = original
		if err := writeFile(MACStatePath, []byte(iface+"\n"+original), 0600); err != nil {
			return fmt.Errorf("save MAC state: %w", err)
		}
	}

	if cfg.RenewIP {
		c.reportStatus("Renewing local IP address")
		iface := cfg.NetworkInterface
		if iface == "" {
			if detected, err := mac.DefaultInterface(); err == nil {
				iface = detected
			}
		}
		if iface != "" {
			if err := netutils.RenewIP(iface); err != nil {
				c.reportStatus("Warning: could not renew IP: " + err.Error())
			}
		}
	}

	if cfg.ForceUTC {
		c.reportStatus("Setting timezone to UTC")
		original, err := timezone.GetCurrent()
		if err == nil && original != "UTC" {
			if err := timezone.Set("UTC"); err != nil {
				return fmt.Errorf("set timezone UTC: %w", err)
			}
			state.timezoneTouched = true
			state.timezoneOriginal = original
			if err := writeFile(TimezoneStatePath, []byte(original), 0600); err != nil {
				return fmt.Errorf("save timezone state: %w", err)
			}
		}
	}

	if cfg.UseProxy && cfg.ProxyBeforeTor {
		if err := c.applyProxy(cfg, state); err != nil {
			return err
		}
	}

	if cfg.UseVPN && cfg.VPNConfigPath != "" {
		c.reportStatus("Connecting VPN")
		if err := vpn.Connect(cfg.VPNConfigPath); err != nil {
			return fmt.Errorf("connect VPN: %w", err)
		}
		state.vpnConnected = true
		state.vpnConfigPath = cfg.VPNConfigPath
	}

	c.reportStatus("Stopping Tor before changing routes")
	if err := systemd.NewTorController().Stop(); err != nil {
		return fmt.Errorf("stop Tor service before routing: %w", err)
	}

	ipv6Manager := ipv6.Manager{}
	c.reportStatus("Disabling IPv6")
	if err := ipv6Manager.Capture(); err != nil {
		return fmt.Errorf("capture IPv6 state: %w", err)
	}
	state.ipv6Touched = true
	if err := ipv6Manager.Disable(); err != nil {
		return err
	}

	manager := firewall.Manager{}
	c.reportStatus("Applying firewall rules")
	if err := manager.Backup(); err != nil {
		return err
	}
	state.firewallTouched = true
	if err := manager.Apply(routing); err != nil {
		return err
	}

	if cfg.DisableCameraMic {
		c.reportStatus("Disabling camera and microphone")
		if err := devices.Disable(); err != nil {
			c.reportStatus("Warning: could not fully disable camera/mic: " + err.Error())
		}
		state.devicesTouched = true
	}

	if cfg.BlockWebRTC {
		c.reportStatus("Blocking WebRTC leaks")
		if err := firewall.BlockWebRTC(); err != nil {
			return fmt.Errorf("block WebRTC: %w", err)
		}
		state.webrtcBlocked = true
	}

	if cfg.BlockTelemetry {
		c.reportStatus("Blocking telemetry endpoints")
		if err := firewall.BlockTelemetry(); err != nil {
			return fmt.Errorf("block telemetry: %w", err)
		}
		state.telemetryBlocked = true
	}
	return nil
}

func (c Controller) readTorRoutingConfig() (tor.RoutingConfig, error) {
	data, err := readFile(tor.TorrcPath)
	if err != nil {
		return tor.RoutingConfig{}, fmt.Errorf("read torrc: %w", err)
	}
	return tor.ParseRoutingConfig(strings.NewReader(string(data)))
}

func (c Controller) applyProxy(cfg config.Config, state *startTransition) error {
	proxyCfg, ok := proxy.FromAnonSurfConfig(
		cfg.UseProxy,
		cfg.Socks5ProxyAddress,
		cfg.Socks5ProxyCredential,
		cfg.HTTPProxyAddress,
		cfg.HTTPProxyCredential,
		cfg.ProxyIsHTTPS,
	)
	if !ok {
		return nil
	}
	c.reportStatus("Configuring proxy")
	if err := proxy.Apply(proxyCfg); err != nil {
		return fmt.Errorf("configure proxy: %w", err)
	}
	state.proxyApplied = true
	return nil
}

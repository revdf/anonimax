package anonsurf

import (
	"errors"
	"fmt"
	"strings"

	"gitlab.com/parrotsec/packages/anonsurf/internal/cleanup"
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

func (c Controller) Stop() error {
	c.reportStatus("Checking managed network state")
	managed, err := c.hasManagedState()
	if err != nil {
		return err
	}
	if !managed {
		if err := c.restoreDNS(); err != nil {
			return err
		}
		c.reportStatus("Anonimax routing is inactive")
		return nil
	}
	// Recovery must not mistake a partially restored session for active routing.
	if err := c.clearActive(); err != nil {
		return err
	}
	restartTor, err := torWasActive()
	if err != nil {
		return err
	}

	var failures errorCollector
	c.reportStatus("Stopping Tor")
	failures.add("stop Tor service", systemd.NewTorController().Stop())
	c.reportStatus("Restoring firewall rules")
	failures.add("restore firewall", (firewall.Manager{}).Restore())
	c.reportStatus("Restoring DNS settings")
	failures.add("restore DNS settings", c.restoreDNS())
	c.reportStatus("Restoring Tor configuration")
	failures.add("restore torrc", c.restoreTorrc())
	c.reportStatus("Restoring IPv6")
	failures.add("restore IPv6 state", (ipv6.Manager{}).Restore())
	if err := failures.err(); err != nil {
		c.reportStatus("Network restoration failed")
		return err
	}
	if restartTor && !systemd.IsSystemStopping() {
		c.reportStatus("Restarting the original Tor service")
		if err := systemd.NewTorController().Start(); err != nil {
			return fmt.Errorf("restart original Tor service: %w", err)
		}
	}
	if restartTor {
		if err := clearTorState(); err != nil {
			return err
		}
	}
	c.reportStatus("Restoring MAC address")
	failures.add("restore MAC address", c.restoreMAC())
	c.reportStatus("Restoring timezone")
	failures.add("restore timezone", c.restoreTimezone())
	if err := failures.err(); err != nil {
		c.reportStatus("Network restoration failed")
		return err
	}
	c.reportStatus("Renewing local IP address")
	if iface, err := mac.DefaultInterface(); err == nil {
		if err := netutils.RenewIP(iface); err != nil {
			c.reportStatus("Warning: could not renew IP: " + err.Error())
		}
	}
	cfg, err := config.Load(config.SystemPath)
	if err != nil {
		cfg = config.Default()
	}
	if cfg.UseVPN && cfg.VPNConfigPath != "" {
		c.reportStatus("Disconnecting VPN")
		if err := vpn.Disconnect(cfg.VPNConfigPath); err != nil {
			c.reportStatus("Warning: could not disconnect VPN: " + err.Error())
		}
	}
	c.reportStatus("Removing proxy configuration")
	failures.add("remove proxy", proxy.Remove())
	if err := failures.err(); err != nil {
		c.reportStatus("Warning: could not remove proxy: " + err.Error())
	}
	c.reportStatus("Re-enabling camera and microphone")
	if err := devices.Enable(); err != nil {
		c.reportStatus("Warning: could not re-enable camera/mic: " + err.Error())
	}
	c.reportStatus("Anonimax routing is inactive")
	if cfg.CleanOnStop {
		c.reportStatus("Cleaning system traces")
		if err := cleanup.CleanTraces(); err != nil {
			c.reportStatus("Warning: some traces could not be cleaned: " + err.Error())
		}
	}
	return nil
}

func (c Controller) restoreMAC() error {
	data, err := readFile(MACStatePath)
	if err != nil {
		// No MAC state means MAC was not changed — nothing to restore.
		return nil
	}
	parts := strings.SplitN(strings.TrimSpace(string(data)), "\n", 2)
	if len(parts) != 2 {
		return fmt.Errorf("malformed MAC state file")
	}
	iface, original := parts[0], parts[1]
	if err := mac.Restore(iface, original); err != nil {
		return err
	}
	return removeFile(MACStatePath)
}

func (c Controller) restoreTimezone() error {
	data, err := readFile(TimezoneStatePath)
	if err != nil {
		return nil
	}
	original := strings.TrimSpace(string(data))
	if original == "" {
		return removeFile(TimezoneStatePath)
	}
	if err := timezone.Set(original); err != nil {
		return err
	}
	return removeFile(TimezoneStatePath)
}

func (c Controller) restoreTorrc() error {
	if err := (tor.TorrcManager{}).Restore(); err != nil && !errors.Is(err, tor.ErrTorrcBackupNotFound) {
		return fmt.Errorf("restore torrc: %w", err)
	}
	return nil
}

package anonsurf

import (
	"fmt"

	"gitlab.com/parrotsec/packages/anonsurf/internal/devices"
	"gitlab.com/parrotsec/packages/anonsurf/internal/firewall"
	"gitlab.com/parrotsec/packages/anonsurf/internal/ipv6"
	"gitlab.com/parrotsec/packages/anonsurf/internal/mac"
	"gitlab.com/parrotsec/packages/anonsurf/internal/proxy"
	"gitlab.com/parrotsec/packages/anonsurf/internal/systemd"
	"gitlab.com/parrotsec/packages/anonsurf/internal/timezone"
	"gitlab.com/parrotsec/packages/anonsurf/internal/vpn"
)

type startTransition struct {
	torWasActive      bool
	torStartAttempted bool
	torrcGenerated    bool
	dnsTouched        bool
	ipv6Touched       bool
	firewallTouched   bool
	macTouched        bool
	macIface          string
	macOriginal       string
	vpnConnected      bool
	vpnConfigPath     string
	webrtcBlocked     bool
	timezoneTouched   bool
	timezoneOriginal  string
	telemetryBlocked  bool
	devicesTouched    bool
	proxyApplied      bool
}

func (c Controller) failStart(operation error, state startTransition) error {
	c.reportStatus("Restoring the previous network state")
	rollback := c.rollbackStart(state)
	if rollback == nil {
		c.reportStatus("Previous network state restored")
		return operation
	}
	c.reportStatus("Network rollback failed")
	return transitionError{operation: operation, rollback: rollback}
}

func (c Controller) rollbackStart(state startTransition) error {
	var failures errorCollector
	if state.torStartAttempted {
		failures.add("stop Tor service", systemd.NewTorController().Stop())
	}
	if state.devicesTouched {
		failures.add("re-enable camera/mic", devices.Enable())
	}
	if state.telemetryBlocked {
		failures.add("unblock telemetry", firewall.UnblockTelemetry())
	}
	if state.webrtcBlocked {
		failures.add("unblock WebRTC", firewall.UnblockWebRTC())
	}
	if state.firewallTouched {
		failures.add("restore firewall", (firewall.Manager{}).Restore())
	}
	if state.ipv6Touched {
		failures.add("restore IPv6 state", (ipv6.Manager{}).Restore())
	}
	if state.dnsTouched {
		failures.add("restore DNS settings", c.restoreDNS())
	}
	if state.torrcGenerated {
		failures.add("restore torrc", c.restoreTorrc())
	}
	if state.macTouched {
		failures.add("restore MAC address", mac.Restore(state.macIface, state.macOriginal))
		failures.add("clear MAC state", removeFile(MACStatePath))
	}
	if state.timezoneTouched {
		failures.add("restore timezone", timezone.Set(state.timezoneOriginal))
		failures.add("clear timezone state", removeFile(TimezoneStatePath))
	}
	if state.vpnConnected {
		failures.add("disconnect VPN", vpn.Disconnect(state.vpnConfigPath))
	}
	if state.proxyApplied {
		failures.add("remove proxy", proxy.Remove())
	}
	if err := failures.err(); err != nil {
		return err
	}
	if state.torWasActive {
		if err := systemd.NewTorController().Start(); err != nil {
			return fmt.Errorf("restart original Tor service: %w", err)
		}
		if err := clearTorState(); err != nil {
			return err
		}
	}
	return nil
}

package vpn

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/parrotsec/packages/anonsurf/internal/command"
)

var runner = command.Runner{}

// Connect starts the VPN from configPath.
// Presence of "[Interface]" selects WireGuard; otherwise OpenVPN is used.
func Connect(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("read VPN config: %w", err)
	}
	if isWireGuard(string(data)) {
		iface := wireGuardIface(configPath)
		if err := runner.Run("wg-quick", "up", configPath); err != nil {
			return fmt.Errorf("start WireGuard %s: %w", iface, err)
		}
		return nil
	}
	if err := runner.Run("openvpn", "--daemon", "--config", configPath); err != nil {
		return fmt.Errorf("start OpenVPN: %w", err)
	}
	return nil
}

// Disconnect stops the VPN associated with configPath.
func Disconnect(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		// Config may already be gone; fall back to stopping known daemons.
		_ = runner.Run("killall", "openvpn")
		return nil
	}
	if isWireGuard(string(data)) {
		iface := wireGuardIface(configPath)
		if err := runner.Run("wg-quick", "down", configPath); err != nil {
			return fmt.Errorf("stop WireGuard %s: %w", iface, err)
		}
		return nil
	}
	if err := runner.Run("pkill", "-f", "openvpn.*"+configPath); err != nil {
		_ = runner.Run("killall", "openvpn")
	}
	return nil
}

// IsConnected reports whether an OpenVPN or WireGuard process is running.
func IsConnected() bool {
	if _, err := runner.QuietOutput("pgrep", "openvpn"); err == nil {
		return true
	}
	if _, err := runner.QuietOutput("pgrep", "wg"); err == nil {
		return true
	}
	return false
}

func isWireGuard(content string) bool {
	return strings.Contains(content, "[Interface]")
}

func wireGuardIface(configPath string) string {
	base := filepath.Base(configPath)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

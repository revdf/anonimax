package mac

import (
	"fmt"
	"strings"

	"gitlab.com/parrotsec/packages/anonsurf/internal/command"
)

var runner = command.Runner{}

// DefaultInterface detects the network interface used for internet traffic.
func DefaultInterface() (string, error) {
	out, err := runner.QuietOutput("ip", "route", "get", "1.1.1.1")
	if err != nil {
		return "", fmt.Errorf("detect network interface: %w", err)
	}
	// Output: "1.1.1.1 via 192.168.x.x dev wlp3s0 src ..."
	fields := strings.Fields(string(out))
	for i, field := range fields {
		if field == "dev" && i+1 < len(fields) {
			return fields[i+1], nil
		}
	}
	return "", fmt.Errorf("could not parse interface from: %s", string(out))
}

// GetCurrent returns the current MAC address of iface.
func GetCurrent(iface string) (string, error) {
	out, err := runner.QuietOutput("ip", "link", "show", iface)
	if err != nil {
		return "", fmt.Errorf("get MAC for %s: %w", iface, err)
	}
	// Look for "link/ether aa:bb:cc:dd:ee:ff"
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "link/ether ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1], nil
			}
		}
	}
	return "", fmt.Errorf("MAC address not found for interface %s", iface)
}

// Randomize sets a random MAC address on iface.
func Randomize(iface string) error {
	if err := runner.Run("ip", "link", "set", iface, "down"); err != nil {
		return fmt.Errorf("bring down %s: %w", iface, err)
	}
	if err := runner.Run("macchanger", "-r", iface); err != nil {
		// Best-effort: try to bring interface back up even if macchanger failed.
		_ = runner.Run("ip", "link", "set", iface, "up")
		return fmt.Errorf("randomize MAC on %s: %w", iface, err)
	}
	if err := runner.Run("ip", "link", "set", iface, "up"); err != nil {
		return fmt.Errorf("bring up %s: %w", iface, err)
	}
	return nil
}

// Restore sets iface back to originalMAC.
func Restore(iface, originalMAC string) error {
	if err := runner.Run("ip", "link", "set", iface, "down"); err != nil {
		return fmt.Errorf("bring down %s: %w", iface, err)
	}
	if err := runner.Run("macchanger", "-m", originalMAC, iface); err != nil {
		_ = runner.Run("ip", "link", "set", iface, "up")
		return fmt.Errorf("restore MAC on %s: %w", iface, err)
	}
	if err := runner.Run("ip", "link", "set", iface, "up"); err != nil {
		return fmt.Errorf("bring up %s: %w", iface, err)
	}
	return nil
}

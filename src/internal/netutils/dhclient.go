package netutils

import (
	"fmt"

	"gitlab.com/parrotsec/packages/anonsurf/internal/command"
)

var runner = command.Runner{}

// RenewIP releases the current DHCP lease (best-effort) and requests a new one.
func RenewIP(iface string) error {
	// Release can fail when no lease exists; continue to request a new address.
	_ = runner.Run("dhclient", "-r", iface)
	if err := runner.Run("dhclient", iface); err != nil {
		return fmt.Errorf("renew IP on %s: %w", iface, err)
	}
	return nil
}

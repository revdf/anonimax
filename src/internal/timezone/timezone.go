package timezone

import (
	"fmt"
	"strings"

	"gitlab.com/parrotsec/packages/anonsurf/internal/command"
)

var runner = command.Runner{}

// GetCurrent returns the system timezone.
func GetCurrent() (string, error) {
	out, err := runner.QuietOutput("timedatectl", "show", "-p", "Timezone", "--value")
	if err != nil {
		return "", fmt.Errorf("get timezone: %w", err)
	}
	tz := strings.TrimSpace(string(out))
	if tz == "" {
		return "", fmt.Errorf("empty timezone from timedatectl")
	}
	return tz, nil
}

// Set changes the system timezone.
func Set(tz string) error {
	if err := runner.Run("timedatectl", "set-timezone", tz); err != nil {
		return fmt.Errorf("set timezone %s: %w", tz, err)
	}
	return nil
}

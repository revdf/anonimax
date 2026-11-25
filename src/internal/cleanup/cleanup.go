package cleanup

import (
	"fmt"
	"strings"

	"gitlab.com/parrotsec/packages/anonsurf/internal/command"
)

var runner = command.Runner{}

// CleanTraces clears caches, logs, temp files and DNS history.
// Each step runs independently; failures are collected and returned together.
func CleanTraces() error {
	var failures []string

	add := func(name string, err error) {
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", name, err))
		}
	}

	add("drop caches", runner.Run("sh", "-c", "sync && echo 3 > /proc/sys/vm/drop_caches"))
	add("journal rotate", runner.Run("journalctl", "--rotate"))
	add("journal vacuum", runner.Run("journalctl", "--vacuum-time=1s"))
	add("clear /tmp", runner.Run("find", "/tmp", "-type", "f", "-delete"))
	add("resolvectl flush", runner.Run("resolvectl", "flush-caches"))
	add("systemd-resolve flush", runner.Run("systemd-resolve", "--flush-caches"))
	add("bleachbit", runner.Run("bleachbit", "--clean",
		"system.cache",
		"system.tmp",
		"system.trash",
	))

	if len(failures) == 0 {
		return nil
	}
	return fmt.Errorf("%s", strings.Join(failures, "; "))
}

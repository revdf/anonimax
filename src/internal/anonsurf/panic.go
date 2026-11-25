package anonsurf

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"gitlab.com/parrotsec/packages/anonsurf/internal/safekill"
)

const panicStepTimeout = 3 * time.Second

// Panic performs an emergency shutdown: cuts the network, clears traces and powers off.
// Each step is best-effort with a short timeout; poweroff always runs last.
func (c Controller) Panic() error {
	c.reportStatus("PANIC: bringing network interfaces down")
	_ = panicBringIfacesDown()

	c.reportStatus("PANIC: stopping Tor")
	_ = runTimed(panicStepTimeout, "systemctl", "stop", "tor@default", "--no-block")

	c.reportStatus("PANIC: flushing firewall")
	_ = panicFlushFirewall()

	c.reportStatus("PANIC: restoring MAC address")
	_ = panicRestoreMAC()

	c.reportStatus("PANIC: clearing traces")
	_ = panicClearTraces()

	c.reportStatus("PANIC: killing sensitive processes")
	_ = panicKillProcesses()

	c.reportStatus("PANIC: powering off")
	if err := runTimed(panicStepTimeout, "systemctl", "poweroff", "--force"); err != nil {
		return fmt.Errorf("poweroff: %w", err)
	}
	return nil
}

func panicBringIfacesDown() error {
	out, err := outputTimed(panicStepTimeout, "ip", "-o", "link", "show")
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// "2: wlp3s0: <BROADCAST,...>"
		parts := strings.SplitN(line, ": ", 3)
		if len(parts) < 2 {
			continue
		}
		iface := strings.TrimSpace(parts[1])
		if iface == "" || iface == "lo" {
			continue
		}
		_ = runTimed(panicStepTimeout, "ip", "link", "set", iface, "down")
	}
	return nil
}

func panicFlushFirewall() error {
	cmds := [][]string{
		{"iptables", "-F"},
		{"iptables", "-t", "nat", "-F"},
		{"iptables", "-t", "mangle", "-F"},
		{"ip6tables", "-F"},
		{"ip6tables", "-t", "nat", "-F"},
	}
	for _, args := range cmds {
		_ = runTimed(panicStepTimeout, args[0], args[1:]...)
	}
	return nil
}

func panicRestoreMAC() error {
	data, err := os.ReadFile(MACStatePath)
	if err != nil {
		return nil
	}
	parts := strings.SplitN(strings.TrimSpace(string(data)), "\n", 2)
	if len(parts) != 2 {
		return fmt.Errorf("malformed MAC state file")
	}
	iface, original := parts[0], parts[1]
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, "macchanger", "-m", original, iface).Run()
}

func panicClearTraces() error {
	_ = runTimed(panicStepTimeout, "sh", "-c", "sync && echo 3 > /proc/sys/vm/drop_caches")
	_ = runTimed(panicStepTimeout, "journalctl", "--rotate")
	_ = runTimed(panicStepTimeout, "journalctl", "--vacuum-time=1s")
	_ = runTimed(panicStepTimeout, "find", "/tmp", "-type", "f", "-delete")
	return nil
}

func panicKillProcesses() error {
	nameTargets := make(map[string]bool, len(safekill.ProcessNames))
	for _, name := range safekill.ProcessNames {
		nameTargets[name] = true
	}
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return err
	}
	self := os.Getpid()
	for _, entry := range entries {
		pid, err := strconv.Atoi(entry.Name())
		if err != nil || pid <= 0 || pid == self {
			continue
		}
		comm, err := os.ReadFile("/proc/" + entry.Name() + "/comm")
		if err != nil {
			continue
		}
		if !nameTargets[strings.TrimSpace(string(comm))] {
			continue
		}
		_ = syscall.Kill(pid, syscall.SIGKILL)
	}
	return nil
}

func runTimed(timeout time.Duration, name string, args ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return exec.CommandContext(ctx, name, args...).Run()
}

func outputTimed(timeout time.Duration, name string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return exec.CommandContext(ctx, name, args...).Output()
}

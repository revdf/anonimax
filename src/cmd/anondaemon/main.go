package main

import (
	"fmt"
	"os"
	"strconv"

	"gitlab.com/parrotsec/packages/anonsurf/internal/anonsurf"
	"gitlab.com/parrotsec/packages/anonsurf/internal/monitor"
	"gitlab.com/parrotsec/packages/anonsurf/internal/systemd"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "anondaemon:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("expected command: start, stop, restart, recover, changeid, monitor, panic, autochange-start, autochange-stop, autochange-loop, proxy-set, or proxy-unset")
	}

	controller := anonsurf.NewController(newStatusReporter())
	switch args[0] {
	case "start":
		if err := controller.Start(); err != nil {
			return err
		}
		return controller.ResumeAutoChange()
	case "stop":
		_ = controller.SuspendAutoChange()
		return controller.Stop()
	case "restart":
		return controller.Restart()
	case "recover":
		return controller.Recover()
	case "changeid":
		return controller.ChangeIdentity()
	case "monitor":
		return monitor.NewNyxMonitor(os.Stdin, os.Stdout, os.Stderr).Run()
	case "panic":
		return controller.Panic()
	case "autochange-start":
		if len(args) != 2 {
			return fmt.Errorf("autochange-start requires an interval in minutes")
		}
		minutes, err := strconv.Atoi(args[1])
		if err != nil || minutes < 1 {
			return fmt.Errorf("interval must be a number >= 1")
		}
		return controller.EnableAutoChange(minutes)
	case "autochange-stop":
		return controller.DisableAutoChange()
	case "autochange-loop":
		if len(args) != 2 {
			return fmt.Errorf("autochange-loop requires an interval in minutes")
		}
		minutes, err := strconv.Atoi(args[1])
		if err != nil || minutes < 1 {
			return fmt.Errorf("interval must be a number >= 1")
		}
		return controller.RunAutoChangeLoop(minutes)
	case "proxy-set":
		if len(args) != 2 {
			return fmt.Errorf("proxy-set requires a proxy URL")
		}
		return controller.SetProxyConfig(args[1])
	case "proxy-unset":
		return controller.UnsetProxyConfig()
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func newStatusReporter() anonsurf.StatusReporter {
	// Notification failure path:
	//   publish phase -> socket error -> log once -> disable later reports
	// Reporting is observability only. Firewall, DNS and Tor transitions keep
	// running because failing them would turn a UI problem into a network or
	// recovery failure.
	available := true
	return func(status string) {
		if !available {
			return
		}
		if err := systemd.NotifyStatus(status); err != nil {
			fmt.Fprintln(os.Stderr, "anondaemon: report service progress:", err)
			available = false
		}
	}
}

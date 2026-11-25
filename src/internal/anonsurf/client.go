package anonsurf

import (
	"fmt"
	"strconv"

	"gitlab.com/parrotsec/packages/anonsurf/internal/command"
)

const BinaryPath = "/usr/lib/anonsurf/anondaemon"

type IdentityClient struct{}

type MonitorClient struct{}

type PanicClient struct{}

type AutoChangeClient struct{}

func (IdentityClient) Change() error {
	executable, args := command.PrivilegedCommand(BinaryPath, "changeid")
	if err := (command.Runner{}).Run(executable, args...); err != nil {
		return fmt.Errorf("request privileged identity change: %w", err)
	}
	return nil
}

func (PanicClient) Run() error {
	executable, args := command.PrivilegedCommand(BinaryPath, "panic")
	if err := (command.Runner{}).Run(executable, args...); err != nil {
		return fmt.Errorf("request privileged panic: %w", err)
	}
	return nil
}

func (AutoChangeClient) Start(minutes int) error {
	executable, args := command.PrivilegedCommand(BinaryPath, "autochange-start", strconv.Itoa(minutes))
	if err := (command.Runner{}).Run(executable, args...); err != nil {
		return fmt.Errorf("request privileged autochange start: %w", err)
	}
	return nil
}

func (AutoChangeClient) Stop() error {
	executable, args := command.PrivilegedCommand(BinaryPath, "autochange-stop")
	if err := (command.Runner{}).Run(executable, args...); err != nil {
		return fmt.Errorf("request privileged autochange stop: %w", err)
	}
	return nil
}

type ProxyClient struct{}

func (ProxyClient) Set(rawURL string) error {
	executable, args := command.PrivilegedCommand(BinaryPath, "proxy-set", rawURL)
	if err := (command.Runner{}).Run(executable, args...); err != nil {
		return fmt.Errorf("request privileged proxy set: %w", err)
	}
	return nil
}

func (ProxyClient) Unset() error {
	executable, args := command.PrivilegedCommand(BinaryPath, "proxy-unset")
	if err := (command.Runner{}).Run(executable, args...); err != nil {
		return fmt.Errorf("request privileged proxy unset: %w", err)
	}
	return nil
}

func (MonitorClient) Run() error {
	executable, args := command.PrivilegedCommand(BinaryPath, "monitor")
	if err := (command.Runner{}).RunInteractive(executable, args...); err != nil {
		return fmt.Errorf("open privileged Tor monitor: %w", err)
	}
	return nil
}

package anonsurf

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"gitlab.com/parrotsec/packages/anonsurf/internal/command"
	"gitlab.com/parrotsec/packages/anonsurf/internal/systemd"
	"gitlab.com/parrotsec/packages/anonsurf/internal/tor"
)

const (
	AutoChangeStatePath = "/run/anonsurf/autochange"
	AutoChangeUnit      = "anonsurf-autochange.service"
)

var ErrNotRunning = errors.New("AnonSurf is not running")

var (
	autoChangeMu   sync.Mutex
	autoChangeDone chan struct{}
)

func (Controller) ChangeIdentity() error {
	active, err := systemd.NewAnonsurfdController().IsActive()
	if err != nil {
		return fmt.Errorf("check AnonSurf service: %w", err)
	}
	if !active {
		return ErrNotRunning
	}
	return (tor.IdentityChanger{}).Change()
}

// StartAutoChange starts a goroutine that calls ChangeIdentity every minutes.
func (c Controller) StartAutoChange(minutes int) {
	autoChangeMu.Lock()
	defer autoChangeMu.Unlock()
	if autoChangeDone != nil {
		close(autoChangeDone)
	}
	done := make(chan struct{})
	autoChangeDone = done
	go func() {
		ticker := time.NewTicker(time.Duration(minutes) * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				_ = c.ChangeIdentity()
			case <-done:
				return
			}
		}
	}()
}

// StopAutoChange stops the auto-change loop in this process.
func StopAutoChange() {
	autoChangeMu.Lock()
	defer autoChangeMu.Unlock()
	if autoChangeDone != nil {
		close(autoChangeDone)
		autoChangeDone = nil
	}
}

// EnableAutoChange persists the interval and starts a long-running autochange unit.
func (c Controller) EnableAutoChange(minutes int) error {
	if minutes < 1 {
		return fmt.Errorf("interval must be a number >= 1")
	}
	active, err := systemd.NewAnonsurfdController().IsActive()
	if err != nil {
		return fmt.Errorf("check AnonSurf service: %w", err)
	}
	if !active {
		return ErrNotRunning
	}
	if err := SaveAutoChangeState(minutes); err != nil {
		return err
	}
	return c.startAutoChangeUnit(minutes)
}

// DisableAutoChange stops the loop and clears persisted state.
func (c Controller) DisableAutoChange() error {
	_ = c.stopAutoChangeUnit()
	StopAutoChange()
	return ClearAutoChangeState()
}

// SuspendAutoChange stops the loop but keeps state for resume after start.
func (c Controller) SuspendAutoChange() error {
	_ = c.stopAutoChangeUnit()
	StopAutoChange()
	return nil
}

// ResumeAutoChange restarts the loop if a persisted interval exists.
func (c Controller) ResumeAutoChange() error {
	minutes, err := LoadAutoChangeState()
	if err != nil || minutes < 1 {
		return nil
	}
	return c.startAutoChangeUnit(minutes)
}

// RunAutoChangeLoop runs the ticker until SIGTERM/SIGINT (used by the systemd unit).
func (c Controller) RunAutoChangeLoop(minutes int) error {
	if minutes < 1 {
		return fmt.Errorf("interval must be a number >= 1")
	}
	c.StartAutoChange(minutes)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
	StopAutoChange()
	return nil
}

func (Controller) startAutoChangeUnit(minutes int) error {
	runner := command.Runner{}
	_ = runner.Run("systemctl", "stop", AutoChangeUnit)
	if err := runner.Run(
		"systemd-run",
		"--unit=anonsurf-autochange",
		"--property=Restart=on-failure",
		BinaryPath, "autochange-loop", strconv.Itoa(minutes),
	); err != nil {
		return fmt.Errorf("start autochange unit: %w", err)
	}
	return nil
}

func (Controller) stopAutoChangeUnit() error {
	return (command.Runner{}).Run("systemctl", "stop", AutoChangeUnit)
}

func SaveAutoChangeState(minutes int) error {
	return writeFile(AutoChangeStatePath, []byte(strconv.Itoa(minutes)+"\n"), 0600)
}

func LoadAutoChangeState() (int, error) {
	data, err := readFile(AutoChangeStatePath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	minutes, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, fmt.Errorf("parse autochange state: %w", err)
	}
	return minutes, nil
}

func ClearAutoChangeState() error {
	return removeFile(AutoChangeStatePath)
}

package status

import (
	"errors"
	"fmt"
	"os"

	"gitlab.com/parrotsec/packages/anonsurf/internal/dns"
	"gitlab.com/parrotsec/packages/anonsurf/internal/systemd"
)

type DNSState string

const (
	DNSMissing       DNSState = "missing"
	DNSEmpty         DNSState = "empty"
	DNSTor           DNSState = "tor"
	DNSTorRedirect   DNSState = "tor-redirect"
	DNSClearnet      DNSState = "clearnet"
	DNSLocalhostOnly DNSState = "localhost-only"
)

type Snapshot struct {
	AnonimaxActive bool
	TorActive      bool
	DNSState       DNSState
	DNSAddresses   []string
}

type Reader struct{}

func (Reader) Snapshot() (Snapshot, error) {
	anoimaxActive, err := systemd.NewAnimaxdController().IsActive()
	if err != nil {
		return Snapshot{}, fmt.Errorf("read Anonimax status: %w", err)
	}
	torActive, err := systemd.NewTorController().IsActive()
	if err != nil {
		return Snapshot{}, fmt.Errorf("read Tor status: %w", err)
	}
	addresses, err := dns.ParseDNSAddresses(dns.ResolvConfPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return Snapshot{}, fmt.Errorf("read DNS status: %w", err)
	}

	return Snapshot{
		AnonimaxActive: anoimaxActive,
		TorActive:      torActive,
		DNSState:       classifyDNS(anoimaxActive, addresses, err),
		DNSAddresses:   addresses,
	}, nil
}

func classifyDNS(anoimaxActive bool, addresses []string, readErr error) DNSState {
	if errors.Is(readErr, os.ErrNotExist) {
		return DNSMissing
	}
	if len(addresses) == 0 {
		return DNSEmpty
	}
	if anoimaxActive {
		if dns.HasOnlyLocalhost(addresses) {
			return DNSTor
		}
		return DNSTorRedirect
	}
	if dns.HasOnlyLocalhost(addresses) {
		return DNSLocalhostOnly
	}
	return DNSClearnet
}

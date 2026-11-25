package firewall

import "gitlab.com/parrotsec/packages/anonsurf/internal/command"

var runner = command.Runner{}

// BlockWebRTC drops STUN traffic on ports 3478 and 5349 (UDP and TCP).
// STUN is used by WebRTC to discover the real IP, bypassing Tor/VPN.
func BlockWebRTC() error {
	rules := [][]string{
		{"iptables", "-A", "OUTPUT", "-p", "udp", "--dport", "3478", "-j", "DROP"},
		{"iptables", "-A", "OUTPUT", "-p", "tcp", "--dport", "3478", "-j", "DROP"},
		{"iptables", "-A", "OUTPUT", "-p", "udp", "--dport", "5349", "-j", "DROP"},
		{"iptables", "-A", "OUTPUT", "-p", "tcp", "--dport", "5349", "-j", "DROP"},
		{"ip6tables", "-A", "OUTPUT", "-p", "udp", "--dport", "3478", "-j", "DROP"},
		{"ip6tables", "-A", "OUTPUT", "-p", "tcp", "--dport", "3478", "-j", "DROP"},
		{"ip6tables", "-A", "OUTPUT", "-p", "udp", "--dport", "5349", "-j", "DROP"},
		{"ip6tables", "-A", "OUTPUT", "-p", "tcp", "--dport", "5349", "-j", "DROP"},
	}
	for _, rule := range rules {
		if err := runner.Run(rule[0], rule[1:]...); err != nil {
			return err
		}
	}
	return nil
}

// UnblockWebRTC removes the STUN DROP rules added by BlockWebRTC.
func UnblockWebRTC() error {
	rules := [][]string{
		{"iptables", "-D", "OUTPUT", "-p", "udp", "--dport", "3478", "-j", "DROP"},
		{"iptables", "-D", "OUTPUT", "-p", "tcp", "--dport", "3478", "-j", "DROP"},
		{"iptables", "-D", "OUTPUT", "-p", "udp", "--dport", "5349", "-j", "DROP"},
		{"iptables", "-D", "OUTPUT", "-p", "tcp", "--dport", "5349", "-j", "DROP"},
		{"ip6tables", "-D", "OUTPUT", "-p", "udp", "--dport", "3478", "-j", "DROP"},
		{"ip6tables", "-D", "OUTPUT", "-p", "tcp", "--dport", "3478", "-j", "DROP"},
		{"ip6tables", "-D", "OUTPUT", "-p", "udp", "--dport", "5349", "-j", "DROP"},
		{"ip6tables", "-D", "OUTPUT", "-p", "tcp", "--dport", "5349", "-j", "DROP"},
	}
	var lastErr error
	for _, rule := range rules {
		// -D fails if the rule doesn't exist — ignore those errors.
		if err := runner.Run(rule[0], rule[1:]...); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

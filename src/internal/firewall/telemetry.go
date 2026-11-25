package firewall

const telemetryChain = "ANONSURF-TELEMETRY"

// telemetryRanges lists known telemetry IP ranges for Google, Microsoft and Amazon.
var telemetryRanges = []string{
	// Google
	"142.250.0.0/15",
	"172.217.0.0/16",
	"216.58.0.0/16",
	"209.85.128.0/17",
	// Microsoft / Windows telemetry
	"13.64.0.0/11",
	"13.96.0.0/13",
	"13.104.0.0/14",
	"52.96.0.0/12",
	"40.76.0.0/14",
	// Amazon telemetry
	"54.239.0.0/16",
	"205.251.192.0/19",
}

// BlockTelemetry creates the ANONSURF-TELEMETRY chain and drops known telemetry ranges.
func BlockTelemetry() error {
	if err := runner.Run("iptables", "-N", telemetryChain); err != nil {
		return err
	}
	if err := runner.Run("iptables", "-I", "OUTPUT", "-j", telemetryChain); err != nil {
		return err
	}
	for _, cidr := range telemetryRanges {
		if err := runner.Run("iptables", "-A", telemetryChain, "-d", cidr, "-j", "DROP"); err != nil {
			return err
		}
	}
	return nil
}

// UnblockTelemetry removes the ANONSURF-TELEMETRY chain.
// If the chain does not exist, the error is ignored.
func UnblockTelemetry() error {
	// Detach from OUTPUT first, then flush and delete the chain.
	_ = runner.Run("iptables", "-D", "OUTPUT", "-j", telemetryChain)
	_ = runner.Run("iptables", "-F", telemetryChain)
	_ = runner.Run("iptables", "-X", telemetryChain)
	return nil
}

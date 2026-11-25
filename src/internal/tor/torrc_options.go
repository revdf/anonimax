package tor

import (
	"fmt"
	"strconv"
	"strings"
)

func renderBridgeSettings(options TorrcOptions) (string, error) {
	addresses := options.BridgeAddresses
	if len(addresses) == 0 && strings.TrimSpace(options.BridgeAddress) != "" {
		addresses = []string{options.BridgeAddress}
	}
	if options.AutoBridge && len(addresses) == 0 {
		return "", fmt.Errorf("auto bridge requires selected bridge addresses")
	}
	if len(addresses) == 0 {
		return "# Bridge setting: No bridge", nil
	}

	var transports []string
	seenTransports := make(map[string]bool)
	for index, address := range addresses {
		addresses[index] = strings.TrimSpace(address)
		if err := validateTorrcSingleLine("bridge address", addresses[index]); err != nil {
			return "", err
		}
		transport, err := bridgeTransport(addresses[index])
		if err != nil {
			return "", err
		}
		if transport == "meek" {
			return "", fmt.Errorf("meek bridges are not supported")
		}
		if transport == "snowflake" {
			if err := validateSnowflakeBridge(addresses[index]); err != nil {
				return "", err
			}
		}
		if transport != "" && !seenTransports[transport] {
			seenTransports[transport] = true
			transports = append(transports, transport)
		}
	}

	var builder strings.Builder
	builder.WriteString("# Enable bridge mode\nUseBridges 1\n")
	for _, transport := range transports {
		switch transport {
		case "obfs4":
			if strings.TrimSpace(options.Obfs4ProxyPath) == "" {
				return "", fmt.Errorf("obfs4 bridge requires obfs4proxy path")
			}
			builder.WriteString("ClientTransportPlugin obfs4 exec ")
			builder.WriteString(strings.TrimSpace(options.Obfs4ProxyPath))
			builder.WriteString(" managed\n")
		case "snowflake":
			if strings.TrimSpace(options.SnowflakeClientPath) == "" {
				return "", fmt.Errorf("snowflake bridge requires snowflake-client path")
			}
			builder.WriteString("ClientTransportPlugin snowflake exec ")
			builder.WriteString(strings.TrimSpace(options.SnowflakeClientPath))
			builder.WriteString("\n")
		}
	}

	for _, address := range addresses {
		builder.WriteString("\nBridge ")
		builder.WriteString(address)
	}
	return builder.String(), nil
}

func bridgeTransport(address string) (string, error) {
	fields := strings.Fields(address)
	if len(fields) == 0 {
		return "", fmt.Errorf("bridge address is empty")
	}

	switch fields[0] {
	case "obfs4", "snowflake", "meek":
		return fields[0], nil
	default:
		return "", nil
	}
}

func renderProxySettings(options TorrcOptions) (string, error) {
	if !options.UseProxy {
		return "", nil
	}

	socks5Address := strings.TrimSpace(options.Socks5ProxyAddress)
	socks5Credential := strings.TrimSpace(options.Socks5ProxyCredential)
	if socks5Address != "" {
		if err := validateTorrcSingleLine("socks5 proxy address", socks5Address); err != nil {
			return "", err
		}
		if err := validateTorrcSingleLine("socks5 proxy credential", socks5Credential); err != nil {
			return "", err
		}
		var builder strings.Builder
		builder.WriteString("# Set SOCKS5 proxy\n")
		builder.WriteString("Socks5Proxy ")
		builder.WriteString(socks5Address)
		builder.WriteString("\n")
		if socks5Credential != "" {
			user, pass, ok := strings.Cut(socks5Credential, ":")
			if !ok {
				user = socks5Credential
			}
			builder.WriteString("Socks5ProxyUsername ")
			builder.WriteString(user)
			builder.WriteString("\n")
			if pass != "" {
				builder.WriteString("Socks5ProxyPassword ")
				builder.WriteString(pass)
				builder.WriteString("\n")
			}
		}
		return strings.TrimRight(builder.String(), "\n"), nil
	}

	address := strings.TrimSpace(options.HTTPProxyAddress)
	credential := strings.TrimSpace(options.HTTPProxyCredential)

	if address == "" {
		if credential != "" {
			return "", fmt.Errorf("proxy credential requires proxy address")
		}
		return "", nil
	}
	if err := validateTorrcSingleLine("proxy address", address); err != nil {
		return "", err
	}
	if credential != "" {
		if err := validateTorrcSingleLine("proxy credential", credential); err != nil {
			return "", err
		}
	}

	var builder strings.Builder
	builder.WriteString("# Set proxy address\n")
	if options.ProxyIsHTTPS {
		builder.WriteString("HTTPSProxy ")
		builder.WriteString(address)
		builder.WriteString("\n")
		if credential != "" {
			builder.WriteString("HTTPSProxyAuthenticator ")
			builder.WriteString(credential)
			builder.WriteString("\n")
		}
		return strings.TrimRight(builder.String(), "\n"), nil
	}

	if credential != "" {
		return "", fmt.Errorf("HTTP proxy credentials require proxy_is_https=true")
	}
	builder.WriteString("HTTPProxy ")
	builder.WriteString(address)
	return builder.String(), nil
}

func validateTorrcSingleLine(name, value string) error {
	if strings.ContainsAny(value, "\r\n") {
		return fmt.Errorf("%s must be a single line", name)
	}
	return nil
}

func torBool(value bool) string {
	if value {
		return "1"
	}
	return "0"
}

func normalizePlaintextPorts(portList string) (string, error) {
	parts := strings.Split(portList, ",")
	ports := make([]string, 0, len(parts))
	seen := make(map[int]bool)

	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			return "", fmt.Errorf("empty plaintext port")
		}

		port, err := strconv.Atoi(value)
		if err != nil || port < 1 || port > 65535 {
			return "", fmt.Errorf("invalid plaintext port %q", value)
		}
		if seen[port] {
			continue
		}
		seen[port] = true
		ports = append(ports, strconv.Itoa(port))
	}
	if len(ports) == 0 {
		return "", fmt.Errorf("no plaintext ports configured")
	}

	return strings.Join(ports, ","), nil
}

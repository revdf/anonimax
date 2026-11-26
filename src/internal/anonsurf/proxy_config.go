package anonsurf

import (
	"fmt"

	"gitlab.com/parrotsec/packages/anonsurf/internal/config"
	"gitlab.com/parrotsec/packages/anonsurf/internal/proxy"
)

// SetProxyConfig parses url and writes proxy settings into anonimax.cfg.
func (Controller) SetProxyConfig(rawURL string) error {
	cfg, err := proxy.ParseURL(rawURL)
	if err != nil {
		return err
	}
	values := map[string]string{
		"use_proxy":               "true",
		"socks5_proxy_address":    "",
		"socks5_proxy_credential": "",
		"http_proxy_address":      "",
		"http_proxy_credential":   "",
		"proxy_is_https":          "false",
	}
	switch cfg.Type {
	case "socks5":
		values["socks5_proxy_address"] = cfg.Address
		values["socks5_proxy_credential"] = cfg.Credential
	case "https":
		values["http_proxy_address"] = cfg.Address
		values["http_proxy_credential"] = cfg.Credential
		values["proxy_is_https"] = "true"
	case "http":
		values["http_proxy_address"] = cfg.Address
		values["http_proxy_credential"] = cfg.Credential
		values["proxy_is_https"] = "false"
	default:
		return fmt.Errorf("unsupported proxy type %q", cfg.Type)
	}
	return config.UpdateKeys(config.SystemPath, values)
}

// UnsetProxyConfig disables proxy settings in anonimax.cfg.
func (Controller) UnsetProxyConfig() error {
	return config.UpdateKeys(config.SystemPath, map[string]string{
		"use_proxy":               "false",
		"socks5_proxy_address":    "",
		"socks5_proxy_credential": "",
		"http_proxy_address":      "",
		"http_proxy_credential":   "",
		"proxy_is_https":          "false",
	})
}

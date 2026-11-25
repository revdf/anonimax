package proxy

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"gitlab.com/parrotsec/packages/anonsurf/internal/fileops"
)

const StatePath = "/run/anonsurf/proxy-active"

type Config struct {
	Type       string // "http", "https", "socks5"
	Address    string // "IP:PORTA"
	Credential string // "user:pass" or ""
}

// ParseURL parses proxy URLs such as:
//
//	socks5://IP:PORTA
//	socks5://user:pass@IP:PORTA
//	http://IP:PORTA
//	https://IP:PORTA
func ParseURL(rawURL string) (Config, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return Config{}, fmt.Errorf("empty proxy URL")
	}
	if !strings.Contains(rawURL, "://") {
		rawURL = "socks5://" + rawURL
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return Config{}, fmt.Errorf("parse proxy URL: %w", err)
	}
	scheme := strings.ToLower(parsed.Scheme)
	switch scheme {
	case "socks5", "socks5h", "http", "https":
	default:
		return Config{}, fmt.Errorf("unsupported proxy scheme %q", parsed.Scheme)
	}
	if scheme == "socks5h" {
		scheme = "socks5"
	}
	host := parsed.Host
	if host == "" {
		return Config{}, fmt.Errorf("proxy URL missing host")
	}
	cfg := Config{Type: scheme, Address: host}
	if parsed.User != nil {
		user := parsed.User.Username()
		pass, _ := parsed.User.Password()
		if user != "" {
			if pass != "" {
				cfg.Credential = user + ":" + pass
			} else {
				cfg.Credential = user
			}
		}
	}
	return cfg, nil
}

// Apply configures proxy environment variables for the current process
// and marks proxy as active under /run/anonsurf.
func Apply(cfg Config) error {
	if cfg.Address == "" {
		return fmt.Errorf("proxy address is empty")
	}
	proxyURL := buildProxyURL(cfg)
	switch cfg.Type {
	case "socks5":
		_ = os.Setenv("ALL_PROXY", proxyURL)
		_ = os.Setenv("all_proxy", proxyURL)
		_ = os.Setenv("SOCKS5_SERVER", cfg.Address)
		_ = os.Unsetenv("http_proxy")
		_ = os.Unsetenv("https_proxy")
		_ = os.Unsetenv("HTTP_PROXY")
		_ = os.Unsetenv("HTTPS_PROXY")
	case "http", "https":
		_ = os.Setenv("http_proxy", proxyURL)
		_ = os.Setenv("https_proxy", proxyURL)
		_ = os.Setenv("HTTP_PROXY", proxyURL)
		_ = os.Setenv("HTTPS_PROXY", proxyURL)
		_ = os.Setenv("ALL_PROXY", proxyURL)
		_ = os.Setenv("all_proxy", proxyURL)
	default:
		return fmt.Errorf("unsupported proxy type %q", cfg.Type)
	}
	content := cfg.Type + "\n" + cfg.Address + "\n" + cfg.Credential + "\n"
	return fileops.WriteAtomicWithDir(StatePath, []byte(content), 0600, 0700)
}

// Remove clears proxy environment variables and state.
// It is a no-op when no proxy was applied.
func Remove() error {
	exists, err := fileops.Exists(StatePath)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	for _, key := range []string{
		"ALL_PROXY", "all_proxy", "SOCKS5_SERVER",
		"http_proxy", "https_proxy", "HTTP_PROXY", "HTTPS_PROXY",
	} {
		_ = os.Unsetenv(key)
	}
	return fileops.RemoveIfExists(StatePath)
}

func buildProxyURL(cfg Config) string {
	if cfg.Credential == "" {
		return cfg.Type + "://" + cfg.Address
	}
	return cfg.Type + "://" + cfg.Credential + "@" + cfg.Address
}

// FromAnonSurfConfig builds a proxy Config from AnonSurf settings.
// ok is false when proxy is disabled or no address is configured.
func FromAnonSurfConfig(useProxy bool, socks5Addr, socks5Cred, httpAddr, httpCred string, https bool) (Config, bool) {
	if !useProxy {
		return Config{}, false
	}
	if socks5Addr != "" {
		return Config{Type: "socks5", Address: socks5Addr, Credential: socks5Cred}, true
	}
	if httpAddr != "" {
		typ := "http"
		if https {
			typ = "https"
		}
		return Config{Type: typ, Address: httpAddr, Credential: httpCred}, true
	}
	return Config{}, false
}

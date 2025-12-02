// Package normalize validates DNS targets, domains, and query types.
// Delegates RFC compliance to miekg/dns and url parsing to stdlib instead of reimplementing validation.
package normalize

import (
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/miekg/dns"
)

const (
	// DefaultDNSPortInt is the default port for DNS (Do53)
	DefaultDNSPortInt = 53
	// DefaultDoTPortInt is the default port for DNS-over-TLS
	DefaultDoTPortInt = 853
	// DefaultHTTPSPortInt is the default port for DNS-over-HTTPS
	DefaultHTTPSPortInt = 443
	// DefaultDoQPortInt is the default port for DNS-over-QUIC
	DefaultDoQPortInt = 853
)

const (
	// SchemeUDP represents UDP protocol
	SchemeUDP = "udp"
	// SchemeTCP represents TCP protocol
	SchemeTCP = "tcp"
	// SchemeTLS represents TLS protocol
	SchemeTLS = "tls"
	// SchemeHTTPS represents HTTPS protocol
	SchemeHTTPS = "https"
	// SchemeQUIC represents QUIC protocol
	SchemeQUIC = "quic"
)

// ProtocolConfig maps DNS protocol schemes to default ports and display names.
type ProtocolConfig struct {
	Scheme       string
	DefaultPort  int
	DisplayName  string
	UsesHostname bool
}

// ProtocolConfigs serves as single source of truth for protocol metadata across the application.
var ProtocolConfigs = map[string]ProtocolConfig{
	SchemeUDP:   {SchemeUDP, DefaultDNSPortInt, "Do53", false},
	SchemeTCP:   {SchemeTCP, DefaultDNSPortInt, "Do53", false},
	SchemeTLS:   {SchemeTLS, DefaultDoTPortInt, "DoT", true},
	SchemeHTTPS: {SchemeHTTPS, DefaultHTTPSPortInt, "DoH", true},
	SchemeQUIC:  {SchemeQUIC, DefaultDoQPortInt, "DoQ", true},
}

// Target validates and normalizes DNS server targets.
// Minimal preprocessing here - delegate port/host handling to AdGuard dnsproxy.
func Target(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("empty target")
	}

	if strings.ContainsAny(raw, "\x00\r\n\t") {
		return "", fmt.Errorf("target contains invalid control characters")
	}

	// Default to UDP for bare addresses
	if !strings.Contains(raw, "://") {
		raw = "udp://" + raw
	}

	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("invalid target URL: %w", err)
	}

	scheme := strings.ToLower(u.Scheme)

	switch scheme {
	case SchemeUDP, SchemeTCP, SchemeTLS, SchemeHTTPS, SchemeQUIC:
	default:
		return "", fmt.Errorf("unsupported scheme '%s' (must be udp, tcp, tls, https, or quic)", scheme)
	}

	host := u.Hostname()
	if host == "" || host == ":" {
		return "", fmt.Errorf("target host cannot be empty")
	}
	if net.ParseIP(host) == nil {
		if _, ok := dns.IsDomainName(host); !ok {
			return "", fmt.Errorf("target host is not a valid domain name: %s", host)
		}
	}

	// DoH requires explicit path
	if scheme == SchemeHTTPS && u.Path == "" {
		raw += "/dns-query"
	}

	// Let AdGuard dnsproxy handle ports, IPv6 brackets, host parsing
	return raw, nil
}

// IsValidIP delegates to net.ParseIP for RFC compliance.
func IsValidIP(s string) bool {
	return net.ParseIP(s) != nil
}

// IsValidDomain delegates RFC 1035 validation to miekg/dns.
// dns.IsDomainName() handles length (253 chars max, 63/label) and character rules.
func IsValidDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain cannot be empty")
	}

	if strings.ContainsAny(domain, "\x00\r\n\t ") {
		return fmt.Errorf("domain contains invalid characters")
	}

	_, ok := dns.IsDomainName(domain)
	if !ok {
		return fmt.Errorf("invalid domain format: %s", domain)
	}

	return nil
}

// Domain lowercases and trims trailing dots before validation.
func Domain(domain string) (string, error) {
	normalized := strings.TrimSpace(strings.ToLower(domain))
	normalized = strings.TrimSuffix(normalized, ".")
	if err := IsValidDomain(normalized); err != nil {
		return "", err
	}
	return normalized, nil
}

// IsValidQType uses miekg/dns type map to avoid maintaining our own list.
func IsValidQType(qtype string) bool {
	_, ok := dns.StringToType[strings.ToUpper(qtype)]
	return ok
}

// QType validates qtype via dns.StringToType, defaults to A.
func QType(qtype string) (string, error) {
	if qtype == "" {
		return "A", nil
	}

	normalized := strings.ToUpper(qtype)
	if !IsValidQType(normalized) {
		return "", fmt.Errorf("invalid query type: %s (must be a valid DNS record type)", qtype)
	}

	return normalized, nil
}

// IPToReverseDNS delegates reverse DNS formatting to dns.ReverseAddr.
func IPToReverseDNS(ip string) (string, error) {
	rev, err := dns.ReverseAddr(ip)
	if err != nil {
		return "", err
	}
	// Trim trailing dot from dns.ReverseAddr output
	return strings.TrimSuffix(rev, "."), nil
}

package normalize

import "testing"

func TestTarget(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
		ok   bool
	}{
		// AdGuard dnsproxy handles port and IPv6 bracketing automatically
		{"plain IPv4", "8.8.8.8", "udp://8.8.8.8", true},
		{"udp explicit", "udp://9.9.9.9:53", "udp://9.9.9.9:53", true},
		{"https default path", "https://dns.google", "https://dns.google/dns-query", true},
		{"https with path", "https://dns.google/dns-query", "https://dns.google/dns-query", true},
		{"quic default port", "quic://dns.adguard.com", "quic://dns.adguard.com", true},
		{"ipv6 plain", "2001:4860:4860::8888", "udp://2001:4860:4860::8888", true},
		{"ipv6 with brackets and port", "[2001:4860:4860::8888]:853", "udp://[2001:4860:4860::8888]:853", true},
		{"unsupported scheme", "ftp://example.com", "", false},
		{"empty input", "", "", false},
		{"tcp scheme", "tcp://8.8.8.8:53", "tcp://8.8.8.8:53", true},
		{"tls with hostname", "tls://dns.quad9.net", "tls://dns.quad9.net", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Target(tt.in)
			if tt.ok && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.ok && err == nil {
				t.Fatalf("expected error, got none")
			}
			if tt.ok && got != tt.want {
				t.Fatalf("got %q want %q", got, tt.want)
			}
		})
	}
}

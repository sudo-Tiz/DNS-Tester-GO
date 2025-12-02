package models

import (
	"testing"
)

func TestDNSServerValidate(t *testing.T) {
	tests := []struct {
		target  string
		wantErr bool
	}{
		{"udp://9.9.9.9:53", false},
		{"https://dns.quad9.net", false},
		{"http://invalid", true},
		{"", true},
	}

	for _, tt := range tests {
		srv := DNSServer{Target: tt.target}
		err := srv.Validate()
		if (err != nil) != tt.wantErr {
			t.Errorf("Validate(%q) error = %v, wantErr %v", tt.target, err, tt.wantErr)
		}
	}
}

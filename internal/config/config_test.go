package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	yamlContent := `
servers:
  - ip: "9.9.9.9"
    port: 53
    tags: ["test"]
    protocols:
      - udp
`
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			t.Logf("Warning: failed to remove temp file: %v", err)
		}
	}()

	if _, err := tmpfile.Write([]byte(yamlContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(cfg.Servers) != 1 || cfg.Servers[0].IP != "9.9.9.9" {
		t.Error("Config not loaded correctly")
	}
}

func TestGetDNSTargets(t *testing.T) {
	cfg := &APIConfig{
		Servers: []DNSServer{
			{IP: "9.9.9.9", Port: 53, Services: []ServiceType{ServiceDo53UDP}},
		},
	}

	targets := cfg.GetDNSTargets()
	if len(targets) == 0 {
		t.Error("Expected at least one target")
	}
}

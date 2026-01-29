package status

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type statusVar struct {
	VarType  string `json:"vartype"`
	VarID    string `json:"varid"`
	VarValue string `json:"varvalue"`
}

func loadFixture(t *testing.T, path string) *os.File {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	t.Cleanup(func() { _ = f.Close() })
	return f
}

func countEntries(t *testing.T, path string) int {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var vars []statusVar
	if err := json.Unmarshal(data, &vars); err != nil {
		t.Fatalf("decode fixture: %v", err)
	}
	return len(vars)
}

func TestParse_NotProvisionedFixture(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "status_not_provisioned.json")
	f := loadFixture(t, path)

	got, err := Parse(f)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(got) != countEntries(t, path) {
		t.Fatalf("map length mismatch: got %d", len(got))
	}

	if got["device_name"] != "Glasfaser-Modem 2" {
		t.Fatalf("device_name mismatch: %q", got["device_name"])
	}
	if got["txpower"] != "--" {
		t.Fatalf("txpower mismatch: %q", got["txpower"])
	}
	if got["rxpower"] != "-16.42" {
		t.Fatalf("rxpower mismatch: %q", got["rxpower"])
	}
	if got["title"] != "Glasfaser-Modem 2 Konfigurationsprogramm" {
		t.Fatalf("title mismatch: %q", got["title"])
	}
}

func TestParse_ProvisionedFixture(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "status_provisioned.json")
	f := loadFixture(t, path)

	got, err := Parse(f)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(got) != countEntries(t, path) {
		t.Fatalf("map length mismatch: got %d", len(got))
	}

	if got["serial_number"] == "" {
		t.Fatalf("serial_number should not be empty")
	}
	if got["txpower"] == "" || got["rxpower"] == "" {
		t.Fatalf("optical power values should not be empty")
	}
}

func TestParse_DuplicateVarID_LastWins(t *testing.T) {
	input := strings.NewReader(`[{"vartype":"value","varid":"rxpower","varvalue":"-10.0"},{"vartype":"value","varid":"rxpower","varvalue":"-9.5"}]`)

	got, err := Parse(input)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got["rxpower"] != "-9.5" {
		t.Fatalf("duplicate varid should use last value, got %q", got["rxpower"])
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	input := strings.NewReader(`not json`)
	if _, err := Parse(input); err == nil {
		t.Fatalf("expected error for invalid JSON")
	}
}

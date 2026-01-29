package status

import (
	"math"
	"testing"
)

func TestParseFloat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantNaN bool
		want    float64
	}{
		{"valid positive", "3.14", false, 3.14},
		{"valid negative", "-16.46", false, -16.46},
		{"valid integer", "42", false, 42.0},
		{"empty string", "", true, 0},
		{"double dash", "--", true, 0},
		{"N/A", "N/A", true, 0},
		{"whitespace only", "   ", true, 0},
		{"leading whitespace", "  -10.5", false, -10.5},
		{"trailing whitespace", "-10.5  ", false, -10.5},
		{"both whitespace", "  -10.5  ", false, -10.5},
		{"zero", "0", false, 0},
		{"negative zero", "-0", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseFloat(tt.input)
			if tt.wantNaN {
				if !math.IsNaN(got) {
					t.Errorf("ParseFloat(%q) = %v, want NaN", tt.input, got)
				}
			} else {
				if math.IsNaN(got) {
					t.Errorf("ParseFloat(%q) = NaN, want %v", tt.input, tt.want)
				} else if got != tt.want {
					t.Errorf("ParseFloat(%q) = %v, want %v", tt.input, got, tt.want)
				}
			}
		})
	}
}

func TestParseUint(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  uint64
	}{
		{"valid positive", "12345", 12345},
		{"zero", "0", 0},
		{"large number", "18446744073709551615", 18446744073709551615},
		{"empty string", "", 0},
		{"invalid string", "abc", 0},
		{"negative number", "-1", 0},
		{"float value", "3.14", 0},
		{"whitespace", "  123  ", 123},
		{"double dash", "--", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseUint(tt.input)
			if got != tt.want {
				t.Errorf("ParseUint(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

package utils

import "testing"

func TestMaskAccountNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1234567890123456", "XXXXXXXXXXXX3456"},
		{"1234", "1234"},
		{"12345", "X2345"},
		{"", ""},
	}

	for _, tt := range tests {
		result := MaskAccountNumber(tt.input)
		if result != tt.expected {
			t.Errorf("MaskAccountNumber(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

package services

import (
	"testing"
	"unicode/utf8"
)

func TestDefineTokenName(t *testing.T) {
	tests := []struct {
		reName   string
		reID     int
		expected string
	}{
		{"Sample Real Estate", 123, "Tooken Sample Real Estate #123"},
		{"Another Real Estate", 456, "Tooken Another Real Estate #456"},
		{"Another Real Estate with a very long name and accents like this ééççç)ààaa and éééé more text to test truncation", 456, "Tooken Another Real Estate with a very long name and accents like this ééççç)ààaa and é #456"},
		{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaéééééééééééééééééééé", 456, "Tooken aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaééé #456"},
		{"👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍❤👍👍", 456, "Tooken 👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍❤ #456"},
		{"   ", 456, "Tooken Unnamed #456"},
	}
	for _, tt := range tests {
		result := defineTokenName(tt.reName, tt.reID)
		if !utf8.ValidString(result) {
			t.Errorf("defineTokenName(%q, %d) produced an invalid UTF-8 string: %q", tt.reName, tt.reID, result)
		}
		if result != tt.expected {
			t.Errorf("defineTokenName(%q, %d) = %q; want %q", tt.reName, tt.reID, result, tt.expected)
		}
	}
}

func TestDefineSymbol(t *testing.T) {
	realEstateID := 123
	expected := "TKN123"
	result := defineSymbol(realEstateID)
	if result != expected {
		t.Errorf("defineSymbol(%d) = %q; want %q", realEstateID, result, expected)
	}
}

func TestDefineSalt(t *testing.T) {
	realEstateID := 123
	expected := "re-123"
	result := defineSalt(realEstateID)
	if result != expected {
		t.Errorf("defineSalt(%d) = %q; want %q", realEstateID, result, expected)
	}
}

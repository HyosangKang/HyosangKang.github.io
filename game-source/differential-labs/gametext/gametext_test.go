package gametext

import (
	"regexp"
	"testing"
)

func TestCleanUsesOnlyGameFontCharacters(t *testing.T) {
	cleaned := Clean("V(t) = sin(t), i1: -0.25 A")
	if !regexp.MustCompile(`^[A-Z0-9 ]+$`).MatchString(cleaned) {
		t.Fatalf("Clean returned unsupported characters: %q", cleaned)
	}
	if cleaned != "V T SIN T I1 0 25 A" {
		t.Fatalf("Clean returned %q", cleaned)
	}
}

func TestScaledValuesAvoidSymbols(t *testing.T) {
	tests := map[string]string{
		Value(1.25, 0.05):      "125",
		Value(-0.8, 0.1):       "NEG 8",
		Signed(0.125, 1000):    "POS 125",
		Signed(-0.125, 1000):   "NEG 125",
		Magnitude(2.449, 1000): "2449",
	}
	for result, want := range tests {
		if result != want {
			t.Fatalf("scaled value = %q, want %q", result, want)
		}
	}
}

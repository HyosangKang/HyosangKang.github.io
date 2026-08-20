package gametext

import (
	"fmt"
	"math"
	"strings"
)

// Clean limits game-font output to uppercase ASCII letters, digits, and spaces.
func Clean(value string) string {
	var cleaned strings.Builder
	for _, character := range strings.ToUpper(value) {
		if character >= 'A' && character <= 'Z' || character >= '0' && character <= '9' {
			cleaned.WriteRune(character)
			continue
		}
		cleaned.WriteByte(' ')
	}
	return strings.Join(strings.Fields(cleaned.String()), " ")
}

// Value displays decimal controls as scaled integers so the game font never
// needs a decimal point or minus symbol.
func Value(value float64, step float64) string {
	scale := 1
	if step < 0.1 {
		scale = 100
	} else if step < 1 {
		scale = 10
	}
	return scaled(value, scale, false)
}

// Signed displays a signed measurement using POS, NEG, or ZERO.
func Signed(value float64, scale int) string {
	return scaled(value, scale, true)
}

// Magnitude displays a nonnegative measurement as a scaled integer.
func Magnitude(value float64, scale int) string {
	return fmt.Sprintf("%d", int(math.Round(math.Abs(value)*float64(scale))))
}

func scaled(value float64, scale int, showPositive bool) string {
	magnitude := int(math.Round(math.Abs(value) * float64(scale)))
	if magnitude == 0 {
		return "ZERO"
	}
	if value < 0 {
		return fmt.Sprintf("NEG %d", magnitude)
	}
	if showPositive {
		return fmt.Sprintf("POS %d", magnitude)
	}
	return fmt.Sprintf("%d", magnitude)
}

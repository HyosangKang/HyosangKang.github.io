package arcade

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type calculatorLab struct {
	display     string
	accumulator float64
	pending     string
	fresh       bool
	message     string
	buttons     []Button
}

func newCalculatorLab() *calculatorLab {
	lab := &calculatorLab{}
	lab.reset()
	labels := [][]string{
		{"C", "+/-", "/", "BACK"},
		{"7", "8", "9", "x"},
		{"4", "5", "6", "-"},
		{"1", "2", "3", "+"},
		{"0", ".", "=", "="},
	}
	for row, items := range labels {
		for column, label := range items {
			if row == 4 && column == 3 {
				continue
			}
			buttonColor := White
			if strings.Contains("+-x/=", label) {
				buttonColor = Gold
			}
			if label == "C" || label == "BACK" {
				buttonColor = Red
			}
			width := float32(92)
			if row == 4 && column == 2 {
				width = 196
			}
			lab.buttons = append(lab.buttons, Button{
				X:     194 + float32(column)*104,
				Y:     156 + float32(row)*58,
				W:     width,
				H:     48,
				Label: label,
				Color: buttonColor,
			})
		}
	}
	return lab
}

func (lab *calculatorLab) reset() {
	lab.display = "0"
	lab.accumulator = 0
	lab.pending = ""
	lab.fresh = true
	lab.message = "EACH  KEY  CHANGES  ONE  STORED  NUMBER"
}

func (lab *calculatorLab) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		for _, button := range lab.buttons {
			if button.Contains(x, y) {
				lab.press(button.Label)
				break
			}
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		lab.reset()
	}
	return nil
}

func (lab *calculatorLab) press(label string) {
	switch label {
	case "C":
		lab.reset()
	case "BACK":
		if !lab.fresh && len(lab.display) > 1 {
			lab.display = lab.display[:len(lab.display)-1]
		} else {
			lab.display = "0"
			lab.fresh = true
		}
	case "+/-":
		if strings.HasPrefix(lab.display, "-") {
			lab.display = strings.TrimPrefix(lab.display, "-")
		} else if lab.display != "0" {
			lab.display = "-" + lab.display
		}
	case "+", "-", "x", "/":
		lab.applyPending()
		lab.pending = label
		lab.fresh = true
		lab.message = "STORED  " + formatNumber(lab.accumulator) + "    NEXT  " + label
	case "=":
		lab.applyPending()
		lab.pending = ""
		lab.fresh = true
		lab.message = "RESULT  STORED  IN  THE  DISPLAY"
	case ".":
		if lab.fresh {
			lab.display = "0."
			lab.fresh = false
		} else if !strings.Contains(lab.display, ".") {
			lab.display += "."
		}
	default:
		if len(lab.display) >= 10 && !lab.fresh {
			return
		}
		if lab.fresh || lab.display == "0" {
			lab.display = label
			lab.fresh = false
		} else {
			lab.display += label
		}
	}
}

func (lab *calculatorLab) applyPending() {
	value, err := strconv.ParseFloat(lab.display, 64)
	if err != nil {
		return
	}
	if lab.pending == "" {
		lab.accumulator = value
		return
	}
	switch lab.pending {
	case "+":
		lab.accumulator += value
	case "-":
		lab.accumulator -= value
	case "x":
		lab.accumulator *= value
	case "/":
		if value == 0 {
			lab.display = "ERROR"
			lab.pending = ""
			lab.fresh = true
			lab.message = "DIVISION  BY  ZERO"
			return
		}
		lab.accumulator /= value
	}
	lab.display = formatNumber(lab.accumulator)
}

func formatNumber(value float64) string {
	if math.Abs(value) < 1e-12 {
		value = 0
	}
	result := strconv.FormatFloat(value, 'g', 10, 64)
	if len(result) > 12 {
		result = fmt.Sprintf("%.6g", value)
	}
	return result
}

func (lab *calculatorLab) Draw(screen *ebiten.Image) {
	DrawText(screen, "ELEMENTARY  CALCULATOR", 28, Width/2, 22, White, text.AlignCenter)
	vector.DrawFilledRect(screen, 194, 72, 404, 64, DarkBlue, false)
	vector.StrokeRect(screen, 194, 72, 404, 64, 2, Cyan, false)
	DrawText(screen, lab.display, 36, 580, 83, Cyan, text.AlignEnd)
	for _, button := range lab.buttons {
		button.Draw(screen)
	}
	DrawText(screen, lab.message, 15, Width/2, 454, Blue, text.AlignCenter)
}

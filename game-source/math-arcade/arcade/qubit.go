package arcade

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type qubitLab struct {
	alpha   float64
	beta    float64
	message string
	buttons []Button
}

func newQubitLab() *qubitLab {
	return &qubitLab{
		alpha:   1,
		message: "START  IN  STATE  0",
		buttons: []Button{
			{X: 90, Y: 394, W: 100, H: 48, Label: "X", Color: Red},
			{X: 204, Y: 394, W: 100, H: 48, Label: "H", Color: Cyan},
			{X: 318, Y: 394, W: 100, H: 48, Label: "Z", Color: Purple},
			{X: 432, Y: 394, W: 162, H: 48, Label: "MEASURE", Color: Gold},
			{X: 608, Y: 394, W: 102, H: 48, Label: "RESET", Color: Green},
		},
	}
}

func (lab *qubitLab) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		for _, button := range lab.buttons {
			if !button.Contains(x, y) {
				continue
			}
			switch button.Label {
			case "X":
				lab.alpha, lab.beta = lab.beta, lab.alpha
				lab.message = "X  SWAPS  THE  TWO  AMPLITUDES"
			case "H":
				a, b := lab.alpha, lab.beta
				lab.alpha = (a + b) / math.Sqrt2
				lab.beta = (a - b) / math.Sqrt2
				lab.message = "H  MAKES  OR  REMOVES  A  SUPERPOSITION"
			case "Z":
				lab.beta = -lab.beta
				lab.message = "Z  FLIPS  THE  PHASE  OF  STATE  1"
			case "MEASURE":
				if rand.Float64() < lab.alpha*lab.alpha {
					lab.alpha, lab.beta = 1, 0
					lab.message = "MEASURED  0"
				} else {
					lab.alpha, lab.beta = 0, 1
					lab.message = "MEASURED  1"
				}
			case "RESET":
				lab.alpha, lab.beta = 1, 0
				lab.message = "START  IN  STATE  0"
			}
			break
		}
	}
	return nil
}

func (lab *qubitLab) Draw(screen *ebiten.Image) {
	DrawText(screen, "ONE  QUBIT", 28, Width/2, 20, White, text.AlignCenter)
	p0 := lab.alpha * lab.alpha
	p1 := lab.beta * lab.beta

	drawProbabilityBar(screen, 150, 118, "|0>", p0, Cyan)
	drawProbabilityBar(screen, 150, 226, "|1>", p1, Gold)

	centerX, centerY := float32(618), float32(212)
	vector.StrokeCircle(screen, centerX, centerY, 108, 2, Blue, false)
	vector.StrokeLine(screen, centerX-108, centerY, centerX+108, centerY, 1, GridBlue, false)
	vector.StrokeLine(screen, centerX, centerY-108, centerX, centerY+108, 1, GridBlue, false)
	angle := 2 * math.Atan2(lab.beta, lab.alpha)
	pointX := centerX + 96*float32(math.Sin(angle))
	pointY := centerY - 96*float32(math.Cos(angle))
	vector.StrokeLine(screen, centerX, centerY, pointX, pointY, 3, White, false)
	vector.DrawFilledCircle(screen, pointX, pointY, 6, Red, false)

	DrawText(screen, fmt.Sprintf("P0  %.2f", p0), 18, 150, 88, Cyan, text.AlignStart)
	DrawText(screen, fmt.Sprintf("P1  %.2f", p1), 18, 150, 196, Gold, text.AlignStart)
	for _, button := range lab.buttons {
		button.Draw(screen)
	}
	DrawText(screen, lab.message, 15, Width/2, 456, Blue, text.AlignCenter)
}

func drawProbabilityBar(
	screen *ebiten.Image,
	x float32,
	y float32,
	label string,
	probability float64,
	itemColor color.Color,
) {
	DrawText(screen, label, 24, float64(x-58), float64(y+14), itemColor, text.AlignCenter)
	vector.StrokeRect(screen, x, y, 300, 42, 2, itemColor, false)
	vector.DrawFilledRect(screen, x+4, y+4, float32(probability)*292, 34, itemColor, false)
}

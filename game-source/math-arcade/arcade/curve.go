package arcade

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type curveLab struct {
	segments int
	minus    Button
	plus     Button
}

func newCurveLab() *curveLab {
	return &curveLab{
		segments: 8,
		minus:    Button{X: 276, Y: 402, W: 92, H: 48, Label: "-", Color: Red},
		plus:     Button{X: 432, Y: 402, W: 92, H: 48, Label: "+", Color: Green},
	}
}

func curvePoint(t float64) (float64, float64) {
	x := 72 + 656*t
	y := 238 - 82*math.Sin(2*math.Pi*t) - 24*math.Sin(5*math.Pi*t)
	return x, y
}

func (lab *curveLab) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if lab.minus.Contains(x, y) && lab.segments > 2 {
			lab.segments /= 2
		}
		if lab.plus.Contains(x, y) && lab.segments < 128 {
			lab.segments *= 2
		}
	}
	return nil
}

func (lab *curveLab) estimate() float64 {
	length := 0.0
	x0, y0 := curvePoint(0)
	for index := 1; index <= lab.segments; index++ {
		x1, y1 := curvePoint(float64(index) / float64(lab.segments))
		length += math.Hypot(x1-x0, y1-y0)
		x0, y0 = x1, y1
	}
	return length
}

func referenceCurveLength() float64 {
	length := 0.0
	x0, y0 := curvePoint(0)
	const samples = 4096
	for index := 1; index <= samples; index++ {
		x1, y1 := curvePoint(float64(index) / samples)
		length += math.Hypot(x1-x0, y1-y0)
		x0, y0 = x1, y1
	}
	return length
}

func (lab *curveLab) Draw(screen *ebiten.Image) {
	DrawText(screen, "CURVE  LENGTH", 28, Width/2, 20, White, text.AlignCenter)
	DrawGrid(screen, 80, 64)

	previousX, previousY := curvePoint(0)
	for index := 1; index <= 600; index++ {
		x, y := curvePoint(float64(index) / 600)
		vector.StrokeLine(
			screen,
			float32(previousX),
			float32(previousY),
			float32(x),
			float32(y),
			2,
			Cyan,
			false,
		)
		previousX, previousY = x, y
	}

	previousX, previousY = curvePoint(0)
	for index := 1; index <= lab.segments; index++ {
		x, y := curvePoint(float64(index) / float64(lab.segments))
		vector.StrokeLine(
			screen,
			float32(previousX),
			float32(previousY),
			float32(x),
			float32(y),
			2,
			Gold,
			false,
		)
		vector.DrawFilledCircle(screen, float32(x), float32(y), 4, Gold, false)
		previousX, previousY = x, y
	}

	estimate := lab.estimate()
	reference := referenceCurveLength()
	DrawText(screen, fmt.Sprintf("SEGMENTS  %03d", lab.segments), 18, 28, 340, White, text.AlignStart)
	DrawText(screen, fmt.Sprintf("POLYGON  %.1f", estimate), 18, 300, 340, Gold, text.AlignStart)
	DrawText(screen, fmt.Sprintf("CURVE  %.1f", reference), 18, 566, 340, Cyan, text.AlignStart)
	lab.minus.Draw(screen)
	lab.plus.Draw(screen)
	DrawText(screen, "MORE  SEGMENTS  FOLLOW  THE  CURVE  MORE  CLOSELY", 15, Width/2, 458, Blue, text.AlignCenter)
}

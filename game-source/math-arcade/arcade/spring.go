package arcade

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type springLab struct {
	x       float64
	v       float64
	time    float64
	damping float64
	force   float64
	trace   []float64
	buttons []Button
}

func newSpringLab() *springLab {
	return &springLab{
		x:       1.3,
		damping: 0.12,
		force:   0.0,
		buttons: []Button{
			{X: 74, Y: 405, W: 104, H: 46, Label: "DAMP -", Color: Red},
			{X: 188, Y: 405, W: 104, H: 46, Label: "DAMP +", Color: Green},
			{X: 508, Y: 405, W: 104, H: 46, Label: "FORCE -", Color: Red},
			{X: 622, Y: 405, W: 104, H: 46, Label: "FORCE +", Color: Green},
		},
	}
}

func (lab *springLab) Update() error {
	const dt = 0.035
	acceleration := -1.15*lab.x - lab.damping*lab.v + lab.force*math.Sin(1.45*lab.time)
	lab.v += acceleration * dt
	lab.x += lab.v * dt
	lab.time += dt
	lab.trace = append(lab.trace, lab.x)
	if len(lab.trace) > 430 {
		lab.trace = lab.trace[len(lab.trace)-430:]
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		for _, button := range lab.buttons {
			if !button.Contains(x, y) {
				continue
			}
			switch button.Label {
			case "DAMP -":
				lab.damping = math.Max(0, lab.damping-0.05)
			case "DAMP +":
				lab.damping = math.Min(1, lab.damping+0.05)
			case "FORCE -":
				lab.force = math.Max(0, lab.force-0.2)
			case "FORCE +":
				lab.force = math.Min(2, lab.force+0.2)
			}
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		lab.x = 1.3
		lab.v = 0
		lab.time = 0
		lab.trace = nil
	}
	return nil
}

func (lab *springLab) Draw(screen *ebiten.Image) {
	DrawText(screen, "SPRING  OSCILLATION", 28, Width/2, 18, White, text.AlignCenter)
	vector.StrokeLine(screen, 92, 118, 92, 246, 4, Blue, false)
	equilibrium := 430.0
	massX := equilibrium + 120*lab.x
	points := 18
	previousX := 92.0
	previousY := 182.0
	for index := 1; index <= points; index++ {
		x := 92 + (massX-92)*float64(index)/float64(points)
		y := 182.0
		if index < points {
			if index%2 == 0 {
				y -= 24
			} else {
				y += 24
			}
		}
		vector.StrokeLine(screen, float32(previousX), float32(previousY), float32(x), float32(y), 2, Cyan, false)
		previousX, previousY = x, y
	}
	vector.StrokeLine(screen, float32(equilibrium), 110, float32(equilibrium), 254, 1, GridBlue, false)
	vector.DrawFilledRect(screen, float32(massX-24), 154, 48, 56, Gold, false)

	vector.StrokeRect(screen, 74, 282, 652, 92, 1, Blue, false)
	for index := 1; index < len(lab.trace); index++ {
		x0 := 75 + float32(index-1)*650/430
		x1 := 75 + float32(index)*650/430
		y0 := 328 - float32(lab.trace[index-1])*26
		y1 := 328 - float32(lab.trace[index])*26
		vector.StrokeLine(screen, x0, y0, x1, y1, 1, White, false)
	}

	DrawText(screen, fmt.Sprintf("DAMP  %.2f", lab.damping), 16, 76, 254, Cyan, text.AlignStart)
	DrawText(screen, fmt.Sprintf("FORCE  %.1f", lab.force), 16, 724, 254, Gold, text.AlignEnd)
	for _, button := range lab.buttons {
		button.Draw(screen)
	}
	DrawText(screen, "R  RESET", 15, Width/2, 420, Blue, text.AlignCenter)
}

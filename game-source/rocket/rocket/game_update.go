package rocket

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	VelocityThreshold = 1
)

func (g *Game) Update() error {
	if g.Scene == SceneNil {
		g.Scene = SceneStart
	}
	switch g.Scene {
	case SceneStart:
		g.Start()
	case ScenePlay:
		g.Play()
	}
	return nil
}

func (g *Game) Start() {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		g.Scene = ScenePlay
	}
}

func (g *Game) Play() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.Clear()
		g.Scene = SceneStart
	}
	switch g.Stage {
	case RocketStopped:
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			g.Stage = RocketLaunching
		}
	case RocketLaunching:
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			g.RocketVelocity[0] = RocketXY[0] - float64(x)
			g.RocketVelocity[1] = RocketXY[1] - float64(y)
			g.Stage = RocketThrust
		}
	case RocketThrust:
		if AbsMax(g.RocketVelocity[0], g.RocketVelocity[1]) < VelocityThreshold {
			g.Stage = RocketStopped
		}
		if Distance(RocketXY, PlanetCenter) < PlanetRadius {
			g.Stage = RocketCrashed
		}
		g.UpdateRocket()
	}
}

func Distance(a, b [2]float64) float64 {
	dx := a[0] - b[0]
	dy := a[1] - b[1]
	return math.Sqrt(dx*dx + dy*dy)
}

const (
	DistanceScale = .1
	DraggingRatio = .05
	Gravitational = 10
)

func (g *Game) UpdateRocket() {
	r := Distance(RocketXY, PlanetCenter)
	for i := 0; i < 2; i++ {
		g.RocketVelocity[i] -= DraggingRatio * g.RocketVelocity[i]
		g.RocketVelocity[i] += Gravitational * (PlanetCenter[i] - RocketXY[i]) / (r * r)
		RocketXY[i] += g.RocketVelocity[i] * DistanceScale
	}
}

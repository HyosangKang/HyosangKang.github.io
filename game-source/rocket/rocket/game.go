package rocket

import (
	"math"
	"math/rand"
)

const (
	RocketStopped int = iota
	RocketLaunching
	RocketThrust
	RocketCrashed
)

const (
	Width  = 600
	Height = 600
)

var (
	PlanetCenter = [2]float64{Width / 2, Height / 2}
)

type Scene int

const (
	SceneNil Scene = iota
	SceneStart
	ScenePlay
)

type Game struct {
	Scene
	Stage          int
	RocketVelocity [2]float64
}

var (
	RocketXY [2]float64
)

func init() {
	r := 100 + 100*rand.Float64()
	t := 2 * math.Pi * rand.Float64()
	RocketXY = [2]float64{Width/2 + r*math.Cos(t), Height/2 + r*math.Sin(t)}
}

func (g *Game) Clear() {
	g.Stage = RocketStopped
	r := 100 + 100*rand.Float64()
	t := 2 * math.Pi * rand.Float64()
	RocketXY = [2]float64{Width/2 + r*math.Cos(t), Height/2 + r*math.Sin(t)}
}

func (g *Game) Layout(int, int) (int, int) {
	return Width, Height
}

package virus

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Agent struct {
	X, Y          float64
	State         int
	MovingTimer   int
	MovingTimeOut int
	ExposedTimer  int
	InfectedTimer int
	Follow        [2]float64
}

const (
	Susceptible int = iota
	Exposed
	Infected
	Recovered
)

const (
	ExposedTimeOut  = 120
	InfectedTimeOut = 120
)

func NewAgent() *Agent {
	return &Agent{
		X:             rand.Float64() * Width,
		Y:             rand.Float64() * Height,
		MovingTimeOut: rand.Intn(100) + 100,
		Follow: [2]float64{
			rand.Float64() * Width,
			rand.Float64() * Height},
	}
}

const AgentTimeOut = 150

func (a *Agent) Update(agents []*Agent) {
	switch a.State {
	case Susceptible:
		a.Move()
		for _, b := range agents {
			if b.State == Exposed || b.State == Infected {
				dx := a.X - b.X
				dy := a.Y - b.Y
				d := math.Sqrt(dx*dx + dy*dy)
				if d < ExposeRadius {
					a.State = Exposed
				}
			}
		}
	case Exposed:
		a.Move()
		a.ExposedTimer++
		if a.ExposedTimer >= ExposedTimeOut {
			a.ExposedTimer = 0
			a.State = Infected
		}
	case Infected:
		a.InfectedTimer++
		if a.InfectedTimer >= InfectedTimeOut {
			a.InfectedTimer = 0
			p := rand.Float64()
			if p < .5 {
				a.State = Recovered
			} else {
				a.State = Susceptible
			}
		}
	case Recovered:
		a.Move()
	}
}

func (a *Agent) ResetFollow() {
	a.Follow = [2]float64{
		rand.Float64() * Width,
		rand.Float64() * Height,
	}
	a.MovingTimer = 0
	a.MovingTimeOut = rand.Intn(100) + 100
}

const (
	AgentSpeed   = 1
	ExposeRadius = 5
)

func (a *Agent) Move() {
	if a.MovingTimer < a.MovingTimeOut {
		a.MovingTimer++
	} else {
		a.ResetFollow()
	}
	dx := a.Follow[0] - a.X
	dy := a.Follow[1] - a.Y
	d := math.Sqrt(dx*dx + dy*dy)
	if d < AgentSpeed {
		a.ResetFollow()
	}
	a.X += AgentSpeed * dx / d
	a.Y += AgentSpeed * dy / d
}

var (
	ColorSusceptible = color.RGBA{60, 179, 113, 255}
	ColorExposed     = color.RGBA{255, 165, 0, 255}
	ColorInfected    = color.RGBA{255, 0, 0, 255}
	ColorRecovered   = color.RGBA{238, 130, 238, 255}
)

func (a *Agent) Draw(screen *ebiten.Image) {
	switch a.State {
	case Susceptible:
		DrawCircle(screen, a.X, a.Y, 5, ColorSusceptible)
	case Exposed:
		DrawCircle(screen, a.X, a.Y, 5, ColorExposed)
	case Infected:
		DrawCircle(screen, a.X, a.Y, 5, ColorInfected)
	case Recovered:
		DrawCircle(screen, a.X, a.Y, 5, ColorRecovered)
	}
}

const (
	Nsub = 100
)

func DrawCircle(screen *ebiten.Image, x, y, r float64, color color.Color) {
	t := Linspace(0, 2*math.Pi, Nsub)
	for i := 0; i < Nsub; i++ {
		x0, y0 := x+r*math.Cos(t[i]), y+r*math.Sin(t[i])
		x1, y1 := x+r*math.Cos(t[(i+1)%Nsub]), y+r*math.Sin(t[(i+1)%Nsub])
		DrawLine(screen, x0, y0, x1, y1, color)
	}
}

func Linspace(a, b float64, n int) []float64 {
	if n < 2 {
		return nil
	}
	t := make([]float64, n)
	for i := 0; i < n; i++ {
		t[i] = a + (b-a)*float64(i)/float64(n-1)
	}
	return t
}

func DrawLine(screen *ebiten.Image, x0, y0, x1, y1 float64, color color.Color) {
	dx := x1 - x0
	dy := y1 - y0
	if dx == 0 && dy == 0 {
		screen.Set(int(x0), int(y0), color)
		return
	}
	n := AbsMax(dx, dy)
	for i := 0; i <= int(n); i++ {
		x := x0 + float64(i)*dx/n
		y := y0 + float64(i)*dy/n
		screen.Set(int(x), int(y), color)
	}
}

func AbsMax(x, y float64) float64 {
	if x < 0 {
		x = -x
	}
	if y < 0 {
		y = -y
	}
	if x < y {
		return y
	}
	return x
}

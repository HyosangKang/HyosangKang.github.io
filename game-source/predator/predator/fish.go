package predator

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Fish struct {
	X, Y   float64
	Prey   *Coral
	Status int
	Timer  int
}

func NewFish() *Fish {
	return &Fish{
		X: rand.Float64() * Width,
		Y: rand.Float64() * Height,
	}
}

const (
	FishSpeed       = 1.5
	FishRadius      = 5
	FishEatingSpeed = .2
	FishEatingTime  = 30
)

const (
	FishHunting int = iota
	FishEating
	FishDead
)

func (f *Fish) Update(corals []*Coral) {
	switch f.Status {
	case FishDead:
		return
	case FishHunting:
		if f.Prey == nil || f.Prey.Status == CoralDead {
			f.Prey = FindAliveCoral(corals)
			if f.Prey == nil {
				f.X += FishSpeed * (rand.Float64()*2 - 1)
				f.Y += FishSpeed * (rand.Float64()*2 - 1)
				return
			}
		} else {
			dx := f.Prey.X - f.X
			dy := f.Prey.Y - f.Y
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist < f.Prey.Radius {
				f.Status = FishEating
			} else {
				f.X += FishSpeed * dx / dist * FishSpeed
				f.Y += FishSpeed * dy / dist * FishSpeed
			}
		}
	case FishEating:
		if f.Timer >= FishEatingTime {
			f.Status = FishHunting
			f.Timer = 0
			return
		}
		f.Timer++
		f.Prey.Radius -= FishEatingSpeed
	}
}

func FindAliveCoral(corals []*Coral) *Coral {
	var alive []*Coral
	for _, c := range corals {
		if c.Status == CoralAlive {
			alive = append(alive, c)
		}
	}
	if len(alive) == 0 {
		return nil
	}
	return alive[rand.Intn(len(alive))]
}

func (f *Fish) Draw(screen *ebiten.Image) {
	if f.Status != FishDead {
		DrawCircle(screen, f.X, f.Y, FishRadius, color.RGBA{60, 179, 113, 255})
	}

}

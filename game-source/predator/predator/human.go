package predator

import (
	"bytes"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Human struct {
	X, Y   float64
	Status int
	Timer  int
	Prey   *Fish
}

func NewHuman() *Human {
	return &Human{
		X: rand.Float64() * Width,
		Y: rand.Float64() * Height,
	}
}

const (
	HumanEating int = iota
	HumanSleeping
	HumanHunting
)

const (
	SleepingTime = 240
	EatingTime   = 120
)

const (
	HumanRadius = 10
	HumanSpeed  = 1
)

func (h *Human) Update(fishes []*Fish) {
	switch h.Status {
	case HumanEating:
		h.Timer++
		if h.Timer >= EatingTime {
			h.Status = HumanSleeping
			h.Timer = 0
		}
	case HumanSleeping:
		h.Timer++
		if h.Timer >= SleepingTime {
			h.Status = HumanHunting
			h.Timer = 0
		}
	case HumanHunting:
		if h.Prey == nil || h.Prey.Status == FishDead {
			h.Prey = FindAliveFish(fishes)
			if h.Prey == nil {
				return
			}
		} else {
			dx := h.Prey.X - h.X
			dy := h.Prey.Y - h.Y
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist < FishRadius {
				h.Prey.Status = FishDead
				h.Status = HumanEating
				h.Timer = 0
			} else {
				h.X += (h.Prey.X - h.X) / dist * HumanSpeed
				h.Y += (h.Prey.Y - h.Y) / dist * HumanSpeed
			}
		}
	}
}

func FindAliveFish(fishes []*Fish) *Fish {
	var alive []*Fish
	for _, f := range fishes {
		if f.Status != FishDead {
			alive = append(alive, f)
		}
	}
	if len(alive) == 0 {
		return nil
	}
	return alive[rand.Intn(len(alive))]
}

func (h *Human) Draw(screen *ebiten.Image) {
	var status string
	switch h.Status {
	case HumanEating:
		status = "E"
	case HumanSleeping:
		status = "S"
	case HumanHunting:
		status = "H"
	}
	DrawCircle(screen, h.X, h.Y, HumanRadius, color.RGBA{244, 197, 66, 255})
	DrawText(screen, status, 15, h.X, h.Y)
}

func init() {
	TitleFontFace, _ = text.NewGoTextFaceSource(bytes.NewReader(ArcadeClassic_ttf))
}

func DrawText(screen *ebiten.Image, msg string, size float64, x, y float64) {
	text.Draw(screen, msg, &text.GoTextFace{
		Source: TitleFontFace,
		Size:   size,
	}, NewTextDrawOption(x, y))
}

func NewTextDrawOption(x, y float64) *text.DrawOptions {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(color.White)
	op.PrimaryAlign = text.AlignStart
	op.SecondaryAlign = text.AlignStart
	return op
}

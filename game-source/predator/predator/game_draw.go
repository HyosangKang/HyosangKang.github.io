package predator

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.scene {
	case sceneStart:
		g.FrontPage(screen)
	case scenePlay:
		g.PlayScene(screen)
	}
}

const (
	Title                = "Predator!"
	TitleFontSize        = 48
	TitleMessage         = "Press  Enter  to  start"
	TitleMessageFontSize = 24
)

var (
	InstructionMessage = []string{
		"Q  to  add  human",
		"W  to  add  fish",
		"E  to  add  coral",
		"Press  Esc  to  exit"}
	TitleFontFace *text.GoTextFaceSource

	//go:embed arcadeclassic.ttf
	ArcadeClassic_ttf []byte
)

func init() {
	TitleFontFace, _ = text.NewGoTextFaceSource(bytes.NewReader(ArcadeClassic_ttf))
}

func drawText(screen *ebiten.Image, msg string, size float64, x, y float64) {
	text.Draw(screen, msg, &text.GoTextFace{
		Source: TitleFontFace,
		Size:   size,
	}, NewTextDrawOption(x, y))
}

func (g *Game) FrontPage(screen *ebiten.Image) {
	drawText(screen, Title, TitleFontSize, Width/2-150, float64(Height)/2-150)
	drawText(screen, TitleMessage, TitleMessageFontSize, Width/2-150, float64(Height)/2)
	for i, msg := range InstructionMessage {
		drawText(screen, msg, TitleMessageFontSize, Width/2-150, Height/2+100+float64(i)*30)
	}
}

func (g *Game) PlayScene(screen *ebiten.Image) {
	for _, c := range g.Corals {
		c.Draw(screen)
	}
	for _, f := range g.Fishes {
		f.Draw(screen)
	}
	for _, h := range g.Humans {
		h.Draw(screen)
	}
	g.DrawCounterBars(screen)
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

func (g *Game) DrawCounterBars(screen *ebiten.Image) {
	// Draw human counter
	DrawBar(screen, 10, 10, len(g.Humans))
	DrawText(screen, fmt.Sprintf("Humans  %d", len(g.Humans)), 15, BarWidth+15, 10)

	// Draw fish counter
	var numFish int
	for _, f := range g.Fishes {
		if f.Status != FishDead {
			numFish++
		}
	}
	DrawBar(screen, 10, 10+BarHeight+BarSkip, numFish)
	DrawText(screen, fmt.Sprintf("Fishes  %d", numFish), 15, BarWidth+15, 10+BarHeight+BarSkip)

	// Draw coral counter
	var numCoral int
	for _, c := range g.Corals {
		if c.Status == CoralAlive {
			numCoral++
		}
	}
	DrawBar(screen, 10, 10+2*BarHeight+2*BarSkip, numCoral)
	DrawText(screen, fmt.Sprintf("Corals  %d", numCoral), 15, BarWidth+15, 10+2*BarHeight+2*BarSkip)
}

const (
	BarWidth  = 100
	BarHeight = 20
	BarSkip   = 10
	MaxNum    = 50
)

func DrawBar(screen *ebiten.Image, xoff, yoff, len int) {
	// Draw frames
	for i := 0; i < BarWidth; i++ {
		screen.Set(xoff+i, yoff, color.White)
		screen.Set(xoff+i, yoff+BarHeight-1, color.White)
	}
	for j := 0; j < BarHeight; j++ {
		screen.Set(xoff, yoff+j, color.White)
		screen.Set(xoff+BarWidth-1, yoff+j, color.White)
	}
	// Draw bars
	for i := 0; i < len; i++ {
		for j := 0; j < BarHeight; j++ {
			screen.Set(xoff+i, yoff+j, color.White)
		}
	}
}

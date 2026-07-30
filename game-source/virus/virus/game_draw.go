package virus

import (
	"bytes"
	_ "embed"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.Scene {
	case SceneStart:
		g.FrontPage(screen)
	case ScenePlay:
		g.PlayScene(screen)
	}
}

func (g *Game) FrontPage(screen *ebiten.Image) {
	drawText(screen, Title, TitleFontSize, Width/2-150, float64(Height)/2-150)
	drawText(screen, TitleMessage, TitleMessageFontSize, Width/2-150, float64(Height)/2)
	for i, msg := range InstructionMessage {
		drawText(screen, msg, TitleMessageFontSize, Width/2-150, Height/2+100+float64(i)*30)
	}
}
func (g *Game) PlayScene(screen *ebiten.Image) {
	for x := 100.0; x < Width; x += 100 {
		DrawLine(screen, x, 0, x, Height, color.RGBA{18, 52, 82, 255})
	}
	for y := 100.0; y < Height; y += 100 {
		DrawLine(screen, 0, y, Width, y, color.RGBA{18, 52, 82, 255})
	}
	for _, a := range g.Agents {
		a.Draw(screen)
	}
}

const (
	Title                = "Virus!"
	TitleFontSize        = 48
	TitleMessage         = "Press  Enter  to  start"
	TitleMessageFontSize = 24
)

var (
	InstructionMessage = []string{
		"Greens  are  healthy",
		"Oranges  are  exposed",
		"Reds  are  infected",
		"Purples  are  recovered(immune)",

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

func NewTextDrawOption(x, y float64) *text.DrawOptions {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(color.White)
	op.PrimaryAlign = text.AlignStart
	op.SecondaryAlign = text.AlignStart
	return op
}

package arcade

import (
	"bytes"
	_ "embed"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	Width  = 800
	Height = 480
)

var (
	Black    = color.RGBA{0, 0, 0, 255}
	White    = color.RGBA{242, 242, 232, 255}
	Blue     = color.RGBA{46, 126, 214, 255}
	DarkBlue = color.RGBA{4, 18, 38, 255}
	GridBlue = color.RGBA{18, 52, 82, 255}
	Cyan     = color.RGBA{77, 215, 229, 255}
	Gold     = color.RGBA{244, 197, 66, 255}
	Red      = color.RGBA{240, 76, 76, 255}
	Green    = color.RGBA{56, 190, 140, 255}
	Purple   = color.RGBA{162, 110, 230, 255}
)

var (
	//go:embed arcadeclassic.ttf
	fontBytes []byte
	fontFace  *text.GoTextFaceSource
)

func init() {
	fontFace, _ = text.NewGoTextFaceSource(bytes.NewReader(fontBytes))
}

type lab interface {
	Update() error
	Draw(*ebiten.Image)
}

type Game struct {
	current lab
}

func NewGame(mode string) *Game {
	var current lab
	switch mode {
	case "curve-length":
		current = newCurveLab()
	case "spring":
		current = newSpringLab()
	case "quantum-simulator":
		current = newQubitLab()
	default:
		current = newCalculatorLab()
	}
	return &Game{current: current}
}

func (g *Game) Update() error {
	return g.current.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(Black)
	g.current.Draw(screen)
}

func (g *Game) Layout(int, int) (int, int) {
	return Width, Height
}

type Button struct {
	X     float32
	Y     float32
	W     float32
	H     float32
	Label string
	Color color.Color
}

func (button Button) Contains(x, y int) bool {
	return float32(x) >= button.X &&
		float32(x) <= button.X+button.W &&
		float32(y) >= button.Y &&
		float32(y) <= button.Y+button.H
}

func (button Button) Draw(screen *ebiten.Image) {
	itemColor := button.Color
	if itemColor == nil {
		itemColor = Blue
	}
	vector.DrawFilledRect(screen, button.X, button.Y, button.W, button.H, DarkBlue, false)
	vector.StrokeRect(screen, button.X, button.Y, button.W, button.H, 2, itemColor, false)
	DrawText(
		screen,
		button.Label,
		20,
		float64(button.X+button.W/2),
		float64(button.Y+button.H/2-10),
		itemColor,
		text.AlignCenter,
	)
}

func DrawText(
	screen *ebiten.Image,
	message string,
	size float64,
	x float64,
	y float64,
	itemColor color.Color,
	align text.Align,
) {
	options := &text.DrawOptions{}
	options.GeoM.Translate(x, y)
	options.ColorScale.ScaleWithColor(itemColor)
	options.PrimaryAlign = align
	options.SecondaryAlign = text.AlignStart
	text.Draw(screen, message, &text.GoTextFace{Source: fontFace, Size: size}, options)
}

func DrawGrid(screen *ebiten.Image, spacing int, top int) {
	for x := spacing; x < Width; x += spacing {
		vector.StrokeLine(screen, float32(x), float32(top), float32(x), Height, 1, GridBlue, false)
	}
	for y := top; y < Height; y += spacing {
		vector.StrokeLine(screen, 0, float32(y), Width, float32(y), 1, GridBlue, false)
	}
}

package labs

import (
	"bytes"
	"differential-labs/gametext"
	_ "embed"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	Red      = color.RGBA{255, 94, 108, 255}
	Green    = color.RGBA{125, 226, 139, 255}
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
	if mode == "double-spring" {
		return &Game{current: NewSpringLab()}
	}
	return &Game{current: NewElectricLab()}
}

func (g *Game) Update() error {
	return g.current.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(Black)
	DrawBackgroundGrid(screen, 24)
	g.current.Draw(screen)
}

func (g *Game) Layout(int, int) (int, int) {
	return Width, Height
}

type Button struct {
	X, Y, W, H float32
	Label      string
	Color      color.Color
	Active     bool
}

func (b Button) Contains(x, y int) bool {
	return float32(x) >= b.X && float32(x) <= b.X+b.W &&
		float32(y) >= b.Y && float32(y) <= b.Y+b.H
}

func (b Button) Draw(screen *ebiten.Image, size float64) {
	itemColor := b.Color
	if itemColor == nil {
		itemColor = Blue
	}
	fill := DarkBlue
	if b.Active {
		fill = color.RGBA{16, 42, 63, 255}
	}
	vector.DrawFilledRect(screen, b.X, b.Y, b.W, b.H, fill, false)
	vector.StrokeRect(screen, b.X, b.Y, b.W, b.H, 1.5, itemColor, false)
	DrawText(screen, b.Label, size, float64(b.X+b.W/2), float64(b.Y+b.H/2)-size/2, itemColor, text.AlignCenter)
}

func DrawText(screen *ebiten.Image, message string, size, x, y float64, itemColor color.Color, align text.Align) {
	options := &text.DrawOptions{}
	options.GeoM.Translate(x, y)
	options.ColorScale.ScaleWithColor(itemColor)
	options.PrimaryAlign = align
	options.SecondaryAlign = text.AlignStart
	text.Draw(screen, gametext.Clean(message), &text.GoTextFace{Source: fontFace, Size: size}, options)
}

func DrawBackgroundGrid(screen *ebiten.Image, spacing int) {
	for x := spacing; x < Width; x += spacing {
		vector.StrokeLine(screen, float32(x), 0, float32(x), Height, 1, color.RGBA{10, 28, 49, 170}, false)
	}
	for y := spacing; y < Height; y += spacing {
		vector.StrokeLine(screen, 0, float32(y), Width, float32(y), 1, color.RGBA{10, 28, 49, 170}, false)
	}
}

func DrawPanel(screen *ebiten.Image, x, y, w, h float32) {
	vector.DrawFilledRect(screen, x, y, w, h, color.RGBA{3, 10, 20, 238}, false)
	vector.StrokeRect(screen, x, y, w, h, 1, Blue, false)
}

func drawLine(screen *ebiten.Image, x0, y0, x1, y1 float64, width float32, itemColor color.Color) {
	vector.StrokeLine(screen, float32(x0), float32(y0), float32(x1), float32(y1), width, itemColor, false)
}

func clamp(value, low, high float64) float64 {
	if value < low {
		return low
	}
	if value > high {
		return high
	}
	return value
}

func formatValue(value float64, step float64) string {
	return gametext.Value(value, step)
}

func justPressedPointer() (int, int, bool) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		return x, y, true
	}
	touchIDs := inpututil.AppendJustPressedTouchIDs(nil)
	if len(touchIDs) > 0 {
		x, y := ebiten.TouchPosition(touchIDs[0])
		return x, y, true
	}
	return 0, 0, false
}

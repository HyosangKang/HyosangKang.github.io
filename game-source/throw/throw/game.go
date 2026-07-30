package throw

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	Width  = 800
	Height = 320
)

type scene int

const (
	sceneTitle scene = iota
	sceneAim
	sceneFlight
	sceneLanded
)

type point struct {
	x float64
	y float64
}

type Game struct {
	scene
	origin   [2]float64
	cursor   [2]float64
	velocity [2]float64
	position [2]float64
	trail    []point
}

func NewGame() *Game {
	g := &Game{scene: sceneTitle}
	g.reset()
	g.scene = sceneTitle
	return g
}

func (g *Game) reset() {
	g.scene = sceneAim
	g.origin = [2]float64{72, 230}
	g.cursor = [2]float64{142, 175}
	g.velocity = [2]float64{}
	g.position = g.origin
	g.trail = []point{{x: g.origin[0], y: g.origin[1]}}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	if g.scene == sceneTitle {
		g.drawTitle(screen)
		return
	}

	vector.DrawFilledRect(screen, 0, 232, Width, Height-232, throwWater, false)
	vector.StrokeLine(screen, 0, 232, Width, 232, 2, throwBlue, false)

	for index := 0; index < len(g.trail); index += 3 {
		item := g.trail[index]
		vector.DrawFilledRect(screen, float32(item.x), float32(item.y), 2, 2, throwWhite, false)
	}

	if g.scene == sceneAim {
		vector.StrokeLine(
			screen,
			float32(g.origin[0]),
			float32(g.origin[1]),
			float32(g.cursor[0]),
			float32(g.cursor[1]),
			2,
			throwCyan,
			false,
		)
		drawThrowText(screen, "DRAG  AND  RELEASE", 17, 22, 24, throwWhite, text.AlignStart)
	} else {
		drawThrowText(
			screen,
			fmt.Sprintf("SPEED  %03.0f", math.Hypot(g.velocity[0], g.velocity[1])),
			16,
			Width-22,
			24,
			throwCyan,
			text.AlignEnd,
		)
	}

	vector.DrawFilledCircle(
		screen,
		float32(g.position[0]),
		float32(g.position[1]),
		7,
		throwGold,
		false,
	)
	if g.scene == sceneLanded {
		drawThrowText(screen, "LANDED", 34, Width/2, 106, throwGold, text.AlignCenter)
		drawThrowText(screen, "CLICK  TO  TRY  AGAIN", 17, Width/2, 158, throwWhite, text.AlignCenter)
	}
}

func (g *Game) drawTitle(screen *ebiten.Image) {
	drawThrowText(screen, "THROW!", 54, Width/2, 62, throwWhite, text.AlignCenter)
	drawThrowText(screen, "A  SMALL  PROJECTILE  GAME", 20, Width/2, 126, throwBlue, text.AlignCenter)
	vector.StrokeLine(screen, 82, 216, Width-82, 216, 2, throwBlue, false)
	for index := 0; index < 7; index++ {
		x := 210.0 + float64(index)*52
		y := 204.0 - 72*math.Sin(math.Pi*float64(index)/6)
		vector.DrawFilledCircle(screen, float32(x), float32(y), 3, throwGold, false)
	}
	drawThrowText(screen, "PRESS  ENTER  OR  CLICK", 22, Width/2, 252, throwGold, text.AlignCenter)
	drawThrowText(screen, "DRAG  TO  CHOOSE  THE  THROW", 16, Width/2, 288, throwWhite, text.AlignCenter)
}

func (g *Game) Layout(int, int) (int, int) {
	return Width, Height
}

var (
	throwWhite = color.RGBA{242, 242, 232, 255}
	throwBlue  = color.RGBA{46, 126, 214, 255}
	throwCyan  = color.RGBA{77, 215, 229, 255}
	throwGold  = color.RGBA{244, 197, 66, 255}
	throwWater = color.RGBA{4, 18, 38, 255}
)

var (
	//go:embed arcadeclassic.ttf
	throwFontBytes []byte
	throwFont      *text.GoTextFaceSource
)

func init() {
	throwFont, _ = text.NewGoTextFaceSource(bytes.NewReader(throwFontBytes))
}

func drawThrowText(
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
	text.Draw(screen, message, &text.GoTextFace{Source: throwFont, Size: size}, options)
}

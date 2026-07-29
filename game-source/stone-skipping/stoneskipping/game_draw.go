package stoneskipping

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

var (
	black     = color.RGBA{0, 0, 0, 255}
	white     = color.RGBA{242, 242, 232, 255}
	blue      = color.RGBA{46, 126, 214, 255}
	darkWater = color.RGBA{4, 18, 38, 255}
	cyan      = color.RGBA{77, 215, 229, 255}
	gold      = color.RGBA{244, 197, 66, 255}
	red       = color.RGBA{240, 76, 76, 255}
)

var (
	//go:embed arcadeclassic.ttf
	arcadeClassic []byte

	arcadeFont *text.GoTextFaceSource
)

func init() {
	arcadeFont, _ = text.NewGoTextFaceSource(bytes.NewReader(arcadeClassic))
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(black)
	if g.scene == sceneTitle {
		g.drawTitle(screen)
		return
	}

	g.drawWater(screen)
	g.drawTrail(screen)
	g.drawRipples(screen)
	g.drawStone(screen)
	g.drawScore(screen)
	g.drawMessage(screen)
}

func (g *Game) drawTitle(screen *ebiten.Image) {
	drawText(screen, "STONE  SKIP!", 52, Width/2, 104, white, text.AlignCenter)
	drawText(screen, "A  SMALL  PHYSICS  GAME", 20, Width/2, 166, blue, text.AlignCenter)

	vector.StrokeLine(screen, 84, 245, Width-84, 245, 2, blue, false)
	for index := 0; index < 4; index++ {
		radius := 18.0 + float64(index)*17
		drawEllipse(screen, Width/2, 245, radius, radius*0.18, cyan)
	}
	drawStoneShape(screen, Width/2-122, 225, white)

	drawText(screen, "PRESS  ENTER  OR  CLICK", 24, Width/2, 318, gold, text.AlignCenter)
	drawText(screen, "DRAG  BACK  AND  RELEASE  TO  THROW", 18, Width/2, 362, white, text.AlignCenter)
	drawText(screen, "ESC  RETURNS  HERE", 15, Width/2, 402, blue, text.AlignCenter)
}

func (g *Game) drawWater(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, waterLineY, Width, Height-waterLineY, darkWater, false)
	vector.StrokeLine(screen, 0, waterLineY, Width, waterLineY, 2, blue, false)
	for x := 100; x < Width; x += 100 {
		vector.StrokeLine(screen, float32(x), waterLineY-4, float32(x), waterLineY+4, 1, blue, false)
	}
}

func (g *Game) drawTrail(screen *ebiten.Image) {
	for index := 0; index < len(g.trail); index += 3 {
		x, y := g.toScreen(g.trail[index].x, g.trail[index].y)
		vector.DrawFilledRect(screen, float32(x), float32(y), 2, 2, white, false)
	}
}

func (g *Game) drawRipples(screen *ebiten.Image) {
	for _, item := range g.ripples {
		alpha := uint8(clamp(item.alpha*255, 0, 255))
		drawEllipse(
			screen,
			item.x,
			waterLineY,
			item.radius,
			item.radius*0.18,
			color.RGBA{cyan.R, cyan.G, cyan.B, alpha},
		)
	}
}

func (g *Game) drawStone(screen *ebiten.Image) {
	x, y := g.toScreen(g.x, g.y)
	drawStoneShape(screen, x, y, gold)

	if g.scene != sceneAim {
		return
	}

	dragX, dragY := g.dragX, g.dragY
	if !g.dragging {
		dragX = x - 72
		dragY = y + 10
	}
	drawLine(screen, x, y, dragX, dragY, white)
	vector.DrawFilledRect(screen, float32(dragX-3), float32(dragY-3), 7, 7, red, false)
}

func drawStoneShape(screen *ebiten.Image, x, y float64, itemColor color.Color) {
	length := 12.0
	dx := length * math.Cos(stonePitch)
	dy := length * math.Sin(stonePitch)
	for offset := -2.0; offset <= 2; offset++ {
		drawLine(screen, x-dx, y+dy+offset, x+dx, y-dy+offset, itemColor)
	}
}

func (g *Game) drawScore(screen *ebiten.Image) {
	drawText(screen, "STONE  SKIP", 22, 22, 20, white, text.AlignStart)
	drawText(screen, fmt.Sprintf("SKIPS  %02d", g.skips), 19, Width-22, 22, cyan, text.AlignEnd)
	drawText(screen, fmt.Sprintf("SPEED  %03.0f", math.Hypot(g.vx, g.vy)), 15, Width-22, 52, white, text.AlignEnd)
}

func (g *Game) drawMessage(screen *ebiten.Image) {
	switch g.scene {
	case sceneAim:
		drawText(screen, "DRAG  BACK  FROM  THE  STONE", 19, 22, Height-66, white, text.AlignStart)
		drawText(screen, "ENTER  DEFAULT  SHOT", 15, 22, Height-35, blue, text.AlignStart)
	case sceneFlight:
		if g.paused {
			drawText(screen, "PAUSED", 36, Width/2, Height/2-26, gold, text.AlignCenter)
		}
		drawText(screen, "SPACE  PAUSE     R  RESET", 15, 22, Height-35, blue, text.AlignStart)
	case sceneFinished:
		drawText(screen, fmt.Sprintf("%d  SKIPS", g.skips), 42, Width/2, Height/2-30, gold, text.AlignCenter)
		drawText(screen, "ENTER  OR  CLICK  TO  TRY  AGAIN", 18, Width/2, Height/2+28, white, text.AlignCenter)
	}
}

func drawText(
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
	text.Draw(screen, message, &text.GoTextFace{Source: arcadeFont, Size: size}, options)
}

func drawLine(screen *ebiten.Image, x0, y0, x1, y1 float64, itemColor color.Color) {
	dx := x1 - x0
	dy := y1 - y0
	steps := math.Max(math.Abs(dx), math.Abs(dy))
	if steps < 1 {
		screen.Set(int(x0), int(y0), itemColor)
		return
	}

	for index := 0; index <= int(steps); index++ {
		ratio := float64(index) / steps
		screen.Set(
			int(x0+ratio*dx),
			int(y0+ratio*dy),
			itemColor,
		)
	}
}

func drawEllipse(
	screen *ebiten.Image,
	centerX float64,
	centerY float64,
	radiusX float64,
	radiusY float64,
	itemColor color.Color,
) {
	const segments = 64
	previousX := centerX + radiusX
	previousY := centerY
	for index := 1; index <= segments; index++ {
		angle := 2 * math.Pi * float64(index) / segments
		x := centerX + radiusX*math.Cos(angle)
		y := centerY + radiusY*math.Sin(angle)
		drawLine(screen, previousX, previousY, x, y, itemColor)
		previousX, previousY = x, y
	}
}

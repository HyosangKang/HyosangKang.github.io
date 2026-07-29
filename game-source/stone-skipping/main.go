package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 960
	screenHeight = 540
	waterLineY   = 354

	timeStep          = 0.005
	stepsPerFrame     = 4
	gravity           = 10.0
	waterDensity      = 1.0
	liftCoefficient   = 0.1
	dragCoefficient   = 0.3
	stoneMass         = 1.0
	stoneLength       = 5.0
	stonePitch        = math.Pi / 200
	verticalScale     = 22.0
	maxTrailPoints    = 2400
	minLaunchSpeed    = 45.0
	defaultHorizontal = 150.0
	defaultVertical   = 0.5
)

var (
	ink       = color.RGBA{22, 41, 66, 255}
	mutedInk  = color.RGBA{72, 94, 118, 255}
	sky       = color.RGBA{231, 242, 246, 255}
	skyBand   = color.RGBA{217, 235, 241, 255}
	water     = color.RGBA{41, 128, 185, 255}
	deepWater = color.RGBA{24, 88, 139, 255}
	foam      = color.RGBA{207, 239, 246, 255}
	coral     = color.RGBA{231, 111, 81, 255}
	gold      = color.RGBA{233, 196, 106, 255}
	teal      = color.RGBA{27, 153, 139, 255}
	white     = color.RGBA{255, 255, 255, 255}
)

type gameState int

const (
	stateAim gameState = iota
	stateFlight
	stateFinished
)

type point struct {
	x float64
	y float64
}

type ripple struct {
	x      float64
	radius float64
	alpha  float64
}

type Game struct {
	state      gameState
	paused     bool
	dragging   bool
	dragX      float64
	dragY      float64
	x          float64
	y          float64
	vx         float64
	vy         float64
	elapsed    float64
	skips      int
	inWater    bool
	trail      []point
	ripples    []ripple
	peakHeight float64
}

func NewGame() *Game {
	g := &Game{}
	g.Reset()
	return g
}

func (g *Game) Reset() {
	g.state = stateAim
	g.paused = false
	g.dragging = false
	g.x = 96
	g.y = 1.8
	g.vx = defaultHorizontal
	g.vy = defaultVertical
	g.elapsed = 0
	g.skips = 0
	g.inWater = false
	g.trail = []point{{g.x, g.y}}
	g.ripples = nil
	g.peakHeight = g.y
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.Reset()
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && g.state == stateFlight {
		g.paused = !g.paused
	}

	switch g.state {
	case stateAim:
		g.updateAim()
	case stateFlight:
		if !g.paused {
			for i := 0; i < stepsPerFrame; i++ {
				g.step(timeStep)
				if g.state != stateFlight {
					break
				}
			}
		}
	case stateFinished:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
			inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.Reset()
		}
	}

	for i := range g.ripples {
		g.ripples[i].radius += 0.8
		g.ripples[i].alpha *= 0.985
	}
	live := g.ripples[:0]
	for _, r := range g.ripples {
		if r.alpha > 0.06 {
			live = append(live, r)
		}
	}
	g.ripples = live
	return nil
}

func (g *Game) updateAim() {
	stoneX, stoneY := g.toScreen(g.x, g.y)
	cursorX, cursorY := ebiten.CursorPosition()

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.launch(defaultHorizontal, defaultVertical)
		return
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		dx := float64(cursorX) - stoneX
		dy := float64(cursorY) - stoneY
		if math.Hypot(dx, dy) < 90 {
			g.dragging = true
			g.dragX = float64(cursorX)
			g.dragY = float64(cursorY)
		}
	}

	if g.dragging && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.dragX = float64(cursorX)
		g.dragY = float64(cursorY)
		g.vx, g.vy = g.velocityFromDrag(stoneX, stoneY, g.dragX, g.dragY)
	}

	if g.dragging && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.dragging = false
		if g.vx >= minLaunchSpeed {
			g.launch(g.vx, g.vy)
		}
	}
}

func (g *Game) velocityFromDrag(stoneX, stoneY, dragX, dragY float64) (float64, float64) {
	vx := (stoneX - dragX) * 1.75
	vy := (dragY - stoneY) / verticalScale * 1.6
	return clamp(vx, 0, 220), clamp(vy, -2.5, 8)
}

func (g *Game) launch(vx, vy float64) {
	g.vx = vx
	g.vy = vy
	g.state = stateFlight
	g.trail = []point{{g.x, g.y}}
}

func (g *Game) step(dt float64) {
	ax := 0.0
	ay := -gravity

	if g.y > 0 {
		g.inWater = false
	} else {
		if !g.inWater {
			g.skips++
			g.ripples = append(g.ripples, ripple{x: g.x, radius: 5, alpha: 0.95})
			g.inWater = true
		}

		if g.y+stoneLength*math.Sin(stonePitch) > 0 {
			submergedLength := math.Min(math.Abs(g.y), stoneLength) / math.Sin(stonePitch)
			speedSquared := g.vx*g.vx + g.vy*g.vy
			ay = -gravity +
				0.5*waterDensity/stoneMass*speedSquared*submergedLength*
					(liftCoefficient*math.Cos(stonePitch)-dragCoefficient*math.Sin(stonePitch))
			ax = -0.5 * waterDensity / stoneMass * speedSquared * submergedLength *
				(liftCoefficient*math.Sin(stonePitch) + dragCoefficient*math.Cos(stonePitch))
		} else {
			ax = -dragCoefficient * g.vx * g.vx
			ay = -gravity + dragCoefficient*g.vy*g.vy
		}
	}

	ax = clamp(ax, -7000, 7000)
	ay = clamp(ay, -7000, 7000)
	g.vx += ax * dt
	g.vy += ay * dt
	g.x += g.vx * dt
	g.y += g.vy * dt
	g.elapsed += dt
	g.peakHeight = math.Max(g.peakHeight, g.y)

	if len(g.trail) == 0 || math.Hypot(g.x-g.trail[len(g.trail)-1].x, g.y-g.trail[len(g.trail)-1].y) > 0.8 {
		g.trail = append(g.trail, point{g.x, g.y})
		if len(g.trail) > maxTrailPoints {
			g.trail = g.trail[len(g.trail)-maxTrailPoints:]
		}
	}

	if g.y > 0.04 {
		g.inWater = false
	}

	if g.x > screenWidth-24 || g.y < -2.2 || g.vx < 3 || g.elapsed > 12 {
		g.state = stateFinished
	}
}

func clamp(value, low, high float64) float64 {
	return math.Max(low, math.Min(high, value))
}

func (g *Game) toScreen(x, y float64) (float64, float64) {
	return x, waterLineY - y*verticalScale
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(sky)
	vector.DrawFilledRect(screen, 0, 0, screenWidth, 150, skyBand, false)
	vector.DrawFilledRect(screen, 0, waterLineY, screenWidth, screenHeight-waterLineY, water, false)
	vector.DrawFilledRect(screen, 0, waterLineY+92, screenWidth, screenHeight-waterLineY-92, deepWater, false)
	vector.StrokeLine(screen, 0, waterLineY, screenWidth, waterLineY, 3, foam, true)

	g.drawScale(screen)
	g.drawTrail(screen)
	g.drawRipples(screen)
	g.drawStone(screen)
	g.drawHud(screen)
	g.drawPrompt(screen)
}

func (g *Game) drawScale(screen *ebiten.Image) {
	for x := float32(80); x < screenWidth; x += 100 {
		vector.StrokeLine(screen, x, waterLineY-7, x, waterLineY+7, 1, foam, true)
	}
	for y := float32(waterLineY - 44); y > 105; y -= 44 {
		vector.StrokeLine(screen, 0, y, screenWidth, y, 1, color.RGBA{190, 215, 224, 110}, true)
	}
}

func (g *Game) drawTrail(screen *ebiten.Image) {
	for i := 1; i < len(g.trail); i++ {
		x0, y0 := g.toScreen(g.trail[i-1].x, g.trail[i-1].y)
		x1, y1 := g.toScreen(g.trail[i].x, g.trail[i].y)
		vector.StrokeLine(screen, float32(x0), float32(y0), float32(x1), float32(y1), 2.2, coral, true)
	}
}

func (g *Game) drawRipples(screen *ebiten.Image) {
	for _, r := range g.ripples {
		alpha := uint8(clamp(r.alpha*255, 0, 255))
		c := color.RGBA{foam.R, foam.G, foam.B, alpha}
		previousX := r.x - r.radius
		previousY := float64(waterLineY)
		for i := 1; i <= 48; i++ {
			t := math.Pi * float64(i) / 48
			x := r.x + r.radius*math.Cos(t)
			y := float64(waterLineY) - 0.18*r.radius*math.Sin(t)
			vector.StrokeLine(screen, float32(previousX), float32(previousY), float32(x), float32(y), 1.5, c, true)
			previousX, previousY = x, y
		}
	}
}

func (g *Game) drawStone(screen *ebiten.Image) {
	x, y := g.toScreen(g.x, g.y)
	halfLength := 13.0
	dx := halfLength * math.Cos(stonePitch)
	dy := halfLength * math.Sin(stonePitch)
	vector.StrokeLine(screen, float32(x-dx), float32(y+dy), float32(x+dx), float32(y-dy), 7, ink, true)
	vector.StrokeLine(screen, float32(x-dx), float32(y+dy-1), float32(x+dx), float32(y-dy-1), 2, gold, true)

	if g.state == stateAim {
		dragX, dragY := x-86, y+8
		if g.dragging {
			dragX, dragY = g.dragX, g.dragY
		}
		vector.StrokeLine(screen, float32(x), float32(y), float32(dragX), float32(dragY), 2, coral, true)
		vector.DrawFilledCircle(screen, float32(dragX), float32(dragY), 7, coral, true)
		g.drawArrow(screen, x, y, x+(x-dragX)*0.82, y-(dragY-y)*0.82, teal)
	}
}

func (g *Game) drawArrow(screen *ebiten.Image, x0, y0, x1, y1 float64, c color.Color) {
	vector.StrokeLine(screen, float32(x0), float32(y0), float32(x1), float32(y1), 3, c, true)
	angle := math.Atan2(y1-y0, x1-x0)
	for _, offset := range []float64{2.55, -2.55} {
		hx := x1 + 11*math.Cos(angle+offset)
		hy := y1 + 11*math.Sin(angle+offset)
		vector.StrokeLine(screen, float32(x1), float32(y1), float32(hx), float32(hy), 3, c, true)
	}
}

func (g *Game) drawHud(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 18, 78, 300, 75, color.RGBA{255, 255, 255, 225}, true)
	vector.StrokeRect(screen, 18, 78, 300, 75, 1, color.RGBA{170, 190, 204, 255}, true)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("speed  %5.1f", math.Hypot(g.vx, g.vy)), 34, 94)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("height %5.2f", g.y), 34, 113)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("skips  %d", g.skips), 180, 94)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("time   %4.2f s", g.elapsed), 180, 113)
	if g.paused {
		ebitenutil.DebugPrintAt(screen, "PAUSED", 128, 134)
	}
}

func (g *Game) drawPrompt(screen *ebiten.Image) {
	var title, detail string
	switch g.state {
	case stateAim:
		title = "DRAG BACK FROM THE STONE, THEN RELEASE"
		detail = fmt.Sprintf("launch speed %.0f   vertical speed %.1f   •   Enter uses the MATLAB defaults", g.vx, g.vy)
	case stateFlight:
		title = "AIR: GRAVITY     WATER: LIFT + DRAG"
		detail = "Space pauses the numerical integration   •   R resets"
	case stateFinished:
		title = fmt.Sprintf("%d SKIPS OVER %.0f UNITS", g.skips, g.x-96)
		detail = "Click or press Enter to try another launch"
	}

	vector.DrawFilledRect(screen, 0, screenHeight-58, screenWidth, 58, color.RGBA{12, 31, 52, 232}, false)
	ebitenutil.DebugPrintAt(screen, title, 26, screenHeight-44)
	ebitenutil.DebugPrintAt(screen, detail, 26, screenHeight-24)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Stone Skipping")
	if err := ebiten.RunGame(NewGame()); err != nil {
		panic(err)
	}
}

package stoneskipping

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.scene = sceneTitle
		g.clearShot()
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) && g.scene != sceneTitle {
		g.Reset()
		return nil
	}

	switch g.scene {
	case sceneTitle:
		g.updateTitle()
	case sceneAim:
		g.updateAim()
	case sceneFlight:
		g.updateFlight()
	case sceneFinished:
		g.updateFinished()
	}

	g.updateRipples()
	return nil
}

func (g *Game) updateTitle() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.Reset()
	}
}

func (g *Game) updateAim() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.launch(defaultHorizontal, defaultVertical)
		return
	}

	stoneX, stoneY := g.toScreen(g.x, g.y)
	cursorX, cursorY := ebiten.CursorPosition()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		dx := float64(cursorX) - stoneX
		dy := float64(cursorY) - stoneY
		if math.Hypot(dx, dy) <= 70 {
			g.dragging = true
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

func (g *Game) updateFlight() {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.paused = !g.paused
	}
	if g.paused {
		return
	}

	for i := 0; i < stepsPerFrame; i++ {
		g.step(timeStep)
		if g.scene != sceneFlight {
			break
		}
	}
}

func (g *Game) updateFinished() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.Reset()
	}
}

func (g *Game) velocityFromDrag(stoneX, stoneY, dragX, dragY float64) (float64, float64) {
	vx := (stoneX - dragX) * 1.8
	vy := (dragY - stoneY) / verticalScale * 1.5
	return clamp(vx, 0, 215), clamp(vy, -2.5, 7.5)
}

func (g *Game) launch(vx, vy float64) {
	g.vx = vx
	g.vy = vy
	g.scene = sceneFlight
	g.trail = []point{{x: g.x, y: g.y}}
}

func (g *Game) step(dt float64) {
	ax := 0.0
	ay := -gravity

	if g.y > 0 {
		g.inWater = false
	} else {
		if !g.inWater {
			g.skips++
			g.ripples = append(g.ripples, ripple{x: g.x, radius: 3, alpha: 1})
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

	g.vx += clamp(ax, -7000, 7000) * dt
	g.vy += clamp(ay, -7000, 7000) * dt
	g.x += g.vx * dt
	g.y += g.vy * dt
	g.elapsed += dt

	if len(g.trail) == 0 ||
		math.Hypot(g.x-g.trail[len(g.trail)-1].x, g.y-g.trail[len(g.trail)-1].y) > 1.2 {
		g.trail = append(g.trail, point{x: g.x, y: g.y})
		if len(g.trail) > maxTrailPoints {
			g.trail = g.trail[len(g.trail)-maxTrailPoints:]
		}
	}

	if g.y > 0.04 {
		g.inWater = false
	}

	if g.x > Width-18 || g.y < -3.5 || g.vx < 3 || g.elapsed > 10 {
		g.scene = sceneFinished
	}
}

func (g *Game) updateRipples() {
	for i := range g.ripples {
		g.ripples[i].radius += 0.65
		g.ripples[i].alpha *= 0.982
	}

	live := g.ripples[:0]
	for _, item := range g.ripples {
		if item.alpha > 0.06 {
			live = append(live, item)
		}
	}
	g.ripples = live
}

func (g *Game) toScreen(x, y float64) (float64, float64) {
	return x, waterLineY - y*verticalScale
}

func clamp(value, low, high float64) float64 {
	return math.Max(low, math.Min(high, value))
}

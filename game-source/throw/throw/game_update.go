package throw

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.reset()
		g.scene = sceneTitle
		return nil
	}

	switch g.scene {
	case sceneTitle:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
			inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.reset()
		}
	case sceneAim:
		g.updateAim()
	case sceneFlight:
		g.updateFlight()
	case sceneLanded:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
			inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.reset()
		}
	}
	return nil
}

func (g *Game) updateAim() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.cursor = [2]float64{float64(x), float64(y)}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.velocity = [2]float64{
			(g.cursor[0] - g.origin[0]) / 5.5,
			(g.cursor[1] - g.origin[1]) / 5.5,
		}
		g.scene = sceneFlight
	}
}

func (g *Game) updateFlight() {
	const dt = 0.45
	g.position[0] += g.velocity[0] * dt
	g.position[1] += g.velocity[1] * dt
	g.velocity[1] += 0.34
	g.trail = append(g.trail, point{x: g.position[0], y: g.position[1]})

	if g.position[0] > Width+12 || g.position[1] >= g.origin[1] {
		g.position[1] = g.origin[1]
		g.scene = sceneLanded
	}
}

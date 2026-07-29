package throw

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Stage              int // 0: Ready 1: Throw 2: Land
	Width, Height      int
	Origin, CursorPos  [2]int
	Velocity, Position [2]float64
}

func (game *Game) Draw(screen *ebiten.Image) {
	switch game.Stage {
	case 0:
		Dot(game.Origin[0], game.Origin[1], screen)
		Line(game.Origin[0], game.Origin[1],
			game.CursorPos[0], game.CursorPos[1], screen)
	default:
		x, y := int(game.Position[0]), int(game.Position[1])
		Dot(x, y, screen)
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	return g.Width, g.Height
}

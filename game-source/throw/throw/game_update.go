package throw

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (game *Game) Update() error {
	switch game.Stage {
	case 0:
		game.StageReady()
	case 1:
		game.StageThrow()
	case 2:
		game.StageLand()
	}
	return nil
}

func (game *Game) StageReady() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		game.Velocity = [2]float64{
			float64(x-game.Origin[0]) / 50,
			float64(y-game.Origin[1]) / 50}
		game.CursorPos = [2]int{x, y}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		game.Stage = 1
	}
}

func (game *Game) StageThrow() {
	dt := 1.0
	for i := 0; i < 2; i++ {
		game.Position[i] += game.Velocity[i] * dt
	}
	game.Velocity[1] += 0.1
	if int(game.Position[0]) > game.Width ||
		int(game.Position[1]) > game.Origin[1] {
		game.Stage = 2
	}
}

func (game *Game) StageLand() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		game.Position = [2]float64{
			float64(game.Origin[0]),
			float64(game.Origin[1])}
		game.Stage = 0
	}
}

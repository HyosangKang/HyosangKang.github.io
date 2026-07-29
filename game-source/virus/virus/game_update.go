package virus

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {
	if g.Scene == SceneNil {
		g.Clear()
		g.Scene = SceneStart
	}
	switch g.Scene {
	case SceneStart:
		g.Start()
	case ScenePlay:
		g.Play()
	}
	return nil
}

func (g *Game) Start() {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		g.Scene = ScenePlay
	}
}

func (g *Game) Play() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.Agents = append(g.Agents, NewAgent())
	}
	for _, a := range g.Agents {
		a.Update(g.Agents)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.Clear()
		g.Scene = SceneStart
	}
}

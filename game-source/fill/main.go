package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	fill "github.com/hyosangkang/multi-game/fill/fill"
)

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.RunGame(&fill.Game{})
}

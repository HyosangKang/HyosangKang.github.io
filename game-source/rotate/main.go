package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hyosangkang/multi-game/rotate/example/torus"
	"github.com/hyosangkang/multi-game/rotate/rotate"
)

func main() {
	ebiten.SetWindowSize(600, 600)
	ebiten.SetWindowTitle("Rotate!")
	if err := ebiten.RunGame(&rotate.Game{
		Sprite: torus.Torus{},
	}); err != nil {
		log.Fatal(err)
	}
}

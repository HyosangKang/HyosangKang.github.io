package main

import (
	"virus/virus"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(virus.Width, virus.Height)
	ebiten.RunGame(&virus.Game{})
}

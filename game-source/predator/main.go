package main

import (
	"predator/predator"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(predator.Width, predator.Height)
	ebiten.RunGame(&predator.Game{})
}

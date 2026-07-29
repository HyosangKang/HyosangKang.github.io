package main

import (
	"rocket/rocket"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(rocket.Width, rocket.Height)
	ebiten.RunGame(&rocket.Game{})
}

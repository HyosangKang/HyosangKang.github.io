package main

import (
	"log"
	"throw/throw"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(throw.Width, throw.Height)
	ebiten.SetWindowTitle("Throw!")
	if err := ebiten.RunGame(throw.NewGame()); err != nil {
		log.Fatal(err)
	}
}

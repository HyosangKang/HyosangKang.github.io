package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"math-projects/lanchester/lanchester"
)

func main() {
	ebiten.SetWindowSize(lanchester.Width, lanchester.Height)
	ebiten.SetWindowTitle("Lanchester Arena")
	if err := ebiten.RunGame(lanchester.NewGame()); err != nil {
		log.Fatal(err)
	}
}

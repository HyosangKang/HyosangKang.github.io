package main

import (
	"differential-labs/labs"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(labs.Width, labs.Height)
	ebiten.SetWindowTitle("Differential Equations Lab")
	if err := ebiten.RunGame(labs.NewGame(browserMode())); err != nil {
		log.Fatal(err)
	}
}

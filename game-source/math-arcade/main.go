package main

import (
	"log"
	"math-arcade/arcade"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(arcade.Width, arcade.Height)
	ebiten.SetWindowTitle("Math Arcade")
	if err := ebiten.RunGame(arcade.NewGame(browserMode())); err != nil {
		log.Fatal(err)
	}
}

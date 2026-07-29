package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"math-projects/stone-skipping/stoneskipping"
)

func main() {
	ebiten.SetWindowSize(stoneskipping.Width, stoneskipping.Height)
	ebiten.SetWindowTitle("Stone Skip!")
	if err := ebiten.RunGame(stoneskipping.NewGame()); err != nil {
		log.Fatal(err)
	}
}

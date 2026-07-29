package throw

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func Dot(x, y int, screen *ebiten.Image) {
	for i := -1; i <= 1; i++ {
		for j := -2; j <= 2; j++ {
			screen.Set(x+i, y+j, color.White)
			screen.Set(x+j, y+i, color.White)
		}
	}
}

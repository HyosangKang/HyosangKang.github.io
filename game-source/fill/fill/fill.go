package rotate

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	margin = 10
)

func (f *Game) square(screen *ebiten.Image) {
	for x := 0; x < f.unitLength; x++ {
		for y := 0; y < f.unitLength; y++ {
			screen.Set(margin+x, h-1-y-margin, color.White)
		}
	}
}

func dots(screen *ebiten.Image, pts [][2]int, c color.Color) {
	for _, p := range pts {
		dot(screen, p[0], p[1], c)
	}
}

func dot(screen *ebiten.Image, x, y int, c color.Color) {
	for i := -1; i <= 1; i++ {
		for j := -2; j <= 2; j++ {
			screen.Set(x+i, y+j, c)
			screen.Set(x+j, y+i, c)
		}
	}
}

func fill(screen *ebiten.Image, pts [][2]int, c color.Color) {
	if len(pts) < 3 {
		return
	}
	for i := 1; i < len(pts)-1; i++ {
		fillTriangle(screen, [3][2]int{pts[0], pts[i], pts[i+1]}, c)
	}
	for i := 0; i < len(pts); i++ {
		for _, p := range line(pts[i], pts[(i+1)%len(pts)]) {
			screen.Set(p[0], p[1], color.White)
		}
	}
}

func fillTriangle(screen *ebiten.Image, tri [3][2]int, c color.Color) {
	var xarr, yarr []int
	for _, p := range tri {
		xarr = append(xarr, p[0])
		yarr = append(yarr, p[1])
	}
	var xb, yb [2]int
	xb[0], xb[1] = minmax(xarr...)
	yb[0], yb[1] = minmax(yarr...)
	dx, dy := xb[1]-xb[0], yb[1]-yb[0]
	tline := make(map[[2]int]struct{})
	for i := 0; i < 3; i++ {
		for _, p := range line(tri[i], tri[(i+1)%3]) {
			tline[[2]int{p[0] - xb[0], p[1] - yb[0]}] = struct{}{}
		}
	}
	for x := 0; x <= dx; x++ {
		var color, crossing bool = false, false
		var ps [][2]int
		for y := 0; y <= dy; y++ {
			if _, ok := tline[[2]int{x, y}]; ok {
				if !crossing {
					for _, p := range ps {
						screen.Set(xb[0]+p[0], yb[0]+p[1], c)
					}
					ps = [][2]int{}
					crossing = true
					color = !color
				}
			} else {
				crossing = false
			}
			if color {
				ps = append(ps, [2]int{x, y})
			}
		}
	}
}

func minmax(a ...int) (int, int) {
	min, max := a[0], a[0]
	for _, x := range a {
		if x < min {
			min = x
		}
		if x > max {
			max = x
		}
	}
	return min, max
}

func line(p, q [2]int) [][2]int {
	dx, dy := q[0]-p[0], q[1]-p[1]
	if dx == 0 && dy == 0 {
		return [][2]int{p}
	}
	var ps [][2]int
	_, n := minmax(abs(dx), abs(dy))
	for k := 0; k <= n; k++ {
		x := int(float64(p[0]) + float64(dx)*float64(k)/float64(n))
		y := int(float64(p[1]) + float64(dy)*float64(k)/float64(n))
		ps = append(ps, [2]int{x, y})
	}
	return ps
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

package torus

import (
	"math"
)

type Torus struct{}

const (
	longitudeRadiusTorus = 1.5
	meridianRadiusTorus  = .5
	numGridTorus         = 20
)

var (
	meshPointsTorus    [][3]float64
	initRotationAngles = [3]float64{
		-math.Pi / 2,
		-3 * math.Pi / 4,
		math.Pi / 6}
	initRotationAxis = [3][2]int{
		{1, 2},
		{2, 0},
		{1, 2},
	}
)

func init() {
	n := numGridTorus * numGridTorus
	meshPointsTorus = make([][3]float64, n)
	for i := 0; i < n; i++ {
		j, k := i/numGridTorus, i%numGridTorus
		u := 2 * math.Pi * float64(j) / float64(numGridTorus)
		v := 2 * math.Pi * float64(k) / float64(numGridTorus)
		meshPointsTorus[i] = parametricTorus(u, v)
	}
}

func parametricTorus(u, v float64) [3]float64 {
	a, b := longitudeRadiusTorus, meridianRadiusTorus
	x := (a + b*math.Cos(u)) * math.Cos(v)
	y := (a + b*math.Cos(u)) * math.Sin(v)
	z := b * math.Sin(u)
	return [3]float64{x, y, z}
}

func (s Torus) Graph() [][2][3]float64 {
	var lines [][2][3]float64
	for i, p := range meshPointsTorus {
		j, k := i/numGridTorus, i%numGridTorus
		j1 := (j + 1) % numGridTorus
		k1 := (k + 1) % numGridTorus
		q1 := meshPointsTorus[j1*numGridTorus+k]
		q2 := meshPointsTorus[j*numGridTorus+k1]
		lines = append(lines, [2][3]float64{p, q1}, [2][3]float64{p, q2})
	}
	return lines
}

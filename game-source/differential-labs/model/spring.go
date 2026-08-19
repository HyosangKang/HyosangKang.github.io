package model

import (
	"errors"
	"math"
)

type SpringParams struct {
	M1, M2, K1, K2     float64
	X10, V10, X20, V20 float64
	T                  float64
}

type SpringSample struct {
	T, X1, V1, X2, V2, Energy float64
}

type coordinateSolution struct {
	cosSlow, sinSlow float64
	cosFast, sinFast float64
	slow, fast       float64
}

type SpringModel struct {
	Params        SpringParams
	first, second coordinateSolution
	Slow, Fast    float64
	Samples       []SpringSample
	MaxAbs        float64
}

func NewSpringModel(params SpringParams) (*SpringModel, error) {
	if params.M1 <= 0 || params.M2 <= 0 || params.K1 <= 0 || params.K2 <= 0 || params.T <= 0 {
		return nil, errors.New("masses, spring constants, and time must be positive")
	}

	polynomialA := params.K2/params.M2 + (params.K1+params.K2)/params.M1
	polynomialB := params.K1 * params.K2 / (params.M1 * params.M2)
	discriminant := math.Max(0, polynomialA*polynomialA-4*polynomialB)
	rootLow := (-polynomialA - math.Sqrt(discriminant)) / 2
	rootHigh := (-polynomialA + math.Sqrt(discriminant)) / 2
	fast := math.Sqrt(math.Max(1e-10, -rootLow))
	slow := math.Sqrt(math.Max(1e-10, -rootHigh))

	acceleration1 := (-(params.K1+params.K2)*params.X10 + params.K2*params.X20) / params.M1
	acceleration2 := (params.K2*params.X10 - params.K2*params.X20) / params.M2
	jerk1 := (-(params.K1+params.K2)*params.V10 + params.K2*params.V20) / params.M1
	jerk2 := (params.K2*params.V10 - params.K2*params.V20) / params.M2

	first, err := newCoordinateSolution(params.X10, params.V10, acceleration1, jerk1, slow, fast)
	if err != nil {
		return nil, err
	}
	second, err := newCoordinateSolution(params.X20, params.V20, acceleration2, jerk2, slow, fast)
	if err != nil {
		return nil, err
	}

	model := &SpringModel{Params: params, first: first, second: second, Slow: slow, Fast: fast}
	model.Samples = make([]SpringSample, 1001)
	for index := range model.Samples {
		t := params.T * float64(index) / float64(len(model.Samples)-1)
		sample := model.Motion(t)
		model.MaxAbs = math.Max(model.MaxAbs, math.Max(math.Abs(sample.X1), math.Abs(sample.X2)))
		model.Samples[index] = sample
	}
	model.MaxAbs = math.Max(0.2, model.MaxAbs*1.12)
	return model, nil
}

func newCoordinateSolution(x0, v0, acceleration, jerk, slow, fast float64) (coordinateSolution, error) {
	slowSquared, fastSquared := slow*slow, fast*fast
	cosSlow, cosFast, err := solve2x2(1, 1, slowSquared, fastSquared, x0, -acceleration)
	if err != nil {
		return coordinateSolution{}, err
	}
	velocitySlow, velocityFast, err := solve2x2(1, 1, slowSquared, fastSquared, v0, -jerk)
	if err != nil {
		return coordinateSolution{}, err
	}
	return coordinateSolution{
		cosSlow: cosSlow,
		sinSlow: velocitySlow / slow,
		cosFast: cosFast,
		sinFast: velocityFast / fast,
		slow:    slow,
		fast:    fast,
	}, nil
}

func (model *SpringModel) Motion(t float64) SpringSample {
	x1, v1 := model.first.at(t)
	x2, v2 := model.second.at(t)
	energy := 0.5*model.Params.M1*v1*v1 +
		0.5*model.Params.M2*v2*v2 +
		0.5*model.Params.K1*x1*x1 +
		0.5*model.Params.K2*(x2-x1)*(x2-x1)
	return SpringSample{T: t, X1: x1, V1: v1, X2: x2, V2: v2, Energy: energy}
}

func (solution coordinateSolution) at(t float64) (float64, float64) {
	position := solution.cosSlow*math.Cos(solution.slow*t) +
		solution.sinSlow*math.Sin(solution.slow*t) +
		solution.cosFast*math.Cos(solution.fast*t) +
		solution.sinFast*math.Sin(solution.fast*t)
	velocity := -solution.cosSlow*solution.slow*math.Sin(solution.slow*t) +
		solution.sinSlow*solution.slow*math.Cos(solution.slow*t) -
		solution.cosFast*solution.fast*math.Sin(solution.fast*t) +
		solution.sinFast*solution.fast*math.Cos(solution.fast*t)
	return position, velocity
}

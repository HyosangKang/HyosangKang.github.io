package model

import (
	"errors"
	"math"
)

type ElectricParams struct {
	L, R1, R2, C, T float64
}

type ElectricSample struct {
	T, I1, I2 float64
}

type quadraticPart struct {
	c, d, a, b, constant float64
}

type ElectricModel struct {
	Params                ElectricParams
	a1, b1, a2, b2        float64
	firstPart, secondPart quadraticPart
	Samples               []ElectricSample
	MaxAbs                float64
	Response              string
}

func solve2x2(m00, m01, m10, m11, b0, b1 float64) (float64, float64, error) {
	determinant := m00*m11 - m01*m10
	if math.Abs(determinant) < 1e-10 {
		return 0, 0, errors.New("singular parameter combination")
	}
	return (b0*m11 - b1*m01) / determinant,
		(m00*b1 - b0*m10) / determinant, nil
}

func NewElectricModel(params ElectricParams) (*ElectricModel, error) {
	if params.L <= 0 || params.R1 <= 0 || params.R2 <= 0 || params.C <= 0 || params.T <= 0 {
		return nil, errors.New("all circuit values must be positive")
	}

	aFirst := params.C * params.L * params.R2
	bFirst := params.L + params.C*params.R1*params.R2
	dFirst := params.R1 + params.R2
	a1, b1, err := solve2x2(
		dFirst-aFirst, bFirst,
		bFirst, aFirst-dFirst,
		0, -1,
	)
	if err != nil {
		return nil, err
	}
	c1 := -aFirst * a1
	d1 := 1 - dFirst*b1

	aSecond := params.L * params.R2
	bSecond := params.L/params.C + params.R1*params.R2
	dSecond := (params.R1 + params.R2) / params.C
	a2, b2, err := solve2x2(
		bSecond, aSecond-dSecond,
		dSecond-aSecond, bSecond,
		params.L-1/params.C, params.R1,
	)
	if err != nil {
		return nil, err
	}
	c2 := -aSecond * a2
	d2 := 1/params.C - dSecond*b2

	model := &ElectricModel{
		Params: params,
		a1:     a1,
		b1:     b1,
		a2:     a2,
		b2:     b2,
		firstPart: quadraticPart{
			c: c1, d: d1, a: aFirst, b: bFirst, constant: dFirst,
		},
		secondPart: quadraticPart{
			c: c2, d: d2, a: aSecond, b: bSecond, constant: dSecond,
		},
	}

	discriminant := bFirst*bFirst - 4*aFirst*dFirst
	switch {
	case math.Abs(discriminant) < 1e-6:
		model.Response = "CRITICAL"
	case discriminant < 0:
		model.Response = "RINGING"
	default:
		model.Response = "OVERDAMPED"
	}

	model.Samples = make([]ElectricSample, 1001)
	for index := range model.Samples {
		t := params.T * float64(index) / float64(len(model.Samples)-1)
		i1, i2 := model.Current(t)
		if math.IsNaN(i1) || math.IsNaN(i2) || math.IsInf(i1, 0) || math.IsInf(i2, 0) {
			return nil, errors.New("circuit response overflow")
		}
		model.MaxAbs = math.Max(model.MaxAbs, math.Max(math.Abs(i1), math.Abs(i2)))
		model.Samples[index] = ElectricSample{T: t, I1: i1, I2: i2}
	}
	model.MaxAbs = math.Max(0.1, model.MaxAbs*1.12)
	return model, nil
}

func (model *ElectricModel) Current(t float64) (float64, float64) {
	i1 := model.a1*math.Cos(t) + model.b1*math.Sin(t) + model.firstPart.value(t)
	i2 := model.a2*math.Cos(t) + model.b2*math.Sin(t) + model.secondPart.value(t)
	return i1, i2
}

func (part quadraticPart) value(t float64) float64 {
	alpha := part.b / (2 * part.a)
	betaSquared := part.constant/part.a - alpha*alpha
	shiftedConstant := part.d - part.c*alpha
	envelope := math.Exp(-alpha * t)

	if math.Abs(betaSquared) < 1e-9 {
		return envelope * (part.c/part.a + shiftedConstant/part.a*t)
	}
	beta := math.Sqrt(math.Abs(betaSquared))
	if betaSquared > 0 {
		return envelope * (part.c/part.a*math.Cos(beta*t) +
			shiftedConstant/(part.a*beta)*math.Sin(beta*t))
	}
	return envelope * (part.c/part.a*math.Cosh(beta*t) +
		shiftedConstant/(part.a*beta)*math.Sinh(beta*t))
}

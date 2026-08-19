package model

import (
	"math"
	"testing"
)

func TestElectricModelInitialValuesAndEquations(t *testing.T) {
	params := ElectricParams{L: 1, R1: 1, R2: 2, C: 1, T: 18}
	model, err := NewElectricModel(params)
	if err != nil {
		t.Fatal(err)
	}
	i1, i2 := model.Current(0)
	if math.Abs(i1) > 1e-9 || math.Abs(i2) > 1e-9 {
		t.Fatalf("initial currents = (%g, %g), want (0, 0)", i1, i2)
	}

	const h = 1e-4
	for _, sampleTime := range []float64{0.4, 1.1, 3.7, 7.2} {
		previousI1, previousI2 := model.Current(sampleTime - h)
		currentI1, currentI2 := model.Current(sampleTime)
		nextI1, nextI2 := model.Current(sampleTime + h)
		i1Prime := (nextI1 - previousI1) / (2 * h)
		i1Second := (nextI1 - 2*currentI1 + previousI1) / (h * h)
		i2Prime := (nextI2 - previousI2) / (2 * h)
		firstResidual := params.L*i1Second + params.R1*i1Prime + (currentI1-currentI2)/params.C
		secondResidual := params.R2*i2Prime + (currentI2-currentI1)/params.C - math.Cos(sampleTime)
		if math.Abs(firstResidual) > 2e-5 || math.Abs(secondResidual) > 2e-7 {
			t.Fatalf("equation residual at t=%g: (%g, %g)", sampleTime, firstResidual, secondResidual)
		}
	}
}

func TestSpringModelModesEquationsAndEnergy(t *testing.T) {
	params := SpringParams{M1: 1, M2: 1, K1: 3, K2: 2, X10: 1, V10: 0.2, X20: -0.3, V20: -0.1, T: 18}
	model, err := NewSpringModel(params)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(model.Slow-1) > 1e-10 || math.Abs(model.Fast-math.Sqrt(6)) > 1e-10 {
		t.Fatalf("normal modes = (%g, %g), want (1, sqrt(6))", model.Slow, model.Fast)
	}
	initial := model.Motion(0)
	if math.Abs(initial.X1-params.X10) > 1e-10 || math.Abs(initial.V1-params.V10) > 1e-10 ||
		math.Abs(initial.X2-params.X20) > 1e-10 || math.Abs(initial.V2-params.V20) > 1e-10 {
		t.Fatalf("initial state = %+v", initial)
	}

	const h = 1e-4
	initialEnergy := initial.Energy
	for _, sampleTime := range []float64{0.4, 1.1, 3.7, 7.2} {
		previous := model.Motion(sampleTime - h)
		current := model.Motion(sampleTime)
		next := model.Motion(sampleTime + h)
		x1Second := (next.X1 - 2*current.X1 + previous.X1) / (h * h)
		x2Second := (next.X2 - 2*current.X2 + previous.X2) / (h * h)
		firstResidual := params.M1*x1Second + params.K1*current.X1 - params.K2*(current.X2-current.X1)
		secondResidual := params.M2*x2Second + params.K2*(current.X2-current.X1)
		if math.Abs(firstResidual) > 2e-6 || math.Abs(secondResidual) > 2e-6 {
			t.Fatalf("equation residual at t=%g: (%g, %g)", sampleTime, firstResidual, secondResidual)
		}
		if math.Abs(current.Energy-initialEnergy) > 1e-9 {
			t.Fatalf("energy changed at t=%g: got %g, want %g", sampleTime, current.Energy, initialEnergy)
		}
	}
}

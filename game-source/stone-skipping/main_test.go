package main

import "testing"

func TestDefaultLaunchProducesARecordedWaterContact(t *testing.T) {
	g := NewGame()
	g.launch(defaultHorizontal, defaultVertical)

	for i := 0; i < 4000 && g.state == stateFlight; i++ {
		g.step(timeStep)
	}

	if g.x <= 96 {
		t.Fatalf("stone did not move forward: x = %.2f", g.x)
	}
	if g.skips < 1 {
		t.Fatalf("expected at least one water contact, got %d", g.skips)
	}
	if len(g.trail) < 2 {
		t.Fatalf("expected a recorded trajectory, got %d points", len(g.trail))
	}
}

func TestGravityActsWhileStoneIsAboveWater(t *testing.T) {
	g := NewGame()
	g.y = 3
	g.vy = 2
	before := g.vy

	g.step(timeStep)

	if g.vy >= before {
		t.Fatalf("expected gravity to reduce vertical velocity: before %.3f after %.3f", before, g.vy)
	}
}

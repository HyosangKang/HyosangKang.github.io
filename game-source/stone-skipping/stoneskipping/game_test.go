package stoneskipping

import "testing"

func TestDefaultLaunchMovesForwardAndSkips(t *testing.T) {
	g := NewGame()
	g.Reset()
	g.launch(defaultHorizontal, defaultVertical)

	for index := 0; index < 4000 && g.scene == sceneFlight; index++ {
		g.step(timeStep)
	}

	if g.x <= launchX {
		t.Fatalf("stone did not move forward: x = %.2f", g.x)
	}
	if g.skips < 2 {
		t.Fatalf("expected the default shot to skip at least twice, got %d", g.skips)
	}
	if len(g.trail) < 2 {
		t.Fatalf("expected a recorded trajectory, got %d points", len(g.trail))
	}
}

func TestGravityActsAboveWater(t *testing.T) {
	g := NewGame()
	g.Reset()
	g.y = 3
	g.vy = 2
	before := g.vy

	g.step(timeStep)

	if g.vy >= before {
		t.Fatalf("expected gravity to reduce vertical velocity: before %.3f after %.3f", before, g.vy)
	}
}

func TestBackwardDragCreatesForwardLaunch(t *testing.T) {
	g := NewGame()
	g.Reset()
	stoneX, stoneY := g.toScreen(g.x, g.y)

	vx, vy := g.velocityFromDrag(stoneX, stoneY, stoneX-80, stoneY+20)

	if vx <= 0 {
		t.Fatalf("expected a forward velocity, got %.2f", vx)
	}
	if vy <= 0 {
		t.Fatalf("expected an upward velocity, got %.2f", vy)
	}
}

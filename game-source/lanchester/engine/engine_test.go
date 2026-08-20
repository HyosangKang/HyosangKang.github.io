package engine

import (
	"math"
	"testing"
)

func TestInitialCountsAndSpawn(t *testing.T) {
	battle := New([]int{12, 19, 7}, 960, 540, 42)
	snapshot := battle.Snapshot()
	if len(snapshot.Counts) != 3 {
		t.Fatalf("got %d teams, want 3", len(snapshot.Counts))
	}
	if len(snapshot.Units) != 38 {
		t.Fatalf("got %d units, want 38", len(snapshot.Units))
	}
	for i, want := range []int{12, 19, 7} {
		if snapshot.Counts[i] != want {
			t.Fatalf("team %d has %d units, want %d", i, snapshot.Counts[i], want)
		}
	}
	for _, unit := range snapshot.Units {
		if unit.X < 0 || unit.X > 960 || unit.Y < 0 || unit.Y > 540 {
			t.Fatalf("unit %d spawned outside the arena", unit.ID)
		}
	}
}

func TestBattleEndsAndCountsOnlyDecrease(t *testing.T) {
	battle := New([]int{24, 20, 16}, 960, 540, 9)
	previous := []int{24, 20, 16}
	sawLaser := false
	for steps := 0; steps < maxBattleTicks && !battle.done; steps++ {
		battle.Step()
		if len(battle.Snapshot().Shots) > 0 {
			sawLaser = true
		}
		for team, count := range battle.counts {
			if count > previous[team] {
				t.Fatalf("team %d count increased from %d to %d", team, previous[team], count)
			}
			previous[team] = count
		}
	}
	if !battle.done {
		t.Fatal("battle did not finish")
	}
	if !sawLaser {
		t.Fatal("battle never produced a ranged shot")
	}
	snapshot := battle.Snapshot()
	if len(snapshot.History) != snapshot.Tick+1 {
		t.Fatalf("history has %d rows at tick %d", len(snapshot.History), snapshot.Tick)
	}
}

func TestUnitsMoveIntoLaserRange(t *testing.T) {
	battle := New([]int{1, 1}, 900, 456, 27)
	initial := battle.Snapshot()
	initialDistance := distanceBetween(initial.Units[0], initial.Units[1])

	for step := 0; step < 500; step++ {
		battle.Step()
		snapshot := battle.Snapshot()
		if len(snapshot.Shots) == 0 {
			continue
		}
		shot := snapshot.Shots[0]
		shotDistance := math.Hypot(shot.ToX-shot.FromX, shot.ToY-shot.FromY)
		if shotDistance > attackRange+0.001 {
			t.Fatalf("laser fired at distance %.2f, beyond range %.2f", shotDistance, attackRange)
		}
		if shotDistance >= initialDistance {
			t.Fatalf("units did not approach: initial %.2f, firing %.2f", initialDistance, shotDistance)
		}
		return
	}
	t.Fatal("units never moved into laser range")
}

func distanceBetween(a, b Unit) float64 {
	return math.Hypot(a.X-b.X, a.Y-b.Y)
}

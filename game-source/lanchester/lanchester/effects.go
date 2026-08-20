package lanchester

import "math-projects/lanchester/engine"

func (g *Game) captureEffects(next engine.Snapshot) {
	nextUnits := make(map[int]engine.Unit, len(next.Units))
	for _, unit := range next.Units {
		nextUnits[unit.ID] = unit
	}
	for id, unit := range g.previousUnits {
		if _, alive := nextUnits[id]; alive {
			continue
		}
		for count := 0; count < 6; count++ {
			vx, vy := randomParticleVelocity()
			g.particles = append(g.particles, particle{x: unit.X, y: unit.Y, vx: vx, vy: vy, life: 1, team: unit.Team})
		}
	}
	g.previousUnits = nextUnits

	for _, shot := range next.Shots {
		g.lasers = append(g.lasers, laser{
			fromX: shot.FromX,
			fromY: shot.FromY,
			toX:   shot.ToX,
			toY:   shot.ToY,
			life:  1,
			team:  shot.Team,
		})
	}
	if len(g.lasers) > 480 {
		g.lasers = g.lasers[len(g.lasers)-480:]
	}
}

package engine

import (
	"math"
	"math/rand"
)

const (
	minTeams       = 2
	maxTeams       = 6
	minForceSize   = 1
	maxForceSize   = 120
	startingHealth = 100.0
	moveSpeed      = 2.15
	attackRange    = 104.0
	attackDamage   = 18.0
	attackDelay    = 7
	maxBattleTicks = 3600
)

// Shot is one ranged attack rendered as a short-lived laser by the game.
type Shot struct {
	FromX float64
	FromY float64
	ToX   float64
	ToY   float64
	Team  int
}

// Unit is the state of one living minion exposed to the game renderer.
type Unit struct {
	ID     int
	Team   int
	X      float64
	Y      float64
	Health float64
	target int
	wait   int
	alive  bool
}

// Snapshot is an immutable view used by the Go game renderer and result graph.
type Snapshot struct {
	Tick    int
	Done    bool
	Winner  int
	Counts  []int
	Units   []Unit
	Shots   []Shot
	History [][]int
}

// Battle contains the complete deterministic simulation state.
type Battle struct {
	units   []Unit
	counts  []int
	history [][]int
	tick    int
	done    bool
	winner  int
	width   float64
	height  float64
	rng     *rand.Rand
	shots   []Shot
}

// New places team clusters at even angles around an arc-shaped perimeter.
func New(forceSizes []int, width, height float64, seed int64) *Battle {
	sizes := sanitiseSizes(forceSizes)
	if width < 320 {
		width = 320
	}
	if height < 240 {
		height = 240
	}

	b := &Battle{
		counts: make([]int, len(sizes)),
		winner: -1,
		width:  width,
		height: height,
		rng:    rand.New(rand.NewSource(seed)),
	}

	total := 0
	for _, size := range sizes {
		total += size
	}
	b.units = make([]Unit, 0, total)

	cx, cy := width/2, height/2
	arenaRadius := math.Min(width*0.37, height*0.35)
	for team, size := range sizes {
		angle := -math.Pi/2 + 2*math.Pi*float64(team)/float64(len(sizes))
		clusterX := cx + math.Cos(angle)*arenaRadius
		clusterY := cy + math.Sin(angle)*arenaRadius
		clusterRadius := math.Min(48, 23+math.Sqrt(float64(size))*2.25)

		for member := 0; member < size; member++ {
			spreadAngle := b.rng.Float64() * 2 * math.Pi
			spreadRadius := math.Sqrt(b.rng.Float64()) * clusterRadius
			b.units = append(b.units, Unit{
				ID:     len(b.units),
				Team:   team,
				X:      clusterX + math.Cos(spreadAngle)*spreadRadius,
				Y:      clusterY + math.Sin(spreadAngle)*spreadRadius,
				Health: startingHealth,
				target: -1,
				wait:   b.rng.Intn(attackDelay),
				alive:  true,
			})
		}
		b.counts[team] = size
	}
	b.recordCounts()
	return b
}

func sanitiseSizes(forceSizes []int) []int {
	if len(forceSizes) < minTeams {
		forceSizes = []int{40, 40}
	}
	if len(forceSizes) > maxTeams {
		forceSizes = forceSizes[:maxTeams]
	}
	sizes := make([]int, len(forceSizes))
	for i, size := range forceSizes {
		if size < minForceSize {
			size = minForceSize
		}
		if size > maxForceSize {
			size = maxForceSize
		}
		sizes[i] = size
	}
	return sizes
}

// Step advances movement and resolves all attacks simultaneously for fairness.
func (b *Battle) Step() {
	if b.done {
		return
	}

	b.shots = b.shots[:0]
	damage := make([]float64, len(b.units))
	for i := range b.units {
		unit := &b.units[i]
		if !unit.alive {
			continue
		}
		if unit.wait > 0 {
			unit.wait--
		}

		if !b.validTarget(unit.Team, unit.target) {
			unit.target = b.randomEnemy(unit.Team)
		}
		if unit.target < 0 {
			continue
		}

		target := &b.units[unit.target]
		dx, dy := target.X-unit.X, target.Y-unit.Y
		distance := math.Hypot(dx, dy)
		if distance > attackRange {
			step := math.Min(moveSpeed, distance-attackRange+0.35)
			if distance > 0 {
				unit.X += dx/distance*step + (b.rng.Float64()-0.5)*0.22
				unit.Y += dy/distance*step + (b.rng.Float64()-0.5)*0.22
			}
			unit.X = clamp(unit.X, 8, b.width-8)
			unit.Y = clamp(unit.Y, 8, b.height-8)
			continue
		}

		if unit.wait == 0 {
			variation := 0.82 + b.rng.Float64()*0.36
			damage[unit.target] += attackDamage * variation
			b.shots = append(b.shots, Shot{
				FromX: unit.X,
				FromY: unit.Y,
				ToX:   target.X,
				ToY:   target.Y,
				Team:  unit.Team,
			})
			unit.wait = attackDelay + b.rng.Intn(4)
		}
	}

	for i, amount := range damage {
		if amount == 0 || !b.units[i].alive {
			continue
		}
		b.units[i].Health -= amount
		if b.units[i].Health <= 0 {
			b.units[i].Health = 0
			b.units[i].alive = false
			b.counts[b.units[i].Team]--
		}
	}

	b.tick++
	b.recordCounts()
	b.updateOutcome()
}

func (b *Battle) validTarget(team, target int) bool {
	return target >= 0 && target < len(b.units) && b.units[target].alive && b.units[target].Team != team
}

func (b *Battle) randomEnemy(team int) int {
	if len(b.units) == 0 {
		return -1
	}
	start := b.rng.Intn(len(b.units))
	stride := 1 + b.rng.Intn(max(1, len(b.units)-1))
	for checked, index := 0, start; checked < len(b.units); checked, index = checked+1, (index+stride)%len(b.units) {
		if b.units[index].alive && b.units[index].Team != team {
			return index
		}
	}
	for i := range b.units {
		if b.units[i].alive && b.units[i].Team != team {
			return i
		}
	}
	return -1
}

func (b *Battle) recordCounts() {
	row := append([]int(nil), b.counts...)
	b.history = append(b.history, row)
}

func (b *Battle) updateOutcome() {
	activeTeams := 0
	winner := -1
	for team, count := range b.counts {
		if count > 0 {
			activeTeams++
			winner = team
		}
	}
	if activeTeams <= 1 {
		b.done = true
		b.winner = winner
		return
	}
	if b.tick >= maxBattleTicks {
		b.done = true
		best := -1
		for team, count := range b.counts {
			if count > best {
				best = count
				b.winner = team
			} else if count == best {
				b.winner = -1
			}
		}
	}
}

// Snapshot returns copies so rendering cannot mutate the simulation.
func (b *Battle) Snapshot() Snapshot {
	living := make([]Unit, 0)
	for _, unit := range b.units {
		if unit.alive {
			living = append(living, unit)
		}
	}
	snapshot := Snapshot{
		Tick:   b.tick,
		Done:   b.done,
		Winner: b.winner,
		Counts: append([]int(nil), b.counts...),
		Units:  living,
		Shots:  append([]Shot(nil), b.shots...),
	}
	if b.done {
		snapshot.History = make([][]int, len(b.history))
		for i, row := range b.history {
			snapshot.History[i] = append([]int(nil), row...)
		}
	}
	return snapshot
}

func clamp(value, low, high float64) float64 {
	return math.Max(low, math.Min(high, value))
}

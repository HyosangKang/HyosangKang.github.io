package lanchester

import (
	"time"

	"math-projects/lanchester/engine"
)

const (
	Width  = 960
	Height = 640

	arenaX      = 30.0
	arenaY      = 112.0
	arenaWidth  = 900.0
	arenaHeight = 456.0
)

type scene int

const (
	sceneSetup scene = iota
	sceneBattle
	sceneResult
)

type rect struct {
	x float32
	y float32
	w float32
	h float32
}

func (r rect) contains(x, y int) bool {
	return float32(x) >= r.x && float32(x) <= r.x+r.w && float32(y) >= r.y && float32(y) <= r.y+r.h
}

type particle struct {
	x, y   float64
	vx, vy float64
	life   float64
	team   int
}

type laser struct {
	fromX, fromY float64
	toX, toY     float64
	life         float64
	team         int
}

type Game struct {
	scene       scene
	teamCount   int
	sizes       [6]int
	editingTeam int
	editBuffer  string
	editFresh   bool

	battle        *engine.Battle
	snapshot      engine.Snapshot
	previousUnits map[int]engine.Unit
	particles     []particle
	lasers        []laser
	paused        bool
	frame         int
}

func NewGame() *Game {
	return &Game{
		scene:       sceneSetup,
		teamCount:   3,
		sizes:       [6]int{52, 44, 36, 32, 32, 32},
		editingTeam: -1,
	}
}

func (g *Game) Layout(int, int) (int, int) {
	return Width, Height
}

func (g *Game) startBattle() {
	g.commitEdit()
	forceSizes := append([]int(nil), g.sizes[:g.teamCount]...)
	g.battle = engine.New(forceSizes, arenaWidth, arenaHeight, time.Now().UnixNano())
	g.snapshot = g.battle.Snapshot()
	g.previousUnits = make(map[int]engine.Unit, len(g.snapshot.Units))
	for _, unit := range g.snapshot.Units {
		g.previousUnits[unit.ID] = unit
	}
	g.particles = nil
	g.lasers = nil
	g.paused = false
	g.frame = 0
	g.scene = sceneBattle
}

func (g *Game) returnToSetup() {
	g.scene = sceneSetup
	g.editingTeam = -1
	g.editBuffer = ""
	g.paused = false
}

func teamCardRect(index int) rect {
	column := index % 3
	row := index / 3
	return rect{
		x: float32(30 + column*310),
		y: float32(126 + row*170),
		w: 280,
		h: 145,
	}
}

func teamMinusRect(index int) rect {
	card := teamCardRect(index)
	return rect{x: card.x + 18, y: card.y + 84, w: 46, h: 42}
}

func teamValueRect(index int) rect {
	card := teamCardRect(index)
	return rect{x: card.x + 72, y: card.y + 84, w: 136, h: 42}
}

func teamPlusRect(index int) rect {
	card := teamCardRect(index)
	return rect{x: card.x + 216, y: card.y + 84, w: 46, h: 42}
}

var (
	teamCountMinus = rect{x: 722, y: 44, w: 42, h: 42}
	teamCountValue = rect{x: 772, y: 44, w: 76, h: 42}
	teamCountPlus  = rect{x: 856, y: 44, w: 42, h: 42}
	startButton    = rect{x: 704, y: 548, w: 226, h: 58}
	editButton     = rect{x: 30, y: 588, w: 154, h: 36}
	pauseButton    = rect{x: 776, y: 588, w: 154, h: 36}
	resultEdit     = rect{x: 566, y: 562, w: 158, h: 48}
	resultAgain    = rect{x: 738, y: 562, w: 192, h: 48}
)

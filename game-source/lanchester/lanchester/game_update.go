package lanchester

import (
	"math"
	"math/rand"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {
	switch g.scene {
	case sceneSetup:
		g.updateSetup()
	case sceneBattle:
		g.updateBattle()
	case sceneResult:
		g.updateResult()
	}
	return nil
}

func (g *Game) updateSetup() {
	if g.editingTeam >= 0 {
		g.updateNumberInput()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.commitEdit()
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if g.editingTeam >= 0 {
			g.commitEdit()
			return
		}
		g.startBattle()
		return
	}
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}

	x, y := ebiten.CursorPosition()
	if teamCountMinus.contains(x, y) {
		g.commitEdit()
		if g.teamCount > 2 {
			g.teamCount--
		}
		return
	}
	if teamCountPlus.contains(x, y) {
		g.commitEdit()
		if g.teamCount < 6 {
			g.teamCount++
		}
		return
	}
	if teamCountValue.contains(x, y) {
		return
	}
	for team := 0; team < g.teamCount; team++ {
		switch {
		case teamMinusRect(team).contains(x, y):
			g.commitEdit()
			g.sizes[team] = clampInt(g.sizes[team]-5, 1, 120)
			return
		case teamValueRect(team).contains(x, y):
			g.beginEdit(team)
			return
		case teamPlusRect(team).contains(x, y):
			g.commitEdit()
			g.sizes[team] = clampInt(g.sizes[team]+5, 1, 120)
			return
		}
	}
	if startButton.contains(x, y) {
		g.startBattle()
		return
	}
	g.commitEdit()
}

func (g *Game) beginEdit(team int) {
	g.commitEdit()
	g.editingTeam = team
	g.editBuffer = strconv.Itoa(g.sizes[team])
	g.editFresh = true
}

func (g *Game) updateNumberInput() {
	for _, typed := range ebiten.AppendInputChars(nil) {
		if typed < '0' || typed > '9' {
			continue
		}
		if g.editFresh {
			g.editBuffer = ""
			g.editFresh = false
		}
		if len(g.editBuffer) < 3 {
			g.editBuffer += string(typed)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		g.editFresh = false
		if len(g.editBuffer) > 0 {
			g.editBuffer = g.editBuffer[:len(g.editBuffer)-1]
		}
	}
}

func (g *Game) commitEdit() {
	if g.editingTeam < 0 {
		return
	}
	if value, err := strconv.Atoi(g.editBuffer); err == nil {
		g.sizes[g.editingTeam] = clampInt(value, 1, 120)
	}
	g.editingTeam = -1
	g.editBuffer = ""
	g.editFresh = false
}

func (g *Game) updateBattle() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.returnToSetup()
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.startBattle()
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.paused = !g.paused
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if editButton.contains(x, y) {
			g.returnToSetup()
			return
		}
		if pauseButton.contains(x, y) {
			g.paused = !g.paused
		}
	}

	if !g.paused {
		g.frame++
		if g.frame%2 == 0 {
			g.battle.Step()
			next := g.battle.Snapshot()
			g.captureEffects(next)
			g.snapshot = next
			if next.Done {
				g.scene = sceneResult
			}
		}
	}
	g.updateEffects()
}

func (g *Game) updateResult() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.startBattle()
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.returnToSetup()
		return
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if resultEdit.contains(x, y) {
			g.returnToSetup()
		} else if resultAgain.contains(x, y) {
			g.startBattle()
		}
	}
}

func (g *Game) updateEffects() {
	particles := g.particles[:0]
	for _, item := range g.particles {
		item.x += item.vx
		item.y += item.vy
		item.vx *= 0.94
		item.vy *= 0.94
		item.life -= 0.065
		if item.life > 0 {
			particles = append(particles, item)
		}
	}
	g.particles = particles

	lasers := g.lasers[:0]
	for _, item := range g.lasers {
		item.life -= 0.22
		if item.life > 0 {
			lasers = append(lasers, item)
		}
	}
	g.lasers = lasers
}

func clampInt(value, low, high int) int {
	return max(low, min(high, value))
}

func randomParticleVelocity() (float64, float64) {
	angle := rand.Float64() * math.Pi * 2
	speed := 0.6 + rand.Float64()*1.6
	return math.Cos(angle) * speed, math.Sin(angle) * speed
}

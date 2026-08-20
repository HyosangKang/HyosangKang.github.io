package lanchester

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	background = color.RGBA{5, 11, 20, 255}
	panel      = color.RGBA{7, 24, 42, 255}
	panelLight = color.RGBA{10, 35, 58, 255}
	grid       = color.RGBA{20, 53, 81, 255}
	white      = color.RGBA{242, 244, 239, 255}
	muted      = color.RGBA{128, 156, 180, 255}
	cyan       = color.RGBA{77, 215, 229, 255}
	gold       = color.RGBA{244, 197, 66, 255}
	teamColors = []color.RGBA{
		{94, 224, 177, 255},
		{255, 111, 145, 255},
		{255, 200, 87, 255},
		{97, 165, 255, 255},
		{184, 140, 255, 255},
		{255, 149, 95, 255},
	}
	teamNames = []string{"JADE", "CORAL", "GOLD", "AZURE", "VIOLET", "EMBER"}
)

var (
	//go:embed arcadeclassic.ttf
	arcadeClassic []byte
	arcadeFont    *text.GoTextFaceSource
)

func init() {
	arcadeFont, _ = text.NewGoTextFaceSource(bytes.NewReader(arcadeClassic))
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(background)
	switch g.scene {
	case sceneSetup:
		g.drawSetup(screen)
	case sceneBattle:
		g.drawBattle(screen)
	case sceneResult:
		g.drawResult(screen)
	}
}

func (g *Game) drawSetup(screen *ebiten.Image) {
	drawText(screen, "01  CONFIGURE", 16, 30, 28, cyan, text.AlignStart)
	drawText(screen, "SHAPE  THE  FORCES", 34, 30, 56, white, text.AlignStart)
	drawText(screen, "CLICK  A  NUMBER  TO  TYPE  AN  EXACT  VALUE", 14, 30, 96, muted, text.AlignStart)

	drawText(screen, "TEAMS", 13, 681, 27, muted, text.AlignStart)
	drawButton(screen, teamCountMinus, "-", cyan, false)
	drawButton(screen, teamCountValue, strconv.Itoa(g.teamCount), white, false)
	drawButton(screen, teamCountPlus, "+", cyan, false)

	for team := 0; team < g.teamCount; team++ {
		g.drawTeamCard(screen, team)
	}
	for team := g.teamCount; team < 6; team++ {
		card := teamCardRect(team)
		vector.DrawFilledRect(screen, card.x, card.y, card.w, card.h, color.RGBA{6, 17, 29, 130}, false)
		vector.StrokeRect(screen, card.x, card.y, card.w, card.h, 1, color.RGBA{20, 53, 81, 120}, false)
	}

	drawText(screen, "2-6  TEAMS     1-120  MINIONS  EACH", 15, 30, 568, muted, text.AlignStart)
	drawText(screen, "TEAM  CLUSTERS  SPAWN  AT  EVEN  ANGLES", 13, 30, 593, muted, text.AlignStart)
	drawButton(screen, startButton, "START  BATTLE", cyan, true)
}

func (g *Game) drawTeamCard(screen *ebiten.Image, team int) {
	card := teamCardRect(team)
	teamColor := teamColors[team]
	vector.DrawFilledRect(screen, card.x, card.y, card.w, card.h, panel, false)
	vector.StrokeRect(screen, card.x, card.y, card.w, card.h, 1, grid, false)
	vector.DrawFilledRect(screen, card.x, card.y, card.w, 3, teamColor, false)
	vector.DrawFilledCircle(screen, card.x+24, card.y+30, 5, teamColor, false)
	drawText(screen, teamNames[team]+"  TEAM", 18, float64(card.x+40), float64(card.y+19), white, text.AlignStart)
	drawText(screen, "MINIONS", 12, float64(card.x+18), float64(card.y+57), muted, text.AlignStart)
	drawButton(screen, teamMinusRect(team), "-", teamColor, false)

	value := strconv.Itoa(g.sizes[team])
	editing := g.editingTeam == team
	if editing {
		value = g.editBuffer + "_"
	}
	drawButton(screen, teamValueRect(team), value, teamColor, editing)
	drawButton(screen, teamPlusRect(team), "+", teamColor, false)
}

func (g *Game) drawBattle(screen *ebiten.Image) {
	drawText(screen, "02  SIMULATE", 15, 30, 20, cyan, text.AlignStart)
	drawText(screen, "BATTLE  IN  PROGRESS", 26, 30, 45, white, text.AlignStart)
	drawText(screen, fmt.Sprintf("TICK  %04d", g.snapshot.Tick), 18, 930, 25, gold, text.AlignEnd)

	chipWidth := float32(860 / max(1, len(g.snapshot.Counts)))
	for team, count := range g.snapshot.Counts {
		x := float32(30) + float32(team)*chipWidth
		vector.DrawFilledRect(screen, x, 78, chipWidth-6, 24, panel, false)
		vector.DrawFilledRect(screen, x, 78, 4, 24, teamColors[team], false)
		drawText(screen, fmt.Sprintf("%s  %d", teamNames[team], count), 13, float64(x+12), 82, white, text.AlignStart)
	}

	vector.DrawFilledRect(screen, arenaX, arenaY, arenaWidth, arenaHeight, panel, false)
	for x := arenaX; x <= arenaX+arenaWidth; x += 45 {
		vector.StrokeLine(screen, float32(x), arenaY, float32(x), arenaY+arenaHeight, 1, grid, false)
	}
	for y := arenaY; y <= arenaY+arenaHeight; y += 45.6 {
		vector.StrokeLine(screen, arenaX, float32(y), arenaX+arenaWidth, float32(y), 1, grid, false)
	}
	g.drawSpawnArcs(screen)

	for _, item := range g.lasers {
		laserColor := withAlpha(teamColors[item.team], uint8(80+item.life*175))
		vector.StrokeLine(
			screen,
			float32(arenaX+item.fromX), float32(arenaY+item.fromY),
			float32(arenaX+item.toX), float32(arenaY+item.toY),
			4, withAlpha(laserColor, 70), false,
		)
		vector.StrokeLine(
			screen,
			float32(arenaX+item.fromX), float32(arenaY+item.fromY),
			float32(arenaX+item.toX), float32(arenaY+item.toY),
			1.4, laserColor, false,
		)
	}

	for _, unit := range g.snapshot.Units {
		itemColor := teamColors[unit.Team]
		x := float32(arenaX + unit.X)
		y := float32(arenaY + unit.Y)
		vector.DrawFilledCircle(screen, x, y, 6, withAlpha(itemColor, 50), false)
		vector.DrawFilledCircle(screen, x, y, 3.2, itemColor, false)
		if unit.Health < 60 {
			vector.DrawFilledRect(screen, x-5, y-8, 10, 1.5, color.RGBA{40, 54, 66, 255}, false)
			vector.DrawFilledRect(screen, x-5, y-8, float32(unit.Health/10), 1.5, white, false)
		}
	}

	for _, item := range g.particles {
		vector.DrawFilledCircle(
			screen,
			float32(arenaX+item.x), float32(arenaY+item.y),
			2, withAlpha(teamColors[item.team], uint8(item.life*255)), false,
		)
	}
	if g.paused {
		vector.DrawFilledRect(screen, arenaX, arenaY, arenaWidth, arenaHeight, color.RGBA{2, 7, 13, 180}, false)
		drawText(screen, "PAUSED", 44, Width/2, Height/2-20, gold, text.AlignCenter)
		drawText(screen, "SPACE  OR  CLICK  RESUME", 16, Width/2, Height/2+34, white, text.AlignCenter)
	}

	drawButton(screen, editButton, "EDIT  TEAMS", cyan, false)
	pauseLabel := "PAUSE  [SPACE]"
	if g.paused {
		pauseLabel = "RESUME  [SPACE]"
	}
	drawButton(screen, pauseButton, pauseLabel, gold, false)
	drawText(screen, "R  RESTART", 13, Width/2, 599, muted, text.AlignCenter)
}

func (g *Game) drawSpawnArcs(screen *ebiten.Image) {
	centerX := float64(arenaX + arenaWidth/2)
	centerY := float64(arenaY + arenaHeight/2)
	radius := math.Min(arenaWidth*0.37, arenaHeight*0.35)
	segment := 2 * math.Pi / float64(len(g.snapshot.Counts))
	for team := range g.snapshot.Counts {
		start := -math.Pi/2 + segment*float64(team) - segment*0.32
		end := start + segment*0.64
		previousX := centerX + math.Cos(start)*radius
		previousY := centerY + math.Sin(start)*radius
		for part := 1; part <= 18; part++ {
			angle := start + (end-start)*float64(part)/18
			x := centerX + math.Cos(angle)*radius
			y := centerY + math.Sin(angle)*radius
			vector.StrokeLine(
				screen, float32(previousX), float32(previousY), float32(x), float32(y),
				2, withAlpha(teamColors[team], 100), false,
			)
			previousX, previousY = x, y
		}
	}
}

func (g *Game) drawResult(screen *ebiten.Image) {
	drawText(screen, "03  COMPARE", 15, 30, 22, cyan, text.AlignStart)
	drawText(screen, "BATTLE  COMPLETE", 32, 30, 49, white, text.AlignStart)
	if g.snapshot.Winner >= 0 {
		winner := g.snapshot.Winner
		drawText(screen, teamNames[winner]+"  SURVIVES", 20, 930, 30, teamColors[winner], text.AlignEnd)
	} else {
		drawText(screen, "DRAW", 20, 930, 30, gold, text.AlignEnd)
	}
	drawText(screen, fmt.Sprintf("%d  TICKS     SAME  COLOURS,  SAME  TEAMS", g.snapshot.Tick), 14, 30, 91, muted, text.AlignStart)

	g.drawGraph(screen)
	g.drawFinalCounts(screen)
	drawButton(screen, resultEdit, "EDIT  TEAMS", cyan, false)
	drawButton(screen, resultAgain, "RUN  AGAIN  [ENTER]", gold, true)
}

func (g *Game) drawGraph(screen *ebiten.Image) {
	panelRect := rect{x: 30, y: 118, w: 900, h: 368}
	vector.DrawFilledRect(screen, panelRect.x, panelRect.y, panelRect.w, panelRect.h, panel, false)
	vector.StrokeRect(screen, panelRect.x, panelRect.y, panelRect.w, panelRect.h, 1, grid, false)
	drawText(screen, "MINIONS  REMAINING  /  SIMULATION  TICK", 15, 52, 136, white, text.AlignStart)

	left, top := float32(76), float32(178)
	plotWidth, plotHeight := float32(818), float32(252)
	maxCount := 1
	for team := 0; team < g.teamCount; team++ {
		maxCount = max(maxCount, g.sizes[team])
	}
	for division := 0; division <= 4; division++ {
		y := top + plotHeight - plotHeight*float32(division)/4
		vector.StrokeLine(screen, left, y, left+plotWidth, y, 1, grid, false)
		value := maxCount * division / 4
		drawText(screen, strconv.Itoa(value), 12, float64(left-12), float64(y-7), muted, text.AlignEnd)
	}

	history := g.snapshot.History
	lastTick := max(1, len(history)-1)
	for team := 0; team < g.teamCount; team++ {
		previousX := left
		previousY := top + plotHeight - plotHeight*float32(history[0][team])/float32(maxCount)
		for tick := 1; tick < len(history); tick++ {
			x := left + plotWidth*float32(tick)/float32(lastTick)
			y := top + plotHeight - plotHeight*float32(history[tick][team])/float32(maxCount)
			vector.StrokeLine(screen, previousX, previousY, x, y, 2.5, teamColors[team], false)
			previousX, previousY = x, y
		}
	}
	for _, tick := range []int{0, lastTick / 2, lastTick} {
		x := left + plotWidth*float32(tick)/float32(lastTick)
		drawText(screen, strconv.Itoa(tick), 12, float64(x), float64(top+plotHeight+15), muted, text.AlignCenter)
	}
}

func (g *Game) drawFinalCounts(screen *ebiten.Image) {
	x := float64(30)
	for team, count := range g.snapshot.Counts {
		vector.DrawFilledCircle(screen, float32(x+5), 522, 5, teamColors[team], false)
		drawText(screen, fmt.Sprintf("%s  %d", teamNames[team], count), 13, x+18, 513, white, text.AlignStart)
		x += 138
	}
}

func drawButton(screen *ebiten.Image, bounds rect, label string, itemColor color.RGBA, strong bool) {
	x, y := ebiten.CursorPosition()
	hovered := bounds.contains(x, y)
	fill := panel
	if strong {
		fill = color.RGBA{12, 53, 63, 255}
	}
	if hovered {
		fill = panelLight
	}
	vector.DrawFilledRect(screen, bounds.x, bounds.y, bounds.w, bounds.h, fill, false)
	lineWidth := float32(1)
	if hovered || strong {
		lineWidth = 2
	}
	vector.StrokeRect(screen, bounds.x, bounds.y, bounds.w, bounds.h, lineWidth, itemColor, false)
	if label == "-" || label == "+" {
		centerX := bounds.x + bounds.w/2
		centerY := bounds.y + bounds.h/2
		vector.StrokeLine(screen, centerX-7, centerY, centerX+7, centerY, 2, itemColor, false)
		if label == "+" {
			vector.StrokeLine(screen, centerX, centerY-7, centerX, centerY+7, 2, itemColor, false)
		}
		return
	}
	drawText(screen, label, 17, float64(bounds.x+bounds.w/2), float64(bounds.y+bounds.h/2-9), itemColor, text.AlignCenter)
}

func drawText(screen *ebiten.Image, message string, size, x, y float64, itemColor color.Color, align text.Align) {
	options := &text.DrawOptions{}
	options.GeoM.Translate(x, y)
	options.ColorScale.ScaleWithColor(itemColor)
	options.PrimaryAlign = align
	options.SecondaryAlign = text.AlignStart
	text.Draw(screen, message, &text.GoTextFace{Source: arcadeFont, Size: size}, options)
}

func withAlpha(item color.RGBA, alpha uint8) color.RGBA {
	item.A = alpha
	return item
}

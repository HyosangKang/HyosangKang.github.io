package labs

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type springPreset struct {
	name   string
	params SpringParams
}

type springField struct {
	label   string
	minimum float64
	maximum float64
	step    float64
	read    func(SpringParams) float64
	write   func(*SpringParams, float64)
}

type SpringLab struct {
	params        SpringParams
	model         *SpringModel
	time          float64
	playing       bool
	tab           int
	selected      int
	preset        int
	systemFields  []springField
	initialFields []springField
	presets       []springPreset
	lastError     string
}

func NewSpringLab() *SpringLab {
	lab := &SpringLab{
		playing: true,
		presets: []springPreset{
			{"TRANSFER", SpringParams{M1: 1, M2: 1, K1: 3, K2: 2, X10: 1, V10: 0, X20: 0, V20: 0, T: 18}},
			{"TOGETHER", SpringParams{M1: 1, M2: 1, K1: 3, K2: 2, X10: 0.8, V10: 0, X20: 0.8, V20: 0, T: 18}},
			{"OPPOSITE", SpringParams{M1: 1, M2: 1, K1: 3, K2: 2, X10: 0.9, V10: 0, X20: -0.9, V20: 0, T: 18}},
		},
	}
	lab.systemFields = []springField{
		{"MASS M1", 0.1, 10, 0.1, func(p SpringParams) float64 { return p.M1 }, func(p *SpringParams, v float64) { p.M1 = v }},
		{"MASS M2", 0.1, 10, 0.1, func(p SpringParams) float64 { return p.M2 }, func(p *SpringParams, v float64) { p.M2 = v }},
		{"SPRING K1", 0.1, 20, 0.1, func(p SpringParams) float64 { return p.K1 }, func(p *SpringParams, v float64) { p.K1 = v }},
		{"SPRING K2", 0.1, 20, 0.1, func(p SpringParams) float64 { return p.K2 }, func(p *SpringParams, v float64) { p.K2 = v }},
		{"TIME RANGE", 4, 40, 1, func(p SpringParams) float64 { return p.T }, func(p *SpringParams, v float64) { p.T = v }},
	}
	lab.initialFields = []springField{
		{"X1 AT ZERO", -3, 3, 0.1, func(p SpringParams) float64 { return p.X10 }, func(p *SpringParams, v float64) { p.X10 = v }},
		{"V1 AT ZERO", -5, 5, 0.1, func(p SpringParams) float64 { return p.V10 }, func(p *SpringParams, v float64) { p.V10 = v }},
		{"X2 AT ZERO", -3, 3, 0.1, func(p SpringParams) float64 { return p.X20 }, func(p *SpringParams, v float64) { p.X20 = v }},
		{"V2 AT ZERO", -5, 5, 0.1, func(p SpringParams) float64 { return p.V20 }, func(p *SpringParams, v float64) { p.V20 = v }},
	}
	lab.applyPreset(0)
	return lab
}

func (lab *SpringLab) fields() []springField {
	if lab.tab == 1 {
		return lab.initialFields
	}
	return lab.systemFields
}

func (lab *SpringLab) Update() error {
	if lab.playing {
		lab.time += 1.0 / 60.0
		if lab.time > lab.params.T {
			lab.time = 0
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		lab.playing = !lab.playing
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		lab.time = 0
		lab.playing = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		lab.tab = 1 - lab.tab
		lab.selected = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		lab.selected = (lab.selected + len(lab.fields()) - 1) % len(lab.fields())
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		lab.selected = (lab.selected + 1) % len(lab.fields())
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		lab.adjustField(lab.selected, -1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		lab.adjustField(lab.selected, 1)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		lab.applyPreset(0)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		lab.applyPreset(1)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		lab.applyPreset(2)
	}
	if x, y, pressed := justPressedPointer(); pressed {
		lab.handleClick(x, y)
	}
	return nil
}

func (lab *SpringLab) handleClick(x, y int) {
	for index := range lab.presets {
		button := Button{X: float32(12 + index*66), Y: 62, W: 61, H: 25}
		if button.Contains(x, y) {
			lab.applyPreset(index)
			return
		}
	}
	if (Button{X: 12, Y: 94, W: 92, H: 27}).Contains(x, y) {
		lab.tab = 0
		lab.selected = 0
		return
	}
	if (Button{X: 112, Y: 94, W: 92, H: 27}).Contains(x, y) {
		lab.tab = 1
		lab.selected = 0
		return
	}
	for index := range lab.fields() {
		rowY := float32(149 + index*48)
		if (Button{X: 122, Y: rowY, W: 32, H: 27}).Contains(x, y) {
			lab.selected = index
			lab.adjustField(index, -1)
			return
		}
		if (Button{X: 172, Y: rowY, W: 32, H: 27}).Contains(x, y) {
			lab.selected = index
			lab.adjustField(index, 1)
			return
		}
		if x >= 10 && x <= 210 && y >= int(rowY-18) && y <= int(rowY+30) {
			lab.selected = index
			return
		}
	}
	if (Button{X: 12, Y: 398, W: 91, H: 31}).Contains(x, y) {
		lab.playing = !lab.playing
		return
	}
	if (Button{X: 113, Y: 398, W: 91, H: 31}).Contains(x, y) {
		lab.time = 0
		lab.playing = true
		return
	}
	if x >= 246 && x <= 774 && y >= 300 && y <= 446 {
		lab.time = clamp(float64(x-246)/528*lab.params.T, 0, lab.params.T)
		lab.playing = false
	}
}

func (lab *SpringLab) adjustField(index, direction int) {
	field := lab.fields()[index]
	value := field.read(lab.params) + float64(direction)*field.step
	field.write(&lab.params, clamp(value, field.minimum, field.maximum))
	lab.preset = -1
	lab.rebuild()
}

func (lab *SpringLab) applyPreset(index int) {
	lab.params = lab.presets[index].params
	lab.preset = index
	lab.time = 0
	lab.playing = true
	lab.rebuild()
}

func (lab *SpringLab) rebuild() {
	model, err := NewSpringModel(lab.params)
	if err != nil {
		lab.lastError = err.Error()
		return
	}
	lab.lastError = ""
	lab.model = model
	if lab.time > lab.params.T {
		lab.time = lab.params.T
	}
}

func (lab *SpringLab) Draw(screen *ebiten.Image) {
	lab.drawControls(screen)
	lab.drawSpringSystem(screen)
	lab.drawGraph(screen)
}

func (lab *SpringLab) drawControls(screen *ebiten.Image) {
	DrawPanel(screen, 7, 8, 207, 464)
	DrawText(screen, "COUPLED MOTION", 22, 14, 15, White, text.AlignStart)
	DrawText(screen, "DOUBLE SPRING LAB", 11, 14, 42, Cyan, text.AlignStart)
	for index, preset := range lab.presets {
		Button{X: float32(12 + index*66), Y: 62, W: 61, H: 25, Label: fmt.Sprintf("%d %s", index+1, preset.name), Color: Cyan, Active: index == lab.preset}.Draw(screen, 8)
	}
	Button{X: 12, Y: 94, W: 92, H: 27, Label: "SYSTEM", Color: Gold, Active: lab.tab == 0}.Draw(screen, 11)
	Button{X: 112, Y: 94, W: 92, H: 27, Label: "INITIAL", Color: Gold, Active: lab.tab == 1}.Draw(screen, 11)

	for index, field := range lab.fields() {
		y := float64(143 + index*48)
		itemColor := Blue
		if index == lab.selected {
			itemColor = Gold
			vector.DrawFilledRect(screen, 11, float32(y-5), 197, 41, color.RGBA{16, 42, 63, 210}, false)
		}
		DrawText(screen, field.label, 11, 15, y, itemColor, text.AlignStart)
		DrawText(screen, formatValue(field.read(lab.params), field.step), 16, 114, y+13, White, text.AlignEnd)
		Button{X: 122, Y: float32(y + 6), W: 32, H: 27, Label: "-", Color: Red}.Draw(screen, 17)
		Button{X: 172, Y: float32(y + 6), W: 32, H: 27, Label: "+", Color: Green}.Draw(screen, 17)
	}

	playLabel := "PAUSE"
	if !lab.playing {
		playLabel = "PLAY"
	}
	Button{X: 12, Y: 398, W: 91, H: 31, Label: playLabel, Color: Gold, Active: lab.playing}.Draw(screen, 12)
	Button{X: 113, Y: 398, W: 91, H: 31, Label: "RESET", Color: Cyan}.Draw(screen, 12)
	DrawText(screen, "TAB PAGE    ARROWS CHANGE", 9, 13, 437, Blue, text.AlignStart)
	DrawText(screen, "SPACE PAUSE    R RESET", 9, 13, 452, Blue, text.AlignStart)
	if lab.lastError != "" {
		DrawText(screen, "INVALID PARAMETERS", 9, 13, 462, Red, text.AlignStart)
	}
}

func (lab *SpringLab) drawSpringSystem(screen *ebiten.Image) {
	DrawPanel(screen, 224, 8, 569, 268)
	DrawText(screen, "TWO MASS SPRING SYSTEM", 17, 240, 18, White, text.AlignStart)
	if lab.model == nil {
		return
	}
	sample := lab.model.Motion(lab.time)
	scale := 52.0 / lab.model.MaxAbs
	x1 := 470 + sample.X1*scale
	x2 := 665 + sample.X2*scale
	x1 = clamp(x1, 390, 530)
	x2 = clamp(x2, 585, 744)

	DrawText(screen, fmt.Sprintf("X1  %+.3f M", sample.X1), 13, 777, 18, Cyan, text.AlignEnd)
	DrawText(screen, fmt.Sprintf("X2  %+.3f M", sample.X2), 13, 777, 37, Red, text.AlignEnd)
	DrawText(screen, fmt.Sprintf("ENERGY  %.3f J", sample.Energy), 11, 777, 57, Gold, text.AlignEnd)

	vector.DrawFilledRect(screen, 252, 82, 18, 146, Blue, false)
	for y := 86; y <= 222; y += 18 {
		drawLine(screen, 246, float64(y+8), 252, float64(y), 1, Blue)
	}
	drawSpring(screen, 270, x1-26, 154, Cyan)
	drawSpring(screen, x1+26, x2-26, 154, Gold)
	vector.DrawFilledRect(screen, float32(x1-26), 126, 52, 56, Cyan, false)
	vector.DrawFilledRect(screen, float32(x2-26), 126, 52, 56, Red, false)
	vector.StrokeRect(screen, float32(x1-26), 126, 52, 56, 2, White, false)
	vector.StrokeRect(screen, float32(x2-26), 126, 52, 56, 2, White, false)
	DrawText(screen, "M1", 14, x1, 145, Black, text.AlignCenter)
	DrawText(screen, "M2", 14, x2, 145, Black, text.AlignCenter)
	drawLine(screen, 240, 205, 770, 205, 2, Blue)
	for x := 250; x < 770; x += 16 {
		drawLine(screen, float64(x), 205, float64(x-8), 215, 1, Blue)
	}
	drawLine(screen, 470, 115, 470, 210, 1, GridBlue)
	drawLine(screen, 665, 115, 665, 210, 1, GridBlue)
	DrawText(screen, "K1", 12, (270+x1-26)/2, 105, Cyan, text.AlignCenter)
	DrawText(screen, "K2", 12, (x1+x2)/2, 105, Gold, text.AlignCenter)
	DrawText(screen, fmt.Sprintf("NORMAL MODES  W1 %.3f   W2 %.3f", lab.model.Slow, lab.model.Fast), 12, 505, 239, White, text.AlignCenter)
	DrawText(screen, "THE EXACT MODAL SOLUTION MOVES BOTH MASSES", 9, 505, 258, Blue, text.AlignCenter)
}

func (lab *SpringLab) drawGraph(screen *ebiten.Image) {
	DrawPanel(screen, 224, 284, 569, 188)
	DrawText(screen, "DISPLACEMENT VS TIME", 13, 240, 290, White, text.AlignStart)
	DrawText(screen, "CLICK GRAPH TO SCRUB", 9, 777, 292, Blue, text.AlignEnd)
	if lab.model == nil {
		return
	}
	const x0, y0, width, height = 246.0, 300.0, 528.0, 146.0
	middleY := y0 + height/2
	drawLine(screen, x0, middleY, x0+width, middleY, 1, GridBlue)
	for index := 0; index <= 4; index++ {
		x := x0 + width*float64(index)/4
		drawLine(screen, x, y0, x, y0+height, 1, GridBlue)
		DrawText(screen, fmt.Sprintf("%.0f", lab.params.T*float64(index)/4), 8, x, 451, Blue, text.AlignCenter)
	}
	for index := 3; index < len(lab.model.Samples); index += 3 {
		previous := lab.model.Samples[index-3]
		current := lab.model.Samples[index]
		xPrevious := x0 + previous.T/lab.params.T*width
		xCurrent := x0 + current.T/lab.params.T*width
		y1Previous := middleY - previous.X1/lab.model.MaxAbs*height/2
		y1Current := middleY - current.X1/lab.model.MaxAbs*height/2
		y2Previous := middleY - previous.X2/lab.model.MaxAbs*height/2
		y2Current := middleY - current.X2/lab.model.MaxAbs*height/2
		drawLine(screen, xPrevious, y1Previous, xCurrent, y1Current, 1.5, Cyan)
		drawLine(screen, xPrevious, y2Previous, xCurrent, y2Current, 1.5, Red)
	}
	sample := lab.model.Motion(lab.time)
	markerX := x0 + lab.time/lab.params.T*width
	drawLine(screen, markerX, y0, markerX, y0+height, 1, Gold)
	vector.DrawFilledCircle(screen, float32(markerX), float32(middleY-sample.X1/lab.model.MaxAbs*height/2), 5, Cyan, false)
	vector.DrawFilledCircle(screen, float32(markerX), float32(middleY-sample.X2/lab.model.MaxAbs*height/2), 5, Red, false)
	DrawText(screen, fmt.Sprintf("T %.2f / %.0f", lab.time, lab.params.T), 10, 777, 454, Gold, text.AlignEnd)
}

func drawSpring(screen *ebiten.Image, x0, x1, y float64, itemColor color.Color) {
	if x1 <= x0 {
		drawLine(screen, x0, y, x1, y, 2, itemColor)
		return
	}
	const segments = 18
	lead := math.Min(14, (x1-x0)*0.12)
	drawLine(screen, x0, y, x0+lead, y, 2, itemColor)
	previousX, previousY := x0+lead, y
	coilWidth := x1 - x0 - 2*lead
	for index := 1; index <= segments; index++ {
		x := x0 + lead + coilWidth*float64(index)/segments
		yOffset := 0.0
		if index < segments {
			yOffset = 13
			if index%2 == 0 {
				yOffset = -13
			}
		}
		drawLine(screen, previousX, previousY, x, y+yOffset, 2, itemColor)
		previousX, previousY = x, y+yOffset
	}
	drawLine(screen, previousX, previousY, x1, y, 2, itemColor)
}

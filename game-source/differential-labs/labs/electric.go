package labs

import (
	"differential-labs/gametext"
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type electricPreset struct {
	name   string
	params ElectricParams
}

type electricField struct {
	label   string
	minimum float64
	maximum float64
	step    float64
	read    func(ElectricParams) float64
	write   func(*ElectricParams, float64)
}

type ElectricLab struct {
	params    ElectricParams
	model     *ElectricModel
	time      float64
	playing   bool
	selected  int
	preset    int
	charge1   float64
	charge2   float64
	fields    []electricField
	presets   []electricPreset
	lastError string
}

func NewElectricLab() *ElectricLab {
	lab := &ElectricLab{
		playing: true,
		presets: []electricPreset{
			{"BALANCED", ElectricParams{L: 1, R1: 1, R2: 2, C: 1, T: 18}},
			{"RINGING", ElectricParams{L: 2.4, R1: 0.25, R2: 0.35, C: 0.8, T: 24}},
			{"RESISTIVE", ElectricParams{L: 0.45, R1: 4.2, R2: 5.5, C: 0.7, T: 14}},
		},
	}
	lab.fields = []electricField{
		{"INDUCTOR L X100", 0.05, 12, 0.05, func(p ElectricParams) float64 { return p.L }, func(p *ElectricParams, v float64) { p.L = v }},
		{"RESISTOR R1 X100", 0.05, 12, 0.05, func(p ElectricParams) float64 { return p.R1 }, func(p *ElectricParams, v float64) { p.R1 = v }},
		{"RESISTOR R2 X100", 0.05, 12, 0.05, func(p ElectricParams) float64 { return p.R2 }, func(p *ElectricParams, v float64) { p.R2 = v }},
		{"CAPACITOR C X100", 0.05, 12, 0.05, func(p ElectricParams) float64 { return p.C }, func(p *ElectricParams, v float64) { p.C = v }},
		{"TIME RANGE", 4, 40, 1, func(p ElectricParams) float64 { return p.T }, func(p *ElectricParams, v float64) { p.T = v }},
	}
	lab.applyPreset(0)
	return lab
}

func (lab *ElectricLab) Update() error {
	if lab.playing {
		const step = 1.0 / 60.0
		if lab.model != nil {
			i1, i2 := lab.model.Current(lab.time)
			lab.charge1 += i1 * step * 0.8
			lab.charge2 += i2 * step * 0.8
		}
		lab.time += step
		if lab.time > lab.params.T {
			lab.time = 0
			lab.charge1 = 0
			lab.charge2 = 0
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		lab.playing = !lab.playing
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		lab.time = 0
		lab.charge1 = 0
		lab.charge2 = 0
		lab.playing = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		lab.selected = (lab.selected + len(lab.fields) - 1) % len(lab.fields)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		lab.selected = (lab.selected + 1) % len(lab.fields)
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

func (lab *ElectricLab) handleClick(x, y int) {
	for index := range lab.presets {
		button := Button{X: float32(12 + index*66), Y: 62, W: 61, H: 25}
		if button.Contains(x, y) {
			lab.applyPreset(index)
			return
		}
	}
	for index := range lab.fields {
		rowY := float32(111 + index*48)
		if (Button{X: 122, Y: rowY, W: 38, H: 29}).Contains(x, y) {
			lab.selected = index
			lab.adjustField(index, -1)
			return
		}
		if (Button{X: 166, Y: rowY, W: 38, H: 29}).Contains(x, y) {
			lab.selected = index
			lab.adjustField(index, 1)
			return
		}
		if x >= 10 && x <= 210 && y >= int(rowY-18) && y <= int(rowY+30) {
			lab.selected = index
			return
		}
	}
	if (Button{X: 12, Y: 371, W: 91, H: 31}).Contains(x, y) {
		lab.playing = !lab.playing
		return
	}
	if (Button{X: 113, Y: 371, W: 91, H: 31}).Contains(x, y) {
		lab.time = 0
		lab.charge1 = 0
		lab.charge2 = 0
		lab.playing = true
		return
	}
	if x >= 246 && x <= 774 && y >= 300 && y <= 446 {
		lab.time = clamp(float64(x-246)/528*lab.params.T, 0, lab.params.T)
		lab.playing = false
	}
}

func (lab *ElectricLab) adjustField(index, direction int) {
	field := lab.fields[index]
	value := field.read(lab.params) + float64(direction)*field.step
	field.write(&lab.params, clamp(value, field.minimum, field.maximum))
	lab.preset = -1
	lab.rebuild()
}

func (lab *ElectricLab) applyPreset(index int) {
	lab.params = lab.presets[index].params
	lab.preset = index
	lab.time = 0
	lab.charge1 = 0
	lab.charge2 = 0
	lab.playing = true
	lab.rebuild()
}

func (lab *ElectricLab) rebuild() {
	model, err := NewElectricModel(lab.params)
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

func (lab *ElectricLab) Draw(screen *ebiten.Image) {
	lab.drawControls(screen)
	lab.drawCircuit(screen)
	lab.drawGraph(screen)
}

func (lab *ElectricLab) drawControls(screen *ebiten.Image) {
	DrawPanel(screen, 7, 8, 207, 464)
	DrawText(screen, "CURRENT QUEST", 25, 14, 13, White, text.AlignStart)
	DrawText(screen, "LAPLACE CIRCUIT LAB", 13, 14, 42, Cyan, text.AlignStart)

	for index := range lab.presets {
		shortName := []string{"BAL", "RING", "RES"}[index]
		button := Button{X: float32(12 + index*66), Y: 62, W: 61, H: 27, Label: fmt.Sprintf("%d %s", index+1, shortName), Color: Cyan, Active: index == lab.preset}
		button.Draw(screen, 11)
	}
	DrawText(screen, "ARROWS SELECT AND CHANGE", 10, 13, 94, Blue, text.AlignStart)

	for index, field := range lab.fields {
		y := float64(105 + index*48)
		itemColor := Blue
		if index == lab.selected {
			itemColor = Gold
			vector.DrawFilledRect(screen, 11, float32(y-5), 197, 41, color.RGBA{16, 42, 63, 210}, false)
		}
		DrawText(screen, field.label, 13, 15, y-1, itemColor, text.AlignStart)
		value := field.read(lab.params)
		DrawText(screen, formatValue(value, field.step), 19, 114, y+13, White, text.AlignEnd)
		Button{X: 122, Y: float32(y + 6), W: 38, H: 29, Label: "DEC", Color: Red}.Draw(screen, 10)
		Button{X: 166, Y: float32(y + 6), W: 38, H: 29, Label: "INC", Color: Green}.Draw(screen, 10)
	}

	playLabel := "PAUSE"
	if !lab.playing {
		playLabel = "PLAY"
	}
	Button{X: 12, Y: 371, W: 91, H: 33, Label: playLabel, Color: Gold, Active: lab.playing}.Draw(screen, 14)
	Button{X: 113, Y: 371, W: 91, H: 33, Label: "RESET", Color: Cyan}.Draw(screen, 14)
	DrawText(screen, "SPACE PAUSE  R RESET", 11, 13, 410, Blue, text.AlignStart)
	DrawText(screen, "SOURCE SINE WAVE", 14, 13, 430, White, text.AlignStart)
	DrawText(screen, "ZERO INITIAL CURRENTS", 11, 13, 451, Cyan, text.AlignStart)
	if lab.lastError != "" {
		DrawText(screen, "INVALID PARAMETERS", 11, 13, 463, Red, text.AlignStart)
	}
}

func (lab *ElectricLab) drawCircuit(screen *ebiten.Image) {
	DrawPanel(screen, 224, 8, 569, 268)
	DrawText(screen, "DOUBLE ELECTRIC CIRCUIT", 20, 240, 16, White, text.AlignStart)
	if lab.model == nil {
		return
	}
	i1, i2 := lab.model.Current(lab.time)
	DrawText(screen, fmt.Sprintf("I1 %s MA", gametext.Signed(i1, 1000)), 16, 777, 18, Cyan, text.AlignEnd)
	DrawText(screen, fmt.Sprintf("I2 %s MA", gametext.Signed(i2, 1000)), 16, 777, 40, Red, text.AlignEnd)

	left, middle, right := 275.0, 505.0, 735.0
	top, bottom := 78.0, 222.0
	drawLine(screen, left, top, 343, top, 3, Cyan)
	drawInductor(screen, 343, top, 429, Cyan)
	drawLine(screen, 429, top, middle, top, 3, Cyan)
	drawLine(screen, middle, top, 563, top, 3, Red)
	drawResistor(screen, 563, top, 650, top, Red)
	drawLine(screen, 650, top, right, top, 3, Red)
	drawLine(screen, left, top, left, 118, 3, Cyan)
	drawVoltageSource(screen, left, 150)
	drawLine(screen, left, 182, left, bottom, 3, Cyan)
	drawLine(screen, left, bottom, middle, bottom, 3, Cyan)
	drawLine(screen, middle, bottom, right, bottom, 3, Red)
	drawLine(screen, right, top, right, 116, 3, Red)
	drawResistor(screen, right, 116, right, 184, Red)
	drawLine(screen, right, 184, right, bottom, 3, Red)
	drawLine(screen, middle, top, middle, 132, 3, Gold)
	drawCapacitor(screen, middle, 132, 168)
	drawLine(screen, middle, 168, middle, bottom, 3, Gold)

	DrawText(screen, "L", 16, 385, 48, Cyan, text.AlignCenter)
	DrawText(screen, "R1", 16, 606, 48, Red, text.AlignCenter)
	DrawText(screen, "R2", 16, 762, 142, Red, text.AlignCenter)
	DrawText(screen, "C", 16, 524, 140, Gold, text.AlignStart)
	DrawText(screen, "V", 16, 249, 139, Gold, text.AlignCenter)
	DrawText(screen, "I1", 16, 383, 233, Cyan, text.AlignCenter)
	DrawText(screen, "I2", 16, 627, 233, Red, text.AlignCenter)

	drawChargeDots(screen, left, top, middle, bottom, lab.charge1, Cyan)
	drawChargeDots(screen, middle, top, right, bottom, lab.charge2+0.08, Red)
	DrawText(screen, lab.model.Response, 14, 777, 248, Gold, text.AlignEnd)
}

func (lab *ElectricLab) drawGraph(screen *ebiten.Image) {
	DrawPanel(screen, 224, 284, 569, 188)
	DrawText(screen, "CURRENT VS TIME", 16, 240, 288, White, text.AlignStart)
	DrawText(screen, "CLICK GRAPH TO SCRUB", 11, 777, 291, Blue, text.AlignEnd)
	if lab.model == nil {
		return
	}
	const x0, y0, width, height = 246.0, 300.0, 528.0, 146.0
	middleY := y0 + height/2
	drawLine(screen, x0, middleY, x0+width, middleY, 1, GridBlue)
	for index := 0; index <= 4; index++ {
		x := x0 + width*float64(index)/4
		drawLine(screen, x, y0, x, y0+height, 1, GridBlue)
		DrawText(screen, fmt.Sprintf("%.0f", lab.params.T*float64(index)/4), 11, x, 449, Blue, text.AlignCenter)
	}
	for index := 3; index < len(lab.model.Samples); index += 3 {
		previous := lab.model.Samples[index-3]
		current := lab.model.Samples[index]
		xPrevious := x0 + previous.T/lab.params.T*width
		xCurrent := x0 + current.T/lab.params.T*width
		y1Previous := middleY - previous.I1/lab.model.MaxAbs*height/2
		y1Current := middleY - current.I1/lab.model.MaxAbs*height/2
		y2Previous := middleY - previous.I2/lab.model.MaxAbs*height/2
		y2Current := middleY - current.I2/lab.model.MaxAbs*height/2
		drawLine(screen, xPrevious, y1Previous, xCurrent, y1Current, 1.5, Cyan)
		drawLine(screen, xPrevious, y2Previous, xCurrent, y2Current, 1.5, Red)
	}
	i1, i2 := lab.model.Current(lab.time)
	markerX := x0 + lab.time/lab.params.T*width
	drawLine(screen, markerX, y0, markerX, y0+height, 1, Gold)
	vector.DrawFilledCircle(screen, float32(markerX), float32(middleY-i1/lab.model.MaxAbs*height/2), 5, Cyan, false)
	vector.DrawFilledCircle(screen, float32(markerX), float32(middleY-i2/lab.model.MaxAbs*height/2), 5, Red, false)
	DrawText(screen, fmt.Sprintf("TIME %d MS END %d MS", int(math.Round(lab.time*1000)), int(math.Round(lab.params.T*1000))), 12, 777, 452, Gold, text.AlignEnd)
}

func drawInductor(screen *ebiten.Image, x0, y, x1 float64, itemColor color.Color) {
	loops := 7
	step := (x1 - x0) / float64(loops)
	for index := 0; index < loops; index++ {
		start := x0 + float64(index)*step
		drawLine(screen, start, y, start+step/2, y-12, 2, itemColor)
		drawLine(screen, start+step/2, y-12, start+step, y, 2, itemColor)
	}
}

func drawResistor(screen *ebiten.Image, x0, y0, x1, y1 float64, itemColor color.Color) {
	segments := 8
	vertical := math.Abs(y1-y0) > math.Abs(x1-x0)
	previousX, previousY := x0, y0
	for index := 1; index <= segments; index++ {
		ratio := float64(index) / float64(segments)
		x := x0 + (x1-x0)*ratio
		y := y0 + (y1-y0)*ratio
		if index < segments {
			offset := 8.0
			if index%2 == 0 {
				offset = -offset
			}
			if vertical {
				x += offset
			} else {
				y += offset
			}
		}
		drawLine(screen, previousX, previousY, x, y, 2, itemColor)
		previousX, previousY = x, y
	}
}

func drawVoltageSource(screen *ebiten.Image, x, y float64) {
	vector.StrokeCircle(screen, float32(x), float32(y), 32, 3, Gold, false)
	drawLine(screen, x-11, y-7, x+11, y-7, 2, Gold)
	drawLine(screen, x, y-18, x, y+4, 2, Gold)
	drawLine(screen, x-10, y+13, x+10, y+13, 2, Gold)
}

func drawCapacitor(screen *ebiten.Image, x, y0, y1 float64) {
	drawLine(screen, x-18, y0+12, x+18, y0+12, 3, Gold)
	drawLine(screen, x-18, y1-12, x+18, y1-12, 3, Gold)
}

func drawChargeDots(screen *ebiten.Image, left, top, right, bottom, phase float64, itemColor color.Color) {
	for index := 0; index < 6; index++ {
		position := math.Mod(float64(index)/6+phase, 1)
		if position < 0 {
			position++
		}
		x, y := pointOnRectangle(left, top, right, bottom, position)
		vector.DrawFilledCircle(screen, float32(x), float32(y), 3.5, itemColor, false)
	}
}

func pointOnRectangle(left, top, right, bottom, position float64) (float64, float64) {
	width, height := right-left, bottom-top
	perimeter := 2 * (width + height)
	distance := position * perimeter
	switch {
	case distance <= width:
		return left + distance, top
	case distance <= width+height:
		return right, top + distance - width
	case distance <= 2*width+height:
		return right - (distance - width - height), bottom
	default:
		return left, bottom - (distance - 2*width - height)
	}
}

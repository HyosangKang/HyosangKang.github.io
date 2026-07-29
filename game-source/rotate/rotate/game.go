package rotate

import (
	"bytes"
	_ "embed"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hyosangkang/multi-game/rotate/example/torus"
)

const (
	width, height    = 600, 600
	maxSize, minSize = 10.0, 0.2
)

type scene int

const (
	sceneNil scene = iota
	sceneStart
	scenePlay
)

type Sprite interface {

	// Graph returns the list of pairs of 3D points.
	// The list represents the edges of the object in 3D space.
	Graph() [][2][3]float64
}

type Game struct {
	Sprite
	scene
	axis, graph matrix
	size        float64
}

func (t *Game) Update() error {
	if t.scene == sceneNil {
		t.scene = sceneStart
	}
	switch t.scene {
	case sceneStart:
		t.start()
	case scenePlay:
		t.play()
	}
	return nil
}

func (t *Game) start() {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		t.scene = scenePlay
		/* initailize axis and graph */
		t.axis = matrix{
			{1, 0, 0},
			{0, 1, 0},
			{0, 0, 1},
		}
		t.size = 2.0
		if t.Sprite == nil {
			t.Sprite = &torus.Torus{}
		}
		t.graph = matrix{}
		for _, l := range t.Sprite.Graph() {
			t.graph = append(t.graph,
				[]float64{l[0][0], l[0][1], l[0][2]},
				[]float64{l[1][0], l[1][1], l[1][2]})
		}
	}
}

const RotationTick = 0.05

var (
	controlKeys = []ebiten.Key{
		ebiten.KeyA, // rotation around z-axis
		ebiten.KeyD,
		ebiten.KeyW, // rotation around x-axis
		ebiten.KeyS,
		ebiten.KeyQ, // rotation around y-axis
		ebiten.KeyE,
	}
	controlAxis = [][2]int{
		{0, 1},
		{0, 1},
		{1, 2},
		{1, 2},
		{2, 0},
		{2, 0},
	}
	controlSign = []float64{
		-1, 1, -1, 1, -1, 1,
	}
)

func (g *Game) play() {
	for i, key := range controlKeys {
		if ebiten.IsKeyPressed(key) {
			a := controlAxis[i]
			s := controlSign[i]
			m := rotation3d(a, s*RotationTick)
			if ebiten.IsKeyPressed(ebiten.KeyShift) {
				g.axis = g.axis.mul(m)
				g.graph = g.graph.mul(m)
			} else {
				m = g.axis.t().mul(m.t()).mul(g.axis)
				g.graph = g.graph.mul(m)
			}
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) && g.size < maxSize {
		g.size += .1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) && g.size > minSize {
		g.size -= .1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.scene = sceneStart
	}
}

const (
	Title                = "Rotate!"
	TitleFontSize        = 48
	TitleMessage         = "Press  Enter  to  start"
	TitleMessageFontSize = 24
)

var (
	InstructionMessage = []string{
		"S  W  for  X  rotation",
		"Q  E  for  Y  rotation",
		"A  D  for  Z  rotation",
		"Hold  Shift  to  rotate  axis",
		"Up  Down  for  zoom",
		"Press  Esc  to  exit"}
	TitleFontFace *text.GoTextFaceSource

	//go:embed arcadeclassic.ttf
	ArcadeClassic_ttf []byte
)

func init() {
	TitleFontFace, _ = text.NewGoTextFaceSource(bytes.NewReader(ArcadeClassic_ttf))
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.scene {
	case sceneStart:
		drawText(screen, Title, TitleFontSize, float64(width)/2, float64(height)/2-150)
		drawText(screen, TitleMessage, TitleMessageFontSize, float64(width)/2, float64(height)/2)
		for i, msg := range InstructionMessage {
			drawText(screen, msg, TitleMessageFontSize, width/2, height/2+100+float64(i)*30)
		}
	case scenePlay:
		g.drawAxis(screen)
		g.drawGraph(screen)
	}
}

func drawText(screen *ebiten.Image, msg string, size float64, x, y float64) {
	text.Draw(screen, msg, &text.GoTextFace{
		Source: TitleFontFace,
		Size:   size,
	}, NewTextDrawOption(x, y))
}

func NewTextDrawOption(x, y float64) *text.DrawOptions {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(color.White)
	op.PrimaryAlign = text.AlignCenter
	op.SecondaryAlign = text.AlignCenter
	return op
}

var (
	colorAxis = []color.Color{
		color.RGBA{255, 51, 51, 255}, // red
		color.RGBA{51, 255, 51, 255}, // green
		color.RGBA{51, 51, 255, 255}, // blue
	}
	stringAxis = [3]string{"X", "Y", "Z"}
)

func (g *Game) drawAxis(screen *ebiten.Image) {
	org := g.pixel(0, 0)
	for i, p := range g.axis {
		pix := g.pixel(p[0], p[1])
		drawLine(screen, org, pix, colorAxis[i])
		drawText(screen, stringAxis[i], 24, float64(pix[0]), float64(pix[1]))
	}
}

func (g *Game) pixel(x, y float64) [2]int {
	i := int((x + g.size) / (2 * g.size) * width)
	j := int((g.size - y) / (2 * g.size) * height)
	return [2]int{i, j}
}

func drawLine(screen *ebiten.Image, p, q [2]int, c color.Color) {
	dx, dy := q[0]-p[0], q[1]-p[1]
	if dx == 0 && dy == 0 {
		screen.Set(p[0], p[1], c)
	}
	n := max(abs(dx), abs(dy))
	for k := 0; k <= n; k++ {
		x := int(float64(p[0]) + float64(dx)*float64(k)/float64(n))
		y := int(float64(p[1]) + float64(dy)*float64(k)/float64(n))
		screen.Set(x, y, c)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func (g *Game) drawGraph(screen *ebiten.Image) {
	for i := 0; i < len(g.graph); i += 2 {
		p, q := g.graph[i], g.graph[i+1]
		pix := g.pixel(p[0], p[1])
		qix := g.pixel(q[0], q[1])
		drawLine(screen, pix, qix, color.White)
	}
}

func (t *Game) Layout(int, int) (int, int) {
	return width, height
}

package rotate

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type scene int

const (
	sceneStart scene = iota
	scenePlay
	sceneEnd
)

const (
	w         = 800
	h         = 600
	maxRounds = 3
)

type Game struct {
	scene
	points        [][2]int
	numVerts      int
	currentPlayer int
	currentRound  int
	unitLength    int
	area, goal    int
	scores        [2]int
}

func (f *Game) Update() error {
	// when the left mouse button is clicked
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		switch f.scene {
		case sceneStart: // at the start scene,
			f.init()                      // initialize the game settings and
			f.scene = scenePlay           // switch to the play scene.
			f.scores[f.currentPlayer] = 0 // reset the player's score.
		case scenePlay: // at the play scene,
			x, y := ebiten.CursorPosition() // get the cursor position
			if len(f.points) < f.numVerts { // if the number of points is less than n,
				f.points = append(f.points, [2]int{x, y}) // add the point to the points slice.
			}
			if len(f.points) == f.numVerts { // if the number of points is equal to n,
				f.computeArea()                                   // calculate the area of the polygon and
				f.scores[f.currentPlayer] += abs(f.area - f.goal) // add the area to the player's score.
				f.scene = sceneEnd                                // switch to the end scene.
			}
		case sceneEnd: // at the end scene,
			f.init() // initialize the game settings and
			if f.currentRound == maxRounds-1 {
				f.currentPlayer = 1 - f.currentPlayer // switch the player and
				f.currentRound = 0                    // reset the round.
				f.scene = sceneStart                  // switch to the start scene.
			} else {
				f.currentRound += 1
				f.scene = scenePlay // switch to the play scene.
			}
		}
	}
	// when the right mouse button is clicked
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		f.scene = sceneStart // switch to the start scene.
	}
	return nil
}

func (f *Game) computeArea() {
	var area float64
	n := len(f.points)
	for i := 0; i < n; i++ {
		p0, p1 := f.points[i], f.points[(i+1)%n]
		x0, y0 := float64(p0[0]), float64(p0[1])
		x1, y1 := float64(p1[0]), float64(p1[1])
		area += (x0 + x1) * (y1 - y0)
	}
	if area < 0 {
		area = -area
	}
	rem := float64(int(area))
	if area-rem >= 0.5 {
		rem += 1
	}
	d := float64(f.unitLength)
	f.area = int(rem / 2 / (d * d))
}

func (f *Game) init() {
	f.numVerts = rand.Intn(5) + 3
	f.points = [][2]int{}
	f.unitLength = rand.Intn(20) + 20
	f.goal = rand.Intn(200) + 3
}

var (
	fillColor     = color.RGBA{150, 0, 0, 150}
	boundaryColor = color.White
)

const (
	Title                = "Fill!"
	TitleFontSize        = 48
	TitleMessage         = "Click  to  start"
	TitleMessageFontSize = 24
)

var (
	InstructionMessage = []string{
		"Left  Click  to  Draw  Points",
		"Right  Click  to  Return",
		"Player  A",
		"Player  B"}
	TitleFontFace *text.GoTextFaceSource

	//go:embed arcadeclassic.ttf
	ArcadeClassic_ttf []byte
)

func init() {
	TitleFontFace, _ = text.NewGoTextFaceSource(bytes.NewReader(ArcadeClassic_ttf))
}

func (f *Game) Draw(screen *ebiten.Image) {
	if f.scene == sceneStart {
		drawText(screen, Title, TitleFontSize, w/2, h/2-150)
		drawText(screen, TitleMessage, TitleMessageFontSize, w/2, h/2)
		for i := 0; i < 4; i++ {
			msg := InstructionMessage[i]
			drawText(screen, msg, TitleMessageFontSize, w/2, h/2+100+float64(i)*30)
		}
		for i := 0; i < 2; i++ {
			drawText(screen, fmt.Sprintf("  %d", f.scores[i]), TitleMessageFontSize, w/2+100, h/2+100+float64(i+2)*30)
		}
	} else {
		f.square(screen)
		drawText(screen, fmt.Sprintf("%d", f.numVerts), TitleFontSize, w-20, h-30)
		dots(screen, f.points, boundaryColor)
		drawText(screen, fmt.Sprintf("%d", f.goal), TitleFontSize, w/2, 20)
		drawText(screen, InstructionMessage[f.currentPlayer+2], TitleMessageFontSize, 50, 20)
		if len(f.points) == f.numVerts {
			fill(screen, f.points, fillColor)
		}
		if f.scene == sceneEnd {
			drawText(screen, fmt.Sprintf("%d", f.area), TitleFontSize, w/2, h-20)
		}
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

func (f *Game) Layout(int, int) (int, int) {
	return w, h
}

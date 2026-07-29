package stoneskipping

const (
	Width  = 800
	Height = 450

	waterLineY    = 292
	verticalScale = 25.0

	timeStep      = 0.005
	stepsPerFrame = 4

	gravity         = 10.0
	waterDensity    = 1.0
	liftCoefficient = 0.1
	dragCoefficient = 0.3
	stoneMass       = 1.0
	stoneLength     = 5.0
	stonePitch      = 3.141592653589793 / 200

	launchX           = 118.0
	launchY           = 0.5
	defaultHorizontal = 150.0
	defaultVertical   = 0.5
	minLaunchSpeed    = 42.0
	maxTrailPoints    = 1500
)

type scene int

const (
	sceneTitle scene = iota
	sceneAim
	sceneFlight
	sceneFinished
)

type point struct {
	x float64
	y float64
}

type ripple struct {
	x      float64
	radius float64
	alpha  float64
}

type Game struct {
	scene
	paused   bool
	dragging bool
	dragX    float64
	dragY    float64

	x       float64
	y       float64
	vx      float64
	vy      float64
	elapsed float64
	skips   int
	inWater bool

	trail   []point
	ripples []ripple
}

func NewGame() *Game {
	g := &Game{scene: sceneTitle}
	g.clearShot()
	return g
}

func (g *Game) Reset() {
	g.scene = sceneAim
	g.clearShot()
}

func (g *Game) clearShot() {
	g.paused = false
	g.dragging = false
	g.dragX = launchX - 72
	g.dragY = waterLineY + 10
	g.x = launchX
	g.y = launchY
	g.vx = defaultHorizontal
	g.vy = defaultVertical
	g.elapsed = 0
	g.skips = 0
	g.inWater = false
	g.trail = []point{{x: g.x, y: g.y}}
	g.ripples = nil
}

func (g *Game) Layout(int, int) (int, int) {
	return Width, Height
}

package predator

type scene int

const (
	sceneNil scene = iota
	sceneStart
	scenePlay
)

type Game struct {
	scene
	Humans []*Human
	Fishes []*Fish
	Corals []*Coral
}

func (g *Game) Clear() {
	g.Humans = []*Human{}
	g.Fishes = []*Fish{}
	g.Corals = []*Coral{}
}

const (
	Width  = 800
	Height = 600
)

func (g *Game) Layout(int, int) (int, int) {
	return Width, Height
}

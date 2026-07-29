package virus

type Scene int

const (
	SceneNil Scene = iota
	SceneStart
	ScenePlay
)

type Game struct {
	Scene
	Agents []*Agent
}

const (
	InitNumAgents = 100
	InitExposed   = 20
)

func (g *Game) Clear() {
	g.Agents = []*Agent{}
	for i := 0; i < InitNumAgents; i++ {
		g.Agents = append(g.Agents, NewAgent())
	}
	for i := 0; i < InitExposed; i++ {
		a := NewAgent()
		a.State = Exposed
		g.Agents = append(g.Agents, a)
	}
}

const (
	Width  = 800
	Height = 600
)

func (g *Game) Layout(int, int) (int, int) {
	return Width, Height
}

package wc

import "yasty.org/peter/wc/ecs"

type Game struct {
	NumWeeks int    `json:"numWeeks"`
	Players  []User `json:"users"`

	Tick int64 `json:"tick"`
	Week int   `json:"week"`

	// These map player index to whether they've drilled/maintained.
	ToDrill    []bool `json:"toDrill"`
	ToMaintain []bool `json:"toMaintain"`
}

func (g *Game) NextWeek() {
	g.Week++

	setAll(g.ToDrill, true)
	setAll(g.ToMaintain, true)
}

func setAll(b []bool, v bool) {
	for i := range b {
		b[i] = v
	}
}

var games []ecs.Entity

func newGame() ecs.Entity {
	return ecs.New()
}

func createGame(users []User) (*Game, error) {
	return &Game{
		Players:    users,
		NumWeeks:   13,
		ToDrill:    make([]bool, len(users)),
		ToMaintain: make([]bool, len(users)),
	}, nil
}

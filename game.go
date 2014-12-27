package wc

import "yasty.org/peter/wc/ecs"

type Game struct {
	NumWeeks int    `json:"numWeeks"`
	Players  []User `json:"users"`

	Map [][]ecs.Entity

	Move int64
	Week int

	// These map player index to whether they've drilled/maintained.
	ToDrill    []bool
	ToMaintain []bool
}

func (g *Game) NextWeek() {
	g.Week++

	g.ToDrill = make([]bool, len(g.Players))
	g.ToMaintain = make([]bool, len(g.Players))
}

var games []ecs.Entity

func newGame() ecs.Entity {
	return ecs.New()
}

func createGame(users []User) (*Game, error) {
	return &Game{Players: users, NumWeeks: 13}, nil
}

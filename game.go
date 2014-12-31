package wc

import "yasty.org/peter/wc/ecs"

// Game is the serializable form of a Wildcatting game.
type Game struct {
	ID       int64  `json:"id,string"`
	NumWeeks int    `json:"numWeeks"`
	Players  []User `json:"players"`

	Tick int `json:"tick"`
	Week int `json:"week"`

	// These map player index to whether they've drilled/maintained.
	ToDrill    []bool `json:"toDrill"`
	ToMaintain []bool `json:"toMaintain"`

	Field Field `json:"field"`

	world  *ecs.World
	oilMap *OilMap
}

type Field struct {
	W     int `json:"w"`
	H     int `json:"h"`
	Sites []Site
}

// Site is the public representation of a site.
type Site struct {
	ID     ecs.Entity `json:"id,string"`
	Point  `json:"point"`
	Survey `json:"survey"`

	*Well `json:"well"`
}

// Survey is a tag applied to site entities once they've been surveyed
// as a player.
type Survey struct {
	DrillCost int
	Tax       int

	// Prob is an integer probability 0..100.
	Prob int8
}

// Well is a tag applied to site entities once they've been drilled.
type Well struct {
	Marker   rune
	OilDepth int8
}

// Point is a coordinate in oilspace.
type Point struct {
	X, Y int
}

func (g *Game) nextWeek() {
	g.Tick++
	g.Week++

	setAll(g.ToDrill, true)
	setAll(g.ToMaintain, true)
}

func (g *Game) player(name string) (int, User) {
	for i, p := range g.Players {
		if p.Username == name {
			return i, p
		}
	}

	return -1, User{}
}

func setAll(b []bool, v bool) {
	for i := range b {
		b[i] = v
	}
}

var games []Game

func createGame(users []User) (*Game, error) {
	game := Game{
		ID:         int64(len(games)),
		Players:    users,
		NumWeeks:   13,
		ToDrill:    make([]bool, len(users)),
		ToMaintain: make([]bool, len(users)),

		Field: Field{W: 80, H: 24},

		world:  ecs.NewWorld(),
		oilMap: newOilMap(80, 24),
	}

	games = append(games, game)

	return &game, nil
}

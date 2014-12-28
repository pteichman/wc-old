package wc

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
}

func (g *Game) nextWeek() {
	g.Week++

	setAll(g.ToDrill, true)
	setAll(g.ToMaintain, true)
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
	}

	games = append(games, game)

	return &game, nil
}

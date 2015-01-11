package wc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Storage struct {
	Static string
	Users  UserStore
}

func NewHandler(s Storage) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(s.Static)))

	mux.HandleFunc("/api/game/new", s.newGame)
	mux.HandleFunc("/api/user/new", s.newUser)

	mux.HandleFunc("/api/state", s.state)
	mux.HandleFunc("/api/move", nil)

	return mux
}

type Response struct {
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

func errorResp(err error) Response {
	return Response{Success: false, Error: err.Error()}
}

func write(w http.ResponseWriter, r Response) {
	j, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		j, err = json.Marshal(Response{Success: false, Error: err.Error()})
	}

	w.Write(j)
	w.Write([]byte("\n"))
}

type newGameResp struct {
	ID    int64  `json:"id,string"`
	State string `json:"state"`
}

func (s Storage) newGame(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		write(w, errorResp(err))
		return
	}

	var users []User
	for _, name := range r.Form["user"] {
		if u, err := s.Users.Get(name); err == nil {
			users = append(users, u)
		}
	}

	game, err := createGame(users)
	if err != nil {
		write(w, errorResp(err))
		return
	}

	game.nextWeek()

	resp := newGameResp{
		ID:    game.ID,
		State: fmt.Sprintf("/api/state?game=%d&user=Alice", game.ID),
	}

	write(w, Response{Success: true, Result: resp})
}

func (s Storage) newUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		write(w, errorResp(err))
		return
	}

	user, err := s.Users.Create(User{Username: r.Form.Get("username")})
	if err != nil {
		write(w, errorResp(err))
		return
	}

	write(w, Response{Success: true, Result: user})
}

func getGame(v url.Values) (*Game, error) {
	gid, err := strconv.ParseInt(v.Get("game"), 10, 64)
	if err != nil {
		return nil, err
	}

	return &games[gid], nil
}

type State struct {
	*Game
	Moves []string `json:"moves"`
}

func (s Storage) state(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		write(w, errorResp(err))
		return
	}

	g, err := getGame(r.Form)
	if err != nil {
		write(w, errorResp(err))
		return
	}

	var moves []string

	if name := r.Form.Get("user"); name != "" {
		// Figure out this user's valid moves.
		i, _ := g.player(name)

		if g.ToDrill[i] {
			moves = append(moves, "action=drill&loc=NNNN")
		} else if g.ToMaintain[i] {
			// Append a sell move for all the user's wells.
		}
	}

	write(w, Response{Success: true, Result: State{Game: g, Moves: moves}})
}

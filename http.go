package wc

import (
	"encoding/json"
	"net/http"
)

type Storage struct {
	Static string
	Users  UserStore
}

func NewMux(s Storage) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(s.Static)))

	mux.HandleFunc("/api/game/new", s.newGame)
	mux.HandleFunc("/api/user/new", s.newUser)

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

	game.NextWeek()

	write(w, Response{Success: true, Result: game})
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

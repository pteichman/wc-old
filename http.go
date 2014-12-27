package wc

import (
	"encoding/json"
	"net/http"
)

func New(staticdir string) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(staticdir)))

	mux.HandleFunc("/api/game/new", newGameHandler)
	mux.HandleFunc("/api/user/new", newUserHandler)

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

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		write(w, errorResp(err))
		return
	}

	var users []User
	for _, name := range r.Form["user"] {
		if u, ok := getUser(name); ok {
			users = append(users, u)
		}
	}

	game, err := createGame(users)
	if err != nil {
		write(w, errorResp(err))
		return
	}

	write(w, Response{Success: true, Result: game})
}

func newUserHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		write(w, errorResp(err))
		return
	}

	user, err := createUser(User{Username: r.Form.Get("username")})
	if err != nil {
		write(w, errorResp(err))
		return
	}

	write(w, Response{Success: true, Result: user})
}

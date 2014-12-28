package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	"yasty.org/peter/wc"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	var (
		addr      = flag.String("addr", "localhost:8080", "HTTP bind address")
		staticdir = flag.String("staticdir", "static", "Path to static files")
	)

	flag.Parse()

	s := wc.Storage{
		Static: *staticdir,
		Users:  &wc.MemUsers{},
	}

	http.Handle("/", logHandler{wc.NewHandler(s)})

	log.Printf("Listening on http://%s/", *addr)
	log.Println(http.ListenAndServe(*addr, nil))
}

type logHandler struct {
	http.Handler
}

func (h logHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s", r.URL)
	h.Handler.ServeHTTP(w, r)
}

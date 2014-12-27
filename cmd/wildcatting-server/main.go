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

	http.Handle("/", wc.New(*staticdir))

	log.Printf("Listening on http://%s/", *addr)
	log.Println(http.ListenAndServe(*addr, nil))
}

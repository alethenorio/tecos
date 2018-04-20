package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/ByteFlinger/tecos/backend/dummy"
	"github.com/ByteFlinger/tecos/v1"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var routes = flag.Bool("routes", false, "Generate router documentation")
var port = flag.Int("port", 8080, "Port to listen to")

func main() {

	flag.Parse()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Mount("/v1", v1.Routes(&dummy.Dummy{}))

	if err := http.ListenAndServe(":"+strconv.Itoa(*port), r); err != nil {
		log.Fatal(err)
	}

}

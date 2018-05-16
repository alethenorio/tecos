package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/ByteFlinger/tecos/backend"
	"github.com/ByteFlinger/tecos/backend/gitmono"
	"github.com/ByteFlinger/tecos/v1"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pkg/errors"
)

var (
	routes *bool
	port   *int
)

func main() {

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	bkndType := os.Args[1]
	be, err := handleBackend(bkndType)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func(be backend.Storage) {
		<-c
		be.Cleanup()
		os.Exit(0)
	}(be)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to init backend - %s\n", err)
		usage()
		os.Exit(1)
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Mount("/v1", v1.Routes(be))

	fmt.Printf("Listening on port %d\n", *port)
	if err := http.ListenAndServe(":"+strconv.Itoa(*port), r); err != nil {
		log.Fatal(err)
	}

}

func handleBackend(backendType string) (backend.Storage, error) {

	var (
		storage backend.Storage
		err     error
	)

	set := flag.NewFlagSet("tecos", flag.ExitOnError)
	addGlobalFlags(set)

	switch backendType {
	case "gitmono":
		c := &gitmono.Config{}
		set.StringVar(&c.RepoURL, "url", "", "URL of the git repository")
		set.StringVar(&c.RepoURL, "modulepath", "", "Relative path to all terraform modules. Wildcards are supported")
		set.Parse(os.Args[2:])
		storage, err = gitmono.New(c)
	default:
		return nil, errors.Errorf("Backend '%s' is invalid", backendType)
	}

	return storage, err
}

func ChooseBackend(backend string) backend.Storage {
	switch backend {
	case "gitmono":

	}

	return nil
}

func addGlobalFlags(set *flag.FlagSet) {
	routes = set.Bool("routes", false, "Generate router documentation")
	port = set.Int("port", 8080, "Port to listen to")
}

func usage() {
	fmt.Fprintln(os.Stderr, `The following backends are supported:

	gitmono

Use -h after the backend type for details on backend specific flags`)
}

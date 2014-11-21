// +build !appengine

package gopi

import (
	"flag"
	"log"
	"net/http"

	"github.com/ludmiloff/gopi/bind"
	"github.com/ludmiloff/gopi/graceful"
)

func init() {
	bind.WithFlag()
	if fl := log.Flags(); fl&log.Ltime != 0 {
		log.SetFlags(fl | log.Lmicroseconds)
	}
}

// Serve starts Goji using reasonable defaults.
func Serve() {
	if !flag.Parsed() {
		flag.Parse()
	}

	DefaultMux.Compile()
	// Install our handler at the root of the standard net/http default mux.
	// This allows packages like expvar to continue working as expected.
	http.Handle("/", DefaultMux)

	listener := bind.Default()
	log.Println("Starting Goji on", listener.Addr())

	graceful.HandleSignals()
	bind.Ready()
	graceful.PreHook(func() { log.Printf("Goji received signal, gracefully stopping") })
	graceful.PostHook(func() { log.Printf("Goji stopped") })

	err := graceful.Serve(listener, http.DefaultServeMux)

	if err != nil {
		log.Fatal(err)
	}

	graceful.Wait()
}

// +build !appengine

package gopi

import (
	"flag"
	"github.com/coopernurse/gorp"
	"github.com/gorilla/sessions"
	"github.com/pelletier/go-toml"
	"html/template"

	"log"
	"net/http"

	"github.com/ludmiloff/gopi/bind"
	"github.com/ludmiloff/gopi/graceful"
	"github.com/ludmiloff/gopi/web"
	"github.com/ludmiloff/gopi/web/middleware"
)

type Application struct {
	DefaultMux  *web.Mux
	Config      *toml.TomlTree
	DB          *gorp.DbMap
	Session     *sessions.FilesystemStore
	CookieStore *sessions.CookieStore
	Render      *Render
	Language    string

	//
	ShutdownFunc func()

	// Application wide parameters
	Params map[string]string
}

var App *Application

func CreateAppliction() *Application {
	app := &Application{}

	if !flag.Parsed() {
		flag.Parse()
	}

	filename := flag.String("config", "config.toml", "Path to configuration file")

	var err error

	app.Config, err = toml.LoadFile(*filename)
	if err != nil {
		log.Fatalf("Config load failed: %s\n", err)
	}

	defaultBind := ":8000"

	if s := bind.Sniff(); s != "" {
		defaultBind = s
	}

	var bindAddress string

	flag.StringVar(&bindAddress, "bind", defaultBind,
		`Address to bind on. If this value has a colon, as in ":8000" or
		"127.0.0.1:9001", it will be treated as a TCP address. If it
		begins with a "/" or a ".", it will be treated as a path to a
		UNIX socket. If it begins with the string "fd@", as in "fd@3",
		it will be treated as a file descriptor (useful for use with
		systemd, for instance). If it begins with the string "einhorn@",
		as in "einhorn@0", the corresponding einhorn socket will be
		used. If an option is not explicitly passed, the implementation
		will automatically select among "einhorn@0" (Einhorn), "fd@3"
		(systemd), and ":8000" (fallback) based on its environment.`)

	bind.WithFlag(bindAddress)
	if fl := log.Flags(); fl&log.Ltime != 0 {
		log.SetFlags(fl | log.Lmicroseconds)
	}

	// Defaults
	app.DefaultMux = web.New()
	app.DefaultMux.Use(middleware.RequestID)
	app.DefaultMux.Use(middleware.Logger)
	app.DefaultMux.Use(middleware.Recoverer)
	app.DefaultMux.Use(middleware.AutomaticOptions)
	app.DefaultMux.Compile()
	// Install our handler at the root of the standard net/http default mux.
	// This allows packages like expvar to continue working as expected.
	http.Handle("/", app.DefaultMux)

	// Static files
	// Setup static files
	static := web.New()
	publicPath := app.Config.Get("general.public_path").(string)
	static.Get("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir(publicPath))))
	http.Handle("/assets/", static)

	// Template renderer
	app.InitRender(template.FuncMap{})

	// Database
	app.InitDB()

	// Cookies
	app.InitCookies()

	// File system session store
	app.InitSessionStore()

	// Parameters
	app.Params = map[string]string{}

	graceful.HandleSignals()
	bind.Ready()
	graceful.PreHook(app.AboutToStop)
	graceful.PostHook(func() { log.Printf("Gopi stopped") })

	App = app

	return app
}

// Start starts Gopi using reasonable defaults.
func (this *Application) Start() {

	listener := bind.Default()
	log.Println("Starting Gopi on", listener.Addr())

	err := graceful.Serve(listener, http.DefaultServeMux)

	if err != nil {
		log.Fatal(err)
	}

	graceful.Wait()
}

func (this *Application) AboutToStop() {
	log.Printf("Gopi received signal, gracefully stopping")
	if this.ShutdownFunc != nil {
		this.ShutdownFunc()
	}
}

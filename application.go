// +build !appengine

package gopi

import (
	"flag"
	"github.com/gorilla/sessions"
	"github.com/pelletier/go-toml"
	"html/template"

	"log"
	"net/http"
	
	"github.com/gorilla/context"
	"github.com/ludmiloff/gopi/bind"
	"github.com/ludmiloff/gopi/graceful"
	"github.com/ludmiloff/gopi/web"
	"github.com/ludmiloff/gopi/web/middleware"
)

type ApplicationVirtualMethods interface {
	Init()
}

type Application struct {
	ApplicationVirtualMethods
	Name		string
	DefaultMux  *web.Mux
	Config      *toml.TomlTree
	Session     *sessions.FilesystemStore
	CookieStore *sessions.CookieStore
	Render      *Render
	Language    string

	//
	ShutdownFunc func()

	// Application wide parameters
	Params map[string]string

	AuthenticateUser func(username, password string)
}

var app *Application

//func CreateAppliction(virt ApplicationVirtualMethods) *Application {
func CreateAppliction(name string) *Application {

	app = &Application{Name: name}
	
	if !flag.Parsed() {
		flag.Parse()
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

	app.initInternal()

	return app
}

func (this *Application) initInternal() {

	filename := flag.String("config", "config.toml", "Path to configuration file")

	var err error

	this.Config, err = toml.LoadFile(*filename)
	if err != nil {
		log.Fatalf("Config load failed: %s\n", err)
	}

	// Defaults
	this.DefaultMux = web.New()
	this.DefaultMux.Use(middleware.RequestID)
	if this.Config.GetDefault("general.logger", true).(bool) {
		this.DefaultMux.Use(middleware.Logger)
	}
	this.DefaultMux.Use(middleware.Recoverer)
	this.DefaultMux.Use(middleware.AutomaticOptions)
	this.DefaultMux.Compile()
	// Install our handler at the root of the standard net/http default mux.
	// This allows packages like expvar to continue working as expected.
	http.Handle("/", this.DefaultMux)

	// Static files
	// Setup static files
	if this.Config.GetDefault("general.handle_assets", true).(bool) {
		static := web.New()
		publicPath := this.Config.Get("general.public_path").(string)
		static.Get("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir(publicPath))))
		http.Handle("/assets/", static)
	}

	// Template renderer
	this.InitRender(template.FuncMap{})

//	if (this.Config.GetDefault("cookies.use_cookies", false).(bool)) {
//		// Cookies
//		this.InitCookies()
//		// Use Middleware
//		this.Use(this.ApplyCookies)
//	}

	// Parameters
	this.Params = map[string]string{}

	graceful.HandleSignals()
	bind.Ready()
	graceful.PreHook(this.AboutToStop)
	graceful.PostHook(func() { log.Printf("Gopi stopped") })

	log.Printf("internal init done ...")
}

// Start starts Gopi using reasonable defaults.
func (this *Application) Start() {

	// User defined initialization
	this.Init()

	// Finalize middleware stack
	this.Use(context.ClearHandler)

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

// DefaultMux actions

// Use appends the given middleware to the default Mux's middleware stack. See
// the documentation for web.Mux.Use for more information.
func (this *Application) Use(middleware interface{}) {
	this.DefaultMux.Use(middleware)
}

// Insert the given middleware into the default Mux's middleware stack. See the
// documentation for web.Mux.Insert for more information.
func (this *Application) Insert(middleware, before interface{}) error {
	return this.DefaultMux.Insert(middleware, before)
}

// Abandon removes the given middleware from the default Mux's middleware stack.
// See the documentation for web.Mux.Abandon for more information.
func (this *Application) Abandon(middleware interface{}) error {
	return this.DefaultMux.Abandon(middleware)
}

// Handle adds a route to the default Mux. See the documentation for web.Mux for
// more information about what types this function accepts.
func (this *Application) Handle(pattern interface{}, handler interface{}) {
	this.DefaultMux.Handle(pattern, handler)
}

// Connect adds a CONNECT route to the default Mux. See the documentation for
// web.Mux for more information about what types this function accepts.
func (this *Application) Connect(pattern interface{}, handler interface{}) {
	this.DefaultMux.Connect(pattern, handler)
}

// Delete adds a DELETE route to the default Mux. See the documentation for
// web.Mux for more information about what types this function accepts.
func (this *Application) Delete(pattern interface{}, handler interface{}) {
	this.DefaultMux.Delete(pattern, handler)
}

// Get adds a GET route to the default Mux. See the documentation for web.Mux for
// more information about what types this function accepts.
func (this *Application) Get(pattern interface{}, handler interface{}) {
	this.DefaultMux.Get(pattern, handler)
}

// Head adds a HEAD route to the default Mux. See the documentation for web.Mux
// for more information about what types this function accepts.
func (this *Application) Head(pattern interface{}, handler interface{}) {
	this.DefaultMux.Head(pattern, handler)
}

// Options adds a OPTIONS route to the default Mux. See the documentation for
// web.Mux for more information about what types this function accepts.
func (this *Application) Options(pattern interface{}, handler interface{}) {
	this.DefaultMux.Options(pattern, handler)
}

// Patch adds a PATCH route to the default Mux. See the documentation for web.Mux
// for more information about what types this function accepts.
func (this *Application) Patch(pattern interface{}, handler interface{}) {
	this.DefaultMux.Patch(pattern, handler)
}

// Post adds a POST route to the default Mux. See the documentation for web.Mux
// for more information about what types this function accepts.
func (this *Application) Post(pattern interface{}, handler interface{}) {
	this.DefaultMux.Post(pattern, handler)
}

// Put adds a PUT route to the default Mux. See the documentation for web.Mux for
// more information about what types this function accepts.
func (this *Application) Put(pattern interface{}, handler interface{}) {
	this.DefaultMux.Put(pattern, handler)
}

// Trace adds a TRACE route to the default Mux. See the documentation for
// web.Mux for more information about what types this function accepts.
func (this *Application) Trace(pattern interface{}, handler interface{}) {
	this.DefaultMux.Trace(pattern, handler)
}

// NotFound sets the NotFound handler for the default Mux. See the documentation
// for web.Mux.NotFound for more information.
func (this *Application) NotFound(handler interface{}) {
	this.DefaultMux.NotFound(handler)
}

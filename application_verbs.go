package gopi

import ()

func init() {
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

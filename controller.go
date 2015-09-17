package gopi

import (
	"github.com/ludmiloff/gopi/web"
	//"log"
	"net/http"
)


type Controller struct {
	Layout string

	// Hooks
	//BeforeActionHook func(c web.C)
	BeforeRenderHook func(c web.C)
}


func (this *Controller) SaveSession(c web.C) {
//	cookies, _ := App.CookieStore.Get(c.Request, "cookies")
//	err := cookies.Save(c.Request, c.W)
//	if err != nil {
//		log.Println("Can't save session: %v", err)
//	}
//
//	// TODO: filestystem store
}

func (this *Controller) End(c web.C) {
	this.SaveSession(c)
}

func (this *Controller) Redirect(c web.C, urlStr string, code int) {
	this.SaveSession(c)
	http.Redirect(c.W, c.Request, urlStr, code)
}

func (this *Controller) Render(c web.C, view string, args RenderArgs, status int) {
	if this.BeforeRenderHook != nil {
		this.BeforeRenderHook(c)
	}

	this.SaveSession(c)

	var r *Render = app.Render
	var data RenderArgs
	if args == nil {
		data = RenderArgs{}
	} else {
		data = args
	}

	data["Env"] = c.Env

	r.RenderHTML(c.W,
		status, view, data,
		HTMLOptions{Layout: this.Layout})
}

func (this *Controller) JSON(c web.C, v interface{}) {
	this.SaveSession(c)
	var r *Render = app.Render
	r.JSON(c.W, http.StatusOK, v)
}

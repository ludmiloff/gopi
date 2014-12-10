package gopi

import (
	"github.com/ludmiloff/gopi/web"
	"log"
	"net/http"
)

type Controller struct {
	Layout string
}

func (this *Controller) SaveSession(c web.C) {
	cookies, _ := App.CookieStore.Get(c.Request, "cookies")
	err := cookies.Save(c.Request, c.W)
	if err != nil {
		log.Println("Can't save session: %v", err)
	}

	// TODO: filestystem store
}

func (this *Controller) End(c web.C) {
	this.SaveSession(c)
}

func (this *Controller) Redirect(c web.C, urlStr string, code int) {
	this.SaveSession(c)
	http.Redirect(c.W, c.Request, urlStr, code)
}

func (this *Controller) BeforeRender(c web.C) {

}

func (this *Controller) Render(c web.C, view string, args RenderArgs, status int) {
	this.BeforeRender(c)
	this.SaveSession(c)

	var r *Render = App.Render
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

package gopi

import (
	"github.com/ludmiloff/gopi/web"
	//"log"
	"net/http"
)

type Controller struct {
	Layout string
}

func (this *Controller) Redirect(c web.C, urlStr string, code int) {
	http.Redirect(c.W, c.Request, urlStr, code)
}

func (this *Controller) Render(c web.C, view string, args RenderArgs, status int) {
	var r *Render = App.Render
	var data RenderArgs
	if args == nil {
		data = RenderArgs{}
	} else {
		data = args
	}
	r.RenderHTML(c.W,
		status, view, data,
		HTMLOptions{Layout: this.Layout})
}

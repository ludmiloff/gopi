package gopi

import (
	"github.com/ludmiloff/gopi/web"
	"net/http"
)

type Controller struct {
	Layout    string
	PageTitle string
}

func (this *Controller) Redirect(c web.C, urlStr string, code int) {
	http.Redirect(c.W, c.Request, urlStr, code)
}

func (this *Controller) Render(c web.C, view string, data interface{}, status int) {
	var r *Render = App.Render
	r.RenderHTML(c.W,
		status, view, data,
		HTMLOptions{Layout: this.Layout})
}

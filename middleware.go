package gopi

import (
	"github.com/ludmiloff/gopi/web"
	"net/http"
)

// Makes sure controllers can have access to session
func (this *Application) ApplySessions(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		cookies, _ := this.CookieStore.Get(r, "cookies")
		c.Env["Cookies"] = cookies
		// TODO: filesystem store
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// func (this *Application) ApplyWebUser(c *web.C, h http.Handler) http.Handler {
// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		session := c.Env["Cookies"].(*sessions.Session)
// 	}
// 	return http.HandlerFunc(fn)
// }

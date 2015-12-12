package gopi

import (
	"crypto/sha256"
	"github.com/gorilla/sessions"
	"io"
)

func (this *Application) InitCookies() {
	hash := sha256.New()
	io.WriteString(hash, this.Config.Get("cookie.secret").(string))
	this.CookieStore = sessions.NewCookieStore(hash.Sum(nil))
	this.CookieStore.Options = &sessions.Options{
		Path: "/",
		HttpOnly: true,
		Secure:   this.Config.Get("cookie.secure").(bool),
	}
}

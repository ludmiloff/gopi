package gopi

import (
	"github.com/ludmiloff/gopi/web"
)

type WebUser struct {
	Id      uint64
	Email   string
	Names   string
	Roles   []string
	IsAdmin bool
}

func (this *WebUser) HasRole() bool {
	// TODO
	return false
}

func (this *Application) GetUser(c web.C) *WebUser {
	if user, exists := c.Env["User"].(*WebUser); exists {
		return user
	}

	new_user := &WebUser{Id: 0, Email: "", Names: "", IsAdmin: false}
	return new_user
}

// func (this *Application) Authenticate(user, password string) bool {

// }

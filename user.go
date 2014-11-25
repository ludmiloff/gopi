package gopi

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

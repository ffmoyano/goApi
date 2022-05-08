package model

type User struct {
	ID       int
	Name     string
	Username string
	Password string
	Email    string
	Roles    []Role
	Active   bool
}

type Role struct {
	ID   int
	Name string
}

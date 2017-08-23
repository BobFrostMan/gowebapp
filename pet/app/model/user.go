package model

type User struct {
	Login string
	Name string
	Password string
	Groups []Group
}
package model

type Permission struct {
	Name string
	Type string
	Value string
	Read bool
	Update bool
	Execute bool
}

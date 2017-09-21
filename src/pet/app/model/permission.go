package model

//TODO: implement Value as Values
//TODO: create simple user permission and permission groups
type Permission struct {
	Name     string `json:"name"`
	Type     string `json:"type"`	// type string from Type constants
	Value    string `json:"value"`
	Read     bool `json:"read"`
	Update   bool `json:"update"`
	Execute  bool `json:"execute"`
}
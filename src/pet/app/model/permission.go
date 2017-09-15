package model

import (
	"gopkg.in/mgo.v2/bson"
)

//TODO: implement Value as Values
//TODO: create simple user permission and permission groups
type Permission struct {
	ObjectID bson.ObjectId `bson:"_id" json:"_id"`
	ID       uint32 `db:"id" json:"id,omitempty" bson:"id,omitempty"` // use PermissionID() instead for consistency with database types
	Name     string `json:"name"`
	Type     string `json:"type"`	// type string from Type constants
	Value    string `json:"value"`
	Read     bool `json:"read"`
	Update   bool `json:"update"`
	Execute  bool `json:"execute"`
}
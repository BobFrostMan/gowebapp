package model

import (
	"gopkg.in/mgo.v2/bson"
)

type Group struct {
	ObjectID    bson.ObjectId `bson:"_id" json:"_id"`
	ID          uint32 `db:"id" json:"id,omitempty" bson:"id,omitempty"` // use GroupID() instead for consistency with database types
	Name        string `bson:"name" json:"name"`
	Permissions []Permission `bson:"permissions" json:"permissions"`
}


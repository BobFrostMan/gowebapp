package model

import (
	"pet/app/shared/database"
	"gopkg.in/mgo.v2/bson"
	"log"
)

// Database tables, collections, fields etc.
const (
	UsersCollection = "Users"
)

type User struct {
	ObjectID bson.ObjectId `bson:"_id" json:"_id"`
	ID       uint32 `db:"id" json:"id,omitempty" bson:"id,omitempty"` // use UserID() instead for consistency with database types
	Login    string `json:"login"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Groups   []string `json:"groups"`
}


// UserById
// Returns user by given _id and error
func UserById(id string) (*User, error) {
	var err error
	var user User
	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C(UsersCollection)
		err = c.FindId(bson.ObjectIdHex(id)).One(&user)
	} else {
		err = NoDBConnection
	}
	if err != nil {
		log.Printf("User '%s' wasn't found", id)
	}
	return &user, err
}


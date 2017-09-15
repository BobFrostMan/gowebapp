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
	Groups   []Group `json:"groups"`
}

// IsAllowed
// Returns true if operation allowed for user object
func (u *User) IsAllowed(operation string) bool {
	for _, group := range u.Groups{
		for _, permission := range group.Permissions{
			if permission.Value == operation{
				return permission.Execute
			}
		}
	}
	return false
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


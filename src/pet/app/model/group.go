package model

import (
	"gopkg.in/mgo.v2/bson"
	"pet/app/shared/database"
	"log"
)

const (
	GroupsCollection = "PermissionGroups"
)

type Group struct {
	ObjectID    bson.ObjectId `bson:"_id" json:"_id"`
	ID          uint32 `db:"id" json:"id,omitempty" bson:"id,omitempty"` // use GroupID() instead for consistency with database types
	Name        string `bson:"name" json:"name"`
	Permissions []Permission `bson:"permissions" json:"permissions"`
}

// UserById
// Returns user by given _id and error
func GetGroups(names []string) ([]Group, error) {
	var err error
	var groups []Group
	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C(GroupsCollection)
		log.Printf("Request: %v", bson.M{ "name" : bson.M{ "$in" : names} })
		err = c.Find(bson.M{ "name" : bson.M{ "$in" : names} }).All(&groups)
	} else {
		err = NoDBConnection
	}
	if err != nil {
		log.Printf("Error occured during receiving permission groups '%s'", err.Error())
	}
	return groups, err
}
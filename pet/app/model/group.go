package model

import (
	"gopkg.in/mgo.v2/bson"
	"gowebapp/pet/app/shared/database"
	"log"
)

// Database tables, collections, fields etc.
const (
	GroupsCollection  = "PermissionGroups"
)

// Messages patterns
const (
	GroupNotFound = "Group '%s' wasn't found"
	GroupNotCreated = "Group '%s' wasn't created"
	GroupCreated = "Group '%s' was successfully created"
)

type Group struct {
	ObjectID  bson.ObjectId `bson:"_id"`
	ID uint32 `db:"id" bson:"id,omitempty"` // use GroupID() instead for consistency with database types
	Name string `bson:"name"`
	Permissions []Permission `bson:"permissions"`
}

// GroupID
// GroupID returns the user id
func (p *Group) GroupId() string {
	r := ""

	switch database.ReadConfig().Type {
	case database.TypeMongoDB:
		r = p.ObjectID.Hex()
	}

	return r
}

// GroupCreate
// Creates permission group with given name, and permissions
func GroupCreate(name string, permissions []Permission) error {
	var err error

	switch database.ReadConfig().Type {
	case database.TypeMongoDB:
		if database.CheckConnection() {
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C(GroupsCollection)

			group := &Group{
				ObjectID:  bson.NewObjectId(),
				Name: name,
				Permissions: permissions,
			}
			err = c.Insert(group)
		} else {
			err = NoDBConnection
		}
	default:
		err = DBNotSelected
	}

	if err != nil {
		log.Printf(GroupNotCreated, name)
	} else {
		log.Printf(GroupCreated, name)
	}

	return err
}

// GroupByName
// Returns group by given name and error
func GroupByName(name string) (*Group, error) {
	var err error
	var group Group
	switch database.ReadConfig().Type {
	case database.TypeMongoDB:
		if database.CheckConnection() {
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C(GroupsCollection)
			err = c.Find(bson.M{"name": name}).One(&group)
		} else {
			err = NoDBConnection
		}
	default:
		err = DBNotSelected
	}

	if err != nil{
		log.Printf(GroupNotFound, name)
	}

	return &group, err
}


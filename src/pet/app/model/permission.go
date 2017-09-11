package model

import (
	"gopkg.in/mgo.v2/bson"
	"log"
	"pet/app/shared/database"
)

// Entity constants
const (
	TypeMethod = "method"
	TypeField = "field"
)

// Database tables, collections, fields etc.
const (
	PermissionsCollection = "Permissions"
)

// Messages patterns
const (
	PermissionNotFound = "Permission '%s' wasn't found"
	PermissionNotCreated = "Permission '%s' wasn't created"
	PermissionCreated = "Permission '%s' was successfully created"
)

//TODO: implement Value as Values
//TODO: create simple user permission and permission groups
type Permission struct {
	ObjectID bson.ObjectId `bson:"_id"`
	ID       uint32 `db:"id" json:"id,omitempty" bson:"id,omitempty"` // use PermissionID() instead for consistency with database types
	Name     string `json:"name"`
	Type     string `json:"type"`	// type string from Type constants
	Value    string `json:"value"`
	Read     bool `json:"read"`
	Update   bool `json:"update"`
	Execute  bool `json:"execute"`
}

// PermissionID
// PermissionID returns the user id
func (p *Permission) PermissionId() string {
	r := ""

	switch database.ReadConfig().Type {
	case database.TypeMongoDB:
		r = p.ObjectID.Hex()
	}

	return r
}

// PermissionCreate
// Creates user permission with given name, type, value, and action types
func PermissionCreate(name string, pType string, value string, read bool, update bool, execute bool) error {
	var err error

	switch database.ReadConfig().Type {
	case database.TypeMongoDB:
		if database.CheckConnection() {
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C(PermissionsCollection)

			permission := &Permission{
				ObjectID:  bson.NewObjectId(),
				Name:  name,
				Type: pType,
				Value: value,
				Read: read,
				Update: update,
				Execute: execute,
			}
			err = c.Insert(permission)
		} else {
			err = NoDBConnection
		}
	default:
		err = DBNotSelected
	}

	if err != nil{
		log.Printf(PermissionNotCreated, name)
	} else {
		log.Printf(PermissionCreated, name)
	}

	return err
}


// PermissionByName
// Returns permission by given name and error
func PermissionByName(name string) (*Permission, error) {
	var err error
	var permission Permission
	switch database.ReadConfig().Type {
	case database.TypeMongoDB:
		if database.CheckConnection() {
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C(PermissionsCollection)
			err = c.Find(bson.M{"name": name}).One(&permission)
		} else {
			err = NoDBConnection
		}
	default:
		err = DBNotSelected
	}

	if err != nil{
		log.Printf(PermissionNotFound, name)
	}

	return &permission, err
}

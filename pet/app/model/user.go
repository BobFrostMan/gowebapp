package model

import (
	"gowebapp/pet/app/shared/database"
	"gopkg.in/mgo.v2/bson"
	"log"
	"gowebapp/pet/app/shared/passhash"
)

// Database tables, collections, fields etc.
const (
	UsersCollection  = "Users"
)

// Messages patterns
const (
	UserNotFound = "User '%s' wasn't found"
	UserNotCreated = "User '%s' wasn't created"
	UserCreated = "User '%s' was successfully created"
)

type User struct {
	ObjectID  bson.ObjectId `bson:"_id"`
	ID uint32 `db:"id" bson:"id,omitempty"` // use UserID() instead for consistency with database types
	Login string
	Name string
	Password string
	Groups []Group
}

// UserID
// UserID returns the user id
func (u *User) UserID() string {
	r := ""

	switch database.ReadConfig().Type {
	case database.TypeMongoDB:
		r = u.ObjectID.Hex()
	}

	return r
}

// UserCreate
// Creates user with given login, name, password, and groups
// Saves password as hash
// User can be created with empty groups value
func UserCreate(login string, name string, password string, groups []Group) error {
	var err error

	switch database.ReadConfig().Type {
	case database.TypeMongoDB:
		if database.CheckConnection() {
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C(UsersCollection)
			hash, er := passhash.HashString(password)
			if er != nil {
				log.Printf("Can't generate hash password for user '%s'", login)
				break
			}
			user := &User{
				ObjectID:  bson.NewObjectId(),
				Login: login,
				Name:  name,
				Password:  hash,
				Groups: groups,
			}
			err = c.Insert(user)
		} else {
			err = NoDBConnection
		}
	default:
		err = DBNotSelected
	}

	if err != nil{
		log.Printf(UserNotCreated, login)
	} else {
		log.Printf(UserCreated, login)
	}

	return err
}

// UserByLogin
// Returns user by given login and error
func UserByLogin(login string) (*User, error) {
	var err error
	var user User
	switch database.ReadConfig().Type {
	case database.TypeMongoDB:
		if database.CheckConnection() {
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C(UsersCollection)
			err = c.Find(bson.M{"login": login}).One(&user)
		} else {
			err = NoDBConnection
		}
	default:
		err = DBNotSelected
	}

	if err != nil{
		log.Printf(UserNotFound, login)
	}

	return &user, err
}

// UserByLogin
// Returns user by given login and error
func UserList() ([]User, error) {
	var err error
	var users []User
	switch database.ReadConfig().Type {
	case database.TypeMongoDB:
		if database.CheckConnection() {
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C(UsersCollection)
			err = c.Find(bson.M{}).All(&users)
		} else {
			err = NoDBConnection
		}
	default:
		err = DBNotSelected
	}

	if err != nil{
		// TODO: implement proper message for this case
		log.Println(UserNotFound)
	}
	return users, err
}


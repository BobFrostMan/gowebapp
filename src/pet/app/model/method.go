package model

import (
	"pet/app/shared/database"
	"gopkg.in/mgo.v2/bson"
	"log"
	"github.com/xenzh/gofsm"
)
// Database tables, collections, fields etc.
const (
	MethodsCollection = "Method"
)

type Method struct {
	ObjectID bson.ObjectId `bson:"_id"`
	Name string `json:"name"`
	Parameters []Parameter `json:"parameters"`
	Fsm simple_fsm.JsonRoot `json:"fsm"`
}

type Parameter struct {
	Name string `json:"name"`
	Required bool `json:"required"`
	Type string `json:"type"`
}

func (m *Method)IsEmpty() bool {
	return m.Name == "" && len(m.Parameters) == 0
}

// GetAllMethods
// Returns all api methods located in database
func GetAllMethods() *[]Method {
	var err error
	var methods []Method
	switch database.ReadConfig().Type {
	case database.TypeMongoDB:
		if database.CheckConnection() {
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C(MethodsCollection)
			err = c.Find(bson.M{}).All(&methods)
		} else {
			err = NoDBConnection
		}
	default:
		err = DBNotSelected
	}

	if err != nil{
		log.Println("Can not receive methods", err.Error())
	}
	return &methods
}


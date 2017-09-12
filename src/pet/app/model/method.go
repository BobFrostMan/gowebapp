package model

import (
	"pet/app/shared/database"
	"gopkg.in/mgo.v2/bson"
	"log"
	"github.com/xenzh/gofsm"
)


// Entity constants
const (
	ParamTypeInt = "int"
	ParamTypeString = "string"
)

// Database tables, collections, fields etc.
const (
	MethodsCollection = "Method"
)

// Messages patterns
const (
	MethodNotFound = "Method '%s' wasn't found"
	MethodNotCreated = "Method '%s' wasn't created"
	MethodCreated = "Method '%s' was successfully created"
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

//TODO: guess that can be done in another way
func (m *Method)IsEmpty() bool {
	return m.Name == "" && len(m.Parameters) == 0
}

// CreateMethod
// Creates api method representation in database
func CreateMethod(name string, parameters []Parameter, fsm interface{}) error {
	var err error

	switch database.ReadConfig().Type {
	case database.TypeMongoDB:
		if database.CheckConnection() {
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C(MethodsCollection)
			method := &Method{
				ObjectID:  bson.NewObjectId(),
				Name: name,
				Parameters: parameters,
				Fsm: fsm.(simple_fsm.JsonRoot),
			}
			err = c.Insert(method)
		} else {
			err = NoDBConnection
		}
	default:
		err = DBNotSelected
	}

	if err != nil{
		log.Printf(MethodNotCreated, name)
	} else {
		log.Printf(MethodCreated, name)
	}
	return err
}

// PermissionByName
// Returns permission by given name and error
func MethodByName(name string) (*Method, error) {
	var err error
	var method Method
	switch database.ReadConfig().Type {
	case database.TypeMongoDB:
		if database.CheckConnection() {
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C(MethodsCollection)
			err = c.Find(bson.M{"name": name}).One(&method)
		} else {
			err = NoDBConnection
		}
	default:
		err = DBNotSelected
	}

	if err != nil{
		log.Printf(MethodNotFound, name)
	}

	return &method, err
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


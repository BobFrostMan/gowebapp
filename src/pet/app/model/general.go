package model

import (
	"pet/app/shared/database"
	"gopkg.in/mgo.v2/bson"
	"log"
)

func create(collection string, )  {
	
}

func MethodByName2(name string) (*Method, error) {
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
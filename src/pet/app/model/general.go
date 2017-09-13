package model

import (
	"pet/app/shared/database"
	"gopkg.in/mgo.v2/bson"
	"log"
	"gopkg.in/mgo.v2"
)

func List(collection string, by bson.M, limit int, resultObj []interface{}) ([]interface{}, error) {
	return EntityArray(collection, by, limit, resultObj)
}

func EntityArray(collection string, by bson.M, limit int, resultObj []interface{}) ([]interface{}, error) {
	var err error
	if err != nil {
		log.Println("Error occured during limit taking. Message: " + err.Error())
	}
	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C(collection)
		if limit == 0 {
			err = c.Find(by).All(&resultObj)
		} else {
			err = c.Find(by).Limit(limit).All(&resultObj)
		}
		log.Printf("Objects found: %v", resultObj)
	} else {
		err = NoDBConnection
	}
	if err != nil {
		log.Printf("Failed to find entity with params %s", resultObj)
	}
	return resultObj, err
}

func CreateEntity(collection string, insert bson.M) (interface{}, error) {
	var err error
	var insertQuery map[string]interface{}
	insertQuery = map[string]interface{}(insert)
	insertQuery["_id"] = bson.NewObjectId()
	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C(collection)
		c.Insert(bson.M(insertQuery))
		log.Printf("Insert: %v", bson.M(insertQuery))
	} else {
		err = NoDBConnection
	}
	if err != nil {
		log.Printf("Failed to find entity with params %s", insertQuery)
	}
	return insertQuery, err
}

func UpdateEntity(collection string, find bson.M, update bson.M, resultObj []interface{}) ([]interface{}, error) {
	var err error
	log.Printf("Find part: %v", find)
	log.Printf("Update part: %v", update)
	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C(collection)
		c.Find(find).Apply(mgo.Change{
			Update: bson.M{
				"$set": update,
			},
			ReturnNew: true,
		}, &resultObj)

	} else {
		err = NoDBConnection
	}
	if err != nil {
		log.Printf("Failed to find entity with find %s and update %s", find, update)
	}
	return resultObj, err
}
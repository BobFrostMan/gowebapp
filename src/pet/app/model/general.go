package model

import (
	"pet/app/shared/database"
	"gopkg.in/mgo.v2/bson"
	"log"
	"gopkg.in/mgo.v2"
)

func Entity(collection string, by bson.M, limit int, resultObj interface{}) (*interface{}, error) {
	var err error
	//count, err := strconv.Atoi(limit)
	if err != nil {
		log.Println("Error occured during limit taking. Message: " + err.Error())
	}
	switch database.ReadConfig().Type {
	case database.TypeMongoDB:
		if database.CheckConnection() {
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C(collection)
			//err = c.Find(bson.M{obj}).Limit(count).All(obj)
			log.Printf("By: %v", by)
			log.Printf("Result Object: %v", resultObj)
			err = c.Find(by).Limit(limit).One(&resultObj)
		} else {
			err = NoDBConnection
		}
	default:
		err = DBNotSelected
	}

	if err != nil {
		log.Printf("Failed to find entity with params %s", resultObj)
	}

	return &resultObj, err
}

func List(collection string, by bson.M, limit int, resultObj []interface{}) ([]interface{}, error) {
	//if limit == 1{
	return EntityArray(collection, by, limit, resultObj)
	/*} else {
		//TODO: handle this
		//return EntityArray(collection, by, limit, resultObj.(*interface{}))
		log.Println("Shouldn't get here!!!!")
		return &resultObj, nil
	}*/
}

func EntityArray(collection string, by bson.M, limit int, resultObj []interface{}) ([]interface{}, error) {
	var err error
	//count, err := strconv.Atoi(limit)
	if err != nil {
		log.Println("Error occured during limit taking. Message: " + err.Error())
	}

	switch database.ReadConfig().Type {
	case database.TypeMongoDB:
		if database.CheckConnection() {
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C(collection)
			//err = c.Find(bson.M{obj}).Limit(count).All(obj)
			//log.Printf("By: %v", by)
			if limit == 0 {
				err = c.Find(by).All(&resultObj)
			} else {
				err = c.Find(by).Limit(limit).All(&resultObj)
			}
			log.Printf("Result Object: %v", resultObj)
		} else {
			err = NoDBConnection
		}
	default:
		err = DBNotSelected
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

	switch database.ReadConfig().Type {
	case database.TypeMongoDB:
		if database.CheckConnection() {
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C(collection)
			c.Insert(bson.M(insertQuery))
			log.Printf("Insert: %v", bson.M(insertQuery))
		} else {
			err = NoDBConnection
		}
	default:
		err = DBNotSelected
	}

	if err != nil {
		log.Printf("Failed to find entity with params %s", insertQuery)
	}

	return insertQuery, err
}


func UpdateEntity(collection string, find bson.M, update bson.M, resultObj []interface{}) ([]interface{}, error) {
	var err error

	switch database.ReadConfig().Type {
	case database.TypeMongoDB:
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
	default:
		err = DBNotSelected
	}

	if err != nil {
		log.Printf("Failed to find entity with find %s and update %s", find,update)
	}

	return resultObj, err
}
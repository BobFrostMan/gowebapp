package database

import (
	"log"
	"gopkg.in/mgo.v2"
	"time"
)

var (
	Mongo *mgo.Session
	databases Info
)

type Info struct {
	Type    Type
	MongoDB MongoDBInfo
}

// Type is the type of database from a Type* constant
type Type string

const (
	TypeMongoDB Type = "MongoDB"
)

type MongoDBInfo struct {
	URL      string
	Database string
}

// Connect to the database
func Connect(d Info) {
	var err error

	// Store the config
	databases = d
	log.Printf("Selected database type is %s", d.Type)
	switch d.Type {

	case TypeMongoDB:
		// Connect to MongoDB
		if Mongo, err = mgo.DialWithTimeout(d.MongoDB.URL, 5 * time.Second); err != nil {
			log.Println("MongoDB Driver Error", err)
			return
		}

		// Prevents these errors: read tcp 127.0.0.1:27017: i/o timeout
		Mongo.SetSocketTimeout(1 * time.Second)

		// Check if is alive
		if err = Mongo.Ping(); err != nil {
			log.Println("Database Error", err)
		}
	default:
		log.Println("No registered database in config")
	}

}

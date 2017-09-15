package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"pet/app/shared/database"
	"log"
)

// Database tables, collections, fields etc.
const (
	TokensCollection = "Token"
)

type Token struct {
	ObjectID bson.ObjectId `bson:"_id" json:"_id"`
	Value string `json:"value"`
	UserId     string        `bson:"userId" json:"userId"`
	Expiration time.Time        `json:"expiration"`
}

// TokenByValue
// Finds token by given value
func TokenByValue(value string) (*Token, error) {
	var err error
	var token Token
	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()

		c := session.DB(database.ReadConfig().MongoDB.Database).C(TokensCollection)
		c.Find(bson.M{"value" : value}).One(&token)
	} else {
		err = NoDBConnection
	}
	if err != nil {
		log.Printf("Token '%s' wasn't found", value)
	}
	return &token, err
}

// CheckToken
// Checks if token exists expiration
func CheckToken(value string) (bool, error) {
	if token, err := TokenByValue(value);err != nil || token.Value == "" {
		return false, err
	} else {
		return !time.Now().After(token.Expiration), nil
	}
}
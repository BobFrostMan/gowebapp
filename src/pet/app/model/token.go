package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"pet/app/shared/database"
	"log"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
)

// Database tables, collections, fields etc.
const (
	TokensCollection = "Token"
	TokenExpirationDefaultInSec = 300
)

// Messages patterns
const (
	TokenNotFoundForUser = "Token wasn't found for user '%s'"
	TokenNotFound = "Token '%s' wasn't found"
	TokenNotCreated = "Token '%s' wasn't created"
	TokenCreated = "Token '%s' was successfully created for user '%s'"
)

type Token struct {
	ObjectID bson.ObjectId `bson:"_id" json:"_id"`
	Value string `json:"value"`
	UserId     string        `bson:"userId" json:"userId"`
	Expiration time.Time        `json:"expiration"`
}

// CreateUserToken
// Generates temporary access token for specified userId
func TokenCreate(userId string) (*Token, error) {
	var err error
	var value string
	var token Token

	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()

		c := session.DB(database.ReadConfig().MongoDB.Database).C(TokensCollection)
		value = uuid.NewV4().String()
		token = Token{
			ObjectID: bson.NewObjectId(),
			UserId: userId,
			Value: value,
			Expiration: time.Now().Add(time.Second * TokenExpirationDefaultInSec),
		}
		c.Insert(&token)
	} else {
		err = NoDBConnection
	}

	if err != nil {
		log.Printf(TokenNotCreated, userId)
	} else {
		log.Printf(TokenCreated, value, userId)
	}

	return &token, err
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
		log.Printf(TokenNotFound, value)
	}
	return &token, err
}

// TokenByUserId
// Finds token by userId
func TokenByUserId(userId string) (*Token, error) {
	var err error
	var token Token
	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C(TokensCollection)
		c.Find(bson.M{"userId" : userId}).One(&token)
	} else {
		err = NoDBConnection
	}
	if err != nil {
		log.Printf(TokenNotFoundForUser, userId)
	}
	return &token, err
}

// TokenUpdate
// Updates existing user token with new random value
func TokenUpdate(userId string) (*Token, error) {
	var err error
	var value string
	var token Token

	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()

		c := session.DB(database.ReadConfig().MongoDB.Database).C(TokensCollection)
		value = uuid.NewV4().String()
		c.Find(bson.M{"userId": userId}).Apply(mgo.Change{
			Update: bson.M{
				"$set": bson.M{
					"value": uuid.NewV4().String(),
					"expiration": time.Now().Add(time.Second * TokenExpirationDefaultInSec),
				},
			},
			ReturnNew: true,
		}, &token)

	} else {
		err = NoDBConnection
	}

	if err != nil {
		log.Printf(TokenNotCreated, userId)
	} else {
		log.Printf(TokenCreated, value, userId)
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
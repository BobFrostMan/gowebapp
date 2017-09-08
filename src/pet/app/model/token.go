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
	TokenNotFound = "Token wasn't found for user '%s'"
	TokenNotCreated = "Token '%s' wasn't created"
	TokenCreated = "Token '%s' was successfully created for user '%s'"
)

type Token struct {
	ID         bson.ObjectId `json:"id",bson:"_id"`
	Value string `json:"value"`
	UserId     string        `json:"userId"`
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
			UserId:userId,
			Value: value,
			Expiration:time.Now().Add(time.Second * TokenExpirationDefaultInSec),
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

// FindUserToken
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
		log.Printf(TokenNotFound, userId)
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
		//c.Update(bson.M{"userId": userId}, bson.M{"token", uuid.NewV4().String()})
		c.Find(bson.M{"userId": userId}).Apply(mgo.Change{
			Update: bson.M{
				"$set": bson.M{
					"token": uuid.NewV4().String(),
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

func TokenSet(userId string) (*Token, error){
	if existingToken, err := TokenByUserId(userId); err == nil && existingToken.Value != "" {
		return TokenUpdate(userId)
	} else {
		return TokenCreate(userId)
	}
}





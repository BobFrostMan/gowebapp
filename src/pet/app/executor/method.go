package executor

import (
	"net/http"
	"log"
	"pet/app/model"
	"pet/app/shared/passhash"
)

// auth
// Checks user existence by given login/pass
func auth(login string, password string) Result{
	log.Printf("Attempt to login as '%s'", login)
	user, err := model.UserByLogin(login)
	if err != nil {
		log.Printf("User '%s' not found", login)
		return Result{ Status: http.StatusForbidden, Data: "Credential data doesn't match to any user"}
	}
	if err = passhash.CompareHashAndPassword(user.Password, password); err != nil{
		log.Printf("User '%s' entered wrong password!", login)
		return Result{ Status: http.StatusForbidden, Data: "Credential data doesn't match to any user"}
	}
	log.Printf("User '%s' successfully logged in", login)
	// TODO: generate token here
	// TODO: save new token to db
	// TODO: set cookie (or think about frontend handling it)
	// TODO: after that return "token" as part of Data
	return Result {Status: http.StatusOK, Data: "{ \"token\" : \"token-value-staub\"}" }
}
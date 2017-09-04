package route

import (
	"pet/app/shared/server"
	"net/http"
	"pet/app/model"
	"log"
	"pet/app/shared/passhash"
	"strconv"
	"encoding/json"
)

// ConfigRoutes
// Registering handlers and binding them to according url
func ConfigRoutes() {
	http.Handle("/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	//TODO: remove unnecessary routes
	http.HandleFunc("/users", users)
	http.HandleFunc("/auth", auth)
}

// StartServer
// Starts server with instance parameters
func StartServer(server *server.Server)  {
	port := strconv.Itoa(server.Port)
	log.Printf("Starting server on port :%s", port)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}

// users
// Returns list of all users as json in payload
func users(w http.ResponseWriter, req *http.Request)  {
	users, _ := model.UserList()
	jsonUsers, err := json.Marshal(&users)
	if err != nil{
		log.Printf("Failed to parse users data:\n%s", jsonUsers)
	}
	w.Write(jsonUsers)
}

// auth
// Checks user existence by given login/pass
func auth(w http.ResponseWriter, req *http.Request){
	req.ParseForm()
	login := req.Form.Get("login")
	password := req.Form.Get("pass")
	log.Printf("Attempt to login as '%s'", login)
	user, err := model.UserByLogin(login)
	if err != nil {
		log.Printf("User '%s' not found", login)
		response(http.StatusForbidden, "Credential data doesn't match to any user", w)
		return
	}
	if err = passhash.CompareHashAndPassword(user.Password, password); err != nil{
		log.Printf("User '%s' entered wrong password!", login)
		response(http.StatusForbidden, "Credential data doesn't match to any user", w)
		return
	}
	log.Printf("User '%s' successfully logged in", login)
	// TODO: save session to db here
	// TODO: set session id here
}

// response
// Returns response, with given code and message, using given response writer
func response(status int, msg string, w http.ResponseWriter)  {
	w.WriteHeader(status)
	resp := ErrorResponse{Code: status, Message: msg}
	response, _ := json.Marshal(resp)
	w.Write(response)
}

type ErrorResponse struct{
	Code int `json:"code"`
	Message string `json:"message"`
}
package route

import (
	"gowebapp/pet/app/shared/server"
	"net/http"
	"gowebapp/pet/app/model"
	"log"
	"gowebapp/pet/app/shared/passhash"
	"io"
	"fmt"
	"strconv"
)

// ConfigRoutes
// Registering handlers and binding them to according url
func ConfigRoutes() {
	http.Handle("/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
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
	for _, user := range users{
		io.WriteString(w, fmt.Sprintf("%s", user))
	}
}

// auth
// Checks user existence by given login/pass
func auth(w http.ResponseWriter, req *http.Request){
	req.ParseForm()
	login := req.Form.Get("login")
	password := req.Form.Get("pass")
	user, err := model.UserByLogin(login)
	if err != nil {
		log.Printf("User '%s' not found", login)
		http.NotFound(w, req);
		return
	}
	if err = passhash.CompareHashAndPassword(user.Password, password); err != nil{
		log.Printf("User '%s' entered wrong password!", login)
		response(http.StatusForbidden, "Credential data doesn't match to any user", w)
		return
	}
	// TODO: save session to db here
	// TODO: set session id here
}

// response
// Returns response, with given code and message, using given response writer
func response(status int, msg string, w http.ResponseWriter)  {
	w.WriteHeader(status)
	io.WriteString(w, fmt.Sprintf("%s - %s", status, msg))
}
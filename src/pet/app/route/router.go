package route

import (
	"pet/app/shared/server"
	"net/http"
	"pet/app/model"
	"log"
	"strconv"
	"encoding/json"
	"pet/app/executor"
)

// ConfigRoutes
// Registering handlers and binding them to according url
func ConfigRoutes() {
	http.Handle("/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handle)
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

func handle(w http.ResponseWriter, req *http.Request)  {
	//Request parsing plus middleware requests logging
	req.ParseForm()
	log.Printf("Processing %s request to %s", req.Method, req.RequestURI)
	result, err := executor.Execute(req.Form)
	if (err != nil){
		log.Printf("[ERROR] Method %s %s executed with error: %s", req.Method, req.RequestURI, err.Error())
		log.Printf("[ERROR] Server response: %s", result)
	} else {
		log.Printf("[INFO] Method %s %s successfully executed", req.Method, req.RequestURI)
		log.Printf("[INFO] Server response: %s", result)
	}
	respond(&result, w)
}

//write result to ResponseWriter, need to be tested
func respond(res executor.Result, w http.ResponseWriter)  {
	w.WriteHeader(res.Status)
	response, res := json.Marshal(res)
	w.Write(response)
}
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
	http.Handle("/", http.StripPrefix("/static/", http.FileServer(http.Dir("src/static"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/api/", handle)
}

// StartServer
// Starts server with instance parameters
func StartServer(server *server.Server)  {
	port := strconv.Itoa(server.Port)
	log.Printf("[INFO] Starting server on port :%s", port)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}

// users
// Returns list of all users as json in payload
func users(w http.ResponseWriter, req *http.Request)  {
	users, _ := model.UserList()
	jsonUsers, err := json.Marshal(&users)
	if err != nil{
		log.Printf("[ERROR] Failed to parse users data:\n%s", jsonUsers)
	}
	w.Write(jsonUsers)
}

func handle(w http.ResponseWriter, req *http.Request)  {
	//Request parsing plus middleware requests logging
	req.ParseForm()
	log.Printf("[INFO] Processing %s request to %s", req.Method, req.RequestURI)
	result, err := executor.Execute(req.URL.Path, req.Form)
	if (err != nil){
		log.Printf("[ERROR] Method %s %s executed with error: %s", req.Method, req.RequestURI, err.Error())
		log.Printf("[ERROR] Server response: %v", result)
	} else {
		log.Printf("[ERROR] Method %s %s successfully executed", req.Method, req.RequestURI)
		log.Printf("[INFO] Server response: %v", result)
	}
	respond(*result, w)
}

//write result to ResponseWriter, need to be tested
func respond(res executor.Result, w http.ResponseWriter)  {
	w.WriteHeader(res.Status)
	response, err := json.Marshal(res)
	if err != nil{
		log.Println(err.Error())
	}
	w.Write(response)
}
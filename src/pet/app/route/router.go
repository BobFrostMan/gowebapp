package route

import (
	"pet/app/shared/server"
	"net/http"
	"log"
	"strconv"
	"encoding/json"
	"pet/app/executor"
	"pet/app/model"
	"fmt"
	"net/http/httputil"
)

var apiExecutor *executor.ApiExecutor

// ConfigRoutes
// Registering handlers and binding them to according url
func ConfigRoutes(methods *[]model.Method) {
	apiExecutor = executor.NewExecutor().LoadMethods(*methods).LoadStructure(*methods)

	http.Handle("/", http.StripPrefix("/static/", http.FileServer(http.Dir("src/static"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/api/", handle)
	http.HandleFunc("/reload-api", reloadApiMethods)
}

// StartServer
// Starts server with instance parameters
func StartServer(server *server.Server) {
	port := strconv.Itoa(server.Port)
	log.Printf("[INFO] Starting server on port :%s", port)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}

// handle
// A primary request handler function
// Contains request parsing plus middleware requests logging
func handle(w http.ResponseWriter, req *http.Request) {
	log.Printf("[INFO] Processing request '%s' metadata:\n%s", req.RequestURI, meta(req))
	request := executor.NewRequest(req)
	log.Printf("[INFO] Processing %s request %s", req.Method, request.MethodName)
	result, err := apiExecutor.Execute(request)
	if (err != nil) {
		log.Printf("[ERROR] Method %s executed with error: %v", request.MethodName, err)
		log.Printf("[ERROR] Server response: %v", result)
	} else {
		log.Printf("[INFO] Method %s successfully executed", request.MethodName)
		log.Printf("[INFO] Server response: %v", result)
	}
	respond(&result, w)
}

func meta(req *http.Request) string {
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	return string(requestDump)
}

// respond
// Write executor result to given response writer ResponseWriter
func respond(res *executor.Result, w http.ResponseWriter) {
	w.WriteHeader(res.Status)
	response, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		log.Println(err.Error())
	}
	w.Write(response)
}

// reloadApiMethods
// Reloads all api methods from database
func reloadApiMethods(w http.ResponseWriter, req *http.Request) {
	//TODO: add security support here, for L3/Admin only
	methods := *model.GetAllMethods()
	apiExecutor.ReloadMethods(methods).LoadStructure(methods)
	result := executor.NewResultMessage(http.StatusAccepted, "Reload methods procedure started")
	respond(&result, w)
}
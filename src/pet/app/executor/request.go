package executor

import (
	"net/http"
	"strings"
)

type Request struct {
	MethodName string        `json:"name"`
	Token      string        `json:"token"`
	Params     map[string]string `json:"params"`
}

func NewRequest(raw *http.Request) *Request {
	var request Request
	raw.ParseForm()

	request.MethodName = methodName(raw)
	request.Token = raw.Form.Get("token")
	request.Params = make(map[string]string)

	for k, _ := range raw.Form {
		//TODO: here we will get all headers even system we should re-think it
		request.Params[k] = raw.Form.Get(k)
	}
	return &request
}

// requestURL
// returns pretty request URL string from given request
func methodName(req *http.Request) string {
	var methodName string

	if strings.Contains(req.URL.Path, "/api/") {
		methodName = strings.Split(req.URL.Path, "/api/")[1]
	}

	return methodName
}



package executor

import (
	"net/http"
	"strings"
	"io/ioutil"
	"log"
	"encoding/json"
)

type Request struct {
	MethodName string        `json:"name"`
	Token      string        `json:"token"`
	Params     map[string]interface{} `json:"params"`
}

func NewRequest(raw *http.Request) *Request {
	var request Request
	raw.ParseForm()
	request.MethodName = getMethodName(raw)
	request.Params = make(map[string]interface{})

	// reading params from headers
	for name, headers := range raw.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request.Params[name] = h
		}
	}

	// reading params from url
	for k, _ := range raw.Form {
		request.Params[k] = raw.Form.Get(k)
	}

	//reading params from body
	if raw.Method == "PUT" || raw.Method == "POST"  {
		body, err := ioutil.ReadAll(raw.Body)
		defer raw.Body.Close()
		if len(body) > 0 {
			if err != nil {
				log.Printf("Cannot read request body! Message:%s", err.Error())
			} else {
				var payload map[string]*json.RawMessage
				err = json.Unmarshal(body, &payload)
				if err != nil {
					log.Printf("Failed to unmarshal object '%v'.\nMessage: %s", body, err.Error())
				} else {
					for k, _ := range payload {
						var str interface{}
						err = json.Unmarshal(*payload[k], &str)
						request.Params[k] = str
					}
				}
			}
		}
	}

	if token, present := request.Params["token"]; present{
		request.Token = token.(string)
	}

	return &request
}

// requestURL
// returns pretty request URL string from given request
func getMethodName(req *http.Request) string {
	var methodName string
	if strings.Contains(req.URL.Path, "/api/") {
		methodName = strings.Split(req.URL.Path, "/api/")[1]
	}
	return methodName
}



package executor

import (
	"github.com/xenzh/gofsm"
	"log"
	"pet/app/model"
	"fmt"
	"net/http"
	"errors"
	"encoding/json"
	"time"
)

type ApiExecutor struct {
	Methods      map[string]model.Method
	StructureMap map[string]*simple_fsm.Structure
	Actions      simple_fsm.ActionMap
}

// New
// Returns new instance of Api executor
func NewExecutor() *ApiExecutor {
	executor := new(ApiExecutor)
	executor.Methods = make(map[string]model.Method)
	executor.StructureMap = make(map[string]*simple_fsm.Structure)
	executor.Actions = make(simple_fsm.ActionMap)
	executor.initActionMap()
	return executor
}

func (a *ApiExecutor) initActionMap() {
	a.AddAction("set_result", setResult)
	a.AddAction("list", list)
	a.AddAction("create", create)
	a.AddAction("auth", authorize)
	a.AddAction("set_to_context", setToContext)
	a.AddAction("update", update)
}

func (a *ApiExecutor) AddAction(name string, action func(ctx simple_fsm.ContextOperator) error) {
	a.Actions[name] = func(ctx simple_fsm.ContextOperator) error {
		log.Printf("'%s' action started", name)
		err := action(ctx)
		if err != nil {
			log.Printf("'%s' action finished with error: %v", name, err)
			return err
		} else {
			log.Printf("'%s' action successfully finished", name)
			return nil
		}
	}
}

// LoadMethods
// Loading api methods to ApiExecutor
func (a *ApiExecutor) LoadStructure(methodsMap []model.Method) *ApiExecutor {
	for _, method := range methodsMap {
		a.Methods[method.Name] = method
		obj, _ := json.MarshalIndent(method.Fsm, "", "    ")
		log.Printf("Method '%s', Finite state machine as json:\n%v", method.Name, string(obj))
		structure, err := simple_fsm.NewBuilder(a.Actions).FromJsonType(method.Fsm).Structure()
		if err != nil {
			log.Fatalf("Failed to construct structure. Message: %s", err.Error())
		}
		a.StructureMap[method.Name] = structure
	}
	return a
}

// LoadMethods
// Loading api methods to ApiExecutor
func (a *ApiExecutor) LoadMethods(methodsMap []model.Method) *ApiExecutor {
	for _, method := range methodsMap {
		a.Methods[method.Name] = method
	}
	return a
}

// ReloadMethods
// Reloads methods to api executor
func (a *ApiExecutor) ReloadMethods(methodsMap []model.Method) *ApiExecutor {
	log.Println("Reloading api methods!")
	newMethods := make(map[string]model.Method)
	for _, method := range methodsMap {
		newMethods[method.Name] = method
	}
	a.Methods = newMethods
	return a
}

// Execute
// Parse parameters from request form, and executes api request, using Finite State Machine
func (a *ApiExecutor) Execute(request *Request) (Result, error) {
	method := a.Methods[request.MethodName]
	if method.IsEmpty() {
		msg := fmt.Sprintf("Method '%s' was not recognized by executor", request.MethodName)
		log.Printf("[ERROR] " + msg)
		return NewResultMessage(http.StatusBadRequest, msg), errors.New(msg)
	}

	ok, err := checkToken(request)
	if err != nil {
		return NewResultMessage(http.StatusBadRequest, err.Error()), err
	}
	if !ok {
		return NewResultMessage(http.StatusForbidden, "Provided token is not valid, or expired. Please provide, valid token or authorize with 'auth'"), nil
	}

	ok, err = validateParams(method, request.Params)
	if err != nil {
		return NewResultMessage(http.StatusBadRequest, err.Error()), err
	}
	if !ok {
		return NewResultMessage(http.StatusBadRequest, "Provided parameters are not valid"), nil
	}

	ok, err = checkPermissions(request)
	if err != nil {
		return NewResultMessage(http.StatusBadRequest, err.Error()), err
	}
	if !ok {
		return NewResultMessage(http.StatusForbidden, "No permissions to perform operation '" + request.MethodName + "'"), nil
	}

	result, err := a.executeRequest(request)
	if err != nil {
		return NewResultMessage(http.StatusInternalServerError, err.Error()), err
	}
	return result, err
}

// checkPermissions
// Checks user permissions to execute method
func checkPermissions(request *Request) (bool, error) {
	if request.MethodName == "auth" {
		return true, nil
	} else {
		if token, err := model.TokenByValue(request.Token); err == nil{
			log.Printf("Token id %v %T", token.UserId, token.UserId)
			if user, err := model.UserById(token.UserId); err == nil{
				return isOperationAllowed(user, request.MethodName), nil
			}
		} else {
			return false, err
		}
		return false, nil
	}
}

// IsAllowed
// Returns true if operation allowed for user object
func isOperationAllowed(user *model.User, operation string) bool {
	if groups, err := model.GetGroups(user.Groups); err == nil {
		for _, group := range groups {
			for _, permission := range group.Permissions {
				if permission.Value ==  operation {
					return permission.Execute
				}
			}
		}
	}
	return false
}

// checkToken
// Checks user token and it's expiration date
func checkToken(request *Request) (bool, error) {
	if request.MethodName == "auth" {
		return true, nil
	} else {
		if valid, err := model.CheckToken(request.Token); valid && err == nil {
			return true, nil
		} else {
			return false, err
		}
	}
}

// validateParams
// Returns true, nil if all required parameters with valid types specified
func validateParams(method model.Method, params map[string]interface{}) (bool, error) {
	var notSpecified []model.Parameter
	for _, param := range method.Parameters {
		value := params[param.Name]
		if param.Required {
			if value != "" && value != nil{
				actualType := getTypeName(value)
				if actualType != param.Type {
					msg := fmt.Sprintf("Wrong argument '%s' for method '%s'. Expected type '%s', but found '%s'", param.Name, method.Name, param.Type, actualType)
					log.Printf("[ERROR] " + msg)
					return false, errors.New(msg)
				}
			} else {
				notSpecified = append(notSpecified, param)
			}
		} else {
			if value != "" && value != nil {
				// optional parameters check
				actualType := getTypeName(value)
				if actualType != param.Type {
					msg := fmt.Sprintf("Wrong argument '%s' for method '%s'. Expected type '%s', but found '%s'", param.Name, method.Name, param.Type, actualType)
					log.Printf("[ERROR] " + msg)
					return false, errors.New(msg)
				}
			}
		}
	}

	if len(notSpecified) != 0 {
		var paramStr string = ""
		for _, param := range notSpecified {
			paramStr += fmt.Sprintf("'%s', ", param.Name)
		}
		msg := fmt.Sprintf("Required parameters are not provided for '%s' method. Please specify: %s", method.Name, paramStr[:len(paramStr) - 2])
		log.Printf("[ERROR] " + msg)
		return false, errors.New(msg)
	}
	return true, nil
}

// executeRequest
// Executes request and returns Result object
func (a *ApiExecutor) executeRequest(req *Request) (Result, error) {
	var fsm *simple_fsm.Fsm
	str := a.StructureMap[req.MethodName]
	fsm = simple_fsm.NewFsm(str)
	fsm.SetInput("methodName", req.MethodName)
	fsm.SetInput("start_date", time.Now())
	fsm.SetInput("failed", false)
	for k, v := range req.Params {
		fsm.SetInput(k, v)
	}
	execRes, err := fsm.Run()
	printFsmDump(fsm)

	if err != nil {
		log.Printf("Error occured during flow execution: %v", err)
	}
	log.Printf("Exec result %v", execRes)
	return NewResultFrom(execRes), nil
}

// printFsmDump
// prints fsm execution params to log in debug format
func printFsmDump(fsm *simple_fsm.Fsm) {
	r, er := fsm.Result()
	log.Printf("FSM state is completed?: %v", fsm.Completed())
	log.Printf("FSM state is fatal?: %v", fsm.Fatal())
	log.Printf("FSM Result: %v", r)
	log.Printf("FSM Error is: %v", er)
	log.Printf("Full FSM dump:\n%s", simple_fsm.Dump(fsm))
}

// getTypeName
// returns type name as string, use type checking to define type
func getTypeName(value interface{}) string {
	if _, ok := value.(int); ok{
		return "int"
	}
	if _, ok := value.(float64); ok{
		return "float64"
	}
	if _, ok := value.(bool); ok{
		return "bool"
	}
	if _, ok := value.(map[string]interface{}); ok{
		return "map[string]interface{}"
	}
	if _, ok := value.(json.RawMessage); ok{
		return "json.RawMessage"
	}
	return "string"
}
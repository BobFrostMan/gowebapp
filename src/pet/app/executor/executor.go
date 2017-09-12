package executor

import (
	"github.com/xenzh/gofsm"
	"log"
	"pet/app/model"
	"fmt"
	"reflect"
	"net/http"
	"errors"
	"encoding/json"
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
	a.AddAction("no_action", noAction)
	a.AddAction("auth", authorize)
}

func (a *ApiExecutor) AddAction(name string, action func(ctx simple_fsm.ContextOperator) error){
	a.Actions[name] = func(ctx simple_fsm.ContextOperator) error{
		log.Printf("'%s' action started", name)
		err := action(ctx)
		if err != nil{
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
		log.Printf("Method '%s' fsm:\n%v", method.Name, method.Fsm)

		obj, _ := json.MarshalIndent(method.Fsm, "", "    ")
		log.Printf("Fsm as json:\n%v", string(obj))

		structure, err := simple_fsm.NewBuilder(a.Actions).FromJsonType(method.Fsm).Structure()
		if err != nil {
			log.Printf("Failed to construct structure. Message: %s", err.Error())
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
		return Result{
			Status: http.StatusBadRequest,
			Data: msg,
		}, errors.New(msg)
	}

	ok, err := validateParams(method, request.Params)
	if err != nil {
		return Result{
			Status: http.StatusBadRequest,
			Data: err.Error(),
		}, err
	}

	if !ok {
		return Result{
			Status: http.StatusBadRequest,
			Data: "Provided parameters are not valid",
		}, nil
	}

	ok, err = checkPermissions(request)

	if err != nil {
		return Result{
			Status: http.StatusBadRequest,
			Data: err.Error(),
		}, err
	}
	if !ok {
		return Result{
			Status: http.StatusForbidden,
			Data: "No permissions to perform operation '" + request.MethodName + "'",
		}, nil
	}

	result, err := a.executeRequest(request)
	if err != nil {
		return Result{
			Status: http.StatusInternalServerError,
			Data: "Probably something happende will fix it later!",
		}, err
	}
	return result, err
}

// checkPermissions
// Checks user permissions to
func checkPermissions(request *Request) (bool, error) {
	if request.MethodName == "auth" {
		return true, nil
	} else {
		if exists, err := model.CheckToken(request.Token); exists && err == nil {
			//TODO: (when first action will be implemented) add exact action permission check
			return true, nil
		} else {
			return false, err
		}
	}
}

// validateParams
// Returns true, nil if all required parameters with valid types specified
func validateParams(method model.Method, params map[string]string) (bool, error) {
	var notSpecified []model.Parameter
	for _, param := range method.Parameters {
		value := params[param.Name]//form.Get(param.Name)
		if param.Required {
			if value != "" {
				actualType := reflect.TypeOf(value).String()
				if actualType != param.Type {
					msg := fmt.Sprintf("Wrong argument '%s' for method '%s'. Expected type '%s', but found '%s'", param.Name, method.Name, param.Type, actualType)
					fmt.Printf("[ERROR] " + msg)
					return false, errors.New(msg)
				}
			} else {
				notSpecified = append(notSpecified, param)
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

	if req.MethodName == "auth0" {
		//TODO: remove temporary condition. Replace it with Method - to action converter
		//TODO: add simple fsm creation here
		fsm = authWithFSM(req.Params["login"], req.Params["pass"])
	} else {
		str := a.StructureMap[req.MethodName]
		log.Printf("Structure map %v", str)
		fsm = simple_fsm.NewFsm(str)
		fsm.SetInput("methodName", req.MethodName)
		fsm.SetInput("failed", false)
		for k, v := range req.Params {
			fsm.SetInput(k, v)
		}
	}
	execRes, err := fsm.Run()
	printFsmDump(fsm)

	if err != nil {
		log.Printf("Error occured during flow execution: %v", err)
	}
	return execRes.(Result), nil
}

// authWithFSM
// Temporary action
func authWithFSM(login string, pass string) *simple_fsm.Fsm {
	fsm := simple_fsm.NewFsm(nil)
	return fsm
}

// printFsmDump
// prints fsm execution params to log in debug format
func printFsmDump(fsm *simple_fsm.Fsm) {
	r, er := fsm.Result()
	log.Printf("FSM state is running?: %v", fsm.Running())
	log.Printf("FSM state is completed?: %v", fsm.Completed())
	log.Printf("FSM state is idle?: %v", fsm.Idle())
	log.Printf("FSM state is fatal?: %v", fsm.Fatal())
	log.Printf("FSM Result: %v", r)
	log.Printf("FSM Error is: %v", er)
	log.Printf("Error kind is: %v", er)
	log.Printf("Error is nill?: %v", er == nil)
	log.Printf("Full FSM dump:\n%s", simple_fsm.Dump(fsm))
}
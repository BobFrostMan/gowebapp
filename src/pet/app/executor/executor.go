package executor

import (
	"github.com/xenzh/gofsm"
	"log"
	"pet/app/model"
	"net/url"
	"pet/app/shared/context"
	"strings"
	"fmt"
	"reflect"
	"net/http"
)


//TODO: executor should take Request entity as input

// loading all api methods to FSM from application context
func LoadFSM(methods []model.Method) {
	actions := createActionMap(methods)
	fsm, err := simple_fsm.NewBuilder(actions).Fsm();
	if (err != nil) {
		log.Fatalf("Error occured during FSM initialization: %s", err.Error())
	} else {
		//TODO: Fill FSM with states somehow (but how? probably smart parsing of FSM object + fsm.AddStates())
		context.AppContext.Put("fsm", fsm)
	}
}

// newAction
// Constructs action using apiMethod function loaded from database
func newAction(apiMethod *model.Method) simple_fsm.ActionFn {
	if apiMethod.Fsm == nil {
		return createSpecificAction(apiMethod)
	} else {
		return createGeneralAction(apiMethod)
	}
}

func createSpecificAction(apiMethod *model.Method) simple_fsm.ActionFn {
	switch apiMethod.Name {
	case "auth":
		return func(ctx simple_fsm.ContextOperator) error {
			res := auth(ctx.Str("login"), ctx.Str("pass"))
			//No fsm -> no transitions -> nothing else to do with context
			//TODO: should we remove values from FSM ctx after method execution?
			if (res.Status != http.StatusOK) {
				return error("Authentification failed! " + res.Data)
			}
			return
		}
	default:
		return func(ctx simple_fsm.ContextOperator) error {
			return error("Method " + apiMethod.Name + " wasn't correctly saved as db object ")
		}
	}
}

func createGeneralAction(apiMethod *model.Method) simple_fsm.ActionFn {
	return func(ctx simple_fsm.ContextOperator) error {
		//TODO: implement list create remove actions creation here
		return
	}
}

func createActionMap(methods []model.Method) map[string]simple_fsm.ActionFn {
	actions := make(map[string]simple_fsm.ActionFn)
	for _, method := range methods {
		actions[method.Name] = newAction(method)
	}
	return actions
}

func getStateInfos(apiMethod *model.Method) []simple_fsm.StateInfo {
	//TODO: somehow generate state infos to add to FSM
	return nil
}

func validateParams(apiUrl string, form url.Values) (bool, error) {
	// method existence check
	methodName := strings.Split(apiUrl, "/")[0]
	method := context.AppContext.GetMethod(methodName)
	if method == nil{
		msg := fmt.Sprintf("Method '%s' was not recognized by executor", methodName)
		log.Printf(msg)
		return false, error.Error(msg)
	}

	// param types checks
	for _, param := range method.Parameters{
		value := form.Get(param.Name)
		if param.Required && value != "" {
			//TODO: it's strong feeling that actual type will always be 'string'
			actualType := reflect.TypeOf(value).String()
			if actualType != param.Type {
				msg := fmt.Sprintf("Wrong argument '%s' for method '%s'. Expected type '%s', but found '%s'", param.Name, method.Name, param.Type, actualType)
				fmt.Printf(msg)
				return false, error.Error(msg)
			}
		}
	}
	return true, nil
}

func checkPermissions(token string) (bool, error) {
	//TODO: validate permissions for method, by token
	return false, nil
}

//TODO: have no idea how to run exact method on fsm
func executeMethod(method model.Method) (*Result, error) {
	var err error
	var result Result
	fsm := context.AppContext.GetFsm("fsm")
	//TODO: feed fsm with parameters from request somehow
	if fsm != nil {
		log.Println("FSM wasn't initialized yet!")
		err = "FSM wasn't initialized yet! Please init it with LoadFSM method first"
		return result, err
	} else {
		execRes, er := fsm.Run()
		if er != nil {
			log.Println("Error occured during flow execution: " + er.Error())
		}
		return execRes.(*Result), er
	}
}

func Execute(url string, form url.Values) (result *Result, err error) {
	var result Result
	var err error

	if ok, err := validateParams(url, form); !ok || err != nil {
		result = Result{
			Status: http.StatusBadRequest,
			Data: err.Error(),
		}
		return
	}

	if ok, err := checkPermissions(form.Get("token")); !ok || err != nil{
		result = Result{
			Status: http.StatusForbidden,
			Data: err.Error(),
		}
		return
	}

	//TODO: put all parameter values somewhere (to some FSM context or how can it be done?)
	//feedFSMWithArguments(form)

	//TODO: if fsm isn't running, run fsm (with parsed params somehow)
	result, err = executeMethod(nil)
	if err != nil {
		result = Result{
			Status: http.StatusInternalServerError,
			Data: err.Error(),
		}
	}
	return &result, err
}
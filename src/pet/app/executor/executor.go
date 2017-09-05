package executor

import (
	"github.com/xenzh/gofsm"
	"log"
	"pet/app/model"
	"net/url"
)

var fsm simple_fsm.Fsm

//TODO: executor should take Request entity as input


func initFsm() {
	//TODO: implement when FromJson will be implemented or remove totally anf use loadApiMethods() instead
	//fsm = simple_fsm.NewBuilder(nil).FromJson()
}

// loading all api methods from database to FSM
func LoadApiMethods() {
	var err simple_fsm.FsmError
	methods := model.GetAllMethods()
	actions := createActionMap(methods)
	fsm, err = simple_fsm.NewBuilder(actions).Fsm();
	if (err != nil) {
		log.Fatalf("Error occured during FSM initialization: %s", err.Error())
	} else {
		//TODO: Fill FSM with states somehow (but how? probably smart parsing of FSM object + fsm.AddStates())
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
			//TODO: should we remove values from ctx after method execution?
			if (res.Status != 200) {
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
	//TODO: somehow henerate state infos to add to FSM
	return nil
}

func validateParams() (bool, error) {
	//TODO: validate parameters (no db access)
	//TODO: validate FSM structure using FSM.validate()
	return false, nil
}

func checkPermissions() (bool, error) {
	//TODO: validate permissions including token
	return false, nil
}

func execute(actions map[string]simple_fsm.ActionFn) (*Result, error) {
	//TODO: execute method return result object
	//TODO: implement interaction with ContextOperator during fsm execution
	//TODO: implement fsm builder
	//TODO: implement fsm execution
	var result interface{}
	var err simple_fsm.FsmError
	fsm, err := simple_fsm.NewBuilder(actions).Fsm()
	if err != nil {
		log.Println("Failed to construct executive state machine: " + err.Error())
	} else {
		result, err = fsm.Run()
		if err != nil {
			log.Println("Error occured during flow execution: " + err.Error())
		}
	}
	return result.(*Result), err
}

func Execute(form url.Values) (*Result, error) {
	//TODO: locate method - else return error
	//TODO: put all parameter values somewhere (to some FSM context or how can it be done?)
	//TODO: validate parameters according to method parameters restriction (somewhere from FSM context?)
	//TODO: validate permissions for method, by token
	//TODO: if fsm isn't running, run fsm (with parsed params)
	//TODO: return general response with error or data inside
	return Result{}, nil
}
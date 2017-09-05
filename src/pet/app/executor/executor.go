package executor

import (
	"pet/app/executor/entity"
	"github.com/xenzh/gofsm"
	"log"
)

//TODO: executor should take Request entity as input

//TODO: define func signature
func constructMethod(req *executor.Request) (func(ctx simple_fsm.ContextOperator), error) {
	//get action name here, create function
	//add actionsMap entry to ctx
	//TODO: implement interaction with ContextOperator during fsm execution
	//TODO: locate a method (from context?) else return error
	//TODO: Method types:
	//TODO: auth
	//TODO: list
	//TODO: create
	//TODO: update
	//TODO: remove
	return nil, nil
}

func validateParams(req *executor.Request) (bool, error) {
	//TODO: validate parameters (no db access)
	//TODO: validate FSM structure using FSM.validate()
	return false, nil
}

func checkPermissions() (bool, error) {
	//TODO: validate permissions including token
	return false, nil
}

func execute(actions map[string]simple_fsm.ActionFn) (*executor.Result, error) {
	//TODO: execute method return result object
	//TODO: implement interaction with ContextOperator during fsm execution
	//TODO: implement fsm builder
	//TODO: implement fsm execution
	var result interface{}
	var err simple_fsm.FsmError
	fsm, err := simple_fsm.NewBuilder(actions).Fsm()
	if err != nil{
		log.Println("Failed to construct executive state machine: " + err.Error())
	} else {
		result, err = fsm.Run()
		if err != nil{
			log.Println("Error occured during flow execution: " + err.Error())
		}
	}
	return result.(*executor.Result), err
}
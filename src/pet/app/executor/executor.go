package executor

import (
	"github.com/xenzh/gofsm"
	"log"
	"pet/app/model"
	"fmt"
	"reflect"
	"net/http"
	"errors"
	"pet/app/shared/passhash"
)

type ApiExecutor struct {
	Methods map[string]model.Method
}

// New
// Returns new instance of Api executor
func NewExecutor() *ApiExecutor {
	executor := new(ApiExecutor)
	executor.Methods = make(map[string]model.Method)
	return executor
}

// LoadMethods
// Loading api methods to ApiExecutor
func (a *ApiExecutor) LoadMethods(methodsMap []model.Method) *ApiExecutor {
	for _, method := range methodsMap{
		a.Methods[method.Name] = method
	}
	return a
}

// ReloadMethods
// Reloads methods to api executor
func (a *ApiExecutor) ReloadMethods(methodsMap []model.Method) *ApiExecutor {
	log.Println("Reloading api methods!")
	newMethods := make(map[string]model.Method)
	for _, method := range methodsMap{
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

	result, err := executeRequest(request)
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
		if exists, err := model.CheckToken(request.Token); exists && err == nil{
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
func executeRequest(req *Request) (Result, error) {
	var fsm *simple_fsm.Fsm
	if req.MethodName == "auth" {
		//TODO: remove temporary condition. Replace it with Method - to action converter
		fsm = authWithFSM(req.Params["login"], req.Params["pass"])
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
	start := simple_fsm.NewState("start",
		simple_fsm.NewTransitionAlways("start-find_user", "find_user",
			func(ctx simple_fsm.ContextOperator) error {
				log.Println("'find_user' action started")
				err := ctx.Put("login", login)
				er := ctx.Put("pass", pass)

				if err != nil || er != nil {
					log.Println("Internal error, can't put values to transition context!: " + err.Error())
					//TODO: how can I return error value when I don't want to?
					return errors.New("Internal error, can't put values to transition context!")
				}
				log.Println("'find_user' action successfully finished")
				return nil
			},

		))
	findUser := simple_fsm.NewState("find_user", []simple_fsm.Transition{
		simple_fsm.NewTransition("find_user-user_found", "user_found",
			func(ctx simple_fsm.ContextAccessor) (bool, error) {
				log.Println("'user_found' guard started")
				if !ctx.Has("login") || !ctx.Has("pass") {
					log.Println("'user_found' can't get login/password from context!")
					return false, errors.New("Login or password were not provided for context")
				}

				login, _ := ctx.Str("login")
				//TODO: two similar requests, fix it
				_, err := model.UserByLogin(login)
				if err != nil {
					log.Printf("'user_found' User '%s' not found", login)
					return false, nil
				}
				log.Println("'user_found' guard successfully finished")
				return true, nil
			},
			func(ctx simple_fsm.ContextOperator) error {
				log.Println("'user_found' action started")
				login, _ := ctx.Str("login")
				//TODO: two similar requests, fix it
				userObj, _ := model.UserByLogin(login)
				ctx.Put("user", userObj)
				log.Println("'user_found' action successfully finished")
				return nil
			}),
		simple_fsm.NewTransition("find_user-user_not_found", "user_not_found",
			func(ctx simple_fsm.ContextAccessor) (bool, error) {
				log.Println("'user_not_found' guard started")
				_, err := model.UserByLogin(login)
				if err != nil {
					log.Printf("User '%s' not found", login)
					return true, nil
				}
				log.Println("'user_not_found' guard successfully finished")
				return false, nil
			},
			func(ctx simple_fsm.ContextOperator) error {
				ctx.PutResult(Result{
					Status: http.StatusForbidden,
					Data: "Credential data doesn't match to any user",
				})
				log.Println("'user_not_found' state aquired")
				return nil
			}),
	})

	userFound := simple_fsm.NewState("user_found", []simple_fsm.Transition{
		simple_fsm.NewTransition("user_found-pass_ok", "pass_ok",
			func(ctx simple_fsm.ContextAccessor) (bool, error) {
				log.Println("'user_found' guard started")
				if !ctx.Has("user") {
					return false, errors.New("User wasn't found inside context")
				}
				log.Println("'user_found' guard is ok")
				return true, nil
			},
			func(ctx simple_fsm.ContextOperator) error {
				log.Println("'user_found' action started")
				rawUser, _ := ctx.Raw("user")
				pass, _ := ctx.Str("pass")
				userObj := rawUser.(*model.User)
				if err := passhash.CompareHashAndPassword(userObj.Password, pass); err != nil {
					log.Printf("User '%s' entered wrong password!", login)
					ctx.PutResult(Result{
						Status: http.StatusForbidden,
						Data: "Wrong password specified",
					})
					return nil
				}
				log.Println("'user_found' action ok")
				return nil
			}),
		simple_fsm.NewTransition("user_found-pass_ok", "pass_not_ok",
			func(ctx simple_fsm.ContextAccessor) (bool, error) {
				log.Println("'pass_not_ok' guard started")
				if !ctx.Has("result") {
					log.Println("'pass_not_ok' state can't be acquired")
					return false, nil
				}
				log.Println("'pass_not_ok' guard successfully finished")
				return true, nil
			},
			func(ctx simple_fsm.ContextOperator) error {
				log.Println("'pass_not_ok' state aquired")
				return nil
			}),
	})
	passOk := simple_fsm.NewState("pass_ok", []simple_fsm.Transition{
		simple_fsm.NewTransition("pass_ok-token_created", "token_created",
			func(ctx simple_fsm.ContextAccessor) (bool, error) {
				log.Println("'pass_ok' guard started")
				if !ctx.Has("result") {
					log.Println("'pass_ok' guard is ok")
					return true, nil
				}
				return false, nil
			},
			func(ctx simple_fsm.ContextOperator) error {
				log.Println("'pass_ok' action started")
				rawUser, _ := ctx.Raw("user")
				userObj := rawUser.(*model.User)
				log.Printf("User pass ok, creating token for %v", userObj)
				token, err := model.TokenSet(userObj.ObjectID.Hex())
				if err != nil {
					log.Printf("Token wasn't created for user'%s' entered wrong password!", userObj.Name)
					ctx.PutResult(Result{
						Status: http.StatusInternalServerError,
						Data: "Token wasn't created for user: " + userObj.Name,
					})
					return errors.New("Token wasn't created for user: " + userObj.Name)
				}
				ctx.Put("token", token)
				log.Printf("Token created for user %v\n", userObj.Name)
				ctx.PutResult(Result{
					Status: http.StatusOK,
					Data: token,
				})
				log.Println("'pass_ok' action successfully finished")
				return nil
			}),
		simple_fsm.NewTransition("pass_ok-failed_to_create_token", "failed_to_create_token",
			func(ctx simple_fsm.ContextAccessor) (bool, error) {
				log.Println("'failed_to_create_token' guard started")
				if !ctx.Has("result") {
					log.Println("'failed_to_create_token' state can't be acquired")
					return false, nil
				}
				log.Println("'failed_to_create_token' guard successfully finished")
				return true, nil
			},
			func(ctx simple_fsm.ContextOperator) error {
				log.Println("'failed_to_create_token' state aquired")
				return nil
			}),
	},
	)
	tokenCreated := simple_fsm.NewState("token_created", nil)
	failedToCreateToken := simple_fsm.NewState("failed_to_create_token", nil)
	passNotOk := simple_fsm.NewState("pass_not_ok", nil)
	userNotFound := simple_fsm.NewState("user_not_found", nil)

	structure := simple_fsm.NewStructure()
	//state -> parent
	structure.AddStartState(start, nil)
	structure.AddState(findUser, start)
	structure.AddState(userFound, findUser)
	structure.AddState(userNotFound, findUser)
	structure.AddState(passOk, userFound)
	structure.AddState(passNotOk, userFound)
	structure.AddState(tokenCreated, passOk)
	structure.AddState(failedToCreateToken, passOk)

	fsm := simple_fsm.NewFsm(structure)
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
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
	"errors"
	"pet/app/shared/passhash"
	"encoding/json"
)
//TODO: executor should take Request entity as input

// loading all api methods to FSM from application context
func LoadFSM(methods []model.Method) {
	//TODO: cannot use builder cause actions map is private
	//actions := createActionMap(methods)
	//fsm, err := simple_fsm.NewBuilder(&actions).Fsm();
	//fsm := simple_fsm.NewFsm();
	//startTransition = simple_fsm.Transition[]{}simple_fsm.NewTransition("finding_user",  )
	//fsm.AddStartState(simple_fsm.NewState("startsub", ), nil)
	/*
	if (err != nil) {
		log.Fatalf("Error occured during FSM initialization: %s", err.Error())
	} else {
	*/
		//TODO: Fill FSM with states somehow (but how? probably smart parsing of FSM object + fsm.AddStates())
		///context.GlobalCtx.Put("fsm", fsm)
	//}
}

func authWithFSM(login string, pass string) *simple_fsm.Fsm {
	fsm := simple_fsm.NewFsm();
	start := simple_fsm.NewState("start",
			simple_fsm.NewTransitionAlways("start-find_user", "find_user",
			func(ctx simple_fsm.ContextOperator) error {
				err := ctx.Put("login", login)
				er := ctx.Put("pass", pass)
				log.Println("Putting logging and pass to context")
				if err != nil || er != nil {
					log.Println("Internal error, can't put values to transition context!: " + err.Error())
					//TODO: how can I return error value when I don't want to?
					return errors.New("Internal error, can't put values to transition context!")
				}
				log.Println("Returning nil instead of error")
				return nil
			},

	))
	findUser := simple_fsm.NewState("find_user", []simple_fsm.Transition{
		simple_fsm.NewTransition("find_user-user_found", "user_found",
			func(ctx simple_fsm.ContextAccessor) (bool, error){
				log.Println("'user_find' guard started")
				if _, err := ctx.Str("login"); err != nil {
					log.Println("'user_find' can't get login from context: " + err.Error())
					return false, err
				}
				//TODO: replace with has
				if _, err := ctx.Str("pass"); err != nil {
					log.Println("'user_find' can't get pass from context: " + err.Error())
					return false, err
				}
				log.Println("'user_find' returning nill - no error")
				return true, nil
			},
			func(ctx simple_fsm.ContextOperator) error {
				login, _ := ctx.Str("login")
				userObj, err := model.UserByLogin(login)
				if err != nil {
					//TODO: not sure if we finish execution here if error returned
					//TODO: or should I do:
					//TODO: simple_fsm.NewTransition("find_user-user_not_found", "user_not_found",
					log.Printf("User '%s' not found", login)
					ctx.PutResult( Result{
						Status: http.StatusForbidden,
						Data: "Credential data doesn't match to any user",
					})
					return err
				}
				log.Printf("Putting user to context '%v'", userObj)
				ctx.Put("user", userObj)
				return nil
			}),
	})

	userFound := simple_fsm.NewState("user_found", []simple_fsm.Transition{
		simple_fsm.NewTransition("user_found-pass_ok", "pass_ok",
			func(ctx simple_fsm.ContextAccessor) (bool, error){
				log.Println("'user_found' guard started")
				if _, err := ctx.Raw("user"); err != nil {
					return false, err
				}

				log.Println("'user_found' guard is ok")
				return true, nil
			},
			func(ctx simple_fsm.ContextOperator) error {
				log.Println("'user_found' action started")
				rawUser, _ := ctx.Raw("user")
				pass, _ := ctx.Str("pass")
				userObj := rawUser.(*model.User)
				log.Printf("'user_found' user object is %v", userObj)
				if err := passhash.CompareHashAndPassword(userObj.Password, pass); err != nil {
					log.Printf("User '%s' entered wrong password!", login)
					ctx.PutResult( Result{
						Status: http.StatusForbidden,
						Data: "Wrong password specified",
					})
					return err
				}
				log.Println("'user_found' action ok")
				return nil
			}),
	})
	passOk := simple_fsm.NewState("pass_ok",
		simple_fsm.NewTransitionAlways("pass_ok-token_created", "token_created",
			func(ctx simple_fsm.ContextOperator) error {
				rawUser, _ := ctx.Raw("user")
				userObj := rawUser.(*model.User)
				log.Printf("User pass ok, creating token for %v", userObj)
				token, err := model.TokenSet(string(userObj.ObjectID))
				if err != nil {
					log.Printf("Token wasn't created for user'%s' entered wrong password!", userObj.Name)
					ctx.PutResult(Result{
						Status: http.StatusInternalServerError,
						Data: "Token wasn't created ",
					})
					return err
				}
				ctx.Put("token", token)
				token_str, _:= json.Marshal(token)
				log.Printf("Token created! for %v\n", userObj)
				ctx.PutResult(Result{
					Status: http.StatusCreated,
					Data: string(token_str),
				})
				return nil
			}),
	)
	tokenCreated:= simple_fsm.NewState("token_created", []simple_fsm.Transition{})
	//global := simple_fsm.NewState(simple_fsm.FsmGlobalStateName, nil)
	//fsm.AddStartState(start, global)
	//fsm.AddState(start, global)
	fsm.AddStartState(start, nil)
	fsm.AddState(findUser, start)
	fsm.AddState(userFound, findUser)
	fsm.AddState(passOk, userFound)
	fsm.AddState(tokenCreated, passOk)

	/*
	//fmt.Printf("%v", fsm)
	step, err := fsm.Advance()
	log.Printf("Step history from/to/transition: %v", step)
	log.Printf("Error: %v", err)

	step, err = fsm.Advance()
	log.Printf("Step history from/to/transition: %v", step)
	log.Printf("Error: %v", err)

	step, err = fsm.Advance()
	log.Printf("Step history from/to/transition: %v", step)
	log.Printf("Error: %v", err)
	*/

	/*
	fsm.AddStates(nil, start,
		findUser,
		userFound,
		passOk,
		tokenCreated,
	)*/

	/*
	//fsm.AddStartState(start, nil)
	fsm.AddStates(nil, start,
		simple_fsm.NewState("find_user", []simple_fsm.Transition{
				simple_fsm.NewTransition("find_user-user_found", "user_found",
				func(ctx simple_fsm.ContextAccessor) (open bool, err error){
					if _, err = ctx.Str("login"); err != nil {
						return false, err
					}
					if _, err = ctx.Str("pass"); err != nil {
						return false, err
					}
					return true, nil
				},
				func(ctx simple_fsm.ContextOperator) error {
					login, _ := ctx.Str("login")
					userObj, err := model.UserByLogin(login)
					if err != nil {
						//TODO: not sure if we finish execution here if error returned
						//TODO: or should I do:
						//TODO: simple_fsm.NewTransition("find_user-user_not_found", "user_not_found",
						log.Printf("User '%s' not found", login)
						ctx.PutResult( Result{
								Status: http.StatusForbidden,
								Data: "Credential data doesn't match to any user",
						})
						return err
					}
					ctx.Put("user", userObj)
					return nil
				}),
		}),

		simple_fsm.NewState("user_found", []simple_fsm.Transition{
				simple_fsm.NewTransition("user_found-pass_ok", "pass_ok",
					func(ctx simple_fsm.ContextAccessor) (open bool, err error){
						if _, err := ctx.Raw("user"); err != nil {
							return false, err
						}
						return true, nil
					},
					func(ctx simple_fsm.ContextOperator) error {
						rawUser, _ := ctx.Raw("user")
						pass, _ := ctx.Str("pass")
						userObj := rawUser.(model.User)

						if err := passhash.CompareHashAndPassword(userObj.Password, pass); err != nil {
							log.Printf("User '%s' entered wrong password!", login)
							ctx.PutResult( Result{
									Status: http.StatusForbidden,
									Data: "Wrong password specified",
							})
							return err
						}
						return nil
					}),
		}),

		//if we already validated user password - there are no reasons not to try create session
		simple_fsm.NewState("pass_ok",
				simple_fsm.NewTransitionAlways("pass_ok-token_created", "token_created",
					func(ctx simple_fsm.ContextOperator) error {
						rawUser, _ := ctx.Raw("user")
						userObj := rawUser.(model.User)
						token, err := model.TokenSet(string(userObj.ID))
						if err != nil {
							log.Printf("Token wasn't created for user'%s' entered wrong password!", userObj.Name)
							ctx.PutResult(Result{
								Status: http.StatusInternalServerError,
								Data: "Token wasn't created ",
							})
							return err
						}
						ctx.Put("token", token)
						return nil
					}),
		),
		simple_fsm.NewState("token_created", []simple_fsm.Transition{}),

	)*/
	return fsm
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
			login, err := ctx.Str("login")
			pass, er := ctx.Str("pass")
			if err != nil || er != nil {
				res := auth(login, pass)
				//No fsm -> no transitions -> nothing else to do with context
				//TODO: should we remove values from FSM ctx after method execution?
				if (res.Status != http.StatusOK) {
					return errors.New("Authentification failed! " + res.Data)
				}
			}
			return nil
		}
	default:
		return func(ctx simple_fsm.ContextOperator) error {
			return errors.New("Method " + apiMethod.Name + " wasn't correctly saved as db object ")
		}
	}
}

func createGeneralAction(apiMethod *model.Method) simple_fsm.ActionFn {
	return func(ctx simple_fsm.ContextOperator) error {
		//TODO: implement list create remove actions creation here
		return errors.New("Not implemented yet!")
	}
}
//TODO: can't return *simple_fsm.actionMap cause it's private
func createActionMap(methods []model.Method) interface {}{
	actions := make(map[string]simple_fsm.ActionFn)
	for _, method := range methods {
		actions[method.Name] = newAction(&method)
	}
	return actions
}

func getStateInfos(apiMethod *model.Method) []simple_fsm.StateInfo {
	//TODO: somehow generate state infos to add to FSM
	return nil
}

func validateParams(method model.Method, form url.Values) (bool, error) {

	// param types checks
	for _, param := range method.Parameters{
		value := form.Get(param.Name)
		if param.Required && value != "" {
			//TODO: it's strong feeling that actual type will always be 'string'
			actualType := reflect.TypeOf(value).String()
			if actualType != param.Type {
				msg := fmt.Sprintf("Wrong argument '%s' for method '%s'. Expected type '%s', but found '%s'", param.Name, method.Name, param.Type, actualType)
				fmt.Printf("[ERROR] " + msg)
				return false, errors.New(msg)
			}
		}
	}
	return true, nil
}

func checkPermissions(methodName string, token string) (bool, error) {
	if methodName == "auth" {
		return true, nil
	} else {
		//TODO: validate permissions for method, by token
		return false, errors.New("Method " + methodName  + " not supported for check permission operation.")
	}
}

//TODO: have no idea how to run exact method on fsm
func executeMethod(methodName string, form url.Values) (Result, error) {
	//fsm := context.GlobalCtx.GetFsm("fsm")
	//TODO: feed fsm with parameters from request somehow
	/*
		if fsm != nil {
			log.Println("FSM wasn't initialized yet!")
			err = "FSM wasn't initialized yet! Please init it with LoadFSM method first"
			return err
		} else {
	*/
	var fsm *simple_fsm.Fsm
	if methodName == "auth" {
		fsm = authWithFSM(form.Get("login"), form.Get("pass"))
	}
	if err := fsm.Validate(); err != nil {
		log.Println("FSM wasn't properly created!")
		return Result{Status:http.StatusInternalServerError, Data: "FSM wasn't properly created!: " +err.Error()}, err
	}
	execRes, err := fsm.Run()

	if err != nil {
		//TODO: handle!
		//log.Println("Error occured during flow execution: " + err.Error())
	}
	return execRes.(Result), err
}

func Execute(url string, form url.Values) (Result, error) {

	// method existence check
	methodName := strings.Split(url, "/")[2]
	method := context.GlobalCtx.GetMethod(methodName)
	if method.IsEmpty() {
		msg := fmt.Sprintf("Method '%s' was not recognized by executor", methodName)
		log.Printf("[ERROR] " + msg)
		return Result{
			Status: http.StatusBadRequest,
			Data: msg,
		}, errors.New(msg)
	}

	if ok, err := validateParams(method, form); !ok || err != nil {
		return Result{
			Status: http.StatusBadRequest,
			Data: err.Error(),
		}, err
	}

	if ok, err := checkPermissions(methodName, form.Get("token")); !ok || err != nil {
		return Result{
			Status: http.StatusForbidden,
			Data: err.Error(),
		}, err
	}

	//TODO: put all parameter values somewhere (to some FSM context or how can it be done?)
	//feedFSMWithArguments(form)

	//TODO: if fsm isn't running, run fsm (with parsed params somehow)
	result, err := executeMethod(methodName, form)
	if err != nil {
		//fmt.Println("Error after method execution " + err.Error())
		return Result{
			Status: http.StatusInternalServerError,
			Data: "Probably something happende will fix it later!", //err.Error(),
		}, err
	}
	return result, err
}
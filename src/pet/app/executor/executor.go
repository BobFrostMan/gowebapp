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
)
//TODO: executor should take Request entity as input

// loading all api methods to FSM from application context
func LoadFSM(methods []model.Method) {
	//TODO: cannot use builder cause actions map is private
	//actions := createActionMap(methods)
	//fsm, err := simple_fsm.NewBuilder(&actions).Fsm();
	fsm := simple_fsm.NewFsm();
	//startTransition = simple_fsm.Transition[]{}simple_fsm.NewTransition("finding_user",  )
	//fsm.AddStartState(simple_fsm.NewState("startsub", ), nil)
	/*
	if (err != nil) {
		log.Fatalf("Error occured during FSM initialization: %s", err.Error())
	} else {
	*/
		//TODO: Fill FSM with states somehow (but how? probably smart parsing of FSM object + fsm.AddStates())
		context.GlobalCtx.Put("fsm", fsm)
	//}
}

func authWithFSM(login string, pass string) *simple_fsm.Fsm {
	fsm := simple_fsm.NewFsm();
	start := simple_fsm.NewState("start",
		simple_fsm.NewTransitionAlways("start-find_user", "find_user",
			func(ctx simple_fsm.ContextOperator) error {
				err := ctx.Put("login", login)
				er := ctx.Put("pass", pass)

				if err != nil || er != nil {
					return errors.New("Internal error, can't put values to transition context!")
				}
				return nil
			},
		))
	findUser := simple_fsm.NewState("find_user", []simple_fsm.Transition{
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
	})

	userFound := simple_fsm.NewState("user_found", []simple_fsm.Transition{
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
	})
	passOk := simple_fsm.NewState("pass_ok",
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
	)
	tokenCreated:= simple_fsm.NewState("token_created", []simple_fsm.Transition{})
	//global := simple_fsm.NewState(simple_fsm.FsmGlobalStateName, nil)
	//fsm.AddStartState(start, global)
	//fsm.AddState(start, global)

	fsm.AddState(start, nil)
	fsm.AddState(findUser, start)
	fsm.AddState(userFound, findUser)
	fsm.AddState(passOk, userFound)
	fsm.AddState(tokenCreated, passOk)

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

func checkPermissions(token string) (bool, error) {
	//TODO: validate permissions for method, by token
	return false, nil
}

//TODO: have no idea how to run exact method on fsm
func executeMethod(methodName string, form url.Values) (*Result, error) {
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
		return &Result{Status:http.StatusInternalServerError, Data: "FSM wasn't properly created!: " +err.Error()}, err
	}
	execRes, err := fsm.Run()

	if err != nil {
		log.Println("Error occured during flow execution: " + err.Error())
	}
	return execRes.(*Result), err
}

func Execute(url string, form url.Values) (result *Result, err error) {

	// method existence check
	methodName := strings.Split(url, "/")[2]
	method := context.GlobalCtx.GetMethod(methodName)
	if method.IsEmpty() {
		msg := fmt.Sprintf("Method '%s' was not recognized by executor", methodName)
		log.Printf("[ERROR] " + msg)
		return &Result{
			Status: http.StatusBadRequest,
			Data: err.Error(),
		}, errors.New(msg)
	}

	if ok, err := validateParams(method, form); !ok || err != nil {
		return &Result{
			Status: http.StatusBadRequest,
			Data: err.Error(),
		}, err
	}

	if methodName != "auth" {
		if ok, err := checkPermissions(form.Get("token")); !ok || err != nil {
			return &Result{
				Status: http.StatusForbidden,
				Data: err.Error(),
			}, err
		}
	}

	//TODO: put all parameter values somewhere (to some FSM context or how can it be done?)
	//feedFSMWithArguments(form)

	//TODO: if fsm isn't running, run fsm (with parsed params somehow)
	result, err = executeMethod(methodName, form)
	if err != nil {
		return &Result{
			Status: http.StatusInternalServerError,
			Data: err.Error(),
		}, err
	}
	return result, err
}
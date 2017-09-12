package executor

import (
	"log"
	"pet/app/model"
	"github.com/xenzh/gofsm"
	"fmt"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"pet/app/shared/passhash"
	"encoding/json"
)

func setResult(ctx simple_fsm.ContextOperator) error {
	failed := get("failure", ctx)
	if failed != nil{
		ctx.PutResult(failed)
	} else {
		if !ctx.Has("result") {
			var msg string
			if name,_ := ctx.Str("methodName"); name != ""{
				msg = fmt.Sprintf("Flow '%s' structure, wasn't formed correctly! No result was set for this flow!", name)
			} else {
				msg = "Current flow structure, wasn't formed correctly! No result was set for this flow!"
			}
			ctx.PutResult(Result{
				Status:http.StatusInternalServerError,
				Data: msg,
			})
			log.Println("Result wasn't set by any action!!!")
		}
		// else already set!
	}
	return nil
}

func noAction (ctx simple_fsm.ContextOperator) error {
	//do nothing
	return nil
}

func list(ctx simple_fsm.ContextOperator) error {
	target := getStr("target", ctx)
	limit := getInt("limit", ctx)
	findQuery := find(ctx)
	entity, er := model.List(target, findQuery, limit, getResultObject("type", limit))
	if er != nil {
		log.Printf("Error during db request: %v", er)
		msg := fmt.Sprintf("Can not find entity by query: '%v'. Message: %s", findQuery, er)
		ctx.Put("failed", true)
		ctx.Put("failure", Result{
			Status: http.StatusInternalServerError,
			Data: msg,
		})
	}
	log.Printf("Entity set to context\n%v", entity)
	ctx.Put("exists", true)
	ctx.Put("entity", entity)
	//TODO: remove tmp
	/*ctx.PutResult(Result{
		Status: http.StatusOK,
		Data: entity,
	})*/
	return nil
}

func create(ctx simple_fsm.ContextOperator) error {
	target := getStr("target", ctx)
	insertQuery := find(ctx)
	entity, er := model.CreateEntity(target, insertQuery)
	if er != nil {
		log.Printf("Error during db request: %v", er)
		ctx.Put("failed", true)
		msg := fmt.Sprintf("Can not insert entity with query: '%v'. Message: %s", insertQuery, er)
		ctx.Put("failure", Result{
			Status: http.StatusInternalServerError,
			Data: msg,
		})
		return nil
	}
	log.Printf("Entity set to context\n%v", entity)
	ctx.Put("entity", entity)
	//TODO: remove tmp
	ctx.PutResult(Result{
		Status: http.StatusOK,
		Data: entity,
	})
	return nil
}

func authorize(ctx simple_fsm.ContextOperator) error {
	rawUser, _ := ctx.Raw("entity")
	pass, _ := ctx.Str("pass")
	log.Printf("User: %v", rawUser)
	usrs := []model.User{}
	usersBytes, _ := json.Marshal(rawUser)
	er := json.Unmarshal(usersBytes, &usrs)
	log.Printf("Unmarshalled err: %v", er)
	log.Printf("Unmarshalled user: %v", usrs)
	userObj := usrs[0]
	if err := passhash.CompareHashAndPassword(userObj.Password, pass); err != nil {
		log.Printf("User '%s' entered wrong password!", userObj.Name)
		ctx.Put("failed", true)
		ctx.Put("failure", Result{
			Status: http.StatusForbidden,
			Data: "Wrong password specified",
		})
		return nil
	}
	//TODO: remove this synthetic result
	ctx.PutResult(Result{
		Status: http.StatusOK,
		Data: usrs,
	})

	return nil
}

func find(ctx simple_fsm.ContextOperator) bson.M  {
	by := findBy(ctx)
	return findRequest(by, ctx)
}

func findRequest(args []string, ctx simple_fsm.ContextOperator) bson.M {
	result := make(map[string]interface{})
	for _, arg := range args {
		//TODO: looks risky
		if arg != "" {
			value, err := ctx.Raw(arg)
			if err != nil {
				log.Printf("Error during form find request. Message: %v", err)
			} else {
				result[arg] = value
			}
		}
	}
	return bson.M(result)
}

//TODO: do we really need this?
func getResultObject(objectType string, limit int) []interface{} {
	res := make([]interface{}, limit)
	switch{
	case objectType == "User":
	//res := make([]model.User, limit)
	//result := &res
	/*for i := range a {
		b[i] = a[i]
	}
	return result.(*[]interface{})*/
	default:
		return nil
	}
	return res
}

func getStr(key string, ctx simple_fsm.ContextOperator) string {
	value, err := ctx.Str(key)
	if err != nil {
		log.Printf("Property '%s' was not found in context! %v", key, err)
	}
	return value
}

func getInt(key string, ctx simple_fsm.ContextOperator) int {
	value, err := ctx.Int(key)
	if err != nil {
		log.Printf("Property '%s' was not found in context! %v", key, err)
	}
	return value
}

func get(key string, ctx simple_fsm.ContextOperator) interface{} {
	value, err := ctx.Raw(key)
	if err != nil {
		log.Printf("Object '%s' wasn't found in context: %v", key, err)
		return nil
	}
	return value
}

func findBy(ctx simple_fsm.ContextOperator) []string {
	byObj := get("by", ctx)
	by := make([]string, len(byObj.([]interface{})))
	for _, v := range byObj.([]interface{}) {
		by = append(by, v.(string))
	}
	log.Printf("By array %v", by)
	return by
}
/*
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
				log.Println("'user_not_found' state acquired")
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
*/


/*
{
    "_id" : ObjectId("59aee26d6ec32f1174db2ba4"),
    "name" : "auth",
    "parameters" : [
        {
            "name" : "login",
            "required" : true,
            "type" : "string"
        },
        {
            "name" : "pass",
            "required" : true,
            "type" : "string"
        }
    ],
    "fsm" : {
        "states" : {
            "start" : {
                "start" : true,
                "transitions" : {
                    "start-find_user" : {
                        "to" : "find_user",
                        "guard" : {
                            "type" : "always"
                        },
                        "action" : {
                            "name" : "list",
                            "params" : {
                                "target" : "Users",
                                "type" : "User",
                                "by" : [
                                    "login"
                                ],
                                "limit" : 1
                            }
                        }
                    }
                }
            },
            "find_user" : {
                "parent" : "start",
                "transitions" : {
                    "find_user-user_found" : {
                        "to" : "user_found",
                        "guard" : {
                            "type" : "context",
                            "key" : "exists",
                            "value" : true
                        },
                        "action" : {
                            "name" : "no_action"
                        }
                    },
                    "find_user-user_not_found" : {
                        "to" : "user_not_found",
                        "guard" : {
                            "type" : "context",
                            "key" : "failed",
                            "value" : true
                        },
                        "action" : {
                            "name" : "set_result"
                        }
                    }
                }
            },
            "user_found" : {
                "parent" : "find_user",
                "transitions" : {
                    "user_found-pass_verified" : {
                        "to" : "pass_verified",
                        "guard" : {
                            "type" : "context",
                            "key" : "failed",
                            "value" : false
                        },
                        "action" : {
                            "name" : "auth"
                        }
                    }
                }
            },
            "user_not_found" : {
                "parent" : "find_user",
                "transitions" : {}
            },
            "pass_verified" : {
                "parent" : "user_found",
                "transitions" : {
                    "pass_verified-end" : {
                        "to" : "end",
                        "guard" : {
                            "type" : "context",
                            "key" : "failed",
                            "value" : false
                        },
                        "action" : {
                            "name" : "set_result"
                        }
                    },
                    "pass_verified-pass_failed" : {
                        "to" : "pass_failed",
                        "guard" : {
                            "type" : "context",
                            "key" : "failed",
                            "value" : true
                        },
                        "action" : {
                            "name" : "set_result"
                        }
                    }
                }
            },
            "pass_failed" : {
                "parent" : "pass_verified",
                "transitions" : {}
            },
            "end" : {
                "parent" : "find_user",
                "transitions" : {}
            }
        }
    }
}


 */
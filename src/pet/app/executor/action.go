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
	"time"
	"github.com/satori/go.uuid"
)

func setResult(ctx simple_fsm.ContextOperator) error {
	failed := get("failure", ctx)
	if failed != nil {
		ctx.PutResult(failed)
	} else {
		if !ctx.Has("result") {
			if !ctx.Has("response") {
				var msg string
				if name, _ := ctx.Str("methodName"); name != "" {
					msg = fmt.Sprintf("Flow '%s' structure, wasn't formed correctly! No result was set for this flow!", name)
				} else {
					msg = "Current flow structure, wasn't formed correctly! No result was set for this flow!"
				}
				ctx.PutResult(Result{
					Status:http.StatusInternalServerError,
					Data: msg,
				})
				log.Println("Result wasn't set by any action!!!")
			} else {
				responseMap := get("response", ctx).(map[string]interface{})
				ctx.PutResult(processData(responseMap, ctx))
				/*
				"response" : {
                                    "code" : 200,
                                    "data" : [
                                        {
                                            "name" : "value",
                                            "value" : "value",
                                            "from" : "entity"
                                        },
                                        {
                                            "name" : "expiration",
                                            "value" : "expiration",
                                            "from" : "entity"
                                        }
                                    ]
                                }
				 */


			}
		}
		// else already set!
	}
	return nil
}

func noAction(ctx simple_fsm.ContextOperator) error {
	//do nothing
	return nil
}
/*
	"set" : {
	    "override" : false,
	    "token_end_date" : {
		"type" : "time",
		"operation" : "add",
		"value" : 500,
		"units" : "seconds"
	    },
	    "uuid" :{
		"type" : "uuid",
		"value" : "new"
	    }
	}
 */
//TODO: strong testing needed
func setToContext(ctx simple_fsm.ContextOperator) error {
	//do nothing
	log.Println("Setting params to context")
	setToContext := get("set", ctx).(map[string]interface{})
	override, _ := setToContext["override"].(bool)
	log.Printf("Set section:\n%v", setToContext)
	for k, v := range setToContext {
		if override {
			if (k == "override"){
				continue
			}
//			log.Printf("Key from set section:\n%v", k)
//			log.Printf("Value from set section:\n%v", v)
			obj := getAsMap(v)
			objType, present := obj["type"]
//			log.Printf("Type value for object present?: %v", present)
//			log.Printf("Type value for object is: %v", objType)
			//log.Printf("Type value for object as string: %v", objType.(string))
			if present {
//				log.Printf("Type specified!\n%v", setToContext)
				switch objType.(string) {
					case "time":
//						log.Printf("Saving TIME:\n%v - %v", k, createTime(obj))
						ctx.Put(k, createTime(obj))
					case "uuid":
						generated :=createUUID(obj)
//						log.Printf("Saving uuid:\n%v - %v", k, generated)
						ctx.Put(k, generated)
					default:
						ctx.Put(k, v)
				}
			} else {
//				log.Printf("Type is not complex type so it will be set 'as is':\n%v = %v", k, v)
				ctx.Put(k, v)
			}
		} else {
			//TODO: implement proper flow for not override
			if !ctx.Has(k) {
				ctx.Put(k, v)
			}
		}
	}
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
		ctx.Put("exists", false)
		ctx.Put("failed", true)
		ctx.Put("failure", Result{
			Status: http.StatusInternalServerError,
			Data: msg,
		})
	}
	if len(entity) > 0 {
		ent1 := entity[0].(bson.M)
//		log.Printf("Entity 1: %v", ent1)
		if ent1["_id"] != nil || ent1["_id"] != "" {
			//log.Println("Entity 1 exists!")
			ctx.Put("exists", true)
		} else {
			//log.Println("Entity 1 doesn't exist!")
			ctx.Put("exists", false)
		}
	} else {
		ctx.Put("exists", false)
	}
	if ctx.Has("entity") {
		previous := get("entity", ctx)
		log.Printf("Found last entity. It will be set as 'last_entity' set to context\n%v", previous)
		ctx.Put("last_entity", previous)
	}

	//log.Printf("Entity set to context\n%v", entity)
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
	log.Printf("Creating new entity in %s", target)
	insertQuery := find(ctx)
	entity, er := model.CreateEntity(target, insertQuery)
	if er != nil {
		log.Printf("Error during db request: %v", er)
		ctx.Put("failed", true)
		ctx.Put("exists", false)
		msg := fmt.Sprintf("Can not insert entity with query: '%v'. Message: %s", insertQuery, er)
		ctx.Put("failure", Result{
			Status: http.StatusInternalServerError,
			Data: msg,
		})
		return nil
	}

	//TODO: handle existence
	/*if len(entity) > 0 {
		ctx.Put("exists", true)
	} else {
		ctx.Put("exists", false)
	}*/
	ctx.Put("exists", true)
	if ctx.Has("entity") {
		previous := get("entity", ctx)
		log.Printf("Found last entity. It will be set as 'last_entity' set to context\n%v", previous)
		ctx.Put("last_entity", previous)
	}
	log.Printf("Entity set to context\n%v", entity)
	ctx.Put("entity", entity)
	//TODO: remove tmp
	/*
	ctx.PutResult(Result{
		Status: http.StatusOK,
		Data: entity,
	})
	*/
	return nil
}

func update(ctx simple_fsm.ContextOperator) error {
	target := getStr("target", ctx)
	log.Printf("Updating new entity in %s", target)
	updateQuery := findForUpdate(get("update_values", ctx).(map[string]interface{}), ctx)
	findQuery := findForUpdate(get("find", ctx).(map[string]interface{}), ctx)
	entity, er := model.UpdateEntity(target, findQuery, updateQuery, getResultObject("type", 0))
	if er != nil {
		log.Printf("Error during db request: %v", er)
		ctx.Put("failed", true)
		ctx.Put("exists", false)
		msg := fmt.Sprintf("Can not insert entity with query: '%v'. Message: %s", updateQuery, er)
		ctx.Put("failure", Result{
			Status: http.StatusInternalServerError,
			Data: msg,
		})
		return nil
	}

	//TODO: handle existence
	if len(entity) > 0 {
		ent1 := entity[0].(bson.M)
		//		log.Printf("Entity 1: %v", ent1)
		if ent1["_id"] != nil || ent1["_id"] != "" {
			//log.Println("Entity 1 exists!")
			ctx.Put("exists", true)
		} else {
			//log.Println("Entity 1 doesn't exist!")
			ctx.Put("exists", false)
		}
	} else {
		ctx.Put("exists", false)
	}
	if ctx.Has("entity") {
		previous := get("entity", ctx)
		log.Printf("Found last entity. It will be set as 'last_entity' set to context\n%v", previous)
		ctx.Put("last_entity", previous)
	}

	//log.Printf("Entity set to context\n%v", entity)
	ctx.Put("entity", entity)
	return nil
}

func authorize(ctx simple_fsm.ContextOperator) error {
	rawUser, _ := ctx.Raw("entity")
	pass, _ := ctx.Str("pass")
	//log.Printf("User: %v", rawUser)
	usrs := []model.User{}
	usersBytes, _ := json.Marshal(rawUser)
	er := json.Unmarshal(usersBytes, &usrs)
	if er != nil {
		log.Printf("Unmarshalled err: %v", er)
	}
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
	ctx.Put("failed", false)
	/*//TODO: remove this synthetic result
	ctx.PutResult(Result{
		Status: http.StatusOK,
		Data: usrs,
	})*/

	return nil
}

func find(ctx simple_fsm.ContextOperator) bson.M {
	by := findBy(ctx)
	if ctx.Has("where") {
//		log.Println("Finding with 'where' conditions!")
		return withWhere(by, ctx)
	} else {
//		log.Println("Simple find, no 'where' were found")
		return findRequest(by, ctx)
	}
}

func findForUpdate(m map[string]interface{}, ctx simple_fsm.ContextOperator) bson.M {
	result := make(map[string]interface{})

	byArr := m["by"].([]interface{})
	by := make([]string, len(byArr))
	for _, v := range byArr {
		by = append(by, v.(string))
	}

	whereSection := m["where"].([]interface{})
	log.Printf("Where section: %v", whereSection)

	// match each 'where' object with 'by'
	for _, w := range whereSection {
		whereObj := newWhereCondition(w.(map[string]interface{}))
		for _, byEl := range by {
			// when 'where.name' == 'by'[i]
			if byEl != "" && whereObj.Name == byEl {
				result[byEl] = from(whereObj, ctx)
				break;
			}
		}

	}
	return bson.M(result)
}

// findRequest
// Matches request arguments with params from Method.fsm then creates search MongoDB search request as json
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

func findBy(ctx simple_fsm.ContextOperator) []string {
	byObj := get("by", ctx)
	by := make([]string, len(byObj.([]interface{})))
	for _, v := range byObj.([]interface{}) {
		by = append(by, v.(string))
	}
//	log.Printf("By array %v", by)
	return by
}

// matchWithWhere
// Matches values for "by" values
func withWhere(by []string, ctx simple_fsm.ContextOperator) bson.M {
	result := make(map[string]interface{})
	whereSection := get("where", ctx)
	//log.Printf("Where section: %v", whereSection)

	// match each 'where' object with 'by'
	for _, w := range whereSection.([]interface{}) {
		whereObj := newWhereCondition(w.(map[string]interface{}))
		for _, byEl := range by {
			// when 'where.name' == 'by'[i]
			if byEl != "" && whereObj.Name == byEl {
				result[byEl] = from(whereObj, ctx)
				break;
			}
		}

	}
	//log.Printf("Where matched objects: %v", result)

	// fill other params that not match to 'where' condition
	for _, arg := range by {
		if _, ok := result[arg]; arg != "" && !ok {
			value, err := ctx.Raw(arg)
			if err != nil {
				log.Printf("Error during mathing 'where'. Message: %v", err)
			} else {
				result[arg] = value
			}
		}
	}
	//log.Printf("Match where result: %v", result)
	return bson.M(result)
}

type where struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	From  string `json:"from"`
}

func processData(responseMap map[string]interface{}, ctx simple_fsm.ContextOperator) Result {
	var result Result
	dataMap := make(map[string]interface{})

	if code, present :=responseMap["code"]; present{
		result.Status = code.(int)
	}
	dataArr := responseMap["data"].([]interface{})

	for _,value := range dataArr {
		whereCond := newWhereCondition(value.(map[string]interface{}))
		dataMap[whereCond.Name] = from(whereCond, ctx)
	}
	result.Data = dataMap
	return result
}

func newWhereCondition(m map[string]interface{}) (where) {
	var res where
	if value, present := m["name"]; present{
		res.Name = value.(string)
	}
	if value, present := m["value"]; present{
		res.Value = value.(string)
	}
	if value, present := m["from"]; present{
		res.From = value.(string)
	}
	return res
}

func from(whereObj where, ctx simple_fsm.ContextOperator) interface{} {
	//log.Println("From analyze started")

	if whereObj.From != "" {
		if whereObj.From == "context" {
			//log.Printf("Taking value from context!\n%v\n Value: %v", whereObj, get(whereObj.Value, ctx))
			//log.Printf("\n Name: %v, Value: %v", whereObj.Name, get(whereObj.Name, ctx))
			//log.Printf("\n Name: %v, Value: %v", whereObj.Name, get(whereObj.Value, ctx))
			//log.Printf("\n From: %v, Value: %v", whereObj.From, get(whereObj.From, ctx))

			// when "context" specified we will take value from root ctx
			return get(whereObj.Value, ctx)
		}
		// when we will take value from specified "from" object in "context"
		rawObj, _ := ctx.Raw(whereObj.From)
		//log.Printf("From Found:\n%v ", rawObj)
		fromParent := getAsMap(rawObj)
		if fromParent != nil {
			//log.Printf("From fromParent is not empty %v", fromParent)
			return fromParent[whereObj.Value]
		}
		//log.Println("From from parent wasn't found " + whereObj.From)
		return nil
	} else {
		// mo from specified -> we will set value as is
		log.Printf("Variable '%s' wasn't found in context, so taking it from global context", whereObj.From)
		return whereObj.Value
	}
}

//TODO: stupidity but it works
func getAsMap(inputObj interface{}) map[string]interface{} {
//	log.Println("getAsMap started")
	var result map[string]interface{}

	converted := getSingleObject(inputObj)
	//log.Printf("getAsMap converted object is:\n%v", converted)
	usersBytes, err := json.Marshal(converted)
	if err != nil {
		log.Printf("Failed to marshal object '%v'.\nMessage: %s", converted, err.Error())
		return nil
	}
	err = json.Unmarshal(usersBytes, &result)
	if err != nil {
		//FIXME: failed to unmarshall 'false'
		log.Printf("Failed to unmarshal object '%v'.\nMessage: %s", converted, err.Error())
		return nil
	}
	//log.Printf("getAsMap result %v", result)
	return result
}

func getSingleObject(inputObj interface{}) interface{} {
	var converted interface{}
	switch obj := inputObj.(type) {
	case []interface{}:
		converted = obj[0]
	case interface{}:
		converted = obj
	default:
		fmt.Printf("Unsupported type: %T\n", converted)
	}
	return converted
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

/*
//TODO: specify start date somehow "start" : "now"
	"token_end_date" : {
		"type" : "time",
		"operation" : "add",
		"value" : 500,
		"units" : "seconds"
	}
 */
func createTime(timeObj map[string]interface{}) interface{} {
	startTime := time.Now()
	var duration time.Duration
	var unitsCount float64
	if value, present := timeObj["units"]; present {
		switch value.(string) {
		case "seconds":
			duration = time.Second
		case "minutes":
			duration = time.Minute
		}
	}

	if value, present := timeObj["value"]; present {
		unitsCount = value.(float64)
	}

	if value, present := timeObj["operation"]; present {
		switch value.(string) {
		case "add":
			return startTime.Add(duration * time.Duration(unitsCount))
		case "subtract":
		//return startTime.Sub(duration * unitsCount)
		}
	}
	return startTime
}

/*
	"uuid" :{
		"type" : "uuid",
		"value" : "new"
	}
 */

func createUUID(timeObj map[string]interface{})  interface{} {
	log.Println("Creating new uuid!")
	if val, present := timeObj["value"]; present && val != "new"{
		return val
	}

	return uuid.NewV4().String()
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
                            "key" : "exists",
                            "value" : false
                        },
                        "action" : {
                            "name" : "set_result",
                            "params" : {
                                "response" : {
                                    "code" : 403,
                                    "data" : [
                                        {
                                            "name" : "message",
                                            "value" : "Provided user wasn't found"
                                        }
                                    ]
                                }
                            }
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
                    "pass_verified-find_token" : {
                        "to" : "find_token",
                        "guard" : {
                            "type" : "context",
                            "key" : "failed",
                            "value" : false
                        },
                        "action" : {
                            "name" : "list",
                            "params" : {
                                "target" : "Token",
                                "type" : "Token",
                                "by" : [
                                    "userId"
                                ],
                                "where" : [
                                    {
                                        "name" : "userId",
                                        "value" : "_id",
                                        "from" : "entity"
                                    }
                                ],
                                "limit" : 1
                            }
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
            "find_token" : {
                "parent" : "pass_verified",
                "transitions" : {
                    "find_token-token_found" : {
                        "to" : "token_found",
                        "guard" : {
                            "type" : "context",
                            "key" : "exists",
                            "value" : true
                        },
                        "action" : {
                            "name" : "set_to_context",
                            "params" : {
                                "set" : {
                                    "override" : true,
                                    "token_end_date" : {
                                        "type" : "time",
                                        "units" : "seconds",
                                        "operation" : "add",
                                        "value" : 500
                                    },
                                    "uuid" : {
                                        "type" : "uuid",
                                        "value" : "new"
                                    }
                                }
                            }
                        }
                    },
                    "find_token-token_not_found" : {
                        "to" : "token_not_found",
                        "guard" : {
                            "type" : "context",
                            "key" : "exists",
                            "value" : false
                        },
                        "action" : {
                            "name" : "set_to_context",
                            "params" : {
                                "set" : {
                                    "override" : true,
                                    "token_end_date" : {
                                        "type" : "time",
                                        "units" : "seconds",
                                        "operation" : "add",
                                        "value" : 500
                                    },
                                    "uuid" : {
                                        "type" : "uuid",
                                        "value" : "new"
                                    }
                                }
                            }
                        }
                    }
                }
            },
            "token_found" : {
                "parent" : "find_token",
                "transitions" : {
                    "token_found-update_token" : {
                        "to" : "update_token",
                        "guard" : {
                            "type" : "context",
                            "key" : "failed",
                            "value" : false
                        },
                        "action" : {
                            "name" : "update",
                            "params" : {
                                "target" : "Token",
                                "type" : "Token",
                                "find" : {
                                    "by" : [
                                        "userId"
                                    ],
                                    "where" : [
                                        {
                                            "name" : "userId",
                                            "value" : "_id",
                                            "from" : "last_entity"
                                        }
                                    ]
                                },
                                "update_values" : {
                                    "by" : [
                                        "expiration",
                                        "value"
                                    ],
                                    "where" : [
                                        {
                                            "name" : "expiration",
                                            "value" : "token_end_date",
                                            "from" : "context"
                                        },
                                        {
                                            "name" : "value",
                                            "value" : "uuid",
                                            "from" : "context"
                                        }
                                    ]
                                }
                            }
                        }
                    }
                }
            },
            "token_not_found" : {
                "parent" : "find_token",
                "transitions" : {
                    "token_not_found-create_token" : {
                        "to" : "create_token",
                        "guard" : {
                            "type" : "context",
                            "key" : "failed",
                            "value" : false
                        },
                        "action" : {
                            "name" : "create",
                            "params" : {
                                "target" : "Token",
                                "type" : "Token",
                                "by" : [
                                    "userId",
                                    "expiration",
                                    "value"
                                ],
                                "where" : [
                                    {
                                        "name" : "userId",
                                        "value" : "_id",
                                        "from" : "last_entity"
                                    },
                                    {
                                        "name" : "expiration",
                                        "value" : "token_end_date",
                                        "from" : "context"
                                    },
                                    {
                                        "name" : "value",
                                        "value" : "uuid",
                                        "from" : "context"
                                    }
                                ]
                            }
                        }
                    }
                }
            },
            "create_token" : {
                "parent" : "token_not_found",
                "transitions" : {
                    "create_token-token_created" : {
                        "to" : "token_created",
                        "guard" : {
                            "type" : "context",
                            "key" : "failed",
                            "value" : false
                        },
                        "action" : {
                            "name" : "set_result",
                            "params" : {
                                "response" : {
                                    "code" : 200,
                                    "data" : [
                                        {
                                            "name" : "value",
                                            "value" : "value",
                                            "from" : "entity"
                                        },
                                        {
                                            "name" : "expiration",
                                            "value" : "expiration",
                                            "from" : "entity"
                                        }
                                    ]
                                }
                            }
                        }
                    }
                }
            },
            "update_token" : {
                "parent" : "token_found",
                "transitions" : {
                    "update_token-token_updated" : {
                        "to" : "token_updated",
                        "guard" : {
                            "type" : "context",
                            "key" : "failed",
                            "value" : false
                        },
                        "action" : {
                            "name" : "set_result",
                            "params" : {
                                "response" : {
                                    "code" : 200,
                                    "data" : [
                                        {
                                            "name" : "value",
                                            "value" : "uuid",
                                            "from" : "context"
                                        },
                                        {
                                            "name" : "expiration",
                                            "value" : "token_end_date",
                                            "from" : "context"
                                        }
                                    ]
                                }
                            }
                        }
                    }
                }
            },
            "token_updated" : {
                "parent" : "update_token",
                "transitions" : {}
            },
            "token_created" : {
                "parent" : "create_token",
                "transitions" : {}
            },
            "pass_failed" : {
                "parent" : "pass_verified",
                "transitions" : {}
            }
        }
    }
}
 */
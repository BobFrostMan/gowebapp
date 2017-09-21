package executor

import (
	"log"
	"pet/app/model"
	"github.com/xenzh/gofsm"
	"fmt"
	"pet/app/shared/passhash"
	"encoding/json"
)

// context keys
const (
	failed = "failed"
	failure = "failure"
	result_key = "result"
	response_key = "response"
	entity_key = "entity"
)
// context action params keys
const (
	override_key = "override"
	target_key = "target"
	limit_key = "limit"
	save_as_key = "save_as"
	find_key = "find"
	update_values_key = "update_values"
	set_key = "set"
)

// list
// Action that stands for mongoDB find query. Result will be set to context as "entity" (slice of interfaces)
// Also set "exists", "failed", and "failure" values to context (depends from request execution result)
// For proper usage, context should contain next values:
// - target - string, Database collection name
// - limit - int, limit of documents obtained,
// - fields - []string, list of fields to query from db
// - where - []where, objects to resolve 'fields' values in find query
func list(ctx simple_fsm.ContextOperator) error {
	target := getStr(target_key, ctx)
	limit := getInt(limit_key, ctx)
	findQuery := findBson(ctx)
	entity, er := model.List(target, findQuery, limit, make([]interface{}, limit))
	if er != nil {
		log.Printf("Error during db request: %v", er)
		msg := fmt.Sprintf("Can not find entity by query: '%v'. Message: %s", findQuery, er)
		setFailureToContext(msg, ctx)
	}
	setExists(entity, ctx)
	saveEntityAs(entity, ctx)
	setEntity(entity, ctx)
	return nil
}

// create
// Action that stands for mongoDB insert query. Result will be set to context as "entity" (slice of interfaces)
// Also set "exists", "failed", and "failure" values to context (depends from request execution result)
// For proper usage, context should contain next values:
// - target - string, Database collection name
// - limit - int, limit of documents obtained,
// - fields - []string, list of fields that will be saved to database document
// - where - []where, objects to resolve 'fields' values
func create(ctx simple_fsm.ContextOperator) error {
	target := getStr(target_key, ctx)
	insertQuery := findBson(ctx)
	entity, er := model.CreateEntity(target, insertQuery)
	if er != nil {
		msg := fmt.Sprintf("Can not insert entity with query: '%v'. Message: %s", insertQuery, er)
		setFailureToContext(msg, ctx)
		return nil
	}
	res := make([]interface{}, 1)
	res[0] = entity
	setExists(res, ctx)
	saveEntityAs(res, ctx)
	setEntity(res, ctx)
	return nil
}

// update
// Action that stands for find and modify action in terms of mongoDB.
// Set modified objects to context as "entity" (slice of interfaces)
// Also set "exists", "failed", and "failure" values to context (depends from request execution result)
// For proper usage, context should contain next values:
// - target - string, Database collection name
// - limit - int, limit of documents obtained,
// - find - object that contains 'fields' and 'where', will be used to build find query
// - update_values - object that contains 'fields' and 'where', will be used to build modify query
// Both 'find' and 'update_values' should have:
// - fields - []string, list of fields that will be saved to database document
// - where - []where, objects to resolve 'fields' values
func update(ctx simple_fsm.ContextOperator) error {
	target := getStr(target_key, ctx)
	log.Printf("Updating new entity in %s", target)
	updateQuery := updateBson(get(update_values_key, ctx).(map[string]interface{}), ctx)
	findQuery := updateBson(get(find_key, ctx).(map[string]interface{}), ctx)
	entity, er := model.UpdateEntity(target, findQuery, updateQuery, make([]interface{}, 0))
	if er != nil {
		msg := fmt.Sprintf("Can not insert entity with query: '%v'. Message: %s", updateQuery, er)
		setFailureToContext(msg, ctx)
		return nil
	}

	setExists(entity, ctx)
	saveEntityAs(entity, ctx)
	setEntity(entity, ctx)
	return nil
}

// authorize
// Context should have model.User object under "entity" key
// Takes password hash from current context "entity", and compares it with "pass" value from context
// Also set "failed", and "failure" values to context base on comparison result
func authorize(ctx simple_fsm.ContextOperator) error {
	rawUser, _ := ctx.Raw(entity_key)
	pass, _ := ctx.Str("pass")
	usrs := []model.User{}
	usersBytes, _ := json.Marshal(rawUser)
	log.Printf("Marshalled: %v", rawUser)
	er := json.Unmarshal(usersBytes, &usrs)
	if er != nil {
		log.Printf("Unmarshalled err: %v", er)
	}
	log.Printf("Unmarshalled user: %v", usrs)
	userObj := usrs[0]
	if err := passhash.CompareHashAndPassword(userObj.Password, pass); err != nil {
		log.Printf("User '%s' entered wrong password!", userObj.Name)
		setFailureToContext("Wrong password specified", ctx)
		return nil
	}
	ctx.PutParent(failed, false)
	return nil
}

/*
   setToContext
    Set data to context, and override it (when override has true value).
    By default, values will be set to context 'as is'.
    Objects with "type" field will be processed in specific way based on value

   Action params json example:
   "action" : {
        "name" : "set_to_context",
        "params" : {
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
	}
    }
*/
func setToContext(ctx simple_fsm.ContextOperator) error {
	setToContext := get(set_key, ctx).(map[string]interface{})
	override, _ := setToContext[override_key].(bool)
	for k, v := range setToContext {
		if (k == override_key) {
			continue
		}
		if override {
			value := resolveValue(v, ctx)
			ctx.PutParent(k, value)
		} else {
			if !ctx.Has(k) {
				value := resolveValue(v, ctx)
				ctx.PutParent(k, value)
			}
		}
	}
	return nil
}

/*
      setResult
    Set "result" to context, based on data that already inside context.
    By default set "entity" as result.
    The customer way to set result context is to use "response" object in action
    By default, values will be set to context 'as is'.
    Objects with "type" field will be processed in specific way based on value
    Action example json:

	"action" : {
		"name" : "set_result",
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
*/
func setResult(ctx simple_fsm.ContextOperator) error {
	failed, _ := ctx.Raw(failure)
	switch {
		case ctx.Has(response_key):
			log.Printf("SET RESULT %v", get(response_key, ctx))
			responseMap := get(response_key, ctx).(map[string]interface{})
			ctx.PutResult(processData(responseMap, ctx))
		case failed != nil:
			ctx.PutResult(failed)
		case !ctx.Has(result_key):
			log.Printf("SET RESULT %v", get(entity_key, ctx))
			ctx.PutResult(get(entity_key, ctx))
		default:
			//do nothing, result already in context
	}
	return nil
}
package executor

import (
	"github.com/satori/go.uuid"
	"github.com/xenzh/gofsm"
	"log"
	"time"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"encoding/json"
)
// general context keys
const (
	exists = "exists"
	last_entity_key = "last_entity"
	fields_key = "fields"
)

// where entity keys
const (
	where_key = "where"
	from_key = "from"
	name_key = "name"
	value_key = "value"
)

// Custom struct for convenient json "where" section handling. "where" is a param of action "params" object
type where struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	From  string `json:"from"`
}

// newWhereCondition
// Creates new where condition based on given map
func newWhereCondition(m map[string]interface{}) (where) {
	var res where
	if value, present := m[name_key]; present {
		res.Name = value.(string)
	}
	if value, present := m[value_key]; present {
		res.Value = value.(string)
	}
	if value, present := m[from_key]; present {
		res.From = value.(string)
	}
	return res
}

// setExists
// Checks entity existence and set according value to context as "exists"
func setExists(entity []interface{}, ctx simple_fsm.ContextOperator){
	if len(entity) > 0 {
		exist := false
		switch entity[0].(type) {
		case bson.M:
			ent1 := entity[0].(bson.M)
			exist = ent1["_id"] != nil || ent1["_id"] != ""
		case map[string]interface{}:
			ent1 := entity[0].(map[string]interface{})
			exist = ent1["_id"] != nil || ent1["_id"] != ""
		}
		ctx.PutParent(exists, exist)
	} else {
		ctx.PutParent(exists, false)
	}
}

// setEntity
// Set given entity to context under "entity" key.
// If "entity" key already contains some data, this data will be stored under "last_entity" key
func setEntity(entity []interface{}, ctx simple_fsm.ContextOperator){
	if ctx.Has(entity_key) {
		previous := get(entity_key, ctx)
		log.Printf("Found last entity. It will be set as 'last_entity' set to context\n%v", previous)
		ctx.PutParent(last_entity_key, previous)
	}
	ctx.PutParent(entity_key, entity)
}

// saveEntityAs
// Set given entity to context under specified name key.
func saveEntityAs(entity []interface{}, ctx simple_fsm.ContextOperator){
	saveAs := getStr(save_as_key, ctx)
	log.Printf("Saved as value is: %v", saveAs)
	if saveAs != "" {
		ctx.PutParent(saveAs, entity)
		ctx.PutParent(save_as_key, "")
	}
}

// resolveValue
// Processing value from valueMap
func resolveValue(valueMap interface{}) interface{} {
	obj := asMap(valueMap)
	objType, present := obj["type"]
	if present {
		switch objType.(string) {
		case "time":
			return createTime(obj)
		case "uuid":
			return createUUID(obj)
		default:
			//do nothing
		}
	}
	return valueMap
}

// findBson
// Return BSON object, from given map. BSON object represents mongoDB find query
func findBson(ctx simple_fsm.ContextOperator) bson.M {
	fields := asStringsSlice(get(fields_key, ctx))
	if ctx.Has(where_key) {
		return withWhere(fields, ctx)
	} else {
		return findRequest(fields, ctx)
	}
}

// updateBson
// Return BSON object, from given map. BSON object represents mongoDB update query
func updateBson(m map[string]interface{}, ctx simple_fsm.ContextOperator) bson.M {
	res := make(map[string]interface{})

	fieldsObj := m[fields_key].([]interface{})
	fields := asStringsSlice(fieldsObj)

	whereSection := m[where_key].([]interface{})
	log.Printf("Where section: %v", whereSection)

	// match each 'where' object with each field from 'fields' array
	for _, w := range whereSection {
		whereObj := newWhereCondition(w.(map[string]interface{}))
		for _, field := range fields {
			// when 'where.name' == field value
			if field != "" && whereObj.Name == field {
				res[field] = handleFrom(whereObj, ctx)
				break;
			}
		}

	}
	return bson.M(res)
}

// findRequest
// Matches request arguments with params from Method.fsm then creates search MongoDB search request as json
func findRequest(args []string, ctx simple_fsm.ContextOperator) bson.M {
	res := make(map[string]interface{})
	for _, arg := range args {
		if arg != "" {
			value, err := ctx.Raw(arg)
			if err != nil {
				log.Printf("Error during form find request. Message: %v", err)
			} else {
				res[arg] = value
			}
		}
	}
	return bson.M(res)
}

// withWhere
// Matches values for "fields" values, with values specified in "where" section
func withWhere(fields []string, ctx simple_fsm.ContextOperator) bson.M {
	res := make(map[string]interface{})
	whereSection := get(where_key, ctx)

	// match each 'where' object with field from 'fields' array
	for _, w := range whereSection.([]interface{}) {
		whereObj := newWhereCondition(w.(map[string]interface{}))
		for _, field := range fields {
			// when 'where.name' == field value
			if field != "" && whereObj.Name == field {
				res[field] = handleFrom(whereObj, ctx)
				break;
			}
		}
	}

	// fill other params that not match to 'where' condition
	for _, arg := range fields {
		if _, ok := res[arg]; arg != "" && !ok {
			value, err := ctx.Raw(arg)
			if err != nil {
				log.Printf("Error during mathing 'where'. Message: %v", err)
			} else {
				res[arg] = value
			}
		}
	}
	return bson.M(res)
}

// processData
// Converts given map to executor.Result object
func processData(responseMap map[string]interface{}, ctx simple_fsm.ContextOperator) Result {
	var res Result
	dataMap := make(map[string]interface{})

	if code, present := responseMap["code"]; present {
		switch code.(type){
		case int:
			res.Status = code.(int)
		default:
			res.Status = int(code.(float64))
		}
	} else {
		res.Status = http.StatusOK
	}
	dataArr := responseMap["data"].([]interface{})
	for _, value := range dataArr {
		whereCond := newWhereCondition(value.(map[string]interface{}))
		dataMap[whereCond.Name] = handleFrom(whereCond, ctx)
	}
	res.Data = dataMap
	return res
}

// handleFrom
// Resolve end value by rules of "from" json session. "from" is a field of "where" section, from action "params"
func handleFrom(whereObj where, ctx simple_fsm.ContextOperator) interface{} {
	if whereObj.From != "" {
		if whereObj.From == "context" {
			// when "context" specified we will take value from root ctx
			return get(whereObj.Value, ctx)
		}
		// when we will take value from specified "from" object in "context"
		rawObj, _ := ctx.Raw(whereObj.From)
		fromParent := asMap(rawObj)
		if fromParent != nil {
			return fromParent[whereObj.Value]
		}
		return nil
	} else {
		// mo from specified -> we will set value as is
		log.Printf("Variable '%s' wasn't found in context, so taking it from global context", whereObj.From)
		return whereObj.Value
	}
}

// asSlice
// Returns obj as slice of strings, if type checking passed
func asStringsSlice(obj interface{}) []string {
	fields := make([]string, len(obj.([]interface{})))
	for _, v := range obj.([]interface{}) {
		fields = append(fields, v.(string))
	}
	return fields
}

// asMap
// Returns interface as map of interface values accessed by string keys
func asMap(inputObj interface{}) map[string]interface{} {
	var res map[string]interface{}

	converted := getSingleObject(inputObj)
	usersBytes, err := json.Marshal(converted)
	if err != nil {
		log.Printf("Failed to marshal object '%v'.\nMessage: %s", converted, err.Error())
		return nil
	}
	err = json.Unmarshal(usersBytes, &res)
	if err != nil {
		log.Printf("Failed to unmarshal object '%v'.\nMessage: %s", converted, err.Error())
		return nil
	}
	return res
}

// getSingleObject
// Returns first object from given collection or object itself, if it can be taken as interface{}
func getSingleObject(inputObj interface{}) interface{} {
	var converted interface{}
	switch obj := inputObj.(type) {
	case []interface{}:
		converted = obj[0]
	case interface{}:
		converted = obj
	default:
		log.Printf("Unsupported type: %T\n", converted)
	}
	return converted
}

// setFailureToContext
// Set flags 'exists' and 'failed'. Also set failure object to context by 'failure' key
func setFailureToContext(msg string, ctx simple_fsm.ContextOperator) {
	ctx.PutParent(exists, false)
	ctx.PutParent(failed, true)
	m := make(map[string]interface{})
	m["message"] = msg
	ctx.PutParent(failure, NewResult(http.StatusInternalServerError, m))
}

// createTime
// Creates new time object base on given map.
// Processed values:
// - "operation" operation to perform with start date (time.Now()). Supported values: ["add"]
// - "units" string value. Supported values: ["seconds", "minutes"]
// - "value" any valid integer value. Units value that will be arithmetically processed
// If no value processed returns time.Now()
func createTime(timeObj map[string]interface{}) interface{} {
	//TODO: specify start date somehow "start" : "now"
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

// createUUID
// Creates new uuid or returns existing from uuidJsonMap (stored by key "value")
// If value is "new" creates new uuid object in string representation
func createUUID(uuidJsonMap map[string]interface{}) interface{} {
	if val, present := uuidJsonMap["value"]; present && val != "new" {
		return val
	}
	return uuid.NewV4().String()
}

// getStr
// wrapper for ctx.Str prints error
func getStr(key string, ctx simple_fsm.ContextOperator) string {
	value, err := ctx.Str(key)
	if err != nil {
		log.Printf("Property '%s' was not found in context! %v", key, err)
	}
	return value
}

// getInt
// wrapper for ctx.Int prints error
func getInt(key string, ctx simple_fsm.ContextOperator) int {
	value, err := ctx.Int(key)
	if err != nil {
		log.Printf("Property '%s' was not found in context! %v", key, err)
	}
	return value
}

// get
// wrapper for ctx.Raw prints error, returns nil on error
func get(key string, ctx simple_fsm.ContextOperator) interface{} {
	value, err := ctx.Raw(key)
	if err != nil {
		log.Printf("Object '%s' wasn't found in context: %v", key, err)
		return nil
	}
	return value
}

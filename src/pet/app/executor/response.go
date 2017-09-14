package executor

import (
	"log"
	"fmt"
)

type Result struct {
	Status int `json:"status"`
	Data   interface{} `json:"data"`
}

func NewResultFrom(inputObj interface{}) Result {
	var result Result
	switch obj := inputObj.(type){
		case []interface{}:
			result.Status = 200
			result.Data = interface{}(inputObj)
		case interface{}:
			resultData := asMap(inputObj)
			if status, ok := resultData["status"]; ok {
				switch status.(type) {
				case float64:
					var val float64 = status.(float64)
					result.Status = int(val)
				default:
					result.Status = status.(int)
				}
			} else {
				result.Status = 200
			}

			if value, ok := resultData["data"]; ok {
				result.Data = value
			} else {
				result.Data = resultData
			}
			log.Printf("data:\n%s", resultData)
		default:
			fmt.Printf("Unsupported type: %T\n", obj)
	}

	log.Printf("Result data:\n%s", result.Data)
	return result
}

func NewResultMessage(status int, msg string) Result {
	return Result{
		Status: status,
		Data:msg,
	}
}

func NewResult(status int, resultData interface{}) Result {
	return Result{
		Status: status,
		Data: resultData,
	}
}

package context

import (
	"pet/app/model"
	"github.com/xenzh/gofsm"
)

const methods = "methods"

type AppContext struct {
	ctx map[string]interface{}
}

func (c *AppContext) InitContext() {
	c.ctx = make(map[string]interface{})
	c.ctx[methods] = make(map[string]model.Method)
}

func (c *AppContext) Put(key string, value interface{}) {
	met, ok := value.(model.Method)
	if !ok {
		c.ctx[key] = value
	} else {
		methodMap := c.ctx[methods].(map[string]model.Method)
		methodMap[met.Name] = met
		c.ctx[methods] = methodMap
	}
}

func (c *AppContext) Get(key string) (interface{}) {
	return &c.ctx[key]
}

func (c *AppContext) GetFsm(key string) (simple_fsm.Fsm) {
	return c.ctx[key].(simple_fsm.Fsm)
}

func (c *AppContext) GetMethod(name string) (model.Method) {
	for methodName, method := range c.ctx[methods].(map[string]model.Method){
		if (name == methodName){
			return method
		}
	}
	return nil
}

func (c *AppContext) GetString(key string) (string) {
	return c.ctx[key].(string)
}

func (c *AppContext) GetBool(key string) (bool) {
	return c.ctx[key].(bool)
}

func (c *AppContext) GetInt(key string) (int) {
	return c.ctx[key].(int)
}

func (c *AppContext) GetAllMethodsMap (map[string]model.Method) {
	return c.ctx[methods]
}

func (c *AppContext) GetAllMethods() ([]model.Method) {
	allMethods := make([]model.Method, 10)
	for _, method := range c.ctx[methods].(map[string]model.Method){
		allMethods = append(allMethods, method)
	}
	return allMethods
}



package executor

type Request struct {
	Name   string
	Token  string
	Params map[string]interface{}
}


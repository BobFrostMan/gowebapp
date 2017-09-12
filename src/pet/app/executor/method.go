package executor

import "github.com/xenzh/gofsm"

//TODO: implement internal executable representation on api method entity


type ApiFunction struct{
	Name string
	Fsm simple_fsm.JsonRoot
}
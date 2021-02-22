package goforit

type BuiltInFunctions interface {
	Execute(funcName string, vm VM, funcDef interface{}) (interface{}, bool)
}

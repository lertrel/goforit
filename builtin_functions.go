package goforit

//BuiltInFunctions an interface provided as an extenstion point
//for client to provide their own built-in functions
type BuiltInFunctions interface {

	//Execute to execute a built-in function as per the given function name
	Execute(funcName string, vm VM, funcDef interface{}) (interface{}, bool)

	//Has to check if the given function name is supported
	//by the current BuiltInFunctions
	Has(funcName string) bool
}

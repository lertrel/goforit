package goforit

//CustomFunctionRepository repository inside Formula
//to maintain custom functions
type CustomFunctionRepository interface {

	//RegisterCustomFunction for registering custom function
	//Ex.
	// r.RegisterCustomFunction(
	// 	"$CIRCLE",
	// 	`
	// 	function $CIRCLE(radius) {
	// 		return $RND(Math.PI * Math.pow(radius, 2), 10);
	// 	}
	// 	`)
	RegisterCustomFunction(funcName string, body string) bool

	//GetCustomFuncBody to get custom function source code
	GetCustomFuncBody(funcName string) string
}

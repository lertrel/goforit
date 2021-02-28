package vm

//CustomFunctionRepository repository inside Formula
//to maintain custom functions
type CustomFunctionRepository interface {

	//RegisterFunction for registering custom function
	//Ex.
	//
	// 		r.RegisterFunction(
	// 			"$CIRCLE",
	// 			`
	// 			function $CIRCLE(radius) {
	// 				return $RND(Math.PI * Math.pow(radius, 2), 10);
	// 			}
	// 			`)
	//
	RegisterFunction(funcName string, body string) bool

	//GetFunctionBody to get custom function source code
	GetFunctionBody(funcName string) string
}

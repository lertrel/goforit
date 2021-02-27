package vm

//DefaultCustomFunctionRepository default implementation CustomFunctionRepository
type DefaultCustomFunctionRepository struct {
	customFuncs map[string]string
}

//NewCustomFunctionRepo is a public function to obtain a deafult
//implementation of CustomFunctionRepository
func NewCustomFunctionRepo() CustomFunctionRepository {
	return DefaultCustomFunctionRepository{customFuncs: make(map[string]string)}
}

//RegisterCustomFunction for registering custom function
//Ex.
//
// 		r.RegisterCustomFunction(
// 			"$CIRCLE",
// 			`
// 			function $CIRCLE(radius) {
// 				return $RND(Math.PI * Math.pow(radius, 2), 10);
// 			}
// 			`)
//
func (r DefaultCustomFunctionRepository) RegisterCustomFunction(funcName string, body string) bool {

	_, found := r.customFuncs[funcName]

	r.customFuncs[funcName] = body

	return found
}

//GetCustomFuncBody to get custom function source code
func (r DefaultCustomFunctionRepository) GetCustomFuncBody(funcName string) string {

	body, found := r.customFuncs[funcName]

	if found {
		return body
	}

	//Falls through
	return ""
}

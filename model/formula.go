package model

//Formula a formula engine for creating Formulacontext
type Formula interface {

	//RegisterCustomFunction for registering custom function
	//Ex.
	//
	// 		f := goforit.NewFormularBuilder().Get()
	// 		f.RegisterCustomFunction(
	// 			"$CIRCLE",
	// 			`
	// 			function $CIRCLE(radius) {
	// 				return $RND(Math.PI * Math.pow(radius, 2), 10);
	// 			}
	// 		`)
	//
	RegisterCustomFunction(funcName string, body string) bool

	//GetCustomFunctionBody to get custom function source code
	GetCustomFunctionBody(funcName string) string

	//NewContext is a method for creating a new FormulaContext
	NewContext(script string) (c FormulaContext, err error)
}

package goforit

import (
	"regexp"

	"github.com/robertkrimen/otto"
)

//Formula a formula engine for creating Formulacontext
type Formula struct {
	r           *regexp.Regexp
	customFuncs map[string]string
	Debug       bool
}

func (f Formula) debug(format string, args ...interface{}) {
	debug(debugFlag || f.Debug, format, args...)
}

//RegisterCustomFunction for registering custom function
//Ex.
// f.RegisterCustomFunction(
// 	"$CIRCLE",
// 	`
// 	function $CIRCLE(radius) {
// 		return $RND(Math.PI * Math.pow(radius, 2), 10);
// 	}
// 	`)
func (f *Formula) RegisterCustomFunction(funcName string, body string) bool {

	_, found := f.customFuncs[funcName]

	f.customFuncs[funcName] = body

	return found
}

//GetCustomFuncBody to get custom function source code
func (f Formula) GetCustomFuncBody(funcName string) string {

	body, found := f.customFuncs[funcName]

	if found {
		return body
	}

	//Falls through
	return ""
}

func (f Formula) extractFunctionListFromFormulaString(formulaStr string) []string {

	matches := f.r.FindAllStringSubmatch(formulaStr, -1)
	dedupMatches := make(map[string]bool)

	for i := 0; i < len(matches); i++ {
		dedupMatches[matches[i][1]] = true
	}

	funArr := make([]string, len(dedupMatches))

	i := 0
	for k := range dedupMatches {
		funArr[i] = k
		i++
	}

	return funArr
}

//LoadContext If context is nil then create a new FormulaContext
//Then preparing a newly created context or a given context
//By loading referred functions (both built-in & custom) into context
func (f Formula) LoadContext(context *FormulaContext, formulaStr string) (c *FormulaContext, err error) {

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		c = nil
	// 		err = r.(error)
	// 	}
	// }()

	if context == nil {
		// context = &	FormulaContext{vm: otto.New()}
		context = &FormulaContext{vm: otto.New(), loadedFuncs: make(map[string]bool), Debug: f.Debug}
	}

	f.debug("Formula.LoadContext() - Extracting function names from %v", formulaStr)
	funcList := f.extractFunctionListFromFormulaString(formulaStr)

	for i := 0; i < len(funcList); i++ {

		f.debug("Formula.LoadContext() - funcList[i]=%v", funcList[i])
		err := f.injectFuncToContext(context, funcList[i])
		if err != nil {
			return nil, err
		}
	}

	return context, nil
}

func (f Formula) injectFuncToContext(context *FormulaContext, funcName string) error {

	f.debug("Formula.injectFuncToContext() started ...")

	if context.isFuncLoaded(funcName) {
		return nil
	}

	if fn := getBuilinFunc(funcName); fn != nil {
		context.markFuncAsLoaded(funcName, false)
		fn(context)
		context.markFuncAsLoaded(funcName, true)

		return nil
	}

	if body := f.GetCustomFuncBody(funcName); body != "" {

		f.debug("Formula.injectFuncToContext() - body=%v", body)
		context.markFuncAsLoaded(funcName, false)
		funcList := f.extractFunctionListFromFormulaString(body)

		for i := 0; i < len(funcList); i++ {

			f.debug("Formula.injectFuncToContext() - funcList[i]=%v", funcList[i])
			subFunc := funcList[i]
			if subFunc != funcName {
				err := f.injectFuncToContext(context, subFunc)
				if err != nil {
					return err
				}
			}
		}

		context.markFuncAsLoaded(funcName, true)

		_, err := context.Run(body)
		if err != nil {
			return err
		}

	}

	f.debug("Formula.injectFuncToContext() ended ...")

	return nil
}

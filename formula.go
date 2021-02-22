package goforit

//Formula a formula engine for creating Formulacontext
type Formula struct {
	driver VMDriver
	// r      *regexp.Regexp
	// customFuncs map[string]string
	customFuncs map[int]CustomFunctionRepository
	// builtInFuncs map[int]BuiltInFunctions
	Debug bool
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

	return f.customFuncs[0].RegisterCustomFunction(funcName, body)
}

//GetCustomFuncBody to get custom function source code
func (f Formula) GetCustomFuncBody(funcName string) string {

	for _, repo := range f.customFuncs {

		if body := repo.GetCustomFuncBody(funcName); body != "" {

			return body
		}

	}

	return ""
}

func (f Formula) extractFunctionListFromFormulaString(formulaStr string) []string {

	return f.driver.ExtractFunctionListFromFormulaString(formulaStr)
}

func (f Formula) newFormulaContext(vm VM) *FormulaContext {

	return &FormulaContext{VM: vm, loadedFuncs: make(map[string]bool), Debug: f.Debug}
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
		// context = &FormulaContext{vm: otto.New(), loadedFuncs: make(map[string]bool), Debug: f.Debug}
		vm, vmErr := f.driver.Get()
		if vmErr != nil {
			return nil, vmErr
		}

		//Falls through
		context = f.newFormulaContext(vm)
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

	// if fn := getBuilinFunc(funcName); fn != nil {
	// 	context.markFuncAsLoaded(funcName, false)
	// 	fn(context)
	// 	context.markFuncAsLoaded(funcName, true)

	// 	return nil
	// }
	if fn := context.VM.GetBuiltInFunc(funcName); fn != nil {
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

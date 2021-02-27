package impl

import (
	"github.com/lertrel/goforit/model"
	"github.com/lertrel/goforit/util"
	"github.com/lertrel/goforit/vm"
)

//DefaultFormula a formula engine for creating Formulacontext
type DefaultFormula struct {
	VM vm.Driver
	// r      *regexp.Regexp
	// CustomFuncs map[string]string
	// CustomFuncs  map[int]CustomFunctionRepository
	CustomFuncs []vm.CustomFunctionRepository
	// BuiltInFuncs map[int]BuiltInFunctions
	BuiltInFuncs []vm.BuiltInFunctions
	Debug        bool
}

//Driver is a method for getting vm.Driver implementation
func (f DefaultFormula) Driver() vm.Driver {
	return f.VM
}

func (f DefaultFormula) debug(format string, args ...interface{}) {
	util.Debug(f.Debug, format, args...)
}

//RegisterCustomFunction for registering custom function
//Ex.
//
// 		f := goforit.GetFormularBuilder().Get()
// 		f.RegisterCustomFunction(
// 			"$CIRCLE",
// 			`
// 			function $CIRCLE(radius) {
// 				return $RND(Math.PI * Math.pow(radius, 2), 10);
// 			}
// 		`)
//
func (f DefaultFormula) RegisterCustomFunction(funcName string, body string) bool {

	return f.CustomFuncs[0].RegisterFunction(funcName, body)
}

//GetCustomFunctionBody to get custom function source code
func (f DefaultFormula) GetCustomFunctionBody(funcName string) string {

	for _, repo := range f.CustomFuncs {

		if body := repo.GetFunctionBody(funcName); body != "" {

			return body
		}

	}

	return ""
}

func (f DefaultFormula) extractFunctionListFromFormulaString(formulaStr string) []string {

	return f.VM.ExtractFunctionNames(formulaStr)
}

//NewContext is a method for creating a new FormulaContext
func (f DefaultFormula) NewContext(script string) (c model.FormulaContext, err error) {

	vm, err := f.VM.Get()
	if err != nil {
		return
	}

	c = DefaultFormulaContext{
		VM:          vm,
		loadedFuncs: make(map[string]bool),
		formula:     f,
		Debug:       f.Debug,
	}

	if script != "" {
		if err = c.Prepare(script); err != nil {
			return nil, err
		}
	}

	return
}

//LoadContext If context is nil then create a new FormulaContext
//Then preparing a newly created context or a given context
//By loading referred functions (both built-in & custom) into context
// func (f Formula) LoadContext(context model.FormulaContext, formulaStr string) (c model.FormulaContext, err error) {

// 	// defer func() {
// 	// 	if r := recover(); r != nil {
// 	// 		c = nil
// 	// 		err = r.(error)
// 	// 	}
// 	// }()

// 	var contextImpl DefaultFormulaContext
// 	c = context
// 	if context == nil {
// 		// context = &	FormulaContext{vm: otto.New()}
// 		// context = &FormulaContext{vm: otto.New(), loadedFuncs: make(map[string]bool), Debug: f.Debug}
// 		// vm, vmErr := f.driver.Get()
// 		// if vmErr != nil {
// 		// 	return nil, vmErr
// 		// }

// 		//Falls through
// 		// context = f.newFormulaContext(vm)
// 		var contextTemp model.FormulaContext
// 		contextTemp, err = f.NewContext()
// 		if err != nil {
// 			return
// 		}
// 		contextImpl = contextTemp.(DefaultFormulaContext)
// 		c = contextImpl
// 	}

// 	f.debug("Formula.LoadContext() - Extracting function names from %v", formulaStr)
// 	funcList := f.extractFunctionListFromFormulaString(formulaStr)

// 	for i := 0; i < len(funcList); i++ {

// 		f.debug("Formula.LoadContext() - funcList[i]=%v", funcList[i])
// 		err := f.injectFuncToContext(&contextImpl, funcList[i])
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	return
// }

func (f DefaultFormula) injectFuncToContext(context *DefaultFormulaContext, funcName string) error {

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
	// if fn := context.VM.GetBuiltInFunc(funcName); fn != nil {
	if fn := context.VM.GetBuiltInFunc(funcName, f.BuiltInFuncs); fn != nil {
		context.markFuncAsLoaded(funcName, false)
		var c model.FormulaContext = context
		fn(&c)
		context.markFuncAsLoaded(funcName, true)

		return nil
	}

	if body := f.GetCustomFunctionBody(funcName); body != "" {

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

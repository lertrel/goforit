package goforit

import (
	"github.com/lertrel/goforit/model"
	"github.com/lertrel/goforit/vm"
)

//DefaultFormulaContext a formula context created by parsing formula string
//Any formula string have to be parsed into a context before using
type DefaultFormulaContext struct {
	// VM          *otto.Otto
	VM          vm.VM
	loadedFuncs map[string]bool
	formula     Formula
	Debug       bool
}

//Prepare If context is nil then create a new FormulaContext
//Then preparing a newly created context or a given context
//By loading referred functions (both built-in & custom) into context
func (c DefaultFormulaContext) Prepare(formulaStr string) error {

	c.debug("DefualtFormulaContext.Prepare() - Extracting function names from %v", formulaStr)
	funcList := c.formula.extractFunctionListFromFormulaString(formulaStr)

	for i := 0; i < len(funcList); i++ {

		c.debug("DefualtFormulaContext.Prepare() - funcList[i]=%v", funcList[i])
		err := c.formula.injectFuncToContext(&c, funcList[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// func (c FormulaContext) GetVM() *otto.Otto {
//
// 	return c.vm
// }

//Run to run a formula (commands given in form of string)
//To formula(s) that are relying on the same set of functions
//can be executed in the same context
//
// Ex.
//
// 		str := `
// 		$IF(true, 1, 2)
// 		`
// 		f := goforit.GetFormularBuilder().Get()
// 		c, err := f.LoadContext(nil, str)
// 		if err != nil {
// 			t.Error(err)
// 		}
//
// 		jsI, runtimeError := c.Run(str)
// 		if runtimeError != nil {
// 			t.Error(runtimeError)
// 		}
//
// 		goI, _ := jsI.ToInteger()
// 		t.Logf("goI = %v\n", goI)
//
func (c DefaultFormulaContext) Run(formulaString string) (model.Value, error) {

	return c.VM.Run(formulaString)
}

// func (c FormulaContext) Run(formulaString string) (JSValue, error) {

// 	value, err := c.vm.Run(formulaString)
// 	if err != nil {
// 		return JSValue{}, err
// 	}

// 	//Falls through
// 	return JSValue{impl: value}, nil
// }

//Get getting a variable inside FormulaContext
//
// 		str := `
// 		i = $SUMI(1, 2, $SUMI(1, $MIN(2,3)), $SUMI(2, 2), 5);
// 		f = $SUMF(1.5, $SUMF($MAX(1.2, 1.1), $ABS(-1.39)), $IF(i == 15, 5.0, 6.0));
// 		`
//
// 		f := goforit.GetFormularBuilder().Get()
// 		c, err := f.LoadContext(nil, str)
// 		if err != nil {
// 			t.Error(err)
// 		}
//
// 		_, runtimeError := c.Run(str)
// 		if runtimeError != nil {
// 			t.Error(runtimeError)
// 		}
//
// 		jsI, _ := c.Get("i")
// 		goI, _ := jsI.ToInteger()
// 		jsF, _ := c.Get("f")
// 		goF, _ := jsF.ToFloat()
//
func (c DefaultFormulaContext) Get(varname string) (model.Value, error) {

	return c.VM.Get(varname)
}

// func (c FormulaContext) Get(varname string) (JSValue, error) {

// 	value, err := c.vm.Get(varname)
// 	if err != nil {
// 		return JSValue{}, err
// 	}

// 	//Falls through
// 	return JSValue{impl: value}, nil
// }

//Set setting a variable inside FormulaContext
//
//Ex.
//
//		f := goforit.GetFormularBuilder().Get()
// 		str := `
// 		area1 = $RND(Math.sqrt($SUMF($RND(a*3,2), $RND(b*4,2), $RND(c*5,2))), 10);
// 		area2 = $CIRCLE(radius);
// 		console.log("a="+a);
// 		console.log("b="+b);
// 		console.log("c="+c);
// 		console.log("radius="+radius);
// 		`
// 		c, err := f.LoadContext(nil, str)
// 		if err != nil {
// 			t.Error(err)
// 		}
//
// 		c.Set("a", 2)
// 		c.Set("b", 3)
// 		c.Set("c", 4)
// 		c.Set("radius", 5)
//
// 		_, runtimeError := c.Run(str)
// 		if runtimeError != nil {
// 			t.Error(runtimeError)
// 		}
//
// 		jsArea1, _ := c.Get("area1")
// 		jsArea2, _ := c.Get("area2")
//
// 		area1, _ := jsArea1.ToFloat()
// 		area2, _ := jsArea2.ToFloat()
//
func (c DefaultFormulaContext) Set(varname string, value interface{}) error {

	err := c.VM.Set(varname, value)
	if err != nil {
		return err
	}

	//Falls through
	return nil
}

func (c DefaultFormulaContext) debug(format string, args ...interface{}) {
	debug(debugFlag || c.Debug, format, args...)
}

func (c DefaultFormulaContext) isFuncLoaded(funcName string) bool {

	_, found := c.loadedFuncs[funcName]

	c.debug("FormulaContext.isFuncLoaded() - c.loadedFunc=%v", c.loadedFuncs)
	c.debug("FormulaContext.isFuncLoaded() - found=%v", found)

	return found
}

func (c *DefaultFormulaContext) markFuncAsLoaded(funcName string, loaded bool) {

	c.debug("FormulaContext.markFuncAsLoaded(before) - c.loadedFunc=%v", c.loadedFuncs)
	c.loadedFuncs[funcName] = loaded
	c.debug("FormulaContext.markFuncAsLoaded(after) - c.loadedFunc=%v", c.loadedFuncs)
}

package goforit

//FormulaContext a formula context created by parsing formula string
//Any formula string have to be parsed into a context before using
type FormulaContext struct {
	// VM          *otto.Otto
	VM          VM
	loadedFuncs map[string]bool
	Debug       bool
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
// str := `
// $IF(true, 1, 2)
// `
// f := Get()
// c, err := f.LoadContext(nil, str)
// if err != nil {
// 	t.Error(err)
// }
//
// jsI, runtimeError := c.Run(str)
// if runtimeError != nil {
// 	t.Error(runtimeError)
// }
//
// goI, _ := jsI.ToInteger()
// t.Logf("goI = %v\n", goI)
func (c FormulaContext) Run(formulaString string) (JSValue, error) {

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
func (c FormulaContext) Get(varname string) (JSValue, error) {

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

//Set getting a variable inside FormulaContext
//
//Ex.
// str := `
// area1 = $RND(Math.sqrt($SUMF($RND(a*3,2), $RND(b*4,2), $RND(c*5,2))), 10);
// area2 = $CIRCLE(radius);
// console.log("a="+a);
// console.log("b="+b);
// console.log("c="+c);
// console.log("radius="+radius);
// `
// 	c, err := f.LoadContext(nil, str)
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	c.Set("a", 2)
// 	c.Set("b", 3)
// 	c.Set("c", 4)
// 	c.Set("radius", 5)
//
// 	_, runtimeError := c.Run(str)
// 	if runtimeError != nil {
// 		t.Error(runtimeError)
// 	}
//
// 	jsArea1, _ := c.Get("area1")
// 	jsArea2, _ := c.Get("area2")
//
// 	area1, _ := jsArea1.ToFloat()
// 	area2, _ := jsArea2.ToFloat()
//
func (c FormulaContext) Set(varname string, value interface{}) error {

	err := c.VM.Set(varname, value)
	if err != nil {
		return err
	}

	//Falls through
	return nil
}

func (c FormulaContext) debug(format string, args ...interface{}) {
	debug(debugFlag || c.Debug, format, args...)
}

func (c FormulaContext) isFuncLoaded(funcName string) bool {

	_, found := c.loadedFuncs[funcName]

	c.debug("FormulaContext.isFuncLoaded() - c.loadedFunc=%v", c.loadedFuncs)
	c.debug("FormulaContext.isFuncLoaded() - found=%v", found)

	return found
}

func (c *FormulaContext) markFuncAsLoaded(funcName string, loaded bool) {

	c.debug("FormulaContext.markFuncAsLoaded(before) - c.loadedFunc=%v", c.loadedFuncs)
	c.loadedFuncs[funcName] = loaded
	c.debug("FormulaContext.markFuncAsLoaded(after) - c.loadedFunc=%v", c.loadedFuncs)
}

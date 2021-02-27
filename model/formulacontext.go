package model

//FormulaContext a formula context created by parsing formula string
//Any formula string have to be parsed into a context before using
type FormulaContext interface {

	//LoadContext If context is nil then create a new FormulaContext
	//Then preparing a newly created context or a given context
	//By loading referred functions (both built-in & custom) into context
	Prepare(formulaString string) error

	//Run to run a formula (commands given in form of string)
	//To formula(s) that are relying on the same set of functions
	//can be executed in the same context
	//
	// Ex.
	//
	// 		str := `
	// 		$IF(true, 1, 2)
	// 		`
	// 		f := goforit.NewFormularBuilder().Get()
	// 		c, err := f.NewContext(str)
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
	Run(formulaString string) (Value, error)

	//Get getting a variable inside FormulaContext
	//
	// 		str := `
	// 		i = $SUMI(1, 2, $SUMI(1, $MIN(2,3)), $SUMI(2, 2), 5);
	// 		f = $SUMF(1.5, $SUMF($MAX(1.2, 1.1), $ABS(-1.39)), $IF(i == 15, 5.0, 6.0));
	// 		`
	//
	// 		f := goforit.NewFormularBuilder().Get()
	// 		c, err := f.NewContext(str)
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
	Get(varname string) (Value, error)

	//Set setting a variable inside FormulaContext
	//
	//Ex.
	//
	//		f := goforit.NewFormularBuilder().Get()
	// 		str := `
	// 		area1 = $RND(Math.sqrt($SUMF($RND(a*3,2), $RND(b*4,2), $RND(c*5,2))), 10);
	// 		area2 = $CIRCLE(radius);
	// 		console.log("a="+a);
	// 		console.log("b="+b);
	// 		console.log("c="+c);
	// 		console.log("radius="+radius);
	// 		`
	// 		c, err := f.NewContext(str)
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
	Set(varname string, value interface{}) error
}

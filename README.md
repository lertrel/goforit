# goforit
A pure go package for enabling runtime customization formulas for Go utilizing JavaScript package. 

With Go-For-It package, developers can externalize program formula(s) (e.g., in text file, csv, excel file, DB, etc.) and load them during runtime. 

The benefit of having formula(s) externalized is, for program that extensively uses formula(s) so numbers of formula(s) can be maintained outside of source code which in turn making the source code more cleaner, adding/changing formula can be done without stopping the running program, as well as complex formula(s) could be handled by specialists (a.k.a. users).  

**NOTE** The current implementation is built on top of "otto"

Here is a very basic example:

```go
package main

import (
	"fmt"
	"github.com/lertrel/goforit"
)

func main() {

	//Creating Formula instance
	f := goforit.NewFormulaBuilder().Get()

//$SUMI is a built-in function
str := `
$SUMI(x, 2*x, 3*x, 4*x) + x
`
  //Loading denpendencies (functions) into context instance
	c, err := f.LoadContext(nil, str)
	if err != nil {
		panic(err)
	}

  //Set value to x
	c.Set("x", 10.0)

  //Run formula & getting result
	jsVal, jsErr := c.Run(str)
	if jsErr != nil {
		panic(jsErr)
	}

  //Converting result into go variable
	goVal, _ := jsVal.ToInteger()

	fmt.Printf("goVal=%v\n", goVal)

}
```

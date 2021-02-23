package main

import (
	"fmt"

	"github.com/lertrel/goforit"
)

func main() {

	f := goforit.GetFormulaBuilder().Get()

	str := `
$SUMI(x, 2*x, 3*x, 4*x) + x
`
	c, err := f.LoadContext(nil, str)
	if err != nil {
		panic(err)
	}

	c.Set("x", 10.0)

	jsVal, jsErr := c.Run(str)
	if jsErr != nil {
		panic(jsErr)
	}

	goVal, _ := jsVal.ToInteger()

	fmt.Printf("goVal=%v\n", goVal)

}

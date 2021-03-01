package main

import (
	"fmt"

	"github.com/lertrel/goforit"
)

func main() {

	f := goforit.NewFormulaBuilder().Get()
	f.RegisterCustomFunction(
		"$CIRCLE",
		`
				function $CIRCLE(radius) {
					return $RND(Math.PI * Math.pow(radius, 2), 10);
				}
			`)

	script := "$CIRCLE(r)"
	fc, err := f.NewContext(script)
	if err != nil {
		panic(err)
	}

	if err := fc.Set("r", 5); err != nil {
		panic(err)
	}

	jsRet, err := fc.Run(script)
	if err != nil {
		panic(err)
	}

	goRet, err := jsRet.ToFloat()
	if err != nil {
		panic(err)
	}

	fmt.Printf("goRet=%v\n", goRet)
}

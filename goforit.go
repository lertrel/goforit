package goforit

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"regexp"
	"strconv"

	"github.com/robertkrimen/otto"
)

var debugFlag = false

//Get An entry point to obtain Formula
func Get() Formula {

	r, _ := regexp.Compile("(\\$[^\\$()\\s]+)\\(")

	return Formula{r: r, customFuncs: make(map[string]string), Debug: false}
}

func debug(b bool, format string, args ...interface{}) {
	if b {
		log.Printf(format, args...)
	}
}

func validateJSFuncArguments(funcName string, cnt int, call otto.FunctionCall) {

	if len(call.ArgumentList) != cnt {
		errMsg := fmt.Sprintf("%s - wrong number of arguments (expecting %d)", funcName, cnt)
		panic(errors.New(errMsg))
	}

}

func getJSFloat(call otto.FunctionCall, index int) float64 {

	v, err := call.ArgumentList[index].ToFloat()
	if err != nil {
		panic(err)
	}

	return v
}

func getJSInt(call otto.FunctionCall, index int) int64 {

	v, err := call.ArgumentList[index].ToInteger()
	if err != nil {
		panic(err)
	}

	return v
}

func getJSString(call otto.FunctionCall, index int) string {

	v, err := call.ArgumentList[index].ToString()
	if err != nil {
		panic(err)
	}

	return v
}

func getJSBoolean(call otto.FunctionCall, index int) bool {

	v, err := call.ArgumentList[index].ToBoolean()
	if err != nil {
		panic(err)
	}

	return v
}

func toJSValue(context *FormulaContext, value interface{}) otto.Value {

	jsResult, err := context.vm.ToValue(value)
	if err != nil {
		panic(err)
	}

	return jsResult
}

func getBuilinFunc(funcName string) func(context *FormulaContext) {

	switch funcName {

	case "$ABS":
		return func(context *FormulaContext) {

			context.vm.Set(funcName, func(call otto.FunctionCall) otto.Value {

				validateJSFuncArguments("$ABS", 1, call)

				v := getJSFloat(call, 0)
				result := math.Abs(v)

				return toJSValue(context, result)
			})
		}
	case "$RND":
		return func(context *FormulaContext) {

			context.vm.Set(funcName, func(call otto.FunctionCall) otto.Value {

				validateJSFuncArguments("$RND", 2, call)

				v := getJSFloat(call, 0)
				p := getJSInt(call, 1)
				if 0 > p || p > 10 {
					panic(errors.New("$RND - second parameter should be between 0 and 10"))
				}

				result := v
				f := math.Pow10(int(p))
				result = math.Round(v*f) / f

				return toJSValue(context, result)
			})
		}
	case "$FLOOR", "$FLR":
		return func(context *FormulaContext) {

			context.vm.Set(funcName, func(call otto.FunctionCall) otto.Value {

				validateJSFuncArguments("$FLOOR", 2, call)

				v := getJSFloat(call, 0)
				p := getJSInt(call, 1)
				if 0 > p || p > 10 {
					panic(errors.New("$RND - second parameter should be between 0 and 10"))
				}

				result := v
				f := math.Pow10(int(p))
				result = math.Floor(v*f) / f

				return toJSValue(context, result)
			})
		}
	case "$CEIL":
		return func(context *FormulaContext) {

			context.vm.Set(funcName, func(call otto.FunctionCall) otto.Value {

				validateJSFuncArguments("$CEIL", 2, call)

				v := getJSFloat(call, 0)
				p := getJSInt(call, 1)
				if 0 > p || p > 10 {
					panic(errors.New("$RND - second parameter should be between 0 and 10"))
				}

				result := v
				f := math.Pow10(int(p))
				result = math.Ceil(v*f) / f

				return toJSValue(context, result)
			})
		}
	case "$IF":
		return func(context *FormulaContext) {

			context.vm.Set(funcName, func(call otto.FunctionCall) otto.Value {

				validateJSFuncArguments("$IF", 3, call)

				b := getJSBoolean(call, 0)
				var jsResult otto.Value

				if b {
					jsResult = call.ArgumentList[1]
				} else {
					jsResult = call.ArgumentList[2]
				}

				return jsResult
			})
		}
	case "$SUMI":
		return func(context *FormulaContext) {

			context.vm.Set(funcName, func(call otto.FunctionCall) otto.Value {

				var result int64
				result = 0

				for i := 0; i < len(call.ArgumentList); i++ {

					v := getJSInt(call, i)
					result = result + v

				}

				return toJSValue(context, result)
			})
		}
	case "$SUMF":
		return func(context *FormulaContext) {

			context.vm.Set(funcName, func(call otto.FunctionCall) otto.Value {

				var result float64
				result = 0.0

				for i := 0; i < len(call.ArgumentList); i++ {

					v := getJSFloat(call, i)
					result = result + v

				}

				return toJSValue(context, result)
			})
		}
	case "$AVG":
		return func(context *FormulaContext) {

			context.vm.Set(funcName, func(call otto.FunctionCall) otto.Value {

				result := big.NewFloat(0.0)
				var i int

				for i = 0; i < len(call.ArgumentList); i++ {

					v := getJSFloat(call, i)
					bigV := big.NewFloat(v)
					result.Add(result, bigV)
				}

				result.Quo(result, big.NewFloat(float64(i)))
				f, err := strconv.ParseFloat(result.String(), 10)
				if err != nil {
					panic(err)
				}

				return toJSValue(context, f)
			})
		}
	case "$MIN":
		return func(context *FormulaContext) {

			context.vm.Set(funcName, func(call otto.FunctionCall) otto.Value {

				hasLast := false
				var last otto.Value = otto.Value{}

				for i := 0; i < len(call.ArgumentList); i++ {

					v := call.ArgumentList[i]

					if !hasLast {
						last = v
						hasLast = true
					} else {
						l, el := last.ToFloat()
						r, er := v.ToFloat()

						if el != nil {
							panic(el)
						}
						if er != nil {
							panic(er)
						}

						if l > r {
							last = v
						}
					}

				}

				return last
			})
		}
	case "$MAX":
		return func(context *FormulaContext) {

			context.vm.Set(funcName, func(call otto.FunctionCall) otto.Value {

				hasLast := false
				var last otto.Value = otto.Value{}

				for i := 0; i < len(call.ArgumentList); i++ {

					v := call.ArgumentList[i]

					if !hasLast {
						last = v
						hasLast = true
					} else {
						l, el := last.ToFloat()
						r, er := v.ToFloat()

						if el != nil {
							panic(el)
						}
						if er != nil {
							panic(er)
						}

						if l < r {
							last = v
						}
					}

				}

				return last
			})
		}
	default:
		return nil

	}
}

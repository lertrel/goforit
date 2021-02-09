package formula

import (
	"fmt"
	"log"
	"errors"
	"strconv"
	"regexp"
	"math"
	"math/big"
	"github.com/robertkrimen/otto"
)

var DEBUG = false

type Formula struct {
	
	r *regexp.Regexp
	customFuncs map[string]string
	Debug bool
}

func Get() Formula {

	r, _ := regexp.Compile("(\\$[^\\$()\\s]+)\\(")

	return Formula{r: r, customFuncs: make(map[string]string), Debug: false}
}

func (f Formula) debug(format string, args ...interface{}) {
	debug(DEBUG || f.Debug, format, args...)
}

func (f *Formula) RegisterCustomFunction(funcName string, body string) bool {

	_, found := f.customFuncs[funcName]

	f.customFuncs[funcName] = body

	return found
}

func (f Formula) GetCustomFuncBody(funcName string) string {

	body, found := f.customFuncs[funcName]

	if found {
		return body
	} else {
		return ""
	}
}

func (f Formula) extractFunctionListFromFormulaString(formulaStr string) []string {

	matches := f.r.FindAllStringSubmatch(formulaStr, -1)
	dedupMatches := make(map[string]bool)

	for i := 0; i < len(matches); i++ {
		dedupMatches[matches[i][1]] = true
	}

	funArr := make([]string, len(dedupMatches))

	i := 0
	for k, _ := range dedupMatches {
		funArr[i] = k
		i++
	}

	return funArr
}

func (f Formula) LoadContext(context *FormulaContext, formulaStr string) (c *FormulaContext, err error) {

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		c = nil
	// 		err = r.(error)
	// 	}
	// }()	
	
	if context == nil {
		// context = &	FormulaContext{vm: otto.New()}
		context = &FormulaContext{vm: otto.New(), loadedFuncs: make(map[string] bool), Debug: f.Debug}
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



type FormulaContext struct {

	vm *otto.Otto
	loadedFuncs map[string]bool
	Debug bool
}

// func (c FormulaContext) GetVM() *otto.Otto {
	
// 	return c.vm
// }

func (c FormulaContext) Run(formulaString string) (JSValue, error) {
	
	value, err := c.vm.Run(formulaString)
	if err != nil {
		return JSValue{}, err
	} else {
		return JSValue{impl: value}, nil
	}
}

func (c FormulaContext) Get(varname string) (JSValue, error) {
	
	value, err := c.vm.Get(varname)
	if err != nil {
		return JSValue{}, err
	} else {
		return JSValue{impl: value}, nil
	}
}

func (c FormulaContext) debug(format string, args ...interface{}) {
	debug(DEBUG || c.Debug, format, args...)
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
		return func (context *FormulaContext) {

			context.vm.Set(funcName, func(call otto.FunctionCall) otto.Value {

				validateJSFuncArguments("$ABS", 1, call)

				v := getJSFloat(call, 0)
				result := math.Abs(v)

				return toJSValue(context, result)
			})
		}
	case "$RND":
		return func (context *FormulaContext) {

			context.vm.Set(funcName, func(call otto.FunctionCall) otto.Value {

				validateJSFuncArguments("$RND", 2, call)

				v := getJSFloat(call, 0)
				p := getJSInt(call, 1)
				if 0 > p || p > 10  {
					panic(errors.New("$RND - second parameter should be between 0 and 10"))
				}

				result := v
				f := math.Pow10(int(p))
				result = math.Round(v*f) / f

				return toJSValue(context, result)
			})
		}
	case "$FLOOR", "$FLR":
		return func (context *FormulaContext) {

			context.vm.Set(funcName, func(call otto.FunctionCall) otto.Value {

				validateJSFuncArguments("$FLOOR", 2, call)

				v := getJSFloat(call, 0)
				p := getJSInt(call, 1)
				if 0 > p || p > 10  {
					panic(errors.New("$RND - second parameter should be between 0 and 10"))
				}

				result := v
				f := math.Pow10(int(p))
				result = math.Floor(v*f) / f

				return toJSValue(context, result)
			})
		}
	case "$CEIL":
		return func (context *FormulaContext) {

			context.vm.Set(funcName, func(call otto.FunctionCall) otto.Value {

				validateJSFuncArguments("$CEIL", 2, call)

				v := getJSFloat(call, 0)
				p := getJSInt(call, 1)
				if 0 > p || p > 10  {
					panic(errors.New("$RND - second parameter should be between 0 and 10"))
				}

				result := v
				f := math.Pow10(int(p))
				result = math.Ceil(v*f) / f

				return toJSValue(context, result)
			})
		}
	case "$IF":
		return func (context *FormulaContext) {

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
		return func (context *FormulaContext) {

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
		return func (context *FormulaContext) {

			context.vm.Set(funcName, func(call otto.FunctionCall) otto.Value {

				var result float64
				result = 0.0

				for i :=0; i < len(call.ArgumentList); i++ {

					v := getJSFloat(call, i)
					result = result + v

				}

				return toJSValue(context, result)
			})
		}
	case "$AVG":
		return func (context *FormulaContext) {

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
		return func (context *FormulaContext) {

			context.vm.Set(funcName, func(call otto.FunctionCall) otto.Value {

				hasLast := false
				var last otto.Value = otto.Value{}

				for i :=0; i < len(call.ArgumentList); i++ {

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
						if (er != nil) {
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
		return func (context *FormulaContext) {

			context.vm.Set(funcName, func(call otto.FunctionCall) otto.Value {

				hasLast := false
				var last otto.Value = otto.Value{}

				for i :=0; i < len(call.ArgumentList); i++ {

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
						if (er != nil) {
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


type JSValue struct {
	impl otto.Value
}

// func (this JSValue) ToValue(value interface{}) (JSValue, error) {
	
// 	impl, err := this.impl.ToValue(value)

// 	return JSValue{impl: impl}, err
// }

func (this JSValue) IsDefined() bool {
	return this.impl.IsDefined()
}

func (this JSValue) IsUndefined() bool {
	return this.impl.IsUndefined()
}

func (this JSValue) IsNull() bool {
	return this.impl.IsNull()
}

func (this JSValue) IsPrimitive() bool {
	return this.impl.IsPrimitive()
}

func (this JSValue) IsBoolean() bool {
	return this.impl.IsBoolean()
}

func (this JSValue) IsNumber() bool {
	return this.impl.IsNumber()
}

func (this JSValue) IsNaN() bool {
	return this.impl.IsNaN()
}

func (this JSValue) IsString() bool {
	return this.impl.IsString()
}

func (this JSValue) IsObject() bool {
	return this.impl.IsObject()
}

func (this JSValue) IsFunction() bool {
	return this.impl.IsFunction()
}

func (this JSValue) ToBoolean() (bool, error) {
	return this.impl.ToBoolean()
}

func (this JSValue) ToFloat() (float64, error) {
	return this.impl.ToFloat()
}

func (this JSValue) ToInteger() (int64, error) {
	return this.impl.ToInteger()
}

func (this JSValue) ToString() (string, error) {
	return this.impl.ToString()
}

func (this JSValue) Export() (interface{}, error) {
	return this.impl.Export()
}

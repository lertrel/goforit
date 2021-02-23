package goforit

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/robertkrimen/otto"
)

var r, _ = regexp.Compile("(\\$[^\\$()\\s]+)\\(")

//OttoVMDriver otto implementation of VMDriver
type OttoVMDriver struct {
	templateVM VM
	// funcs      map[int]BuiltInFunctions
}

//Get getting VM implementation of this driver
func (d OttoVMDriver) Get() (VM, error) {

	if d.templateVM != nil {

		nativeVM := d.templateVM.(OttoVM)

		// return OttoVM{vm: nativeVM.vm.Copy(), funcs: nativeVM.funcs}, nil
		return OttoVM{vm: nativeVM.vm.Copy()}, nil
	}

	//Falls through
	// return OttoVM{vm: otto.New(), funcs: d.funcs}, nil
	return OttoVM{vm: otto.New()}, nil
}

//SetTemplate set template VM, this VM will be used for cloning
//another VM to reduce constructing overhead
func (d *OttoVMDriver) SetTemplate(template VM) {
	d.templateVM = template
}

//ExtractFunctionListFromFormulaString extracting functions names
//from the given script/formular so that the function can be loaded
//by Formula.LoadContext() before being executed, otherwise the
//unloaded functions will not be known to the scripting/VM engine
func (d OttoVMDriver) ExtractFunctionListFromFormulaString(formulaStr string) []string {

	matches := r.FindAllStringSubmatch(formulaStr, -1)
	dedupMatches := make(map[string]bool)

	for i := 0; i < len(matches); i++ {
		dedupMatches[matches[i][1]] = true
	}

	funArr := make([]string, len(dedupMatches))

	i := 0
	for k := range dedupMatches {
		funArr[i] = k
		i++
	}

	return funArr
}

//OttoVM otto implementation of VM (abstract layer)
type OttoVM struct {
	vm *otto.Otto
	// funcs map[int]BuiltInFunctions
}

//Run for running a given script/formual
func (v OttoVM) Run(formulaString string) (JSValue, error) {

	value, err := v.vm.Run(formulaString)
	if err != nil {
		return NewJSValue(nil, nil), err
	}

	//Falls through
	// return JSValue{impl: value}, nil
	return NewJSValue(v, value), nil
}

//Get getting value of a variable out of scripting context
func (v OttoVM) Get(varname string) (JSValue, error) {

	value, err := v.vm.Get(varname)
	if err != nil {
		// return JSValue{}, err
		return NewJSValue(nil, nil), err
	}

	//Falls through
	// return JSValue{impl: value}, nil
	return NewJSValue(v, value), nil
}

//Set setting value of a valiable inside scripting context
func (v OttoVM) Set(varname string, value interface{}) error {

	err := v.vm.Set(varname, value)
	if err != nil {
		return err
	}

	//Falls through
	return nil
}

//ValidateFuncArguments validting if a given scripting function has
//number of arguments equal to the given cnt
func (v OttoVM) ValidateFuncArguments(funcName string, cnt int, funcDef interface{}) {

	validateJSFuncArguments(funcName, cnt, funcDef.(otto.FunctionCall))
}

//GetFuncArgsCount getting function arguments count
func (v OttoVM) GetFuncArgsCount(funcDef interface{}) int {

	call := funcDef.(otto.FunctionCall)

	return len(call.ArgumentList)
}

//GetFuncArgAsFloat getting function argument refering by the index as float64
func (v OttoVM) GetFuncArgAsFloat(funcDef interface{}, index int) float64 {

	return getJSFloat(funcDef.(otto.FunctionCall), index)
}

//GetFuncArgAsInt getting function argument refering by the index as int64
func (v OttoVM) GetFuncArgAsInt(funcDef interface{}, index int) int64 {

	return getJSInt(funcDef.(otto.FunctionCall), index)
}

//GetFuncArgAsString getting function argument refering by the index as string
func (v OttoVM) GetFuncArgAsString(funcDef interface{}, index int) string {

	return getJSString(funcDef.(otto.FunctionCall), index)
}

//GetFuncArgAsBoolean getting function argument refering by the index as boolean
func (v OttoVM) GetFuncArgAsBoolean(funcDef interface{}, index int) bool {

	return getJSBoolean(funcDef.(otto.FunctionCall), index)
}

//GetFuncArgAsIs getting function argument refering by the index as raw type
//depending on each scripting/VM engine e.g., otto.Value for otto
func (v OttoVM) GetFuncArgAsIs(funcDef interface{}, index int) interface{} {

	return funcDef.(otto.FunctionCall).ArgumentList[index]
}

//ToVMValue converting go variable into scriing/VM value e.g., otto.Value for otto
func (v OttoVM) ToVMValue(goValue interface{}) interface{} {

	jsResult, err := v.vm.ToValue(goValue)
	if err != nil {
		panic(err)
	}

	return jsResult
}

//GetBuiltInFunc getting a function for registering / connecting
func (v OttoVM) GetBuiltInFunc(funcName string, funcs map[int]BuiltInFunctions) func(context *FormulaContext) {

	// for _, f := range v.funcs {
	for _, f := range funcs {

		if f.Has(funcName) {

			return func(context *FormulaContext) {

				// ctxVM := context.VM.(OttoVM)

				// if &ctxVM != &v {
				// 	panic(errors.New("The given FormulaContext does not connect to the current OttoVM"))
				// }

				v.vm.Set(funcName, func(call otto.FunctionCall) otto.Value {

					// for _, f := range v.funcs {
					for _, f := range funcs {
						result, found := f.Execute(funcName, v, call)

						if found {
							return v.ToVMValue(result).(otto.Value)
						}
					}

					panic(errors.New("Runtime error: Function not found - " + funcName))
				})
			}
		}
	}

	//Falls through
	return nil
}

//IsDefined check if the current JS value is defined
func (v OttoVM) IsDefined(vmValue interface{}) bool {

	return vmValue.(otto.Value).IsDefined()
}

//IsUndefined check if the current JS value is undefined
func (v OttoVM) IsUndefined(vmValue interface{}) bool {

	return vmValue.(otto.Value).IsUndefined()
}

//IsNull check if the current JS value is null
func (v OttoVM) IsNull(vmValue interface{}) bool {

	return vmValue.(otto.Value).IsNull()
}

//IsPrimitive check if the current JS value is primitive
func (v OttoVM) IsPrimitive(vmValue interface{}) bool {

	return vmValue.(otto.Value).IsPrimitive()
}

//IsBoolean check if the current JS value is boolean
func (v OttoVM) IsBoolean(vmValue interface{}) bool {

	return vmValue.(otto.Value).IsBoolean()
}

//IsNumber check if the current JS value is number
func (v OttoVM) IsNumber(vmValue interface{}) bool {

	return vmValue.(otto.Value).IsNumber()
}

//IsNaN check if the current JS value is nan
func (v OttoVM) IsNaN(vmValue interface{}) bool {

	return vmValue.(otto.Value).IsNaN()
}

//IsString check if the current JS value is string
func (v OttoVM) IsString(vmValue interface{}) bool {

	return vmValue.(otto.Value).IsString()
}

//IsObject check if the current JS value is object
func (v OttoVM) IsObject(vmValue interface{}) bool {

	return vmValue.(otto.Value).IsObject()
}

//IsFunction check if the current JS value is function
func (v OttoVM) IsFunction(vmValue interface{}) bool {

	return vmValue.(otto.Value).IsFunction()
}

//ToBoolean get value as bool
func (v OttoVM) ToBoolean(vmValue interface{}) (bool, error) {

	return vmValue.(otto.Value).ToBoolean()
}

//ToFloat get value as float
func (v OttoVM) ToFloat(vmValue interface{}) (float64, error) {

	return vmValue.(otto.Value).ToFloat()
}

//ToInteger get value as integer
func (v OttoVM) ToInteger(vmValue interface{}) (int64, error) {

	return vmValue.(otto.Value).ToInteger()
}

//ToString get value as string
func (v OttoVM) ToString(vmValue interface{}) (string, error) {

	return vmValue.(otto.Value).ToString()
}

//Export js value to Go value type as specified in a given interface
func (v OttoVM) Export(vmValue interface{}) (interface{}, error) {

	return vmValue.(otto.Value).Export()
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

// func toJSValue(context *FormulaContext, value interface{}) otto.Value {

// 	jsResult, err := context.VM.ToValue(value)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return jsResult
// }

// func getBuilinFunc(funcName string) func(context *FormulaContext) {

// 	switch funcName {

// 	case "$ABS":
// 		return func(context *FormulaContext) {

// 			context.VM.Set(funcName, func(call otto.FunctionCall) otto.Value {

// 				validateJSFuncArguments("$ABS", 1, call)

// 				v := getJSFloat(call, 0)
// 				result := math.Abs(v)

// 				return toJSValue(context, result)
// 			})
// 		}
// 	case "$RND":
// 		return func(context *FormulaContext) {

// 			context.VM.Set(funcName, func(call otto.FunctionCall) otto.Value {

// 				validateJSFuncArguments("$RND", 2, call)

// 				v := getJSFloat(call, 0)
// 				p := getJSInt(call, 1)
// 				if 0 > p || p > 10 {
// 					panic(errors.New("$RND - second parameter should be between 0 and 10"))
// 				}

// 				result := v
// 				f := math.Pow10(int(p))
// 				result = math.Round(v*f) / f

// 				return toJSValue(context, result)
// 			})
// 		}
// 	case "$FLOOR", "$FLR":
// 		return func(context *FormulaContext) {

// 			context.VM.Set(funcName, func(call otto.FunctionCall) otto.Value {

// 				validateJSFuncArguments("$FLOOR", 2, call)

// 				v := getJSFloat(call, 0)
// 				p := getJSInt(call, 1)
// 				if 0 > p || p > 10 {
// 					panic(errors.New("$RND - second parameter should be between 0 and 10"))
// 				}

// 				result := v
// 				f := math.Pow10(int(p))
// 				result = math.Floor(v*f) / f

// 				return toJSValue(context, result)
// 			})
// 		}
// 	case "$CEIL":
// 		return func(context *FormulaContext) {

// 			context.VM.Set(funcName, func(call otto.FunctionCall) otto.Value {

// 				validateJSFuncArguments("$CEIL", 2, call)

// 				v := getJSFloat(call, 0)
// 				p := getJSInt(call, 1)
// 				if 0 > p || p > 10 {
// 					panic(errors.New("$RND - second parameter should be between 0 and 10"))
// 				}

// 				result := v
// 				f := math.Pow10(int(p))
// 				result = math.Ceil(v*f) / f

// 				return toJSValue(context, result)
// 			})
// 		}
// 	case "$IF":
// 		return func(context *FormulaContext) {

// 			context.VM.Set(funcName, func(call otto.FunctionCall) otto.Value {

// 				validateJSFuncArguments("$IF", 3, call)

// 				b := getJSBoolean(call, 0)
// 				var jsResult otto.Value

// 				if b {
// 					jsResult = call.ArgumentList[1]
// 				} else {
// 					jsResult = call.ArgumentList[2]
// 				}

// 				return jsResult
// 			})
// 		}
// 	case "$SUMI":
// 		return func(context *FormulaContext) {

// 			context.VM.Set(funcName, func(call otto.FunctionCall) otto.Value {

// 				var result int64
// 				result = 0

// 				for i := 0; i < len(call.ArgumentList); i++ {

// 					v := getJSInt(call, i)
// 					result = result + v

// 				}

// 				return toJSValue(context, result)
// 			})
// 		}
// 	case "$SUMF":
// 		return func(context *FormulaContext) {

// 			context.VM.Set(funcName, func(call otto.FunctionCall) otto.Value {

// 				var result float64
// 				result = 0.0

// 				for i := 0; i < len(call.ArgumentList); i++ {

// 					v := getJSFloat(call, i)
// 					result = result + v

// 				}

// 				return toJSValue(context, result)
// 			})
// 		}
// 	case "$AVG":
// 		return func(context *FormulaContext) {

// 			context.VM.Set(funcName, func(call otto.FunctionCall) otto.Value {

// 				result := big.NewFloat(0.0)
// 				var i int

// 				for i = 0; i < len(call.ArgumentList); i++ {

// 					v := getJSFloat(call, i)
// 					bigV := big.NewFloat(v)
// 					result.Add(result, bigV)
// 				}

// 				result.Quo(result, big.NewFloat(float64(i)))
// 				f, err := strconv.ParseFloat(result.String(), 10)
// 				if err != nil {
// 					panic(err)
// 				}

// 				return toJSValue(context, f)
// 			})
// 		}
// 	case "$MIN":
// 		return func(context *FormulaContext) {

// 			context.VM.Set(funcName, func(call otto.FunctionCall) otto.Value {

// 				hasLast := false
// 				var last otto.Value = otto.Value{}

// 				for i := 0; i < len(call.ArgumentList); i++ {

// 					v := call.ArgumentList[i]

// 					if !hasLast {
// 						last = v
// 						hasLast = true
// 					} else {
// 						l, el := last.ToFloat()
// 						r, er := v.ToFloat()

// 						if el != nil {
// 							panic(el)
// 						}
// 						if er != nil {
// 							panic(er)
// 						}

// 						if l > r {
// 							last = v
// 						}
// 					}

// 				}

// 				return last
// 			})
// 		}
// 	case "$MAX":
// 		return func(context *FormulaContext) {

// 			context.VM.Set(funcName, func(call otto.FunctionCall) otto.Value {

// 				hasLast := false
// 				var last otto.Value = otto.Value{}

// 				for i := 0; i < len(call.ArgumentList); i++ {

// 					v := call.ArgumentList[i]

// 					if !hasLast {
// 						last = v
// 						hasLast = true
// 					} else {
// 						l, el := last.ToFloat()
// 						r, er := v.ToFloat()

// 						if el != nil {
// 							panic(el)
// 						}
// 						if er != nil {
// 							panic(er)
// 						}

// 						if l < r {
// 							last = v
// 						}
// 					}

// 				}

// 				return last
// 			})
// 		}
// 	default:
// 		return nil

// 	}
// }

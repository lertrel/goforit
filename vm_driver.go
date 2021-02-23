package goforit

//VMDriver An interface for implementing driver
//for specific VM/scripting implementation
type VMDriver interface {

	//Get getting VM implementation of this driver
	Get() (VM, error)

	//SetTemplate set template VM, this VM will be used for cloning
	//another VM to reduce constructing overhead
	SetTemplate(vm VM)

	//ExtractFunctionListFromFormulaString extracting functions names
	//from the given script/formular so that the function can be loaded
	//by Formula.LoadContext() before being executed, otherwise the
	//unloaded functions will not be known to the scripting/VM engine
	ExtractFunctionListFromFormulaString(formulaStr string) []string
}

//VM an interface acting as an abstract layer of VM implementation
//so most parts of goforit will interact with the abstract VM
//instead of the actual implemnetation
//
//The actual implementation will be wrapped up by driver
//Ex. OttoVM is a driver for otto
type VM interface {

	//Run for running a given script/formual
	Run(formulaString string) (JSValue, error)

	//Get getting value of a variable out of scripting context
	Get(varname string) (JSValue, error)

	//Set setting value of a valiable inside scripting context
	Set(varname string, value interface{}) error

	//ValidateFuncArguments validting if a given scripting function has
	//number of arguments equal to the given cnt
	ValidateFuncArguments(funcName string, cnt int, funcDef interface{})

	//GetFuncArgsCount getting function arguments count
	GetFuncArgsCount(funcDef interface{}) int

	//GetFuncArgAsFloat getting function argument refering by the index as float64
	GetFuncArgAsFloat(funcDef interface{}, index int) float64

	//GetFuncArgAsInt getting function argument refering by the index as int64
	GetFuncArgAsInt(funcDef interface{}, index int) int64

	//GetFuncArgAsString getting function argument refering by the index as string
	GetFuncArgAsString(funcDef interface{}, index int) string

	//GetFuncArgAsBoolean getting function argument refering by the index as boolean
	GetFuncArgAsBoolean(funcDef interface{}, index int) bool

	//GetFuncArgAsIs getting function argument refering by the index as raw type
	//depending on each scripting/VM engine e.g., otto.Value for otto
	GetFuncArgAsIs(funcDef interface{}, index int) interface{}

	//ToVMValue converting go variable into scriing/VM value e.g., otto.Value for otto
	ToVMValue(goValue interface{}) interface{}

	//GetBuiltInFunc getting a function for registering / connecting
	//a built-in function referred by the given name to the scripting/VM
	//
	//An abstract implementation of built-in function will be provided
	//by the given BuiltInFunctions, the native VM has to execute the built-in
	//function thtough them
	//
	//Below is an example of the implemenation for otto
	//
	//		var vm otto.Otto = ...
	//
	// 		for _, f := range funcs {
	//
	// 			if f.Has(funcName) {
	//
	// 				return func(context *FormulaContext) {
	//
	// 					vm.Set(funcName, func(call otto.FunctionCall) otto.Value {
	//
	// 						for _, f := range funcs {
	// 							result, found := f.Execute(funcName, v, call)
	//
	// 							if found {
	// 								return v.ToVMValue(result).(otto.Value)
	// 							}
	// 						}
	//
	// 						panic(errors.New("Runtime error: Function not found - " + funcName))
	// 					})
	// 				}
	// 			}
	// 		}
	//
	// 		return nil
	//
	GetBuiltInFunc(funcName string, funcs []BuiltInFunctions) func(context *FormulaContext)

	//IsDefined check if the current JS value is defined
	IsDefined(vmValue interface{}) bool

	//IsUndefined check if the current JS value is undefined
	IsUndefined(vmValue interface{}) bool

	//IsNull check if the current JS value is null
	IsNull(vmValue interface{}) bool

	//IsPrimitive check if the current JS value is primitive
	IsPrimitive(vmValue interface{}) bool

	//IsBoolean check if the current JS value is boolean
	IsBoolean(vmValue interface{}) bool

	//IsNumber check if the current JS value is number
	IsNumber(vmValue interface{}) bool

	//IsNaN check if the current JS value is nan
	IsNaN(vmValue interface{}) bool

	//IsString check if the current JS value is string
	IsString(vmValue interface{}) bool

	//IsObject check if the current JS value is object
	IsObject(vmValue interface{}) bool

	//IsFunction check if the current JS value is function
	IsFunction(vmValue interface{}) bool

	//ToBoolean get value as bool
	ToBoolean(vmValue interface{}) (bool, error)

	//ToFloat get value as float
	ToFloat(vmValue interface{}) (float64, error)

	//ToInteger get value as integer
	ToInteger(vmValue interface{}) (int64, error)

	//ToString get value as string
	ToString(vmValue interface{}) (string, error)

	//Export js value to Go value type as specified in a given interface
	Export(vmValue interface{}) (interface{}, error)
}

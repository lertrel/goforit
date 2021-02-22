package goforit

//VMDriver An interface for implementing driver
//for specific VM/scripting implementation
type VMDriver interface {
	Get() (VM, error)

	SetTemplate(vm VM)

	ExtractFunctionListFromFormulaString(formulaStr string) []string
}

type VM interface {
	Run(formulaString string) (JSValue, error)
	Get(varname string) (JSValue, error)
	Set(varname string, value interface{}) error
	ValidateFuncArguments(funcName string, cnt int, funcDef interface{})
	GetFuncArgsCount(funcDef interface{}) int
	GetFuncArgAsFloat(funcDef interface{}, index int) float64
	GetFuncArgAsInt(funcDef interface{}, index int) int64
	GetFuncArgAsString(funcDef interface{}, index int) string
	GetFuncArgAsBoolean(funcDef interface{}, index int) bool
	GetFuncArgAsIs(funcDef interface{}, index int) interface{}
	ToVMValue(goValue interface{}) interface{}
	GetBuiltInFunc(funcName string) func(context *FormulaContext)

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

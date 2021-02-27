package model

//Value is a placeholder for a VM/script value
type Value interface {
	//IsDefined check if the current JS value is defined
	IsDefined() bool
	//IsUndefined check if the current JS value is undefined
	IsUndefined() bool
	//IsNull check if the current JS value is null
	IsNull() bool
	//IsPrimitive check if the current JS value is primitive
	IsPrimitive() bool
	//IsBoolean check if the current JS value is boolean
	IsBoolean() bool
	//IsNumber check if the current JS value is number
	IsNumber() bool
	//IsNaN check if the current JS value is nan
	IsNaN() bool
	//IsString check if the current JS value is string
	IsString() bool
	//IsObject check if the current JS value is object
	IsObject() bool
	//IsFunction check if the current JS value is function
	IsFunction() bool
	//ToBoolean get value as bool
	ToBoolean() (bool, error)
	//ToFloat get value as float
	ToFloat() (float64, error)
	//ToInteger get value as integer
	ToInteger() (int64, error)
	//ToString get value as string
	ToString() (string, error)
	//Export js value to Go value type
	Export() (interface{}, error)
}

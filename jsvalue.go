package goforit

import "github.com/robertkrimen/otto"

//JSValue a placeholder for a java script value
type JSValue struct {
	impl otto.Value
}

// func (this JSValue) ToValue(value interface{}) (JSValue, error) {

// 	impl, err := this.impl.ToValue(value)

// 	return JSValue{impl: impl}, err
// }

//IsDefined check if the current JS value is defined
func (j JSValue) IsDefined() bool {
	return j.impl.IsDefined()
}

//IsUndefined check if the current JS value is undefined
func (j JSValue) IsUndefined() bool {
	return j.impl.IsUndefined()
}

//IsNull check if the current JS value is null
func (j JSValue) IsNull() bool {
	return j.impl.IsNull()
}

//IsPrimitive check if the current JS value is primitive
func (j JSValue) IsPrimitive() bool {
	return j.impl.IsPrimitive()
}

//IsBoolean check if the current JS value is boolean
func (j JSValue) IsBoolean() bool {
	return j.impl.IsBoolean()
}

//IsNumber check if the current JS value is number
func (j JSValue) IsNumber() bool {
	return j.impl.IsNumber()
}

//IsNaN check if the current JS value is nan
func (j JSValue) IsNaN() bool {
	return j.impl.IsNaN()
}

//IsString check if the current JS value is string
func (j JSValue) IsString() bool {
	return j.impl.IsString()
}

//IsObject check if the current JS value is object
func (j JSValue) IsObject() bool {
	return j.impl.IsObject()
}

//IsFunction check if the current JS value is function
func (j JSValue) IsFunction() bool {
	return j.impl.IsFunction()
}

//ToBoolean get value as bool
func (j JSValue) ToBoolean() (bool, error) {
	return j.impl.ToBoolean()
}

//ToFloat get value as float
func (j JSValue) ToFloat() (float64, error) {
	return j.impl.ToFloat()
}

//ToInteger get value as integer
func (j JSValue) ToInteger() (int64, error) {
	return j.impl.ToInteger()
}

//ToString get value as string
func (j JSValue) ToString() (string, error) {
	return j.impl.ToString()
}

//Export js value to Go value type as specified in a given interface
func (j JSValue) Export() (interface{}, error) {
	return j.impl.Export()
}

package goforit

//JSValue a placeholder for a java script value
type JSValue struct {
	// impl otto.Value
	vm   VM
	impl interface{}
}

//NewJSValue create a new JSValue
func NewJSValue(vm VM, impl interface{}) JSValue {

	return JSValue{vm: vm, impl: impl}
}

// func (this JSValue) ToValue(value interface{}) (JSValue, error) {

// 	impl, err := this.impl.ToValue(value)

// 	return JSValue{impl: impl}, err
// }

//IsDefined check if the current JS value is defined
func (j JSValue) IsDefined() bool {
	// return j.impl.IsDefined()
	return j.vm.IsDefined(j.impl)
}

//IsUndefined check if the current JS value is undefined
func (j JSValue) IsUndefined() bool {
	// return j.impl.IsUndefined()
	return j.vm.IsUndefined(j.impl)
}

//IsNull check if the current JS value is null
func (j JSValue) IsNull() bool {
	// return j.impl.IsNull()
	return j.vm.IsNull(j.impl)
}

//IsPrimitive check if the current JS value is primitive
func (j JSValue) IsPrimitive() bool {
	// return j.impl.IsPrimitive()
	return j.vm.IsPrimitive(j.impl)
}

//IsBoolean check if the current JS value is boolean
func (j JSValue) IsBoolean() bool {
	// return j.impl.IsBoolean()
	return j.vm.IsBoolean(j.impl)
}

//IsNumber check if the current JS value is number
func (j JSValue) IsNumber() bool {
	// return j.impl.IsNumber()
	return j.vm.IsNumber(j.impl)
}

//IsNaN check if the current JS value is nan
func (j JSValue) IsNaN() bool {
	// return j.impl.IsNaN()
	return j.vm.IsNaN(j.impl)
}

//IsString check if the current JS value is string
func (j JSValue) IsString() bool {
	// return j.impl.IsString()
	return j.vm.IsString(j.impl)
}

//IsObject check if the current JS value is object
func (j JSValue) IsObject() bool {
	// return j.impl.IsObject()
	return j.vm.IsObject(j.impl)
}

//IsFunction check if the current JS value is function
func (j JSValue) IsFunction() bool {
	// return j.impl.IsFunction()
	return j.vm.IsFunction(j.impl)
}

//ToBoolean get value as bool
func (j JSValue) ToBoolean() (bool, error) {
	// return j.impl.ToBoolean()
	return j.vm.ToBoolean(j.impl)
}

//ToFloat get value as float
func (j JSValue) ToFloat() (float64, error) {
	// return j.impl.ToFloat()
	return j.vm.ToFloat(j.impl)
}

//ToInteger get value as integer
func (j JSValue) ToInteger() (int64, error) {
	// return j.impl.ToInteger()
	return j.vm.ToInteger(j.impl)
}

//ToString get value as string
func (j JSValue) ToString() (string, error) {
	// return j.impl.ToString()
	return j.vm.ToString(j.impl)
}

//Export js value to Go value type
func (j JSValue) Export() (interface{}, error) {
	// return j.impl.Export()
	return j.vm.Export(j.impl)
}

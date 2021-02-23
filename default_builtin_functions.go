package goforit

import (
	"errors"
	"math"
	"math/big"
	"strconv"
)

//DefaultBuiltInFunctions providing built-in functions shipped with goforit
type DefaultBuiltInFunctions struct {
}

//Has to check if the given function name is supported
//by the current BuiltInFunctions
func (fs DefaultBuiltInFunctions) Has(funcName string) bool {

	switch funcName {

	case "$ABS":
		return true
	case "$RND":
		return true
	case "$FLOOR", "$FLR":
		return true
	case "$CEIL":
		return true
	case "$IF":
		return true
	case "$SUMI":
		return true
	case "$SUMF":
		return true
	case "$AVG":
		return true
	case "$MIN":
		return true
	case "$MAX":
		return true
	default:
		return false

	}

}

//Execute to execute a built-in function as per the given function name
func (fs DefaultBuiltInFunctions) Execute(funcName string, vm VM, funcDef interface{}) (interface{}, bool) {

	switch funcName {

	case "$ABS":
		return fAbs(vm, funcDef), true
	case "$RND":
		return fRnd(vm, funcDef), true
	case "$FLOOR", "$FLR":
		return fFloor(vm, funcDef), true
	case "$CEIL":
		return fCeil(vm, funcDef), true
	case "$IF":
		return fIf(vm, funcDef), true
	case "$SUMI":
		return fSumi(vm, funcDef), true
	case "$SUMF":
		return fSumf(vm, funcDef), true
	case "$AVG":
		return fAvg(vm, funcDef), true
	case "$MIN":
		return fMin(vm, funcDef), true
	case "$MAX":
		return fMax(vm, funcDef), true
	default:
		return nil, false

	}

}

func fAbs(vm VM, funcDef interface{}) interface{} {

	vm.ValidateFuncArguments("$ABS", 1, funcDef)

	v := vm.GetFuncArgAsFloat(funcDef, 0)
	result := math.Abs(v)

	return vm.ToVMValue(result)
}

func fRnd(vm VM, funcDef interface{}) interface{} {

	vm.ValidateFuncArguments("$ABS", 2, funcDef)

	v := vm.GetFuncArgAsFloat(funcDef, 0)
	p := vm.GetFuncArgAsInt(funcDef, 1)
	if 0 > p || p > 10 {
		panic(errors.New("$RND - second parameter should be between 0 and 10"))
	}

	result := v
	f := math.Pow10(int(p))
	result = math.Round(v*f) / f

	return vm.ToVMValue(result)
}

func fFloor(vm VM, funcDef interface{}) interface{} {

	vm.ValidateFuncArguments("$FLOOR", 2, funcDef)

	v := vm.GetFuncArgAsFloat(funcDef, 0)
	p := vm.GetFuncArgAsInt(funcDef, 1)
	if 0 > p || p > 10 {
		panic(errors.New("$RND - second parameter should be between 0 and 10"))
	}

	result := v
	f := math.Pow10(int(p))
	result = math.Floor(v*f) / f

	return vm.ToVMValue(result)
}

func fCeil(vm VM, funcDef interface{}) interface{} {

	vm.ValidateFuncArguments("$CEIL", 2, funcDef)

	v := vm.GetFuncArgAsFloat(funcDef, 0)
	p := vm.GetFuncArgAsInt(funcDef, 1)
	if 0 > p || p > 10 {
		panic(errors.New("$RND - second parameter should be between 0 and 10"))
	}

	result := v
	f := math.Pow10(int(p))
	result = math.Ceil(v*f) / f

	return vm.ToVMValue(result)
}

func fIf(vm VM, funcDef interface{}) interface{} {

	vm.ValidateFuncArguments("$IF", 3, funcDef)

	b := vm.GetFuncArgAsBoolean(funcDef, 0)
	var jsResult interface{}

	if b {
		jsResult = vm.GetFuncArgAsIs(funcDef, 1)
	} else {
		jsResult = vm.GetFuncArgAsIs(funcDef, 2)
	}

	return jsResult
}

func fSumi(vm VM, funcDef interface{}) interface{} {

	var result int64
	result = 0

	for i := 0; i < vm.GetFuncArgsCount(funcDef); i++ {

		v := vm.GetFuncArgAsInt(funcDef, i)
		result = result + v

	}

	return vm.ToVMValue(result)
}

func fSumf(vm VM, funcDef interface{}) interface{} {

	var result float64
	result = 0.0

	for i := 0; i < vm.GetFuncArgsCount(funcDef); i++ {

		v := vm.GetFuncArgAsFloat(funcDef, i)
		result = result + v

	}

	return vm.ToVMValue(result)
}

func fAvg(vm VM, funcDef interface{}) interface{} {

	result := big.NewFloat(0.0)
	var i int

	for i = 0; i < vm.GetFuncArgsCount(funcDef); i++ {

		v := vm.GetFuncArgAsFloat(funcDef, i)
		bigV := big.NewFloat(v)
		result.Add(result, bigV)
	}

	result.Quo(result, big.NewFloat(float64(i)))
	f, err := strconv.ParseFloat(result.String(), 10)
	if err != nil {
		panic(err)
	}

	return vm.ToVMValue(f)
}

func fMin(vm VM, funcDef interface{}) interface{} {

	hasLast := false
	var last interface{} = nil

	for i := 0; i < vm.GetFuncArgsCount(funcDef); i++ {

		v := vm.GetFuncArgAsIs(funcDef, i)

		if !hasLast {
			last = v
			hasLast = true
		} else {
			l, el := vm.ToFloat(last)
			r, er := vm.ToFloat(v)

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
}

func fMax(vm VM, funcDef interface{}) interface{} {

	hasLast := false
	var last interface{} = nil

	for i := 0; i < vm.GetFuncArgsCount(funcDef); i++ {

		v := vm.GetFuncArgAsIs(funcDef, i)

		if !hasLast {
			last = v
			hasLast = true
		} else {
			l, el := vm.ToFloat(last)
			r, er := vm.ToFloat(v)

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
}

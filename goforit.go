package goforit

import (
	"github.com/lertrel/goforit/builder"
	"github.com/lertrel/goforit/trigger"
	"github.com/lertrel/goforit/util"
)

var debugFlag = false

//Get An entry point to obtain Formula
// func Get() Formula {

// 	r, _ := regexp.Compile("(\\$[^\\$()\\s]+)\\(")

// 	return Formula{r: r, customFuncs: make(map[string]string), Debug: false}
// }
func debug(b bool, format string, args ...interface{}) {
	// if b {
	// 	log.Printf(format, args...)
	// }
	util.Debug(b, format, args...)
}

//NewFormulaBuilder To obtain formula builder
//
//Ex.
//
//		builder := NewFormulaBuilder()
//		formula := builder.Get()
//
func NewFormulaBuilder() builder.FormulaBuilder {

	return builder.NewFormulaBuilder()
	// return FormulaBuilder{
	// 	Debug:  false,
	// 	repos:  make(map[vm.CustomFunctionRepository]vm.CustomFunctionRepository),
	// 	Driver: nil,
	// }
}

//NewTriggersBuilder to get TriggersBuilder
func NewTriggersBuilder() trigger.TriggersBuilder {
	return trigger.TriggersBuilder{}
}

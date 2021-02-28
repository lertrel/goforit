package goforit

import (
	"github.com/lertrel/goforit/builder"
	"github.com/lertrel/goforit/trigger"
)

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

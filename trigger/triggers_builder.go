package trigger

import (
	"errors"

	"github.com/lertrel/goforit/model"
)

//TriggersBuilder a builder for Formulas
type TriggersBuilder struct {
	formula            model.Formula
	isFormulaSet       bool
	triggerLookup      Lookup
	isTriggerLookupSet bool
	formulaLookup      FormulaLookup
	isFormulaLookupSet bool
}

//SetFormula setting FormulaBuilder
func (b TriggersBuilder) SetFormula(formula model.Formula) TriggersBuilder {

	return TriggersBuilder{
		formula:            formula,
		isFormulaSet:       true,
		triggerLookup:      b.triggerLookup,
		isTriggerLookupSet: b.isTriggerLookupSet,
		formulaLookup:      b.formulaLookup,
		isFormulaLookupSet: b.isFormulaLookupSet,
	}
}

//SetTriggerLookup setting TriggerLookup
func (b TriggersBuilder) SetTriggerLookup(lookup Lookup) TriggersBuilder {

	return TriggersBuilder{
		formula:            b.formula,
		isFormulaSet:       b.isFormulaSet,
		triggerLookup:      lookup,
		isTriggerLookupSet: true,
		formulaLookup:      b.formulaLookup,
		isFormulaLookupSet: b.isFormulaLookupSet,
	}
}

//SetFormulaLookup setting FormulaLookup
func (b TriggersBuilder) SetFormulaLookup(lookup FormulaLookup) TriggersBuilder {

	return TriggersBuilder{
		formula:            b.formula,
		isFormulaSet:       b.isFormulaSet,
		triggerLookup:      b.triggerLookup,
		isTriggerLookupSet: b.isTriggerLookupSet,
		formulaLookup:      lookup,
		isFormulaLookupSet: true,
	}
}

//Get get Triggers
func (b TriggersBuilder) Get() Triggers {

	fb := b.formula
	if !b.isFormulaSet {
		panic(errors.New("A Formula is yet to be defined"))
	}

	if !b.isTriggerLookupSet {
		panic(errors.New("A TriggerLookup is yet to be defined"))
	}

	if !b.isFormulaLookupSet {
		panic(errors.New("A FormulaLookup is yet to be defined"))
	}

	return SimpleTriggers{
		triggerLookup: b.triggerLookup,
		formulaLookup: b.formulaLookup,
		formula:       fb,
	}
}

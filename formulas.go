package goforit

import "fmt"

//Triggers layer to help executing external formula related to
//a pre-defined trigger point
type Triggers interface {
	Execute(trigger string, context map[string]string) (map[string]JSValue, error)
}

//TriggersBuilder a builder for Formulas
type TriggersBuilder struct {
	formulaBuilder      FormulaBuilder
	isFormulaBuilderSet bool
	triggerLookup       TriggerLookup
	isTriggerLookupSet  bool
	formulaLookup       FormulaLookup
	isFormulaLookupSet  bool
}

func (b TriggersBuilder) SetFormulaBuilder(builder FormulaBuilder) TriggersBuilder {

	return TriggersBuilder{
		formulaBuilder:      builder,
		isFormulaBuilderSet: true,
		triggerLookup:       b.triggerLookup,
		isTriggerLookupSet:  b.isTriggerLookupSet,
		formulaLookup:       b.formulaLookup,
		isFormulaLookupSet:  b.isFormulaLookupSet,
	}
}

func (b TriggersBuilder) SetTriggerLookup(lookup TriggerLookup) TriggersBuilder {

	return TriggersBuilder{
		formulaBuilder:      b.formulaBuilder,
		isFormulaBuilderSet: b.isFormulaBuilderSet,
		triggerLookup:       lookup,
		isTriggerLookupSet:  true,
		formulaLookup:       b.formulaLookup,
		isFormulaLookupSet:  b.isFormulaLookupSet,
	}
}

func (b TriggersBuilder) SetFormulaLookup(lookup FormulaLookup) TriggersBuilder {

	return TriggersBuilder{
		formulaBuilder:      b.formulaBuilder,
		isFormulaBuilderSet: b.isFormulaBuilderSet,
		triggerLookup:       b.triggerLookup,
		isTriggerLookupSet:  b.isTriggerLookupSet,
		formulaLookup:       lookup,
		isFormulaLookupSet:  true,
	}
}

func (b TriggersBuilder) Get() Triggers {

	fb := b.formulaBuilder
	if !b.isFormulaBuilderSet {
		fb = FormulaBuilder{}
		fmt.Println(fb)
	}

	panic("Not yet implemented")
}

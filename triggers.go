package goforit

import "errors"

//TriggersBuilder a builder for Formulas
type TriggersBuilder struct {
	formulaBuilder      FormulaBuilder
	isFormulaBuilderSet bool
	triggerLookup       TriggerLookup
	isTriggerLookupSet  bool
	formulaLookup       FormulaLookup
	isFormulaLookupSet  bool
}

//SetFormulaBuilder setting FormulaBuilder
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

//SetTriggerLookup setting TriggerLookup
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

//SetFormulaLookup setting FormulaLookup
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

//Get get Triggers
func (b TriggersBuilder) Get() Triggers {

	fb := b.formulaBuilder
	if !b.isFormulaBuilderSet {
		fb = NewFormulaBuilder()
		debug(true, "No specific FormulaBuilder is defined, using default")
	}

	if !b.isTriggerLookupSet {
		panic(errors.New("A TriggerLookup is yet to be defined"))
	}

	if !b.isFormulaLookupSet {
		panic(errors.New("A FormulaLookup is yet to be defined"))
	}

	return SimpleTriggers{
		formulaBuilder: fb,
		triggerLookup:  b.triggerLookup,
		formulaLookup:  b.formulaLookup,
		formula:        fb.Get(),
	}
}

//Triggers so-called a controller layer to help executing external formula
//related to a pre-defined trigger point
type Triggers interface {

	//Execute executing formulas for a given trigger point with the states
	//provided in context
	//
	//Triggers will find Trigger definition by looking up into TriggerLookup
	//
	//Then, will find formula(s) matches with the trigger definitions by looking up
	//FormulaLookup
	//
	//Then executing the matches formula(s) under the FormulaContext created by
	//FormulaBuilder
	//
	Execute(trigger string, context map[string]interface{}) (map[string]interface{}, error)
}

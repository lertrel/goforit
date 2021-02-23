package goforit

import "errors"

//GetTriggersBuilder to get TriggersBuilder
func GetTriggersBuilder() TriggersBuilder {
	return TriggersBuilder{}
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
		fb = GetFormulaBuilder()
		debug(true, "No specific FormulaBuilder is defined, using default")
	}

	if !b.isTriggerLookupSet {
		panic(errors.New("A TriggerLookup is yet to be defined"))
	}

	if !b.isFormulaLookupSet {
		panic(errors.New("A FormulaLookup is yet to be defined"))
	}

	return Triggers{
		formulaBuilder: fb,
		triggerLookup:  b.triggerLookup,
		formulaLookup:  b.formulaLookup,
	}
}

//Triggers so-called a controller layer to help executing external formula
//related to a pre-defined trigger point
type Triggers struct {
	formulaBuilder FormulaBuilder
	triggerLookup  TriggerLookup
	formulaLookup  FormulaLookup
	formula        Formula
}

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
func (t Triggers) Execute(trigger string, context map[string]string) (map[string]JSValue, error) {

	triggerDef, lookupTriggerErr := t.triggerLookup.GetTrigger(trigger)
	if lookupTriggerErr != nil {
		return nil, lookupTriggerErr
	}

	formulas, lookupFormulaErr := t.formulaLookup.GetFormulars(triggerDef, context)
	if lookupFormulaErr != nil {
		return nil, lookupFormulaErr
	}

	if !formulas.HasNext() {
		return nil, errors.New("No matched formula found for trigger - " + trigger)
	}

	formulaDef := formulas.Next()

	f, _ := t.getFormula(triggerDef)
	fc, parseScriptErr := f.LoadContext(nil, formulaDef.Body)
	if parseScriptErr != nil {
		return nil, parseScriptErr
	}

	t.mapInputs(*fc, context, formulaDef)
	ret, runtimeErr := fc.Run(formulaDef.Body)
	if runtimeErr != nil {
		return nil, runtimeErr
	}

	arr := make(map[string]JSValue)

	t.mapOutputs(*fc, arr, formulaDef)
	arr["return"] = ret

	return arr, nil
}

func (t Triggers) getFormula(trigger Trigger) (Formula, error) {

	panic("Not yet implemented")
}

func (t Triggers) mapInputs(f FormulaContext, context map[string]string, c FormulaConfig) {

	panic("Not yet implemented")
}

func (t Triggers) mapOutputs(f FormulaContext, result map[string]JSValue, c FormulaConfig) {

	panic("Not yet implemented")
}

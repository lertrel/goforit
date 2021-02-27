package goforit

import (
	"errors"

	"github.com/lertrel/goforit/model"
)

//SimpleTriggers simple implementation of Triggers
type SimpleTriggers struct {
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
func (t SimpleTriggers) Execute(trigger string, context map[string]interface{}) (result map[string]interface{}, err error) {

	//Obtaining trigger definition of the given trigger ID
	triggerDef, err := t.triggerLookup.GetTrigger(trigger)
	if err != nil {
		return
	}

	//Searching matched formula definition
	formulas, err := t.formulaLookup.GetFormulars(triggerDef, context)
	if err != nil {
		return
	}

	//If no formula is matched, return with error
	if !formulas.HasNext() {
		return nil, errors.New("No matched formula found for trigger - " + trigger)
	}

	//Obtaining formula definition
	formulaDef := formulas.Next()

	//Obtaining formula engine
	f, _ := t.getFormula(triggerDef)
	//Creating a new formula context to run this formula
	fc, err := f.LoadContext(nil, formulaDef.Body)
	if err != nil {
		return
	}

	//Mapping context states as formula inputs
	//by following input mapping rule provided
	//by the trigger definition
	if err = t.mapInputs(&fc, context, triggerDef); err != nil {
		return
	}

	//Running the formula, and obtaining result
	jsRet, err := fc.Run(formulaDef.Body)
	if err != nil {
		return
	}

	//Mapping script variables (in VM) as outputs
	//by following input mapping rule provided
	//by the trigger definition
	result, err = t.mapOutputs(&fc, triggerDef)
	if err != nil {
		return
	}

	result["_return"], err = jsRet.Export()

	return
}

func (t SimpleTriggers) getFormula(trigger Trigger) (Formula, error) {

	return t.formula, nil
}

func (t SimpleTriggers) mapInputs(f *model.FormulaContext, context map[string]interface{}, c Trigger) (err error) {

	if c.InputMapping == "" {
		return
	}

	if err = (*f).Set(c.ContextVarName, context); err != nil {
		return
	}

	*f, err = t.formula.LoadContext(*f, c.InputMapping)
	if err != nil {
		return
	}

	//Falls through
	_, err = (*f).Run(c.InputMapping)

	return
}

func (t SimpleTriggers) mapOutputs(f *model.FormulaContext, c Trigger) (result map[string]interface{}, err error) {

	result = make(map[string]interface{})

	if c.OuputMapping == "" {
		return
	}

	if err = f.Set(c.OutputVarName, result); err != nil {
		return
	}

	f, err = t.formula.LoadContext(f, c.OuputMapping)
	if err != nil {
		return
	}

	//Falls through
	_, err = f.Run(c.OuputMapping)

	return
}

package trigger

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

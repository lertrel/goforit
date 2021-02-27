package trigger

//FormulaLookup a formula config repository
//
//An underlying mechanism for formula config to be stored
//is varied from implementation to implementation.
//E.g., RDB, in-memory, file, API, etc.
type FormulaLookup interface {

	//GetFormula getting FormulaConfig by name
	GetFormula(id string) (FormulaConfig, error)

	//GetAllFormulas getting all FormulaConfig(s)
	GetAllFormulas() (FormulaIterator, error)

	//GetAllFormulas search all FormulaConfig(s) that matches the given Trigger
	GetFormulars(trigger Trigger, context map[string]interface{}) (FormulaIterator, error)
}

//FormulaIterator a interator of FormulaConfig
type FormulaIterator interface {
	HasNext() bool
	Next() FormulaConfig
}

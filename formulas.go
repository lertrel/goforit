package goforit

//Formulas layer to help executing external formula related to
//a pre-defined trigger point
type Formulas interface {
	Execute(trigger string, context map[string]string) (map[string]JSValue, error)
}

//FormulasBuilder a builder for Formulas
type FormulasBuilder struct {
	formularBuilder FormulaBuilder
}

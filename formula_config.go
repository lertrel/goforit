package goforit

//FormulaConfig representing (probably persisted) formula
//and its custom attributes (for matching with triggers)
type FormulaConfig struct {
	ID          string
	Description string
	Body        string
	Attributes  map[string]string
	Enabled     bool
	// Inputs      []string
	// Outputs     []string
}

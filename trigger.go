package goforit

//Trigger represting a trigger setup
type Trigger struct {
	ID          string
	Description string
	Filter      string
	// Orders      map[int]string
	// DynamicOrders        string
	// DynamicOrdersVarName string
	ContextVarName string
	OutputVarName  string
	InputMapping   string
	OuputMapping   string
}

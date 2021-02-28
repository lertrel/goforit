package parse

//Parser is a script parser
type Parser interface {

	//ExtractFunctionNames extracting functions names
	//from the given script/formular so that the function can be loaded
	//by Formula.LoadContext() before being executed, otherwise the
	//unloaded functions will not be known to the scripting/VM engine
	ExtractFunctionNames(formulaStr string) []string
}

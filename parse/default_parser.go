package parse

import "regexp"

//New will return a default implementation of Parser
func New() Parser {

	r, _ := regexp.Compile("(\\$[^\\$()\\s]+)\\(")

	return DefaultParser{r}
}

//DefaultParser is a default implementation of Parser
type DefaultParser struct {
	r *regexp.Regexp
}

//ExtractFunctionNames extracting functions names
//from the given script/formular so that the function can be loaded
//by Formula.LoadContext() before being executed, otherwise the
//unloaded functions will not be known to the scripting/VM engine
func (p DefaultParser) ExtractFunctionNames(formulaStr string) []string {

	matches := p.r.FindAllStringSubmatch(formulaStr, -1)
	dedupMatches := make(map[string]bool)

	for i := 0; i < len(matches); i++ {
		dedupMatches[matches[i][1]] = true
	}

	funArr := make([]string, len(dedupMatches))

	i := 0
	for k := range dedupMatches {
		funArr[i] = k
		i++
	}

	return funArr
}

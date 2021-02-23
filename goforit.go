package goforit

import (
	"log"
)

var debugFlag = false

//Get An entry point to obtain Formula
// func Get() Formula {

// 	r, _ := regexp.Compile("(\\$[^\\$()\\s]+)\\(")

// 	return Formula{r: r, customFuncs: make(map[string]string), Debug: false}
// }
func debug(b bool, format string, args ...interface{}) {
	if b {
		log.Printf(format, args...)
	}
}

//GetFormularBuilder To obtain formula builder
//
//Ex.
//
//		builder := GetFormularBuilder()
//		formula := builder.Get()
//
func GetFormularBuilder() FormulaBuilder {

	return FormulaBuilder{
		Debug:  false,
		repos:  make(map[CustomFunctionRepository]CustomFunctionRepository),
		Driver: nil,
	}
}

//FormulaBuilder a formula builder
type FormulaBuilder struct {
	Debug  bool
	repos  map[CustomFunctionRepository]CustomFunctionRepository
	funcs  map[BuiltInFunctions]BuiltInFunctions
	Driver VMDriver
}

//SetDebug setting debug flag (if yes log wll be printed)
func (b FormulaBuilder) SetDebug(debug bool) FormulaBuilder {

	return FormulaBuilder{
		Debug:  debug,
		repos:  b.repos,
		funcs:  b.funcs,
		Driver: b.Driver,
	}
}

//AddCustomFunctionRepository adding a custom function repository
//before getting a new formula, so the custom functions will be also
//being looked up in the given repository if it's not found in the
//default repository.
//
//With this, it's allow client to have freedom to store and provide
//customer functions from various sources e.g., DB, files, etc.
func (b FormulaBuilder) AddCustomFunctionRepository(repo CustomFunctionRepository) FormulaBuilder {

	b2 := FormulaBuilder{
		Debug:  b.Debug,
		repos:  b.repos,
		funcs:  b.funcs,
		Driver: b.Driver,
	}

	_, found := b2.repos[repo]

	if !found {
		b2.repos[repo] = repo
	}

	return b2
}

//AddBuiltInFunctions adding a BuiltInFunctions implementation
//before getting a new formula, so the BuiltInFunction will be also
//being looked up if a desirable function is not found in the default
//repository.
//
//With this, it's allow client to provide additional built-in functions
//customer functions from various sources e.g., DB, files, etc.
//
//*NOTE* that built-in functions are different from custom functions that
//the built-in functions are implemented in Go and loaded as static packages
//together with client program, so not need for VM to intrepret this fucntions
//so basically it should run faster than the custom functions
func (b FormulaBuilder) AddBuiltInFunctions(funcs BuiltInFunctions) FormulaBuilder {

	b2 := FormulaBuilder{
		Debug:  b.Debug,
		repos:  b.repos,
		funcs:  b.funcs,
		Driver: b.Driver,
	}

	_, found := b2.funcs[funcs]

	if !found {
		b2.funcs[funcs] = funcs
	}

	return b2
}

//SetDriver allow client to use another VM rather than the default one
func (b FormulaBuilder) SetDriver(driver VMDriver) FormulaBuilder {

	return FormulaBuilder{
		Debug:  b.Debug,
		repos:  b.repos,
		funcs:  b.funcs,
		Driver: driver,
	}
}

//Get to obtain a new Formula
func (b FormulaBuilder) Get() Formula {

	repos := make([]CustomFunctionRepository, len(b.repos)+1)

	i := 0
	repos[i] = DefaultCustomFunctionRepository{customFuncs: make(map[string]string)}

	for r := range b.repos {
		i++
		repos[i] = r
	}

	funcs := make([]BuiltInFunctions, len(b.funcs)+1)

	i = 0
	funcs[i] = DefaultBuiltInFunctions{}

	for r := range b.funcs {
		i++
		funcs[i] = r
	}

	var driver VMDriver

	if b.Driver != nil {
		driver = b.Driver
	} else {
		// driver = GetVMDriver(funcs)
		driver = GetVMDriver()
	}

	return Formula{
		driver:       driver,
		customFuncs:  repos,
		builtInFuncs: funcs,
		Debug:        b.Debug,
	}
}

// func (b FormulaBuilder) Get() Formula {

// 	r, _ := regexp.Compile("(\\$[^\\$()\\s]+)\\(")

// 	return Formula{r: r, customFuncs: b.repo, Debug: b.debug}
// }

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

func GetFormularBuilder() FormulaBuilder {

	return FormulaBuilder{
		Debug:  false,
		repos:  make(map[CustomFunctionRepository]CustomFunctionRepository),
		Driver: nil,
	}
}

type FormulaBuilder struct {
	Debug  bool
	repos  map[CustomFunctionRepository]CustomFunctionRepository
	funcs  map[BuiltInFunctions]BuiltInFunctions
	Driver VMDriver
}

func (b FormulaBuilder) SetDebug(debug bool) FormulaBuilder {

	return FormulaBuilder{
		Debug:  debug,
		repos:  b.repos,
		funcs:  b.funcs,
		Driver: b.Driver,
	}
}

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

func (b FormulaBuilder) SetDriver(driver VMDriver) FormulaBuilder {

	return FormulaBuilder{
		Debug:  b.Debug,
		repos:  b.repos,
		funcs:  b.funcs,
		Driver: driver,
	}
}

func (b FormulaBuilder) Get() Formula {

	repos := make(map[int]CustomFunctionRepository)

	i := 0
	repos[i] = DefaultCustomFunctionRepository{customFuncs: make(map[string]string)}

	for r := range b.repos {
		i++
		repos[i] = r
	}

	funcs := make(map[int]BuiltInFunctions)

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
		driver = GetVMDriver(funcs)
	}

	return Formula{
		driver:      driver,
		customFuncs: repos,
		Debug:       b.Debug,
	}
}

// func (b FormulaBuilder) Get() Formula {

// 	r, _ := regexp.Compile("(\\$[^\\$()\\s]+)\\(")

// 	return Formula{r: r, customFuncs: b.repo, Debug: b.debug}
// }

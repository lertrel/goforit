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
		debug:  false,
		repos:  make(map[int]CustomFunctionRepository),
		driver: nil,
	}
}

type FormulaBuilder struct {
	debug  bool
	repos  map[CustomFunctionRepository]CustomFunctionRepository
	driver VMDriver
}

func (b *FormulaBuilder) SetDebug(debug bool) FormulaBuilder {

	return FormulaBuilder{
		debug:  debug,
		repos:  b.repos,
		driver: b.driver,
	}
}

func (b FormulaBuilder) AetCustomFunctionRepository(repo CustomFunctionRepository) FormulaBuilder {

	_, found := b.repos[repo]

	if !found {
		b.repos[repo] = repo
	}

	return FormulaBuilder{
		debug:  b.debug,
		repos:  b.repos,
		driver: b.driver,
	}
}

func (b *FormulaBuilder) SetDriver(driver VMDriver) FormulaBuilder {

	return FormulaBuilder{
		debug:  b.debug,
		repos:  b.repos,
		driver: driver,
	}
}

func (b FormulaBuilder) Get() Formula {

	var driver VMDriver
	var repo CustomFunctionRepository

	if b.driver != nil {
		driver = b.driver
	} else {
		driver = &OttoVMDriver{}
	}

	if b.repos != nil {
		repo = b.repos
	} else {
		repo = DefaultCustomFunctionRepository{}
	}

	return Formula{driver: driver, customFuncs: repo, Debug: b.debug}
}

// func (b FormulaBuilder) Get() Formula {

// 	r, _ := regexp.Compile("(\\$[^\\$()\\s]+)\\(")

// 	return Formula{r: r, customFuncs: b.repo, Debug: b.debug}
// }

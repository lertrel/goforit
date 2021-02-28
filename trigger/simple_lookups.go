package trigger

import (
	"errors"
	"fmt"
	"sort"

	"github.com/lertrel/goforit/model"
)

//NewSimpleLookup returns a new simple (in-memory) trigger.Lookup
func NewSimpleLookup(triggerList []Trigger) Lookup {

	sort.Sort(byTriggerID(triggerList))

	return SimpleLookup{triggerList}
}

//NewSimpleFormulaLookup returns a new simple (in-memory) trigger.FormulaLookup
func NewSimpleFormulaLookup(configs []FormulaConfig, f model.Formula) FormulaLookup {

	sort.Sort(byConfigID(configs))

	return SimpleFormulaLookup{configs, f}
}

type byTriggerID []Trigger

func (a byTriggerID) Len() int           { return len(a) }
func (a byTriggerID) Less(i, j int) bool { return a[i].ID < a[j].ID }
func (a byTriggerID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type byConfigID []FormulaConfig

func (a byConfigID) Len() int           { return len(a) }
func (a byConfigID) Less(i, j int) bool { return a[i].ID < a[j].ID }
func (a byConfigID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

//SimpleLookup is a simple (in-memory) implementation of trigger.Lookup
type SimpleLookup struct {
	triggerList []Trigger
}

// GetTrigger getting Trigger by ID
func (tl SimpleLookup) GetTrigger(triggerName string) (Trigger, error) {

	size := len(tl.triggerList)

	index := sort.Search(size, func(i int) bool {
		return tl.triggerList[i].ID >= triggerName
	})

	var t Trigger
	if index < size && tl.triggerList[index].ID == triggerName {
		t = tl.triggerList[index]
	}

	return t, nil
}

//Triggers getting all Trigger(s)
func (tl SimpleLookup) Triggers() (Iterator, error) {

	return &SimpleIterator{0, tl.triggerList}, nil
}

//SimpleIterator is a simple (in-memory) implementation of trigger.Iterator
type SimpleIterator struct {
	index       int
	triggerList []Trigger
}

//HasNext tells if there's more element in the current iterator
func (si SimpleIterator) HasNext() bool {

	return si.index < len(si.triggerList)
}

//Next returns the next element in the current iterator
func (si *SimpleIterator) Next() Trigger {

	t := si.triggerList[si.index]
	si.incIndex()

	return t
}

func (si *SimpleIterator) incIndex() {
	si.index++
}

//SimpleFormulaLookup is a simple (in-memory) implementation of FormulaLookup
type SimpleFormulaLookup struct {
	configs []FormulaConfig
	formula model.Formula
}

//GetFormula getting FormulaConfig by name
func (fl SimpleFormulaLookup) GetFormula(id string) (FormulaConfig, error) {

	size := len(fl.configs)

	index := sort.Search(size, func(i int) bool {
		return fl.configs[i].ID >= id
	})

	var t FormulaConfig
	if index < size && fl.configs[index].ID == id {
		t = fl.configs[index]
	}

	return t, nil
}

//Formulas getting all FormulaConfig(s)
func (fl SimpleFormulaLookup) Formulas() (FormulaIterator, error) {
	return &SimpleFormulaIterator{0, fl.configs}, nil
}

//GetFormulars search all FormulaConfig(s) that matches the given Trigger
func (fl SimpleFormulaLookup) GetFormulars(trigger Trigger, context map[string]interface{}) (FormulaIterator, error) {

	if fl.formula == nil {
		return nil, errors.New("No model.Formula assigned for the current FormulaLookup")
	}

	filter := trigger.Filter
	fc, err := fl.formula.NewContext(filter)
	if err != nil {
		return nil, err
	}

	temp := make([]FormulaConfig, 0, len(fl.configs))

	for _, t := range fl.configs {

		if err = fc.Set("config", t); err != nil {
			return nil, err
		}

		if trigger.ContextVarName != "" {
			if err = fc.Set(trigger.ContextVarName, context); err != nil {
				return nil, err
			}
		}

		ret, err := fc.Run(filter)
		if err != nil {
			return nil, err
		}

		if !ret.IsBoolean() {
			return nil, fmt.Errorf("A result returned from %s is not a boolean", filter)
		}

		matched, err := ret.ToBoolean()
		if err != nil {
			return nil, err
		}

		if matched {
			temp = append(temp, t)
		}
	}

	return &SimpleFormulaIterator{0, temp}, nil
}

//SimpleFormulaIterator a iterator of FormulaConfig
type SimpleFormulaIterator struct {
	index   int
	configs []FormulaConfig
}

//HasNext tells if there's more element in the current iterator
func (fi SimpleFormulaIterator) HasNext() bool {

	return fi.index < len(fi.configs)
}

//Next returns the next element in the current iterator
func (fi *SimpleFormulaIterator) Next() FormulaConfig {

	c := fi.configs[fi.index]
	fi.incIndex()

	return c
}

func (fi *SimpleFormulaIterator) incIndex() {
	fi.index++
}

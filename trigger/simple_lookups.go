package trigger

import (
	"fmt"
	"sort"

	"github.com/lertrel/goforit/model"
)

//NewSimpleLookup returns a new simple (in-memory) trigger.Lookup
func NewSimpleLookup(triggerList []Trigger) Lookup {

	sort.Sort(byTriggerID(triggerList))

	return SimpleLookup{triggerList}
}

type byTriggerID []Trigger

func (a byTriggerID) Len() int           { return len(a) }
func (a byTriggerID) Less(i, j int) bool { return a[i].ID < a[j].ID }
func (a byTriggerID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

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

//GetTriggers search all Trigger(s) that matches the given filter (script)
func (tl SimpleLookup) GetTriggers(filter string) (Iterator, error) {
	var f model.Formula

	fc, err := f.NewContext(filter)
	if err != nil {
		return nil, err
	}

	temp := make([]Trigger, 0, len(tl.triggerList))

	for _, t := range tl.triggerList {

		if err = fc.Set("trigger", t); err != nil {
			return nil, err
		}

		ret, err := fc.Run(filter)
		if err != nil {
			return nil, err
		}

		if !ret.IsBoolean() {
			return nil, fmt.Errorf("A result returned from %s is not a boolean", filter)
		}

		//Falls through
		temp = append(temp, t)
	}

	return &SimpleIterator{0, temp}, nil
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

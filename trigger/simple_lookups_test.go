package trigger

import (
	"fmt"
	"testing"

	"github.com/lertrel/goforit/builder"
)

func triggerListForTest() []Trigger {

	return []Trigger{
		{ID: "Trigger 2"},
		{ID: "Trigger 4"},
		{ID: "Trigger 1"},
		{ID: "Trigger 3"},
	}
}

func TestSimpleIteratorIncIndex(t *testing.T) {

	i := SimpleIterator{0, triggerListForTest()}

	expected := 0
	actual := i.index
	if expected != actual {
		t.Errorf("i.index - expected %v but %v", expected, actual)
	}

	i.incIndex()
	expected, actual = 1, i.index
	if expected != actual {
		t.Errorf("i.index - expected %v but %v", expected, actual)
	}

	i.incIndex()
	expected, actual = 2, i.index
	if expected != actual {
		t.Errorf("i.index - expected %v but %v", expected, actual)
	}

}

func TestSimpleIteratorHasNextOK(t *testing.T) {

	i := SimpleIterator{0, triggerListForTest()}

	expected := true
	actual := i.HasNext()
	if expected != actual {
		t.Errorf("i.HasNext() - expected %v but %v", expected, actual)
	}

	i.index++
	expected, actual = true, i.HasNext()
	if expected != actual {
		t.Errorf("i.HasNext() - expected %v but %v", expected, actual)
	}

	i.index++
	expected, actual = true, i.HasNext()
	if expected != actual {
		t.Errorf("i.HasNext() - expected %v but %v", expected, actual)
	}

	i.index++
	expected, actual = true, i.HasNext()
	if expected != actual {
		t.Errorf("i.HasNext() - expected %v but %v", expected, actual)
	}

	i.index++
	expected, actual = false, i.HasNext()
	if expected != actual {
		t.Errorf("i.HasNext() - expected %v but %v", expected, actual)
	}

}

func TestSimpleIteratorHasNextNoElement(t *testing.T) {

	i := SimpleIterator{0, []Trigger{}}

	expected := false
	actual := i.HasNext()
	if expected != actual {
		t.Errorf("i.HasNext() - expected %v but %v", expected, actual)
	}

}

func TestSimpleIteratorNextNoElement(t *testing.T) {

	i := SimpleIterator{0, []Trigger{}}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("TestSimpleIteratorNextNoElement -", r)
		}
	}()

	_ = i.Next()

	t.Errorf("This line should not be executed")

}

func TestSimpleIteratorNextOK(t *testing.T) {

	i := SimpleIterator{0, triggerListForTest()}

	expected := "Trigger 2"
	actual := i.Next().ID
	if expected != actual {
		t.Errorf("i.Next().ID - expected %v but %v", expected, actual)
	}

	expected = "Trigger 4"
	actual = i.Next().ID
	if expected != actual {
		t.Errorf("i.Next().ID - expected %v but %v", expected, actual)
	}

	expected = "Trigger 1"
	actual = i.Next().ID
	if expected != actual {
		t.Errorf("i.Next().ID - expected %v but %v", expected, actual)
	}

	expected = "Trigger 3"
	actual = i.Next().ID
	if expected != actual {
		t.Errorf("i.Next().ID - expected %v but %v", expected, actual)
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("TestSimpleIteratorNextOK -", r)
		}
	}()

	_ = i.Next()

	t.Errorf("This line should not be executed")

}

func TestNewSimpleLookup(t *testing.T) {

	s := NewSimpleLookup(triggerListForTest()).(SimpleLookup)

	i := 0
	expected := "Trigger 1"
	actual := s.triggerList[i].ID
	if expected != actual {
		t.Errorf("s.triggerList[%d].ID - expected %v but %v", i, expected, actual)
	}

	i++
	expected = "Trigger 2"
	actual = s.triggerList[i].ID
	if expected != actual {
		t.Errorf("s.triggerList[%d].ID - expected %v but %v", i, expected, actual)
	}

	i++
	expected = "Trigger 3"
	actual = s.triggerList[i].ID
	if expected != actual {
		t.Errorf("s.triggerList[%d].ID - expected %v but %v", i, expected, actual)
	}

	i++
	expected = "Trigger 4"
	actual = s.triggerList[i].ID
	if expected != actual {
		t.Errorf("s.triggerList[%d].ID - expected %v but %v", i, expected, actual)
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("TestNewSimpleLookup -", r)
		}
	}()

	i++
	_ = s.triggerList[i].ID

	t.Errorf("This line should not be executed")

}

func TestSimpleLookupGetTriggerOK(t *testing.T) {

	s := NewSimpleLookup(triggerListForTest()).(SimpleLookup)

	expected := "Trigger 1"
	trigger, err := s.GetTrigger(expected)
	if err != nil {
		t.Error(err)
	}
	actual := trigger.ID
	if expected != actual {
		t.Errorf("s.GetTrigger('%s') - expected %v but %v", expected, expected, actual)
	}

	expected = "Trigger 3"
	trigger, err = s.GetTrigger(expected)
	if err != nil {
		t.Error(err)
	}
	actual = trigger.ID
	if expected != actual {
		t.Errorf("s.GetTrigger('%s') - expected %v but %v", expected, expected, actual)
	}

	if trigger, _ = s.GetTrigger("ABCD"); trigger.ID != "" {
		t.Error("s.GetTrigger('ABCD') should return nil")
	}
}

func TestSimpleLookupGetTriggerEmpty(t *testing.T) {

	s := NewSimpleLookup([]Trigger{}).(SimpleLookup)

	if trigger, _ := s.GetTrigger("Trigger 1"); trigger.ID != "" {
		t.Error("s.GetTrigger('Trigger 1') should return nil")
	}
}

func TestSimpleLookupTriggersOK(t *testing.T) {

	s := NewSimpleLookup(triggerListForTest()).(SimpleLookup)

	i, err := s.Triggers()
	if err != nil {
		t.Error(err)
	}

	if _, ok := i.(*SimpleIterator); !ok {
		t.Errorf("Expected *SimpleIterator but %T", i)
	}

	j := 0
	for i.HasNext() {

		if i.Next().ID != s.triggerList[j].ID {
			t.Errorf("i.Next().ID - Expected %s but %s", s.triggerList[j].ID, i.Next().ID)
		}

		j++
	}

	if j != len(s.triggerList) {
		t.Errorf("j - Expected %d but %d", len(s.triggerList), j)
	}
}

func TestSimpleLookupTriggersEmpty(t *testing.T) {

	s := NewSimpleLookup([]Trigger{}).(SimpleLookup)

	i, err := s.Triggers()
	if err != nil {
		t.Error(err)
	}

	if _, ok := i.(*SimpleIterator); !ok {
		t.Errorf("Expected *SimpleIterator but %T", i)
	}

	if i.HasNext() {
		t.Error("i.HasNext() Expected false but true")
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("TestSimpleLookupTriggersEmpty -", r)
		}
	}()

	_ = i.Next()

	t.Errorf("This line should not be executed")

}

func newConfigForTest() []FormulaConfig {

	return []FormulaConfig{
		{ID: "Formula 1"},
		{ID: "Formula 3"},
		{ID: "Formula 5"},
		{ID: "Formula 2"},
		{ID: "Formula 4"},
	}
}

func TestNewSimpleFormulaLookup(t *testing.T) {

	s := NewSimpleFormulaLookup(
		newConfigForTest(),
		builder.NewFormulaBuilder().Get(),
	).(SimpleFormulaLookup)

	i := 0
	expected := "Formula 1"
	actual := s.configs[i].ID

	if expected != actual {
		t.Errorf("s.configs[%d].ID - Expected %v but %v", i, expected, actual)
	}

	i++
	expected = "Formula 2"
	actual = s.configs[i].ID

	if expected != actual {
		t.Errorf("s.configs[%d].ID - Expected %v but %v", i, expected, actual)
	}

	i++
	expected = "Formula 3"
	actual = s.configs[i].ID

	if expected != actual {
		t.Errorf("s.configs[%d].ID - Expected %v but %v", i, expected, actual)
	}

	i++
	expected = "Formula 4"
	actual = s.configs[i].ID

	if expected != actual {
		t.Errorf("s.configs[%d].ID - Expected %v but %v", i, expected, actual)
	}

	i++
	expected = "Formula 5"
	actual = s.configs[i].ID

	if expected != actual {
		t.Errorf("s.configs[%d].ID - Expected %v but %v", i, expected, actual)
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("TestNewSimpleFormulaLookup -", r)
		}
	}()

	i++
	_ = s.configs[i]

	t.Errorf("This line should not be executed")

}

func TestSimpleFormulaLookupGetFormulas(t *testing.T) {

	s := NewSimpleFormulaLookup(
		newConfigForTest(),
		builder.NewFormulaBuilder().Get(),
	).(SimpleFormulaLookup)

	trigger := Trigger{Filter: "config.ID == 'Formula 5' || config.ID == 'Formula 3'"}

	i, err := s.GetFormulars(trigger, nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(i)

	if !i.HasNext() {
		t.Error("i.HasNext() - should be true")
	}

	actual := i.Next().ID
	expected := "Formula 3"
	if expected != actual {
		t.Errorf("i.HasNext() - Expected %v but %v", expected, actual)
	}

	if !i.HasNext() {
		t.Error("i.HasNext() - should be true")
	}

	actual = i.Next().ID
	expected = "Formula 5"
	if expected != actual {
		t.Errorf("i.HasNext() - Expected %v but %v", expected, actual)
	}

	if i.HasNext() {
		t.Error("i.HasNext() - should be false")
	}

}

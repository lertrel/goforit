package goforit

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/lertrel/goforit/model"
)

var f model.Formula

func NewTriggerLookup(t Trigger) TriggerLookup {

	return MockTriggerLookup{t}
}

type MockTriggerLookup struct {
	trigger Trigger
}

// GetTrigger getting Trigger by name
func (tl MockTriggerLookup) GetTrigger(triggerName string) (Trigger, error) {

	return tl.trigger, nil
}

//GetAllTriggers getting all Trigger(s)
func (tl MockTriggerLookup) GetAllTriggers() (TriggerIterator, error) {
	panic("")
}

//GetTriggers search all Trigger(s) that matches the given filter (script)
func (tl MockTriggerLookup) GetTriggers(filter string) (TriggerIterator, error) {
	panic("")
}

func NewFormulaLookup(c FormulaConfig) FormulaLookup {

	return MockFormulaLookup{c}
}

type MockFormulaLookup struct {
	c FormulaConfig
}

//GetAllFormulas search all FormulaConfig(s) that matches the given Trigger
func (ml MockFormulaLookup) GetFormulars(trigger Trigger, context map[string]interface{}) (FormulaIterator, error) {

	return NewFormulaIterator(ml.c), nil
}

//GetFormula getting FormulaConfig by name
func (ml MockFormulaLookup) GetFormula(id string) (FormulaConfig, error) {
	panic("")
}

//GetAllFormulas getting all FormulaConfig(s)
func (ml MockFormulaLookup) GetAllFormulas() (FormulaIterator, error) {
	panic("")
}

func NewFormulaIterator(t FormulaConfig) FormulaIterator {

	return MockFormulaIterator{configs: []FormulaConfig{t}}
}

type MockFormulaIterator struct {
	index   int
	configs []FormulaConfig
}

func (it MockFormulaIterator) HasNext() bool {
	return it.index < len(it.configs)
}
func (it MockFormulaIterator) Next() (config FormulaConfig) {

	config = it.configs[it.index]
	it.index++

	return
}

func newTriggers() (triggers SimpleTriggers) {

	f = NewFormulaBuilder().Get()
	tl := NewTriggerLookup(newTrigger(
		`
	p = context['principal']
	d = context['dow']
	r = context['rate']
	t = context['term']
	v = context['vat']
	`,
		``))

	fl := NewFormulaLookup(FormulaConfig{
		ID:      "Formula 1",
		Body:    "$LOAN(p, d, r, t, v)",
		Enabled: true,
	})

	triggers = SimpleTriggers{
		formula:       f,
		formulaLookup: fl,
		triggerLookup: tl,
	}

	f.RegisterCustomFunction(
		"$CIRCLE",
		`
		function $CIRCLE(r, pi, context) {
			console.log("JS: context['r']="+context['r']);
			console.log("JS: context['pi']="+context['pi']);
			console.log("JS: r="+r);
			console.log("JS: pi="+pi);
			area = pi*r*r;
			console.log("JS: area="+area);
		
			return area;
		}
		`)

	f.RegisterCustomFunction(
		"$LOAN",
		`
		function $LOAN(principal, dow_payment, interest_rate, term, is_vat) {
			
			initial_loan = principal - dow_payment
			i1 = $RND(initial_loan * interest_rate * term / 12, 2)
			i2 = $IF(is_vat, i1 * 1.07, i1)
			gross_loan = $SUMF(initial_loan, i2)
			
			return gross_loan
		}
		`)

	return
}

func newTrigger(in string, out string) (trigger Trigger) {

	trigger = Trigger{}

	trigger.ContextVarName = "context"
	trigger.OutputVarName = "output"
	trigger.ID = "Test"
	trigger.InputMapping = in
	trigger.OuputMapping = out

	return
}

func TestGetFormula(t *testing.T) {

	triggers := newTriggers()

	actual, _ := triggers.getFormula(Trigger{})
	if !reflect.DeepEqual(f, actual) {
		t.Errorf("Expected %v but %v", f, actual)
	}
}

func TestMapInputsHappy(t *testing.T) {

	triggers := newTriggers()
	trigger := newTrigger("r = context['r']; pi = context['pi'];", "")
	f, _ := triggers.getFormula(trigger)
	fc, _ := f.NewContext("")
	context := make(map[string]interface{})

	context["r"] = 0.5
	context["pi"] = math.Pi

	triggers.mapInputs(&fc, context, trigger)

	script := "$CIRCLE(r, pi, context)"
	if err := fc.Prepare(script); err != nil {
		t.Error(err)
	}

	jsRet, err := fc.Run(script)
	if err != nil {
		t.Error(err)
	}

	goRet, err := jsRet.ToFloat()
	fmt.Println("GO: area=", goRet)

	jsR, err := fc.Get("r")
	if err != nil {
		t.Error(err)
	}

	goR, err := jsR.ToFloat()
	if err != nil {
		t.Error(err)
	}

	if goR != 0.5 {
		t.Errorf("Expected 0.5 but %v", goR)
	}

	jsPi, err := fc.Get("pi")
	if err != nil {
		t.Error(err)
	}

	goPi, err := jsPi.ToFloat()
	if err != nil {
		t.Error(err)
	}

	if goPi != math.Pi {
		t.Errorf("Expected %v but %v", math.Pi, goPi)
	}

}

func TestMapInputsNoMapping(t *testing.T) {

	triggers := newTriggers()
	trigger := newTrigger("", "")
	f, _ := triggers.getFormula(trigger)
	fc, _ := f.NewContext("")
	context := make(map[string]interface{})

	context["r"] = 0.5
	context["pi"] = math.Pi

	triggers.mapInputs(&fc, context, trigger)

	script := "$CIRCLE(r, pi, context)"
	if err := fc.Prepare(script); err != nil {
		t.Error(err)
	}

	if _, err := fc.Run(script); err == nil {
		t.Error(err)
	}

}

func TestMapOutputsHappy(t *testing.T) {

	triggers := newTriggers()
	trigger := newTrigger("r = context['r']; pi = context['pi'];", "output['area'] = area;")
	f, _ := triggers.getFormula(trigger)
	fc, _ := f.NewContext("")
	context := make(map[string]interface{})

	context["r"] = 0.5
	context["pi"] = math.Pi

	triggers.mapInputs(&fc, context, trigger)

	script := "area = pi * r * r; console.log('JS: area='+area); result = area;"
	if err := fc.Prepare(script); err != nil {
		t.Error(err)
	}

	jsRet, err := fc.Run(script)
	if err != nil {
		t.Error(err)
	}

	goRet, err := jsRet.ToFloat()
	if err != nil {
		t.Error(err)
	}
	fmt.Println("GO: goRet=", goRet)

	output, err := triggers.mapOutputs(&fc, trigger)

	goArea, ok := output["area"]
	if !ok {
		t.Error("No value assigned to context['area']")
		t.FailNow()
	}
	fmt.Println("GO: goArea=", goArea)

	if goArea != goRet {
		t.Errorf("Expected %v but %v", goRet, goArea)
	}

}

func TestMapOutputsNoMapping(t *testing.T) {

	triggers := newTriggers()
	trigger := newTrigger("r = context['r']; pi = context['pi'];", "")
	f, _ := triggers.getFormula(trigger)
	fc, _ := f.NewContext("")
	context := make(map[string]interface{})

	context["r"] = 0.5
	context["pi"] = math.Pi

	triggers.mapInputs(&fc, context, trigger)

	script := "area = pi * r * r;"
	if err := fc.Prepare(script); err != nil {
		t.Error(err)
	}

	jsRet, err := fc.Run(script)
	if err != nil {
		t.Error(err)
	}

	goRet, err := jsRet.ToFloat()
	if err != nil {
		t.Error(err)
	}
	fmt.Println("GO: goRet=", goRet)

	output, err := triggers.mapOutputs(&fc, trigger)

	goArea, ok := output["area"]
	if ok {
		t.Errorf("Expected no value assigned to context['area'] but found %v\n", goArea)
		t.FailNow()
	}

}

func TestExecuteHappy(t *testing.T) {

	triggers := newTriggers()
	context := make(map[string]interface{})

	context["principal"] = 1000000
	context["dow"] = 100000
	context["rate"] = 0.99 / 100
	context["term"] = 12
	context["vat"] = false

	result, err := triggers.Execute("loan", context)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(result["_return"])

	expected := 908910.0
	actual := result["_return"].(float64)
	if expected != actual {
		t.Errorf("Expected %v but found %v\n", expected, actual)
	}

	context["term"] = 24
	result, err = triggers.Execute("loan", context)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(result["_return"])

	expected = 917820.0
	actual = result["_return"].(float64)
	if expected != actual {
		t.Errorf("Expected %v but found %v\n", expected, actual)
	}

}

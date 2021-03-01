package main

import (
	"fmt"

	"github.com/lertrel/goforit"
	"github.com/lertrel/goforit/trigger"
)

var triggers trigger.Triggers

func init() {

	//*************************************************************************
	// LOOKING COMPLEX ISN'T IT?
	// *YES* it is.
	//
	// As Triggers' design is targeting a medium-to-large system, it would be
	// overkill for small one as it's providing a mechanism for configuring
	// formulas and trigger points for a larger system where formulas changed
	// frequently.
	//
	// If you are looking for a quick win, trying *Formula* instead.
	//*************************************************************************

	//To user Triggers mechanism, construction of triggers system is required
	//as following (one-time setup at program starts)
	fmt.Println("Initializing ...")

	f := goforit.NewFormulaBuilder().Get()

	//Basic concept of (custom) function (like an excel function)
	//that can be invoked by formula
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

	//Basic concept for formula i.g., a formula executes one/more function(s)
	fl := trigger.NewSimpleFormulaLookup([]trigger.FormulaConfig{
		{
			ID:      "LOAN_FORMULA1",
			Body:    "$LOAN(p, d, r, t, v)",
			Enabled: true,
		},
	}, f)

	//Basic concept of trigger i.e., a trigger point where the program
	//asking for result(s) from a particular formula by supplying required
	//parameters as input(s)
	tl := trigger.NewSimpleLookup([]trigger.Trigger{
		{
			ID:             "Calculating Loan for Product 1",
			Description:    "Calculating a loan amount from given parameters",
			Filter:         "config.ID == 'LOAN_FORMULA1'",
			ContextVarName: "context",
			InputMapping: `
			p = context['principal']
			d = context['dow']
			r = context['rate']
			t = context['term']
			v = context['vat']
			`,
		},
	})

	//Creating a new Triggers to execute a trigger
	triggers = goforit.NewTriggersBuilder().
		SetFormula(f).
		SetFormulaLookup(fl).
		SetTriggerLookup(tl).
		Get()

}

func main() {

	fmt.Println("Trigger Examples")

	//Providing required inputs through a context
	context := make(map[string]interface{})

	context["principal"] = 1000000
	context["dow"] = 100000
	context["rate"] = 0.99 / 100
	context["term"] = 12
	context["vat"] = false

	//Triggering a trigger to get result
	result, err := triggers.Execute("Calculating Loan for Product 1", context)
	if err != nil {
		panic(err)
	}

	//Retreiving a result
	loan := result["_return"].(float64)

	fmt.Println(loan)

}

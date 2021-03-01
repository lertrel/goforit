# goforit
## Introduction
A pure go package for enabling runtime customization formulas for Go utilizing JavaScript package. 

With Go-For-It package, developers can externalize program formula(s) (e.g., in text file, csv, excel file, DB, etc.) and load them during runtime. 

The benefit of having formula(s) externalized is, for program that extensively uses formula(s) so numbers of formula(s) can be maintained outside of source code which in turn making the source code more cleaner, adding/changing formula can be done without stopping the running program, as well as complex formula(s) could be handled by specialists (a.k.a. users).  

**NOTE** The current implementation is built on top of "otto"

## Examples

Here is a very basic example:

```go
package main

import (
	"fmt"
	"github.com/lertrel/goforit"
)

func main() {

	//Creating Formula instance
	f := goforit.GetFormulaBuilder().Get()

//$SUMI is a built-in function
str := `
$SUMI(x, 2*x, 3*x, 4*x) + x
`
  //Loading denpendencies (functions) into context instance
	c, err := f.LoadContext(nil, str)
	if err != nil {
		panic(err)
	}

  //Set value to x
	c.Set("x", 10.0)

  //Run formula & getting result
	jsVal, jsErr := c.Run(str)
	if jsErr != nil {
		panic(jsErr)
	}

  //Converting result into go variable
	goVal, _ := jsVal.ToInteger()

	fmt.Printf("goVal=%v\n", goVal)

}
```

Custom formula example:

```go
package main

import (
	"fmt"

	"github.com/lertrel/goforit"
)

func main() {

	f := goforit.NewFormulaBuilder().Get()
	f.RegisterCustomFunction(
		"$CIRCLE",
		`
				function $CIRCLE(radius) {
					return $RND(Math.PI * Math.pow(radius, 2), 10);
				}
			`)

	script := "$CIRCLE(r)"
	fc, err := f.NewContext(script)
	if err != nil {
		panic(err)
	}

	if err := fc.Set("r", 5); err != nil {
		panic(err)
	}

	jsRet, err := fc.Run(script)
	if err != nil {
		panic(err)
	}

	goRet, err := jsRet.ToFloat()
	if err != nil {
		panic(err)
	}

	fmt.Printf("goRet=%v\n", goRet)
}
```

Triggers example:

```go
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
```

## Overview Concepts

Following sections are general idea to the main components of Go-For-It.

**Abstract Model**

[_Fundamental_]
- Functions (custom and built-in)
- In-line / adhoc formula

[_Intermediate_]
- FormulaConfig
- Trigger

**Programming Components**

[_Basic_]
- Formula
- FormulaContext
- Value

[_Intermediate_]
- Triggers

**Extention Points** (sorting for rather simple -> complex elements)

[_Intermediate_]
- CustomFunctionRepository
- FormulaLookup
- TriggerLookup

[_Advance_]
- BuiltinFunctions
- VM
- VMDriver

## [Basic]
### 1. Functions
Thinking of a function as an excel function like SUM() for examples. There are 2 types for functions.
#### 1.1 Built-in function

    Built-in functions are functions implemented in Go and registered to the underlying scripting VM/engine so they can be accessed by scripting language.

	The advantange of using built-in functions is, they are run faster and do not required to be loaded and parsed by the underlying scripting VM/engine.

	**Full list of supported built-in functions** are provided at the end of this document.

	To create additional built-in functions, see BuiltInFunctions section 

	To implement a VM, developer has to directly interact with the native API of the underlying scripting VM/engine (e.g., Otto) instead of abstract layers provided by Go-For-It.

#### 1.2 Custom function

    In contrast to built-in function, a custom function is written in scripting language (e.g, JavaScript), registered into Formula, and lastly loaded into FormulaContext at runtime.

	As custom function is done by scripting language, so it is neitgher required to be compiled or shipped together with your go program. Moreover any developers with some scripting language (e.g., JavaScript) could easily create new custom functions to satisfy every changing business requirements of your customer or organization.

### 2. Formula

    Here a tricky one! As a center of Go-For-It core concepts, the term formula could refer to several things.

#### 2.1 In-line / Adhoc formula

    In a small program, it's often that formulas will be part of the source codes or at best stored in text files shipped as part of program delivery.

	These script will be passed to Go-For-It to executes them at program runtime.

#### 2.2 FormulaConfig (intermediate)

    In contrast to in-line/adhoc formula, formulas can be externalized from the program and made configurable in form of FormulaConfig.

	FormulaConfig(s) are pre-configured in advance of running the program, and will be loaded during runtime through FormulaLookup.

#### 2.3 Formula

    An API component "Formula" (obtained from FormulaBuilder) for executing formulas and functions, and also used for creating a FormulaContext

> **Formula _VS_ Function**
> 
> (What's the difference???)
>
> Function characteristics
> - is a complex or long script -- for computing or creating desirable result -- that's not easy to be re-written thus it's should be prepared in separately advance (e.g., finding list of states/provinces in the country from which selling some your products is forbidden)
> - could be common enough to be shared among programs or modules (e.g., SUM, rounding, etc)
> - some of the functions, could not be implemented or it's not easly to be implemented using scripting language (such as excessing DB or API). This type of functions are required to be implemented built-in functions (as there's no built-in formula in Go-For-It)
>
> Formula characteristics
> - is a concise script usually describing steps to produce a result in regards business requirement in question like calculating product pricing baused on given inputs, calculating age of a person, calculating applicable tax etc.
> - usually calls function(s)
> - can be called by a trigger
>
> Program -- calls --> Formula -- calls --> functions

### 3. FormulaContext

	A FormulaContext is holding process states and loaded functions.

    It's representing an execution context of formula(s) in a particula ...
	- Program Thread
      As it's holding states of the current exectuion process, FormulaContext is not thread-safe. For concurrent environment, a dedicated FormulaContext should be assigned to each program thread.

	- Program Module
	  As processes under the same program module usually depend on some set of functions, using the same FormulaContext of (it's clone) would save time of re-loading and re-parsing of same functions again and again

	- Program Process
	  As it's holding states of the current exectuion process, states of unfinish process would be lost if completing each step of the same process under different FormulaContext(s).

	  An exception is, if your system is very large or intensively focusing on concurrency and scalibility and that the system design suggests that steps in underthe same process could be executed in different timing.

	  If that's the case, there'd be a chance of externalizing process states by persisting and re-loading them into a different FormulaContext

### 4. Value

    Having dynamic formula could be useless if it's unable to obtain the result. Value is a wrapper of the result produced by formula (inside scripting VM/engine). So developer can retrieve for the Value and convert it into Go data type

## [Intermediate]

### 5. CustomFunctionRepository
### 6. Triggers System
#### 6.1 Trigger
#### 6.2 Triggers
#### 6.3 TriggerLookup
#### 6.4 FormulaLookup

## [Advance]

### 7. BuiltinFunctions
### 8. VMDriver / VM

## Appendix

### List of Built-in Functions

- $ABS
- $RND
- $FLOOR or $FLR
- $CEIL
- $IF
- $SUMI
- $SUMF
- $AVG
- $MIN
- $MAX

#### $ABS(value)
    Returning and absolute (positive only) of the given value

	Ex.
	$ABS(1) // 1
	$ABS(-1.0) // 1.0
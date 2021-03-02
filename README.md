# goforit
## Introduction

**What is Go-For-It?**

A pure go package for enabling runtime customization formulas for Go utilizing JavaScript package. 

With Go-For-It package, developers can externalize program formula(s) (e.g., in text file, csv, excel file, DB, etc.) and load them during runtime. 

**Why would it might be needed?**

The benefit of having formula(s) externalized is, for program that extensively uses formula(s) so numbers of formula(s) can be maintained outside of source code which in turn making the source code more cleaner, adding/changing formula can be done without stopping the running program, as well as complex formula(s) could be handled by specialists (a.k.a. users).

**When should be considering using (Go-For-)it?**

- Having many formulas to manage / discover, or
- Formulas can be frequently adding or changing, or
- Needing non-programmer business specialist to develop some formulars, or
- Wanting to share common formulas across multiple programs, or
- Bascially, any other good reasons to have formulas / calculation logics externalized from program

**When shouldn't?**

- _Program is too small_ to bother or logic is too straight forward or rarely changed
- _Performance_ of scripting language (e.g., JavaScript) is not acceptable
- _Extreme data analytic_, though it's something to do with formual, Go-For-It was designed for formula management (configuring, editing, dicovering, and exeucting) but not intended to be providing powerful formulas as part of the shipment
- _Existing JavaScript libraries_ are needed. There are plenty of sophisticated and powerful JavaScript libraries out there, but loading and parsing entire or multiple libraries in Go are something you should think twice, and you might want to consider using Nodes.js instead.

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
Thinking of a function as an excel function like SUM() for examples.

For performance sake, function _must_ start with '$', and conventionally all letters should be captital, word separater is optional, but should be '_'

Ex.

function $PRICE(_product, qty, discountType, payType, orderDate_)

function $DAMAGE_POINTS(_attakcer, target, attack, times, distance_)

There are 3 types for functions.
#### 1.1 Built-in function

    Built-in functions are functions implemented in Go and registered to the underlying scripting VM/engine so they can be accessed by scripting language.

	The advantange of using built-in functions is, they are run faster and do not required to be loaded and parsed by the underlying scripting VM/engine.

	And as they are written in Go, so they are abstract logic that expected to be working on different scripting launuages.

	**Full list of supported built-in functions** are provided at the end of this document.

	To create additional built-in functions, see BuiltInFunctions section 

#### 1.2 Custom function

    In contrast to built-in function, a custom function is written in a specific scripting language (e.g, JavaScript), registered into Formula component, and lastly loaded into FormulaContext at runtime.

	As custom function is done by scripting language, so it is neitgher required to be compiled or shipped together with your go program. Moreover any developers with some scripting language skill (e.g., JavaScript) could easily create new custom functions to satisfy every changing business requirements of your customer or organization.

#### 1.3 Pre-defined JavaScript function
    What is it?

	Pre-defined JavaScript functions are also custom functions as they were written in scripting language (i.e., JavaScript) thus they will only work for JavaScript engine whereas built-in functions are more abstract and will work for any scripting VM/engine of choice.

	Why?

	Though quite rarely but some kind of COMMON functions are easier to be written in JavaScript.

	For examples $SUM(), since Go is a type-safe language, so to write $SUM() that support both int/float and returning int or float depending on the given parameters in Go will be too complex comparing to its benefit, but writting this kind of function in JavaScript is pretty easy as JavaScript has no explicit type.

	Some functions are a wrapper of JavaScript existing functions like Math.*, etc. so they are provided as interim solutions before having their built-in functions counterparts.

	Lastly rich set of pre-defined functions would save developers/specialists/users times for creating common stuff so they can utilize their time to develop business specific logic instead.

	How?

	To use pre-defined JavaScript functions, one just has to add concrete implementation of (pre-built) CustomFunctionRepo(s) shipped with Go-For-It to a Formula component

	When?

	Before creating FormulaContext 
	
	Where?
	
	The pre-defined CustomFunctionRepo(s) are Under the subpackage "js"
	 
### 2. Formula

    Here a tricky one! As a center of Go-For-It core concepts, the term formula could refer to several things.

#### 2.1 In-line / Adhoc formula

    In a small program, it's often that formulas will be part of the source codes or at best stored in text files shipped as part of program delivery.

	These script will be passed to Go-For-It to executes them at program runtime.

#### 2.2 FormulaConfig (intermediate)

    In contrast to in-line/adhoc formula, formulas can be externalized from the program and made configurable in form of FormulaConfig.

	FormulaConfig(s) are pre-configured in advance of running the program, and will be loaded during runtime through FormulaLookup.

#### 2.3 Formula Component

    An API component "Formula" (obtained from FormulaBuilder) for executing formulas and functions, and also used for creating a FormulaContext

> **Formula _VS_ Function**
> 
> (What's the difference???)
>
> Function characteristics
> - is a **COMPLEX** or **LONG** script -- **FOR DETAILS COMPUTING** or creating desirable result -- that's not easy to be re-written thus it's should be prepared separately in advance (e.g., finding list of states/provinces in the country from which selling some your products are forbidden)
> - could be **COMMON** enough to be shared among programs or modules (e.g., SUM, rounding, etc)
> - some of the functions, **COULD NOT BE IMPLEMENTED** or it's not easly to be implemented **USING SCRIPTING LANGUAGE** (such as excessing DB, API, or low level operations). This type of functions are required to be implemented as built-in functions (as there's no built-in formula in Go-For-It)
> - **BEING CALLED BY FORMULA(s)**
>
> Formula characteristics
> - is a more **CONCISE** script usually **DESCRIBING STEPS** to produce a result in regards **SPECIFIC** business requirement in question like calculating product pricing baused on given inputs, calculating age of a person, calculating applicable tax etc.
> - usually **CALLS FUNCTION(s)**
> - can **BE CALLED BY A TRIGGER**
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
	To implement a VM, developer has to directly interact with the native API of the underlying scripting VM/engine (e.g., Otto) instead of abstract layers provided by Go-For-It.

## Appendix

### List of Built-in Functions

- $ABS
- $AVG
- $CEIL
- $FLOOR or $FLR
- $IF
- $MAX
- $MIN
- $RND
- $SUMF
- $SUMI
- ...

##### $ABS ( _value_ )
    Returning and absolute (positive only) of the given value

	Ex.
	$ABS(1) // 1
	$ABS(-1.0) // 1.0
##### $AVG ( _float1, float2, ..._ )
    Returning a average of the given floats

	Ex.
	$AVG(1.0, 2.5, 3.0, 4.5, 5.5) //3.3
	$AVG(price1, price2, price3, price4)
##### $CEIL ( _value, precision_ )
    Returning a round-up value the given value with the precision as given

	Ex.
	$CEIL(1.5, 0) // 2
	$CEIL(1.4, 0) // 2
	$CEIL(1.445, 2) // 1.45
	$CEIL(1.445, 1) // 1.5
##### $FLR ( _value, precision_ )
##### $FLOOR ( _value, precision_ )
    Returning a round-down value the given value with the precision as given

	Ex.
	$FLR(1.5, 0) // 1
	$FLR(1.4, 0) // 1
	$FLR(1.445, 2) // 1.44
	$FLR(1.445, 1) // 1.4
##### $IF ( _condition, value1, value2_ )
    Returning a value1 if the given condition is true otherwise returning value2

	Ex.
	$IF(gender == "M", 200, 150)
	$IF(year > 2020, $NEWPRICE(), $OLDPRICE())
##### $MAX ( _float1, float2, ..._ )
    Returning the maximum value among the given floats

	Ex.
	$MIN(1.0, 2.5, 3.0, 4.5, 5.5) //5.5
	$MIN(price1, price2, price3, price4)
##### $MIN ( _float1, float2, ..._ )
    Returning the minimum value among the given floats

	Ex.
	$MIN(1.0, 2.5, 3.0, 4.5, 5.5) //1.0
	$MIN(price1, price2, price3, price4)
##### $RND ( _value, precision_ )
    Returning a (normal) rounded value the given value with the precision as given

	Ex.
	$RND(1.5, 0) // 2
	$RND(1.4, 0) // 1
	$RND(1.445, 2) // 1.45
	$RND(1.445, 1) // 1.4
##### $SUMF ( _float1, float2, ..._ )
    Returning a result of summing the given floats

	Ex.
	$SUMI(1.0, 2.5, 3.0, 4.5, 5.5) //16.5
	$SUMI(price1, price2, price3, price4)
##### $SUMI ( _integer1, integer2, ..._ )
    Returning a result of summing the given integers

	Ex.
	$SUMI(1, 2, 3, 4, 5) //15
	$SUMI(count1, count2, count3, count4)

**<< MORE TO COME >>**
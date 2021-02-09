package goforit

import (
	"testing"
)

func TestSimpleFormula(t *testing.T) {

str := `
abc = 2 + 2;
console.log("The value of abc is " + abc); // 4
`

	f := Get()
	c, err := f.LoadContext(nil, str)
	if err != nil {
		t.Error(err)
	}

	_, runtimeError := c.Run(str)
	if runtimeError != nil {
		t.Error(runtimeError)
	}
}

func TestBuiltInFormula(t *testing.T) {

str := `
i = $SUMI(1, 2, $SUMI(1, $MIN(2,3)), $SUMI(2, 2), 5);
f = $SUMF(1.5, $SUMF($MAX(1.2, 1.1), $ABS(-1.39)), $IF(i == 15, 5.0, 6.0));
console.log("i = " + i);
console.log("f = " + f);
`

	f := Get()
	c, err := f.LoadContext(nil, str)
	if err != nil {
		t.Error(err)
	}

	_, runtimeError := c.Run(str)
	if runtimeError != nil {
		t.Error(runtimeError)
	}

	jsI, err2 := c.Get("i")
	if err2 != nil {
		t.Error(err2)
	}
	goI, err3 := jsI.ToInteger()
	if err3 != nil {
		t.Error(err3)
	} else if goI != 15 {
		t.Errorf("Expect %v but got %v\n", 15, goI)
	}

	jsF, err4 := c.Get("f")
	if err4 != nil {
		t.Error(err4)
	}
	goF, err5 := jsF.ToFloat()
	if err5 != nil {
		t.Error(err5)
	} else if goF != 9.09 {
		t.Errorf("Expect %v but got %v\n", 9.09, goF)
	}
}

func TestSUMI(t *testing.T) {

str := `
$SUMI(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
`

	f := Get()
	c, err := f.LoadContext(nil, str)
	if err != nil {
		t.Error(err)
	}

	jsI, runtimeError := c.Run(str)
	if runtimeError != nil {
		t.Error(runtimeError)
	}

	var expected int64 = 55
	goI, err3 := jsI.ToInteger()
	t.Logf("goI = %v\n", goI)
	if err3 != nil {
		t.Error(err3)
	} else if goI != expected {
		t.Errorf("Expect %v but got %v\n", expected, goI)
	}
	
}

func TestSUMF(t *testing.T) {

	str := `
	$SUMF(1.0, 2.0, 3.0, 4.0, 5.1, 6.1, 7.1, 8.2, 9, 10)
	`
	
		f := Get()
		c, err := f.LoadContext(nil, str)
		if err != nil {
			t.Error(err)
		}
	
		jsI, runtimeError := c.Run(str)
		if runtimeError != nil {
			t.Error(runtimeError)
		}
	
		var expected float64 = 55.5
		goI, err3 := jsI.ToFloat()
		t.Logf("goI = %v\n", goI)
		if err3 != nil {
			t.Error(err3)
		} else if goI != expected {
			t.Errorf("Expect %v but got %v\n", expected, goI)
		}
		
}
	
func TestMIN(t *testing.T) {

str := `
$MIN(6, 7, 8, 9, 10, 5.2, 5.5, 5.1)
`

	f := Get()
	c, err := f.LoadContext(nil, str)
	if err != nil {
		t.Error(err)
	}

	jsI, runtimeError := c.Run(str)
	if runtimeError != nil {
		t.Error(runtimeError)
	}

	var expected float64 = 5.1
	goI, err3 := jsI.ToFloat()
	t.Logf("goI = %v\n", goI)
	if err3 != nil {
		t.Error(err3)
	} else if goI != expected {
		t.Errorf("Expect %v but got %v\n", expected, goI)
	}
	
}

func TestMAX(t *testing.T) {

str := `
$MAX(6, 7, 8, 9, 10, 5.2, 5.5, 5.1)
`

	f := Get()
	c, err := f.LoadContext(nil, str)
	if err != nil {
		t.Error(err)
	}

	jsI, runtimeError := c.Run(str)
	if runtimeError != nil {
		t.Error(runtimeError)
	}

	var expected float64 = 10.0
	goI, err3 := jsI.ToFloat()
	t.Logf("goI = %v\n", goI)
	if err3 != nil {
		t.Error(err3)
	} else if goI != expected {
		t.Errorf("Expect %v but got %v\n", expected, goI)
	}
	
}

func TestAVG(t *testing.T) {

	str := `
	$AVG(6, 7, 8, 9, 10, 5.2, 5.5, 5.1)
	`
	
	f := Get()
	c, err := f.LoadContext(nil, str)
	if err != nil {
		t.Error(err)
	}

	jsI, runtimeError := c.Run(str)
	if runtimeError != nil {
		t.Error(runtimeError)
	}

	var expected float64 = 6.975
	goI, err3 := jsI.ToFloat()
	t.Logf("goI = %v\n", goI)
	if err3 != nil {
		t.Error(err3)
	} else if goI != expected {
		t.Errorf("Expect %v but got %v\n", expected, goI)
	}
		
}
	
func TestABS(t *testing.T) {

	str := `
	$ABS(-65.2285)
	`
	
	f := Get()
	c, err := f.LoadContext(nil, str)
	if err != nil {
		t.Error(err)
	}

	jsI, runtimeError := c.Run(str)
	if runtimeError != nil {
		t.Error(runtimeError)
	}

	var expected float64 = 65.2285
	goI, err3 := jsI.ToFloat()
	t.Logf("goI = %v\n", goI)
	if err3 != nil {
		t.Error(err3)
	} else if goI != expected {
		t.Errorf("Expect %v but got %v\n", expected, goI)
	}
		
}
	
func TestRND(t *testing.T) {

	str := `
	$RND(-65.2285, 3)
	`
	
	f := Get()
	c, err := f.LoadContext(nil, str)
	if err != nil {
		t.Error(err)
	}

	jsI, runtimeError := c.Run(str)
	if runtimeError != nil {
		t.Error(runtimeError)
	}

	var expected float64 = -65.229
	goI, err3 := jsI.ToFloat()
	t.Logf("goI = %v\n", goI)
	if err3 != nil {
		t.Error(err3)
	} else if goI != expected {
		t.Errorf("Expect %v but got %v\n", expected, goI)
	}
		
	str = `
	$RND(-65.2285, 2)
	`
	
	jsI, runtimeError = c.Run(str)
	if runtimeError != nil {
		t.Error(runtimeError)
	}

	expected = -65.23
	goI, err3 = jsI.ToFloat()
	t.Logf("goI = %v\n", goI)
	if err3 != nil {
		t.Error(err3)
	} else if goI != expected {
		t.Errorf("Expect %v but got %v\n", expected, goI)
	}
		
	str = `
	$RND(-65.2285, 0)
	`
	
	jsI, runtimeError = c.Run(str)
	if runtimeError != nil {
		t.Error(runtimeError)
	}

	expected = -65
	goI, err3 = jsI.ToFloat()
	t.Logf("goI = %v\n", goI)
	if err3 != nil {
		t.Error(err3)
	} else if goI != expected {
		t.Errorf("Expect %v but got %v\n", expected, goI)
	}
		
}
	
func TestCEIL(t *testing.T) {

	str := `
	$CEIL(-65.2244, 3)
	`
	
	f := Get()
	c, err := f.LoadContext(nil, str)
	if err != nil {
		t.Error(err)
	}

	jsI, runtimeError := c.Run(str)
	if runtimeError != nil {
		t.Error(runtimeError)
	}

	var expected float64 = -65.224
	goI, err3 := jsI.ToFloat()
	t.Logf("goI = %v\n", goI)
	if err3 != nil {
		t.Error(err3)
	} else if goI != expected {
		t.Errorf("Expect %v but got %v\n", expected, goI)
	}
		
	str = `
	$CEIL(-65.2244, 2)
	`
	
	jsI, runtimeError = c.Run(str)
	if runtimeError != nil {
		t.Error(runtimeError)
	}

	expected = -65.22
	goI, err3 = jsI.ToFloat()
	t.Logf("goI = %v\n", goI)
	if err3 != nil {
		t.Error(err3)
	} else if goI != expected {
		t.Errorf("Expect %v but got %v\n", expected, goI)
	}
		
	str = `
	$CEIL(-65.2285, 0)
	`
	
	jsI, runtimeError = c.Run(str)
	if runtimeError != nil {
		t.Error(runtimeError)
	}

	expected = -65
	goI, err3 = jsI.ToFloat()
	t.Logf("goI = %v\n", goI)
	if err3 != nil {
		t.Error(err3)
	} else if goI != expected {
		t.Errorf("Expect %v but got %v\n", expected, goI)
	}
		
}
	
func TestFLOOR(t *testing.T) {

	str := `
	$FLOOR(-65.2244, 3)
	`
	
	f := Get()
	c, err := f.LoadContext(nil, str)
	if err != nil {
		t.Error(err)
	}

	jsI, runtimeError := c.Run(str)
	if runtimeError != nil {
		t.Error(runtimeError)
	}

	var expected float64 = -65.225
	goI, err3 := jsI.ToFloat()
	t.Logf("goI = %v\n", goI)
	if err3 != nil {
		t.Error(err3)
	} else if goI != expected {
		t.Errorf("Expect %v but got %v\n", expected, goI)
	}
		
	str = `
	$FLOOR(-65.2244, 2)
	`
	
	jsI, runtimeError = c.Run(str)
	if runtimeError != nil {
		t.Error(runtimeError)
	}

	expected = -65.23
	goI, err3 = jsI.ToFloat()
	t.Logf("goI = %v\n", goI)
	if err3 != nil {
		t.Error(err3)
	} else if goI != expected {
		t.Errorf("Expect %v but got %v\n", expected, goI)
	}
		
	str = `
	$FLOOR(-65.2285, 0)
	`
	
	jsI, runtimeError = c.Run(str)
	if runtimeError != nil {
		t.Error(runtimeError)
	}

	expected = -66
	goI, err3 = jsI.ToFloat()
	t.Logf("goI = %v\n", goI)
	if err3 != nil {
		t.Error(err3)
	} else if goI != expected {
		t.Errorf("Expect %v but got %v\n", expected, goI)
	}
		
}
	
func TestIF(t *testing.T) {

	str := `
	$IF(true, 1, 2)
	`
	
	f := Get()
	c, err := f.LoadContext(nil, str)
	if err != nil {
		t.Error(err)
	}

	jsI, runtimeError := c.Run(str)
	if runtimeError != nil {
		t.Error(runtimeError)
	}

	var expected int64 = 1
	goI, err3 := jsI.ToInteger()
	t.Logf("goI = %v\n", goI)
	if err3 != nil {
		t.Error(err3)
	} else if goI != expected {
		t.Errorf("Expect %v but got %v\n", expected, goI)
	}
		
	str = `
	$IF(false, 1, 2)
	`	
	jsI, runtimeError = c.Run(str)
	if runtimeError != nil {
		t.Error(runtimeError)
	}

	expected = 2
	goI, err3 = jsI.ToInteger()
	t.Logf("goI = %v\n", goI)
	if err3 != nil {
		t.Error(err3)
	} else if goI != expected {
		t.Errorf("Expect %v but got %v\n", expected, goI)
	}
		
}
	
func TestExtractFunctionListFromFormulaString(t * testing.T) {

src := `
    // Sample xyzzy example
    $function(){
        if (3.14159 > 0) {
            console.log("Hello, World.");
            return;
        }

        var xyzzy = NaN;
        console.log("Nothing happens.");
        return xyzzy;
	};
	$sum(5, $length(8) * 15);
	$if(a > 5, $if( b < 2, $int($left(x, 2)) + 5, $sum(5, 6, 7)), -1)
	// Sample xyzzy example
		function(){
			if (3.14159 > 0) {
				console.log("Hello, World.");
				return;
			}
	
			var xyzzy = NaN;
			console.log("Nothing happens.");
			return xyzzy;
		};	
`

	f := Get()
	listOfFuncs := f.extractFunctionListFromFormulaString(src)

	if len(listOfFuncs) != 6 {
		t.Error("Expected arror of length = 6")
	}

	for i := 0; i < len(listOfFuncs); i++ {
		t.Logf("Dedup: %s\n", listOfFuncs[i])
	}

}

func TestCustomFormulaSimple(t *testing.T) {

	f := Get()
	f.RegisterCustomFunction(
		"$PREMIUM1",
		`
		function $PREMIUM1(gender, age, rates, sa) {

			console.log("gender="+gender);
			console.log("age="+age);
			console.log("rates="+rates);
			console.log("sa="+sa);

			r = rates[age];
			console.log("r="+r);

			if (gender == "M") {
				r *= 2;
			}
			console.log("r="+r);
			console.log("sa / 1000 * r="+(sa / 1000 * r));

			return sa / 1000 * r;
		}
		`)


str := `
$PREMIUM1("M", 5, [1.1, 1.2, 1.3, 1.4, 1.5, 1.615], 100000)
`

	c, err := f.LoadContext(nil, str)
	if err != nil {
		t.Error(err)
	}

	jsRet, runtimeError := c.Run(str)
	if runtimeError != nil {
		t.Error(runtimeError)
	}

	goRet, err2 := jsRet.ToFloat()
	var expected float64 = 323.0
	if err2 != nil {
		t.Error(err2)
	} else if goRet != expected {
		t.Errorf("Expect %v but got %v\n", expected, goRet)
	}

str = `
$PREMIUM1("F", 5, [1.1, 1.2, 1.3, 1.4, 1.5, 1.615], 100000)
`

	c, err = f.LoadContext(nil, str)
	if err != nil {
		t.Error(err)
	}

	jsRet, runtimeError = c.Run(str)
	if runtimeError != nil {
		t.Error(runtimeError)
	}

	goRet, err2 = jsRet.ToFloat()
	expected = 161.5
	if err2 != nil {
		t.Error(err2)
	} else if goRet != expected {
		t.Errorf("Expect %v but got %v\n", expected, goRet)
	}
}

func TestCustomFormulaComplex(t *testing.T) {

	f := Get()
	f.RegisterCustomFunction(
		"$PREMIUM1",
		`
		function $PREMIUM1(gender, age, rates, sa) {

			console.log("$PREMIUM1 - gender="+gender);
			console.log("$PREMIUM1 - age="+age);
			console.log("$PREMIUM1 - rates="+rates);
			console.log("$PREMIUM1 - sa="+sa);

			r = rates[age];
			console.log("$PREMIUM1 - r="+r);

			r = $IF(gender=="M", $RND(r*2, 2), r)
			console.log("$PREMIUM1 - r="+r);
			console.log("$PREMIUM1 - sa / 1000 * r="+$RND(sa / 1000 * r, 2));

			return $RND(sa / 1000 * r, 2);
		}
		`)


	f.RegisterCustomFunction(
		"$PREMIUM2",
		`
		function $PREMIUM2(gender, age, rates1, sa1, rates2, sa2, rates3, sa3) {

			console.log("$PREMIUM2 - gender="+gender);
			console.log("$PREMIUM2 - age="+age);
			console.log("$PREMIUM2 - rates1="+rates1);
			console.log("$PREMIUM2 - sa1="+sa1);
			console.log("$PREMIUM2 - rates2="+rates2);
			console.log("$PREMIUM2 - sa2="+sa2);
			console.log("$PREMIUM2 - rates3="+rates3);
			console.log("$PREMIUM2 - sa3="+sa3);

			r1 = rates1[age];
			console.log("$PREMIUM2 - r1="+r1);
			r2 = rates2[age];
			console.log("$PREMIUM2 - r2="+r2);

			rates1[age] = $IF(gender=="M", $RND(r1*2, 2), r1)
			rates2[age] = $IF(gender=="F", $RND(r2*3, 2), r2)

			sa4 = $MAX(sa1, sa2, sa3)
			sa5 = $FLOOR($AVG(sa1, sa2, sa3), 2)

			p = $RND($SUMF(
					$PREMIUM1(gender, age, rates1, sa1), 
					$PREMIUM1(gender, age, rates2, sa2), 
					$PREMIUM1(gender, age, rates3, sa3), 
					$PREMIUM1(gender, age, rates3, sa4), 
					$PREMIUM1(gender, age, rates3, sa5) 
				), 2)

			console.log("$PREMIUM2 - p="+p);

			return p;
		}
		`)


str := `
$PREMIUM2("M", 1, [1.1, 1.2, 1.3], 100000, [2.1, 2.2, 2.3], 50000, [3.1, 3.2, 3.3], 4555)
`
	c, err := f.LoadContext(nil, str)
	if err != nil {
		t.Error(err)
	}

	jsRet, runtimeError := c.Run(str)
	if runtimeError != nil {
		t.Error(runtimeError)
	}

	goRet, err2 := jsRet.ToFloat()
	var expected float64 = 1698.87
	if err2 != nil {
		t.Error(err2)
	} else if goRet != expected {
		t.Errorf("Expect %v but got %v\n", expected, goRet)
	}


str = `
$PREMIUM2("F", 1, [1.1, 1.2, 1.3], 100000, [2.1, 2.2, 2.3], 50000, [3.1, 3.2, 3.3], 4555)
`
	c, err = f.LoadContext(c, str)
	if err != nil {
		t.Error(err)
	}

	jsRet, runtimeError = c.Run(str)
	if runtimeError != nil {
		t.Error(runtimeError)
	}

	goRet, err2 = jsRet.ToFloat()
	expected = 949.44
	if err2 != nil {
		t.Error(err2)
	} else if goRet != expected {
		t.Errorf("Expect %v but got %v\n", expected, goRet)
	}

}

func TestCustomMultipleFormula(t *testing.T) {

	f := Get()
	f.RegisterCustomFunction(
		"$CIRCLE",
		`
		function $CIRCLE(radius) {

			return $RND(Math.PI * Math.pow(radius, 2), 10); 
		}
		`)


str := `
area1 = $RND(Math.sqrt($SUMF($RND(a*3,2), $RND(b*4,2), $RND(c*5,2))), 10);
area2 = $CIRCLE(radius);
console.log("a="+a);
console.log("b="+b);
console.log("c="+c);
console.log("radius="+radius);
`
	c, err := f.LoadContext(nil, str)
	if err != nil {
		t.Error(err)
	}

	c.Set("a", 2)
	c.Set("b", 3)
	c.Set("c", 4)
	c.Set("radius", 5)

	_, runtimeError := c.Run(str)
	if runtimeError != nil {
		t.Error(runtimeError)
	}

	jsArea1, _ := c.Get("area1")
	jsArea2, _ := c.Get("area2")

	area1, _ := jsArea1.ToFloat()
	area2, _ := jsArea2.ToFloat()

	var expected1 float64 = 6.164414003
	var expected2 float64 = 78.5398163397
	if area1 != expected1 {
		t.Errorf("Expect %v but got %v\n", expected1, area1)
	}

	if area2 != expected2 {
		t.Errorf("Expect %v but got %v\n", expected2, area2)
	}

}

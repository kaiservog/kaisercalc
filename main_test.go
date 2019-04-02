package main

import (
	"testing"
	//"fmt"
)

func TestIsNextCharacterNumber(t *testing.T) {

	if !IsNextCharacterIsNumber(0, "12a3") {
		t.Errorf("Error")
	}

	if IsNextCharacterIsNumber(1, "12a3") {
		t.Errorf("Error")
	}

	if IsNextCharacterIsNumber(3, "12a3") {
		t.Errorf("Error")
	}
}

func TestToPostfix(t *testing.T) {
	s := ConvertToPostfix("5*(6+2)-12/4")

	st := ""
	for s.Len() > 0 {
		st += s.Pop().(string)
	}

	if st != "562+*124/-" {
		t.Errorf("Expression in stack is wrong")
	}
}

func TestToPostfixWithVariable(t *testing.T) {
	s := ConvertToPostfix("5*(6+pi)-12/4")

	st := ""
	for s.Len() > 0 {
		st += s.Pop().(string)
	}

	if st != "56pi+*124/-" {
		t.Errorf("Expression in stack is wrong")
	}
}

func TestResolveExpression(t *testing.T) {
	defs := make(map[string]CalcExp)
	r, err := Resolve(ConvertToPostfix("5*(6+2)-12/4"), &defs)

	if err != nil {
		t.Errorf(err.Error())
	}

	if r != "37" {
		t.Errorf("Expression wrong resolution")
	}
}

func TestResolveExpressionWithDefinitions(t *testing.T) {
	defs := make(map[string]CalcExp)
	defs["pi"] = CalcExp{"3"}

	s := ConvertToPostfix("2*pi+5")
	r, err := Resolve(s, &defs)

	if err != nil {
		t.Errorf(err.Error())
	}

	if r != "11" {
		t.Errorf("Expression wrong resolution")
	}
}

func TestExpression(t *testing.T) {
	r := NewRegexpCalc()
	exp := CalcExp{"1+2-3*4/5"}

	if exp.IsSpecialExpression(r) {
		t.Errorf("it's not a special expression")
	}

	exp = CalcExp{"1+2-3*4/var"}
	if !exp.IsSpecialExpression(r) {
		t.Errorf("it's a special expression")
	}
}

func TestFunctionCall(t *testing.T) {
	s := ConvertToPostfix("print(3+5)")
	expected := "35+print"

	st := ""
	for s.Len() > 0 {
		st += s.Pop().(string)
	}

	if st != expected {
		t.Errorf("stacked wrong using function call")
	}

}

func TestFloatNumbers(t *testing.T) {
	defs := make(map[string]CalcExp)
	r, err := Resolve(ConvertToPostfix("0.5+1.6"), &defs)

	if err != nil {
		t.Errorf(err.Error())
	}

	if r != "2.1" {
		t.Errorf("Expression wrong resolution shound be " + "2.1" + " it was " + r)
	}
}

func TestArgExpression(t *testing.T) {
	processExpression("a=2+2; b=a + 1")
}

func TestStackFunction(t *testing.T) {
	cc := NewCalcCompiler()
	cc.CompileLine("mysum(x, y)=x+y")

	if cc.Funcs["mysum"].Exp != "x+y" {
		t.Error("stack for functions failed")
	}

	if len(cc.Funcs["mysum"].Args) != 2 {
		t.Error("stack for functions failed")
	}

}

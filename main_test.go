package main

import (
	"testing"
	//"fmt"
)

func TestIsNextCharacterNumber(t *testing.T) {

	if !isNextCharacterIsNumber(0, "12a3") {
		t.Errorf("Error")
	}

	if isNextCharacterIsNumber(1, "12a3") {
		t.Errorf("Error")
	}

	if isNextCharacterIsNumber(3, "12a3") {
		t.Errorf("Error")
	}
}

func TestToPostfix(t *testing.T) {
	s := convertToPostfix("5*(6+2)-12/4", nil)

	st := ""
	for s.Len() > 0 {
		st += s.Pop().(string)
	}

	if st != "562+*124/-" {
		t.Errorf("Expression in stack is wrong")
	}
}

func TestToPostfixWithVariable(t *testing.T) {
	s := convertToPostfix("5*(6+pi)-12/4", nil)

	st := ""
	for s.Len() > 0 {
		st += s.Pop().(string)
	}

	if st != "56pi+*124/-" {
		t.Errorf("Expression in stack is wrong")
	}
}

func TestResolveExpression(t *testing.T) {
	defs := make(map[string]exp)
	r, err := resolve(convertToPostfix("5*(6+2)-12/4", nil), &defs, nil)

	if err != nil {
		t.Errorf(err.Error())
	}

	if r != "37" {
		t.Errorf("Expression wrong resolution")
	}
}

func TestResolveExpressionWithDefinitions(t *testing.T) {
	defs := make(map[string]exp)
	defs["pi"] = exp{"3"}

	s := convertToPostfix("2*pi+5", nil)
	r, err := resolve(s, &defs, nil)

	if err != nil {
		t.Errorf(err.Error())
	}

	if r != "11" {
		t.Errorf("Expression wrong resolution")
	}
}

func TestExpression(t *testing.T) {
	p := newPattern()
	expr := exp{"1+2-3*4/5"}

	if expr.isSpecialExpression(p) {
		t.Errorf("it's not a special expression")
	}

	expr = exp{"1+2-3*4/var"}
	if !expr.isSpecialExpression(p) {
		t.Errorf("it's a special expression")
	}
}

func TestFunctionCall(t *testing.T) {
	s := convertToPostfix("print(3+5)", nil)
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
	vars := make(map[string]exp)
	r, err := resolve(convertToPostfix("0.5+1.6", nil), &vars, nil)

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
	cc := newCompiler()
	cc.CompileLine("mysum(x, y)=x+y")

	if (*cc.Funcs)["mysum"].Exp != "x+y" {
		t.Error("stack for functions failed")
	}

	if len((*cc.Funcs)["mysum"].Args) != 2 {
		t.Error("stack for functions failed")
	}

}

func TestDeclaredFunctionCall(t *testing.T) {
	cc := newCompiler()
	cc.CompileLine("mydiv(x, y)=x/y")
	cc.CompileLine("result=mydiv(4, 2)+1")
	cc.CompileLine("print(result)") //should be 6
}

func TestDeclaredFunctionCall2(t *testing.T) {
	cc := newCompiler()
	cc.CompileLine("mysum(x, y)=x+y")
	cc.CompileLine("age=5")
	cc.CompileLine("result=mysum(3+age, 2)+1")
	cc.CompileLine("print(result)") //should be 6
}

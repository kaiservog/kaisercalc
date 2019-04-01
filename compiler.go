package main

import (
	"regexp"
//	"fmt"
	"errors"
	"strings"
)

type RegexpCalc struct {
	ReDefinitions *regexp.Regexp
	ReVariableExpression *regexp.Regexp
	ReExpression *regexp.Regexp
	ReFuncCall *regexp.Regexp
}

func NewRegexpCalc() *RegexpCalc {
	re := &RegexpCalc{}

	c, _ := regexp.Compile("\\w[\\w]*=.*")
	re.ReDefinitions = c

	c, _ = regexp.Compile("[a-zA-Z_]+")
	re.ReVariableExpression = c

	c, _ = regexp.Compile("[0-9]+")
	re.ReExpression = c

	c, _ = regexp.Compile("[a-zA-Z_]+\\(.*\\)")
	re.ReFuncCall = c

	return re
}

type CalcExp struct {
	Exp string
}

func (ce *CalcExp) IsSpecialExpression(rc *RegexpCalc) bool {
	return rc.ReVariableExpression.MatchString(ce.Exp)
}

func (ce *CalcExp) IsExpression(rc *RegexpCalc) bool {
	return rc.ReVariableExpression.FindString(ce.Exp) != ""
}


func (ce *CalcExp) IsFunctionCall(rc *RegexpCalc) bool {
	return !rc.ReFuncCall.MatchString(ce.Exp)
}

type CalcCompiler struct {
	data string
	Defs map[string]CalcExp
	rc *RegexpCalc
}

func NewCalcCompiler() *CalcCompiler {
	cc := &CalcCompiler{}
	cc.Defs = make(map[string]CalcExp)
	cc.rc = NewRegexpCalc()

	return cc
}

func (cc *CalcCompiler) CompileLine(line string) error {
	//TODO improve that part

	if strings.Contains(line, "=") {
		e := strings.Split(line, "=")

		if _, ok := cc.Defs[e[0]]; ok {
			return errors.New("variable already defined")
		}

		dom := CalcExp{e[1]}

		if dom.IsExpression(cc.rc) {
			s := ConvertToPostfix(dom.Exp)
			result, err := Resolve(s, &cc.Defs)
			if err != nil {
				return err
			}
			
			cc.Defs[e[0]] = CalcExp{result}

		} else {
			cc.Defs[e[0]] = dom
		}
	} else {
		s := ConvertToPostfix(line)
		_, err := Resolve(s, &cc.Defs)
		if err != nil {
			return err
		}
	}

	return nil
}
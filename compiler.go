package main

import (
	"regexp"

	//	"fmt"
	"errors"
	"strings"
)

type RegexpCalc struct {
	ReDefinitions        *regexp.Regexp
	ReVariableExpression *regexp.Regexp
	ReExpression         *regexp.Regexp
	ReFuncCall           *regexp.Regexp
	ReFuncArgs           *regexp.Regexp
}

func NewRegexpCalc() *RegexpCalc {
	re := &RegexpCalc{}

	c, _ := regexp.Compile("\\w[\\w]*=.*")
	re.ReDefinitions = c

	c, _ = regexp.Compile("[a-zA-Z_]+")
	re.ReVariableExpression = c

	c, _ = regexp.Compile("[0-9\\.]+")
	re.ReExpression = c

	c, _ = regexp.Compile("[a-zA-Z_]+\\(.*\\)")
	re.ReFuncCall = c

	c, _ = regexp.Compile(`\((.*)\)`)
	re.ReFuncArgs = c

	return re
}

type CalcExp struct {
	Exp string
}

type CalcFuncExp struct {
	Exp  string
	Args []string
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
	data  string
	Defs  map[string]CalcExp
	Funcs map[string]CalcFuncExp
	rc    *RegexpCalc
}

func NewCalcCompiler() *CalcCompiler {
	cc := &CalcCompiler{}
	cc.Defs = make(map[string]CalcExp)
	cc.Funcs = make(map[string]CalcFuncExp)
	cc.rc = NewRegexpCalc()

	return cc
}

func NewCalcFuncExp(name, rs string, cc *CalcCompiler) CalcFuncExp {
	args := cc.rc.ReFuncArgs.FindStringSubmatch(name)
	if len(args) > 0 {
		args = strings.Split(args[1], ",") //its group is 0?
	} else {
		args = make([]string, 0)
	}

	return CalcFuncExp{
		Exp:  rs,
		Args: args}
}

func (cc *CalcCompiler) checkDuplicate(name string) error {
	if _, ok := cc.Funcs[name]; ok {
		return errors.New("name already defined (" + name + ")")
	}
	if _, ok := cc.Defs[name]; ok {
		return errors.New("name already defined (" + name + ")")
	}

	return nil
}

func (cc *CalcCompiler) CompileLine(line string) error {
	//TODO improve

	line = cleanup(line)

	if strings.Contains(line, "=") {
		e := strings.Split(line, "=")
		leftSide := e[0]
		rightSide := e[1]

		if strings.Contains(leftSide, "(") { //It's a function
			name := strings.Split(leftSide, "(")[0]
			err := cc.checkDuplicate(name)
			if err != nil {
				return err
			}

			cc.Funcs[name] = NewCalcFuncExp(leftSide, rightSide, cc)

		} else { //It's a variable
			err := cc.checkDuplicate(leftSide)
			if err != nil {
				return err
			}

			exp := CalcExp{rightSide}
			s := ConvertToPostfix(exp.Exp)
			result, err := Resolve(s, &cc.Defs)
			if err != nil {
				return err
			}

			cc.Defs[leftSide] = CalcExp{result}
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

func cleanup(line string) string {
	return strings.ReplaceAll(line, " ", "")
}

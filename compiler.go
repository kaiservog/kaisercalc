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

type calcExp struct {
	Exp string
}

type calcFuncExp struct {
	Exp  string
	Args []string
}

func (ce *calcExp) isSpecialExpression(rc *RegexpCalc) bool {
	return rc.ReVariableExpression.MatchString(ce.Exp)
}

func (ce *calcExp) isExpression(rc *RegexpCalc) bool {
	return rc.ReVariableExpression.FindString(ce.Exp) != ""
}

func (ce *calcExp) isFunctionCall(rc *RegexpCalc) bool {
	return !rc.ReFuncCall.MatchString(ce.Exp)
}

type CalcCompiler struct {
	data  string
	Defs  *map[string]calcExp
	Funcs *map[string]calcFuncExp
	rc    *RegexpCalc
}

func newCalcCompiler() *CalcCompiler {
	cc := &CalcCompiler{}

	defs := make(map[string]calcExp)
	funcs := make(map[string]calcFuncExp)

	cc.Defs = &defs
	cc.Funcs = &funcs

	cc.rc = NewRegexpCalc()

	return cc
}

func newCalcFuncExp(name, rs string, cc *CalcCompiler) calcFuncExp {
	args := cc.rc.ReFuncArgs.FindStringSubmatch(name)
	if len(args) > 0 {
		args = strings.Split(args[1], ",") //its group is 0?
	} else {
		args = make([]string, 0)
	}

	return calcFuncExp{
		Exp:  rs,
		Args: args}
}

func (cc *CalcCompiler) checkDuplicate(name string) error {
	if _, ok := (*cc.Funcs)[name]; ok {
		return errors.New("name already defined (" + name + ")")
	}
	if _, ok := (*cc.Defs)[name]; ok {
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

			(*cc.Funcs)[name] = newCalcFuncExp(leftSide, rightSide, cc)

		} else { //It's a variable
			err := cc.checkDuplicate(leftSide)
			if err != nil {
				return err
			}

			exp := calcExp{rightSide}
			s := ConvertToPostfix(exp.Exp, cc.Funcs)
			result, err := resolve(s, cc.Defs, cc.Funcs)
			if err != nil {
				return err
			}

			(*cc.Defs)[leftSide] = calcExp{result}
		}

	} else {
		s := ConvertToPostfix(line, nil)
		_, err := resolve(s, cc.Defs, cc.Funcs)
		if err != nil {
			return err
		}
	}

	return nil
}

func cleanup(line string) string {
	return strings.ReplaceAll(line, " ", "")
}

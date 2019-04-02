package main

import "strings"

type exp struct {
	Exp string
}

type funcExp struct {
	Exp  string
	Args []string
}

func (ce *exp) isSpecialExpression(p *pattern) bool {
	return p.variableExpression.MatchString(ce.Exp)
}

func (ce *exp) isExpression(p *pattern) bool {
	return p.variableExpression.FindString(ce.Exp) != ""
}

func (ce *exp) isFunctionCall(p *pattern) bool {
	return !p.funcCall.MatchString(ce.Exp)
}

func newFuncExp(name, rs string, comp *compiler) funcExp {
	args := comp.rc.funcArgs.FindStringSubmatch(name)
	if len(args) > 0 {
		args = strings.Split(args[1], ",")
	} else {
		args = make([]string, 0)
	}

	return funcExp{
		Exp:  rs,
		Args: args}
}

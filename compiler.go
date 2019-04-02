package main

import (
	//	"fmt"
	"errors"
	"strings"
)

type compiler struct {
	data  string
	Vars  *map[string]exp
	Funcs *map[string]funcExp
	rc    *pattern
}

func newCompiler() *compiler {
	cc := &compiler{}

	vars := make(map[string]exp)
	funcs := make(map[string]funcExp)

	cc.Vars = &vars
	cc.Funcs = &funcs

	cc.rc = newPattern()

	return cc
}

func (cc *compiler) checkDuplicate(name string) error {
	if _, ok := (*cc.Funcs)[name]; ok {
		return errors.New("name already defined (" + name + ")")
	}
	if _, ok := (*cc.Vars)[name]; ok {
		return errors.New("name already defined (" + name + ")")
	}

	return nil
}

func (cc *compiler) CompileLine(line string) error {
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

			(*cc.Funcs)[name] = newFuncExp(leftSide, rightSide, cc)

		} else { //It's a variable
			err := cc.checkDuplicate(leftSide)
			if err != nil {
				return err
			}

			expr := exp{rightSide}
			s := ConvertToPostfix(expr.Exp, cc.Funcs)
			result, err := resolve(s, cc.Vars, cc.Funcs)
			if err != nil {
				return err
			}

			(*cc.Vars)[leftSide] = exp{result}
		}

	} else {
		s := ConvertToPostfix(line, nil)
		_, err := resolve(s, cc.Vars, cc.Funcs)
		if err != nil {
			return err
		}
	}

	return nil
}

func cleanup(line string) string {
	return strings.ReplaceAll(line, " ", "")
}

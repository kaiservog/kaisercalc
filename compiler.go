package main

import (
	//	"fmt"
	"errors"
	"fmt"
	"strings"
)

type compiler struct {
	data      string
	Vars      *map[string]exp
	Funcs     *map[string]funcExp
	rc        *pattern
	simpleExp bool
	root      string //local where user call this program
}

func newCompiler(root string, simpleExp bool) *compiler {
	cc := &compiler{}

	vars := make(map[string]exp)
	funcs := make(map[string]funcExp)

	cc.Vars = &vars
	cc.Funcs = &funcs

	cc.rc = newPattern()
	cc.root = root
	cc.simpleExp = simpleExp
	return cc
}

func (cc *compiler) checkDuplicateName(name string) error {
	if _, ok := (*cc.Funcs)[name]; ok {
		return errors.New("name already defined (" + name + ")")
	}
	if _, ok := (*cc.Vars)[name]; ok {
		return errors.New("name already defined (" + name + ")")
	}

	return nil
}

func (cc *compiler) CompileLine(line string) error {
	if strings.Contains(line, "=") {
		line = cleanup(line)
		e := strings.Split(line, "=")
		leftSide := e[0]
		rightSide := e[1]

		if strings.Contains(leftSide, "(") { //It's a function
			name := strings.Split(leftSide, "(")[0]
			err := cc.checkDuplicateName(name)
			if err != nil {
				return err
			}

			f := newFuncExp(leftSide, rightSide, cc)
			(*cc.Funcs)[name] = f

		} else { //It's a variable
			err := cc.checkDuplicateName(leftSide)
			if err != nil {
				return err
			}

			expr := exp{rightSide}
			s := convertToPostfix(expr.Exp, cc.Funcs)
			result, err := resolve(s, cc.Vars, cc.Funcs)
			if err != nil {
				return err
			}
			(*cc.Vars)[leftSide] = exp{result}
		}
	} else if isImport(line) {
		names := strings.Split(line, " ")
		if len(names) != 3 {
			return errors.New("import syntax must be 'import alias path/to/file'")
		}

		//tree := strings.Split(names[2], "/")
		importPath := names[2]
		if cc.root != "" {
			importPath = cc.root + "/" + names[2]
		}

		comp := processFile(importPath)
		mixVarsAndFuncs(comp.Vars, cc.Vars, comp.Funcs, cc.Funcs, names[1])
	} else {
		line = cleanup(line)
		s := convertToPostfix(line, nil)
		r, err := resolve(s, cc.Vars, cc.Funcs)
		if err != nil {
			return err
		}

		if cc.simpleExp {
			fmt.Println(r)
		}
	}

	return nil
}

func mixVarsAndFuncs(vvSrc, vvTgt *map[string]exp, ffSrc, ffTgt *map[string]funcExp, importName string) {
	for key, value := range *vvSrc {
		(*vvTgt)[importName+"."+key] = value
		//REMOVE THIS
	}

	for key, value := range *ffSrc {
		(*ffTgt)[importName+"."+key] = value
		//REMOVE THIS
	}
}

func cleanup(line string) string {
	line = removeComments(line)
	insideStr := false
	l := len(line) - 1

	for i := l; i >= 0; i-- { //cannot remove spaces from strings
		c := string([]rune(line)[i])

		if c == " " && !insideStr {
			line = line[:i] + line[i+1:]
		}

		if c == "'" {
			insideStr = !insideStr
		}
	}

	return line
}

func isImport(line string) bool {
	p := newPattern()
	return p.importSyntx.MatchString(line)
}

func removeComments(line string) string {
	p := newPattern()
	return p.comments.ReplaceAllString(line, "")
}

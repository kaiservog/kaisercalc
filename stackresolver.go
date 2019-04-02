package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/golang-collections/collections/stack"
)

func resolve(s *stack.Stack, defs *map[string]calcExp, funcs *map[string]calcFuncExp) (string, error) {
	rs := stack.New()

	for s.Len() > 0 {
		elm := s.Pop().(string)

		if isOperator(elm) {
			a, _ := strconv.ParseFloat(rs.Pop().(string), 64)
			b, _ := strconv.ParseFloat(rs.Pop().(string), 64)

			if elm == "*" {
				f := fmt.Sprintf("%g", b*a)
				rs.Push(f)
			} else if elm == "/" {
				f := fmt.Sprintf("%g", b/a)
				rs.Push(f)
			} else if elm == "+" {
				f := fmt.Sprintf("%g", b+a)
				rs.Push(f)
			} else if elm == "-" {
				f := fmt.Sprintf("%g", b-a)
				rs.Push(f)
			}
		} else if isNumber(elm) {
			rs.Push(elm)
		} else if isSystemFunction(elm) {
			fmt.Println(rs.Pop())
		} else if isDefinedFunction(elm, funcs) {
			r, err := callDefinedFunction(elm, rs, funcs)
			if err != nil {
				return "", nil
			}
			rs.Push(r)
		} else { //is definition
			if def, ok := (*defs)[elm]; ok {
				rs.Push(def.Exp)
			} else {
				return "", errors.New("variable not defined " + elm)
			}
		}
	}
	if rs.Len() == 1 {
		return rs.Pop().(string), nil
	}

	return "", nil
}

func getReverseArr(a []string) []string {
	r := make([]string, len(a))

	for i := len(a) - 1; i >= 0; i-- {
		r[len(a)-i-1] = a[i]
	}
	return r
}

func callDefinedFunction(name string, s *stack.Stack, funcs *map[string]calcFuncExp) (string, error) {
	fDefs := make(map[string]calcExp)
	reversed_args := getReverseArr((*funcs)[name].Args)
	for _, name := range reversed_args {
		fDefs[name] = calcExp{s.Pop().(string)}
	}

	newStack := ConvertToPostfix((*funcs)[name].Exp, funcs)
	return resolve(newStack, &fDefs, funcs)
}

func showStack(s *stack.Stack) {
	t := stack.New()

	fmt.Println()
	fmt.Println("stack")

	for s.Len() > 0 {
		value := s.Pop()
		fmt.Print(value, "-")
		t.Push(value)
	}
	fmt.Println()

	for t.Len() > 0 {
		value := t.Pop()
		s.Push(value)
	}
}

func ConvertToPostfix(exp string, funcs *map[string]calcFuncExp) *stack.Stack {
	s := stack.New()
	temp := stack.New()
	buffer := ""

	for i := 0; i < len(exp); i++ {
		c := string([]rune(exp)[i])

		if isNumber(c) || isVariablePart(c) {
			if IsNextCharNumberOrDefinitions(i, exp) {
				buffer += c
			} else {
				if isSystemFunction(buffer + c) {
					temp.Push(buffer + c)
					buffer = ""
				} else if isDefinedFunction(buffer+c, funcs) {
					temp.Push(buffer + c)
					buffer = ""
				} else {
					s.Push(buffer + c)
					buffer = ""
				}
			}
		} else if c == "," {
			continue
		} else if c == "(" {
			temp.Push(c)
		} else if c == ")" {
			for temp.Len() > 0 && temp.Peek() != "(" {
				s.Push(temp.Pop())
			}
			temp.Pop()
		} else if !isNumber(c) {
			if temp.Len() == 0 || temp.Peek() == "(" {
				temp.Push(c)
			} else {
				for temp.Len() > 0 && temp.Peek() != "(" && isPrecedenceHigher(temp.Peek().(string), c) {
					s.Push(temp.Pop())
				}
				temp.Push(c)
			}
		}
	}

	for temp.Len() > 0 {
		s.Push(temp.Pop())
	}

	//invertendo
	for s.Len() > 0 {
		temp.Push(s.Pop())
	}

	//showStack(temp)
	return temp
}

func isVariablePart(t string) bool {
	//TODO pool of compiled regexp
	c, _ := regexp.Compile("[a-zA-Z_]")
	return c.MatchString(t)
}

func isOperator(v string) bool {
	return v == "+" || v == "-" || v == "*" || v == "/"
}

func isSystemFunction(s string) bool {
	return s == "print"
}

func isDefinedFunction(name string, funcs *map[string]calcFuncExp) bool {
	if funcs != nil {
		_, ok := (*funcs)[name]
		return ok
	}
	return false
}

func operatorPrecedence(op string) int {
	if op == "*" {
		return 4
	}

	if op == "/" {
		return 4
	}

	if op == "+" {
		return 3
	}

	if op == "-" {
		return 3
	}

	return 5
}

func IsNextCharacterIsNumber(i int, s string) bool {
	return i+1 < len(s) && isNumber(string([]rune(s)[i+1]))
}

func IsNextCharNumberOrDefinitions(i int, s string) bool {
	if i+1 >= len(s) {
		return false
	}

	cre, _ := regexp.Compile("[a-zA-Z_]")
	c := string([]rune(s)[i+1])

	return isNumber(c) || cre.MatchString(c)
}

func isPrecedenceHigher(x, y string) bool {
	return operatorPrecedence(x) >= operatorPrecedence(y)
}

//is not one character only
func isNumber(v string) bool {
	//TODO compile regexp pool
	rc := NewRegexpCalc()
	return rc.ReExpression.MatchString(v)
}

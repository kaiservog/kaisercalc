package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang-collections/collections/stack"
)

func resolve(s *stack.Stack, vars *map[string]exp, funcs *map[string]funcExp) (string, error) {
	rs := stack.New()
	p := newPattern()

	for s.Len() > 0 {
		elm := s.Pop().(string)
		if isOperator(elm) {
			a, _ := strconv.ParseFloat(rs.Pop().(string), 64)

			if elm == "*" {
				b, _ := strconv.ParseFloat(rs.Pop().(string), 64)
				f := fmt.Sprintf("%g", b*a)
				rs.Push(f)
			} else if elm == "/" {
				b, _ := strconv.ParseFloat(rs.Pop().(string), 64)
				f := fmt.Sprintf("%g", b/a)
				rs.Push(f)
			} else if elm == "+" {
				b, _ := strconv.ParseFloat(rs.Pop().(string), 64)
				f := fmt.Sprintf("%g", b+a)
				rs.Push(f)
			} else if elm == "-" {
				b, _ := strconv.ParseFloat(rs.Pop().(string), 64)
				f := fmt.Sprintf("%g", b-a)
				rs.Push(f)
			} else if elm == "^" {
				b, _ := strconv.ParseFloat(rs.Pop().(string), 64)
				f := fmt.Sprintf("%g", elevation(b, int(a)))
				rs.Push(f)
			} else if elm == "!" {
				f := fmt.Sprintf("%g", factorial(a))
				rs.Push(f)
			}
		} else if isNumber(elm, p) {
			rs.Push(elm)
		} else if isSystemFunction(elm) {
			reverseStack(rs)
			for rs.Len() > 0 {
				arg := rs.Pop()
				if elm == "print" {
					fmt.Print(arg.(string) + " ")
				} else {
					fmt.Println(arg)
				}
			}
		} else if isDefinedFunction(elm, funcs) {
			r, err := callDefinedFunction(elm, rs, funcs)
			if err != nil {
				return "", nil
			}
			rs.Push(r)
		} else { //is definition
			if def, ok := (*vars)[elm]; ok {
				rs.Push(def.Exp)
			} else {
				rs.Push(elm)
				//if not var so str? good idea?
				//return "", errors.New("variable not defined " + elm)
			}
		}
	}
	if rs.Len() == 1 {
		return rs.Pop().(string), nil
	}

	return "", nil
}

func elevation(x float64, y int) float64 {
	for i := 0; i < y; i++ {
		x = x * x
	}

	return x
}

func factorial(x float64) float64 {
	if x == 0 {
		return 1
	}

	return x * factorial(x-1)
}

func getReverseArr(a []string) []string {
	r := make([]string, len(a))

	for i := len(a) - 1; i >= 0; i-- {
		r[len(a)-i-1] = a[i]
	}
	return r
}

func callDefinedFunction(name string, s *stack.Stack, funcs *map[string]funcExp) (string, error) {
	funcVars := make(map[string]exp)
	reversedArgs := getReverseArr((*funcs)[name].Args)
	for _, name := range reversedArgs {
		funcVars[name] = exp{s.Pop().(string)}
	}

	newStack := convertToPostfix((*funcs)[name].Exp, funcs)
	return resolve(newStack, &funcVars, funcs)
}

func showStack(s *stack.Stack) {
	t := stack.New()

	fmt.Println()
	fmt.Println("stack")

	fmt.Print("-")
	for s.Len() > 0 {
		value := s.Pop()
		fmt.Print(value, " ")
		t.Push(value)
	}
	fmt.Print("-")
	fmt.Println()

	for t.Len() > 0 {
		value := t.Pop()
		s.Push(value)
	}
}

//returns first string in expression from start idx
//return last index from string returned
func readStr(exp string, idx int) (string, error) {
	if string(exp[idx]) != "'" {
		return "", errors.New("First character must be '")
	}

	subexp := exp[idx:]

	end := strings.Index(subexp[1:], "'")
	return subexp[1 : end+1], nil
}

func convertToPostfix(exp string, funcs *map[string]funcExp) *stack.Stack {
	s := stack.New()
	temp := stack.New()
	buffer := ""
	pattn := newPattern()

	for i := 0; i < len(exp); i++ {
		c := string([]rune(exp)[i])
		if isNumber(c, pattn) || isVariablePart(c, pattn) {
			if isNextCharNumberOrDefinitions(i, exp, pattn) {
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
		} else if c == "'" {
			r, _ := readStr(exp, i)
			s.Push(r)
			i = i + len(r) + 1
		} else if c == "," {
			continue
		} else if c == "(" {
			temp.Push(c)
		} else if c == ")" {
			for temp.Len() > 0 && temp.Peek() != "(" {
				s.Push(temp.Pop())
			}
			temp.Pop()
		} else if !isNumber(c, pattn) {
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

	//showStack(temp)
	reverseStack(s)
	return s
}

func isVariablePart(t string, p *pattern) bool {
	//alterei aqui
	return p.variableExpression.MatchString(t)
}

func isOperator(v string) bool {
	return v == "+" || v == "-" || v == "*" || v == "/" || v == "^" || v == "!"
}

func isSystemFunction(s string) bool {
	return s == "print" || s == "println"
}

func isDefinedFunction(name string, funcs *map[string]funcExp) bool {
	if funcs != nil {
		_, ok := (*funcs)[name]
		return ok
	}
	return false
}

func operatorPrecedence(op string) int {
	if op == "!" {
		return 7
	}

	if op == "^" {
		return 6
	}

	if op == "*" {
		return 5
	}

	if op == "/" {
		return 4
	}

	if op == "+" {
		return 3
	}

	if op == "-" {
		return 2
	}

	//functions will return 5
	return 5
}

func isNextCharacterIsNumber(i int, s string, p *pattern) bool {
	return i+1 < len(s) && isNumber(string([]rune(s)[i+1]), p)
}

func isNextCharNumberOrDefinitions(i int, s string, p *pattern) bool {
	if i+1 >= len(s) {
		return false
	}

	c := string([]rune(s)[i+1])
	return isNumber(c, p) || p.variableExpression.MatchString(c)
}

func isPrecedenceHigher(x, y string) bool {
	return operatorPrecedence(x) >= operatorPrecedence(y)
}

func isNumber(val string, p *pattern) bool {
	return p.expression.MatchString(val)
}

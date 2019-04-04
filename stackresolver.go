package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/golang-collections/collections/stack"
)

func resolve(s *stack.Stack, vars *map[string]exp, funcs *map[string]funcExp) (string, error) {
	rs := stack.New()

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
		} else if isNumber(elm) {
			rs.Push(elm)
		} else if isSystemFunction(elm) {
			arg := rs.Pop()
			fmt.Println(arg)
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
				return "", errors.New("variable not defined " + elm)
			}
		}
	}
	if rs.Len() == 1 {
		return rs.Pop().(string), nil
	}

	return "", nil
}

func elevation(x float64, y int) float64 {
	for i :=0; i < y; i++ {
		x = x*x
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

func convertToPostfix(exp string, funcs *map[string]funcExp) *stack.Stack {
	s := stack.New()
	temp := stack.New()
	buffer := ""

	for i := 0; i < len(exp); i++ {
		c := string([]rune(exp)[i])

		if isNumber(c) || isVariablePart(c) {
			if isNextCharNumberOrDefinitions(i, exp) {
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

	//reversing
	for s.Len() > 0 {
		temp.Push(s.Pop())
	}

	showStack(temp)
	return temp
}

func isVariablePart(t string) bool {
	//TODO pool of compiled regexp
	c, _ := regexp.Compile("[a-zA-Z_.]")
	return c.MatchString(t)
}

func isOperator(v string) bool {
	return v == "+" || v == "-" || v == "*" || v == "/" || v == "^" || v== "!"
}

func isSystemFunction(s string) bool {
	return s == "print"
}

func isDefinedFunction(name string, funcs *map[string]funcExp) bool {
	if funcs != nil {
		_, ok := (*funcs)[name]
		return ok
	}
	return false
}

func operatorPrecedence(op string) int {
	if op == "^" {
		return 6
	}

	if op == "*" {
		return 5
	}

	if op == "/" {
		return 4
	}

	if op == "!" {
		return 5
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

func isNextCharacterIsNumber(i int, s string) bool {
	return i+1 < len(s) && isNumber(string([]rune(s)[i+1]))
}

func isNextCharNumberOrDefinitions(i int, s string) bool {
	if i+1 >= len(s) {
		return false
	}

	cre, _ := regexp.Compile("[a-zA-Z_\\.]")
	c := string([]rune(s)[i+1])

	return isNumber(c) || cre.MatchString(c)
}

func isPrecedenceHigher(x, y string) bool {
	return operatorPrecedence(x) >= operatorPrecedence(y)
}

func isNumber(val string) bool {
	//TODO compile regexp pool
	p := newPattern()
	return p.expression.MatchString(val)
}

package main

import (
	"strconv"
	"fmt"
	"github.com/golang-collections/collections/stack"
	"regexp"
	"errors"

)

func Resolve(s *stack.Stack, defs *map[string]CalcExp) (string, error) {
	rs := stack.New()

	for s.Len() > 0 {
		elm := s.Pop().(string)

		if IsOperator(elm) {
			a, _ := strconv.Atoi(rs.Pop().(string))
			b, _ := strconv.Atoi(rs.Pop().(string))

			if elm == "*" {
				rs.Push( strconv.Itoa(b * a) )
			} else if elm == "/" {
				rs.Push( strconv.Itoa(b / a) )
			} else if elm == "+" {
				rs.Push( strconv.Itoa(b + a) )				
			} else if elm == "-" {
				rs.Push( strconv.Itoa(b - a) )
			}
		} else if isNumber(elm) {
			rs.Push(elm)
		} else if IsSystemFunction(elm) {
			fmt.Println(rs.Pop())
		} else { //is definition
			if def, ok := (*defs)[elm]; ok {
				rs.Push(def.Exp)
			} else {
				return "", errors.New("variable not defined " + elm)
			}
		}
	}
	if rs.Len() == 1{
		return rs.Pop().(string), nil
	}

	return "", nil
}

func ConvertToPostfix(exp string) *stack.Stack {
	s := stack.New()
	temp := stack.New()
	buffer := ""

	for i:=0; i < len(exp); i++ {
		c := string([]rune(exp)[i])

		if isNumber(c) || isVariablePart(c) {
			if IsNextCharNumberOrDefinitions(i, exp) {
				buffer += c
			} else {
				if IsSystemFunction(buffer + c) {
					temp.Push(buffer + c)
					buffer = ""
				} else {
					s.Push(buffer + c)
					buffer = ""
				}
			}
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
				for temp.Len() > 0 && temp.Peek() != "(" && isPrecedenceHigher(temp.Peek().(string), c)  {
					s.Push(temp.Pop())
				}
				temp.Push(c)
			}
		}
	}

	for temp.Len() > 0  {
		s.Push(temp.Pop())
	}

	//invertendo
	for s.Len() > 0 {
		temp.Push(s.Pop())
	}

	return temp
}

func isVariablePart(t string) bool{
	//TODO pool of compiled regexp
	c, _ := regexp.Compile("[a-zA-Z_]")
	return c.MatchString(t)
}

func IsOperator(v string) bool {
	return v == "+" || v == "-" || v == "*" || v == "/"
}

func IsSystemFunction(s string) bool {
	return s == "print"
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

	return 0
}

func IsNextCharacterIsNumber(i int, s string) bool {
	return i+1 < len(s) && isNumber(string([]rune(s)[i+1])) 
}

func IsNextCharNumberOrDefinitions(i int, s string) bool {
	if  i+1 >= len(s) {
		return false
	}

	cre, _ := regexp.Compile("[a-zA-Z_]")
	c := string([]rune(s)[i+1])

	return isNumber(c) || cre.MatchString(c)
}

func isPrecedenceHigher(x, y string) bool {
	return operatorPrecedence(x) >= operatorPrecedence(y)
}

func isNumber(v string) bool {
	if _, err := strconv.Atoi(v); err == nil {
		return true
	} else {
		return false
	}
}
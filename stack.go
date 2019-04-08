package main

import "github.com/golang-collections/collections/stack"

func reverseStack(s *stack.Stack) {
	var tmp []string

	for s.Len() > 0 {
		tmp = append(tmp, s.Pop().(string))
	}

	for i := 0; i < len(tmp); i++ {
		s.Push(tmp[i])
	}
}

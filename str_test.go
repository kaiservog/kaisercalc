package main

import (
	"fmt"
	"testing"

	"github.com/golang-collections/collections/stack"
)

func TestPrintStr(t *testing.T) {
	s := stack.New()
	s.Push("print")
	s.Push("teste")
	vars := make(map[string]exp)
	vars["teste"] = exp{"hello world!"}
	_, err := resolve(s, &vars, nil)

	if err != nil {
		t.Error(err)
	}
}
func TestConvertPosfixStr(t *testing.T) {
	funcs := make(map[string]funcExp)
	s := convertToPostfix("'hello world!'", &funcs)

	r := s.Pop()
	fmt.Println(r)
	if r != "hello world!" {
		t.Error("not correctly converted")
	}
}

func TestPrintCallStr(t *testing.T) {
	c := newCompiler("", false)
	c.CompileLine("print('result is', ' ', 42)")
}

func TestCleanLine(t *testing.T) {
	s := cleanup("print('result is', ' ', 42)")
	fmt.Println(s)
	if s != "print('result is',' ',42)" {
		t.Error("spaces are not correctlt removed")
	}
}

func TestReadStr(t *testing.T) {
	str, _ := readStr("'teste'1234'test'", 0)
	if str != "teste" {
		t.Error("readStr wrong behavior")
	}

	str, _ = readStr("'teste'1234'test'", 11)
	if str != "test" {
		t.Error("readStr wrong behavior")
	}

}

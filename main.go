package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func main() {
	if contains(os.Args, "-e") {
		processExpression(os.Args[2])
	} else {
		p := os.Args[1]
		processFile(p)
	}
}

func processExpression(exp string) {
	exps := strings.Split(exp, ";")
	if len(exps) > 0 {
		cc := NewCalcCompiler()

		for _, ln := range exps {
			process(ln, cc)
		}
	}
}

func processFile(p string) {
	file, err := os.Open(p)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	cc := NewCalcCompiler()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ln := scanner.Text()
		process(ln, cc)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func process(ln string, cc *CalcCompiler) {
	err := cc.CompileLine(ln)
	if err != nil {
		log.Fatal(err)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

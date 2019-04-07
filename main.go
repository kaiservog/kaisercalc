package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	if contains(os.Args, "-e") {
		processExpression(os.Args[2])
	} else if len(os.Args) > 1 {
		p := os.Args[1]
		processFile(p)
	} else {
		showHelp()
	}
}

func processExpression(exp string) {
	exps := strings.Split(exp, ";")
	if len(exps) > 0 {
		cc := newCompiler("", true)

		for _, ln := range exps {
			process(ln, cc)
		}
	}
}

func processFile(p string) *compiler {
	root := extractRoot(p)

	file, err := os.Open(p)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	cc := newCompiler(root, false)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ln := scanner.Text()
		process(ln, cc)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return cc
}

func extractRoot(p string) string {
	idx := strings.LastIndex(p, "/")
	if idx != -1 {
		return p[:idx-1]
	}

	idx = strings.LastIndex(p, "\\")
	if idx != -1 {
		return p[:idx]
	}

	return ""
}

func process(ln string, cc *compiler) {
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

func showHelp() {
	help := `Kaisercalc Help

	kaisercalc -e "exp"      to resolve simple expression
	kaisercalc filename      to resolve expression in file
	
	examples
	kaisercalc -e "(3+1)/(2*3+1)"
	kaisercalc mymath.txt
	`

	fmt.Println(help)
}

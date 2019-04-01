package main

import (
	"bufio"
	"os"
//	"fmt"
	"log"
)

func main() {
	p := os.Args[1]
	process(p)
}

func process(p string) {
	file, err := os.Open(p)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	cc := NewCalcCompiler()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ln := scanner.Text()
		err = cc.CompileLine(ln)
		if err != nil {
			log.Fatal(err)	
		}
	}
	
	if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}
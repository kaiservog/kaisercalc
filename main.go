package main

import (
	"bufio"
	"os"
//	"fmt"
	"log"
)

func main() {
	file, err := os.Open(os.Args[1])
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

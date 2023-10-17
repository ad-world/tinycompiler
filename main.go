package main

import (
	"log"
	"os"
	"tinycompiler/lexer"
	"tinycompiler/parser"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Missing parameter, please provide filename.")
		return
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	contents := string(data)

	l := lexer.NewLexer(contents)
	p := parser.NewParser(l)

	parser.Program(p)
}

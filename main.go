package main

import (
	"log"
	"os"
	"strings"
	"tinycompiler/emitter"
	"tinycompiler/lexer"
	"tinycompiler/parser"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Missing parameter, please provide filename.")
		return
	}
	filename := os.Args[1]
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	splitFileName := strings.Split(filename, "/")
	lastToken := splitFileName[len(splitFileName)-1]
	splitLastToken := strings.Split(lastToken, ".")
	resultingFilename := splitLastToken[0]

	contents := string(data)

	l := lexer.NewLexer(contents)
	e := emitter.NewEmitter("results/" + resultingFilename + ".c")
	p := parser.NewParser(l, e)

	parser.Program(p)
	emitter.WriteFile(e)
}

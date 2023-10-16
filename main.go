package main

import (
	"fmt"
	"tinycompiler/lexer"
)

func main() {
	source := "IF+-123 foo*THEN/"
	l := lexer.NewLexer(source)

	token := lexer.GetToken(l)

	for token.TokenType != lexer.EOF {
		fmt.Println(token.TokenType)
		token = lexer.GetToken(l)
	}
}

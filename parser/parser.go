package parser

import (
	"fmt"
	"log"
	"tinycompiler/lexer"
)

type Parser struct {
	Symbols        map[string]bool
	LabelsDeclared map[string]bool
	LabelsGotoed   map[string]bool
	lexer          *lexer.Lexer
	CurToken       lexer.Token
	PeekToken      lexer.Token
}

func NewParser(l *lexer.Lexer) *Parser {
	parser := &Parser{
		Symbols:        make(map[string]bool),
		LabelsDeclared: make(map[string]bool),
		LabelsGotoed:   make(map[string]bool),
		lexer:          l,
		CurToken: lexer.Token{
			TokenType: lexer.UNKNOWN,
		},
		PeekToken: lexer.Token{
			TokenType: lexer.UNKNOWN,
		},
	}

	NextToken(parser)
	NextToken(parser)

	return parser
}

func CheckToken(p *Parser, kind lexer.TokenType) bool {
	return kind == p.CurToken.TokenType
}

func CheckPeek(p *Parser, kind lexer.TokenType) bool {
	return kind == p.PeekToken.TokenType
}

func Match(p *Parser, kind lexer.TokenType) {
	if !CheckToken(p, kind) {
		Abort(p, "Expected "+lexer.TokenTypeName(kind)+", got "+lexer.TokenTypeName(p.CurToken.TokenType))
	}
	NextToken(p)
}

func NextToken(p *Parser) {
	p.CurToken = p.PeekToken
	p.PeekToken = lexer.GetToken(p.lexer)
}

func Abort(p *Parser, message string) {
	log.Fatalln("Parsing error. " + message)
}

// The following functions are language parsers

func Program(p *Parser) {
	fmt.Println("PROGRAM")

	for CheckToken(p, lexer.NEWLINE) {
		NextToken(p)
	}

	for !CheckToken(p, lexer.EOF) {
		Statement(p)
	}

	for label := range p.LabelsGotoed {
		_, ok := p.LabelsDeclared[label]
		if ok == false {
			Abort(p, "Attempting to GOTO undeclared label: "+label)
		}
	}

}

func Statement(p *Parser) {
	if CheckToken(p, lexer.PRINT) {
		fmt.Println("STATEMENT-PRINT")
		NextToken(p)

		if CheckToken(p, lexer.STRING) {
			NextToken(p)
		} else {
			Expression(p)
		}
	} else if CheckToken(p, lexer.IF) {
		fmt.Println("STATEMENT-IF")
		NextToken(p)

		Comparison(p)
		Match(p, lexer.THEN)
		Newline(p)

		for !CheckToken(p, lexer.ENDIF) {
			Statement(p)
		}

		Match(p, lexer.ENDIF)
	} else if CheckToken(p, lexer.WHILE) {
		fmt.Println("STATEMENT-WHILE")
		NextToken(p)
		Comparison(p)
		Match(p, lexer.REPEAT)
		Newline(p)

		for !CheckToken(p, lexer.ENDWHILE) {
			Statement(p)
		}

		Match(p, lexer.ENDWHILE)
	} else if CheckToken(p, lexer.LABEL) {
		fmt.Println("STATEMENT-LABEL")
		NextToken(p)

		if p.LabelsDeclared[p.CurToken.TokenText] {
			Abort(p, "Label already exists: "+p.CurToken.TokenText)
		}

		p.LabelsDeclared[p.CurToken.TokenText] = true

		Match(p, lexer.IDENT)
	} else if CheckToken(p, lexer.GOTO) {
		fmt.Println("STATEMENT-GOTO")
		NextToken(p)
		p.LabelsGotoed[p.CurToken.TokenText] = true
		Match(p, lexer.IDENT)
	} else if CheckToken(p, lexer.LET) {
		fmt.Println("STATEMENT-LET")
		NextToken(p)

		_, ok := p.Symbols[p.CurToken.TokenText]
		if ok == false {
			p.Symbols[p.CurToken.TokenText] = true
		}

		Match(p, lexer.IDENT)
		Match(p, lexer.EQ)
		Expression(p)
	} else if CheckToken(p, lexer.INPUT) {
		fmt.Println("STATEMENT-INPUT")
		NextToken(p)

		_, ok := p.Symbols[p.CurToken.TokenText]
		if ok == false {
			p.Symbols[p.CurToken.TokenText] = true
		}

		Match(p, lexer.IDENT)
	} else {
		Abort(p, "Invalid statement at "+p.CurToken.TokenText+" ("+lexer.TokenTypeName(p.CurToken.TokenType)+")")
	}
	Newline(p)
}

func Newline(p *Parser) {
	fmt.Println("NEWLINE")
	Match(p, lexer.NEWLINE)

	for CheckToken(p, lexer.NEWLINE) {
		NextToken(p)
	}
}

func Comparison(p *Parser) {
	fmt.Println("COMPARISON")

	Expression(p)

	if IsComparisonOperator(p) {
		NextToken(p)
		Expression(p)
	} else {
		Abort(p, "Expected comparison operator at: "+p.CurToken.TokenText)
	}

	for IsComparisonOperator(p) {
		NextToken(p)
		Expression(p)
	}
}

func IsComparisonOperator(p *Parser) bool {
	return CheckToken(p, lexer.GT) ||
		CheckToken(p, lexer.GTEQ) ||
		CheckToken(p, lexer.LT) ||
		CheckToken(p, lexer.LTEQ) ||
		CheckToken(p, lexer.EQEQ) ||
		CheckToken(p, lexer.NOTEQ)
}

func Expression(p *Parser) {
	fmt.Println("EXPRESSION")

	Term(p)
	for CheckToken(p, lexer.PLUS) || CheckToken(p, lexer.MINUS) {
		NextToken(p)
		Term(p)
	}
}

func Term(p *Parser) {
	fmt.Println("TERM")
	Unary(p)

	for CheckToken(p, lexer.SLASH) || CheckToken(p, lexer.ASTERISK) {
		NextToken(p)
		Unary(p)
	}
}

func Unary(p *Parser) {
	fmt.Println("UNARY")

	if CheckToken(p, lexer.PLUS) || CheckToken(p, lexer.MINUS) {
		NextToken(p)
	}
	Primary(p)
}

func Primary(p *Parser) {
	fmt.Println("PRIMARY (" + p.CurToken.TokenText + ")")

	if CheckToken(p, lexer.NUMBER) {
		NextToken(p)
	} else if CheckToken(p, lexer.IDENT) {
		_, ok := p.Symbols[p.CurToken.TokenText]
		if ok == false {
			Abort(p, "Referencing variable before assignment: "+p.CurToken.TokenText)
		}
		NextToken(p)
	} else {
		Abort(p, "Unexpected token at "+p.CurToken.TokenText)
	}
}

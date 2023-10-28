package parser

import (
	"log"
	"tinycompiler/emitter"
	"tinycompiler/lexer"
)

type Parser struct {
	Symbols        map[string]bool
	LabelsDeclared map[string]bool
	LabelsGotoed   map[string]bool
	emitter        *emitter.Emitter
	lexer          *lexer.Lexer
	CurToken       lexer.Token
	PeekToken      lexer.Token
}

func NewParser(l *lexer.Lexer, e *emitter.Emitter) *Parser {
	parser := &Parser{
		Symbols:        make(map[string]bool),
		LabelsDeclared: make(map[string]bool),
		LabelsGotoed:   make(map[string]bool),
		lexer:          l,
		emitter:        e,
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
	emitter.HeaderLine(p.emitter, "#include <stdio.h>")
	emitter.HeaderLine(p.emitter, "int main(void) {")

	for CheckToken(p, lexer.NEWLINE) {
		NextToken(p)
	}

	for !CheckToken(p, lexer.EOF) {
		Statement(p)
	}

	emitter.EmitLine(p.emitter, "return 0;")
	emitter.EmitLine(p.emitter, "}")

	for label := range p.LabelsGotoed {
		_, ok := p.LabelsDeclared[label]
		if ok == false {
			Abort(p, "Attempting to GOTO undeclared label: "+label)
		}
	}

}

func Statement(p *Parser) {
	if CheckToken(p, lexer.PRINT) {
		NextToken(p)

		if CheckToken(p, lexer.STRING) {
			emitter.EmitLine(p.emitter, "printf(\""+p.CurToken.TokenText+"\\n\");")
			NextToken(p)
		} else {
			emitter.Emit(p.emitter, "printf(\"%"+".2f\\n\", (float)(")
			Expression(p)
			emitter.EmitLine(p.emitter, "));")
		}
	} else if CheckToken(p, lexer.IF) {
		NextToken(p)
		emitter.Emit(p.emitter, "if(")

		Comparison(p)
		Match(p, lexer.THEN)
		Newline(p)
		emitter.EmitLine(p.emitter, "){")

		for !CheckToken(p, lexer.ENDIF) {
			Statement(p)
		}

		Match(p, lexer.ENDIF)
		emitter.EmitLine(p.emitter, "}")
	} else if CheckToken(p, lexer.WHILE) {
		NextToken(p)
		emitter.Emit(p.emitter, "while(")
		Comparison(p)
		Match(p, lexer.REPEAT)
		Newline(p)
		emitter.EmitLine(p.emitter, "){")

		for !CheckToken(p, lexer.ENDWHILE) {
			Statement(p)
		}

		Match(p, lexer.ENDWHILE)
		emitter.EmitLine(p.emitter, "}")
	} else if CheckToken(p, lexer.LABEL) {
		NextToken(p)

		if p.LabelsDeclared[p.CurToken.TokenText] {
			Abort(p, "Label already exists: "+p.CurToken.TokenText)
		}

		p.LabelsDeclared[p.CurToken.TokenText] = true
		emitter.EmitLine(p.emitter, p.CurToken.TokenText+":")

		Match(p, lexer.IDENT)
	} else if CheckToken(p, lexer.GOTO) {
		NextToken(p)
		p.LabelsGotoed[p.CurToken.TokenText] = true
		emitter.EmitLine(p.emitter, "goto "+p.CurToken.TokenText+";")
		Match(p, lexer.IDENT)
	} else if CheckToken(p, lexer.LET) {
		NextToken(p)

		_, ok := p.Symbols[p.CurToken.TokenText]
		if ok == false {
			p.Symbols[p.CurToken.TokenText] = true
			emitter.HeaderLine(p.emitter, "float "+p.CurToken.TokenText+";")
		}

		emitter.Emit(p.emitter, p.CurToken.TokenText+" = ")
		Match(p, lexer.IDENT)
		Match(p, lexer.EQ)
		Expression(p)
		emitter.EmitLine(p.emitter, ";")
	} else if CheckToken(p, lexer.INPUT) {
		NextToken(p)

		_, ok := p.Symbols[p.CurToken.TokenText]
		if ok == false {
			p.Symbols[p.CurToken.TokenText] = true
			emitter.HeaderLine(p.emitter, "float "+p.CurToken.TokenText+";")
		}

		emitter.EmitLine(p.emitter, "if(0 == scanf(\"%"+"f\", &"+p.CurToken.TokenText+")) {")
		emitter.EmitLine(p.emitter, p.CurToken.TokenText+" = 0;")
		emitter.Emit(p.emitter, "scanf(\"%")
		emitter.EmitLine(p.emitter, "*s\");")
		emitter.EmitLine(p.emitter, "}")
		Match(p, lexer.IDENT)
	} else {
		Abort(p, "Invalid statement at "+p.CurToken.TokenText+" ("+lexer.TokenTypeName(p.CurToken.TokenType)+")")
	}
	Newline(p)
}

func Newline(p *Parser) {
	Match(p, lexer.NEWLINE)

	for CheckToken(p, lexer.NEWLINE) {
		NextToken(p)
	}
}

func Comparison(p *Parser) {
	Expression(p)

	if IsComparisonOperator(p) {
		emitter.Emit(p.emitter, p.CurToken.TokenText)
		NextToken(p)
		Expression(p)
	} else {
		Abort(p, "Expected comparison operator at: "+p.CurToken.TokenText)
	}

	for IsComparisonOperator(p) {
		emitter.Emit(p.emitter, p.CurToken.TokenText)
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
	Term(p)
	for CheckToken(p, lexer.PLUS) || CheckToken(p, lexer.MINUS) {
		emitter.Emit(p.emitter, p.CurToken.TokenText)
		NextToken(p)
		Term(p)
	}
}

func Term(p *Parser) {
	Unary(p)

	for CheckToken(p, lexer.SLASH) || CheckToken(p, lexer.ASTERISK) {
		emitter.Emit(p.emitter, p.CurToken.TokenText)
		NextToken(p)
		Unary(p)
	}
}

func Unary(p *Parser) {
	if CheckToken(p, lexer.PLUS) || CheckToken(p, lexer.MINUS) {
		emitter.Emit(p.emitter, p.CurToken.TokenText)
		NextToken(p)
	}
	Primary(p)
}

func Primary(p *Parser) {
	if CheckToken(p, lexer.NUMBER) {
		emitter.Emit(p.emitter, p.CurToken.TokenText)
		NextToken(p)
	} else if CheckToken(p, lexer.IDENT) {
		_, ok := p.Symbols[p.CurToken.TokenText]
		if ok == false {
			Abort(p, "Referencing variable before assignment: "+p.CurToken.TokenText)
		}
		emitter.Emit(p.emitter, p.CurToken.TokenText)
		NextToken(p)
	} else {
		Abort(p, "Unexpected token at "+p.CurToken.TokenText)
	}
}

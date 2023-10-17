package lexer

import (
	"log"
	"unicode"
)

type Lexer struct {
	Source  string
	CurChar byte
	CurPos  int
}

type Token struct {
	TokenText string
	TokenType
}

type TokenType int64

const (
	EOF      TokenType = -1
	NEWLINE  TokenType = 0
	NUMBER   TokenType = 1
	IDENT    TokenType = 2
	STRING   TokenType = 3
	LABEL    TokenType = 101
	GOTO     TokenType = 102
	PRINT    TokenType = 103
	INPUT    TokenType = 104
	LET      TokenType = 105
	IF       TokenType = 106
	THEN     TokenType = 107
	ENDIF    TokenType = 108
	WHILE    TokenType = 109
	REPEAT   TokenType = 110
	ENDWHILE TokenType = 111
	EQ       TokenType = 201
	PLUS     TokenType = 202
	MINUS    TokenType = 203
	ASTERISK TokenType = 204
	SLASH    TokenType = 205
	EQEQ     TokenType = 206
	NOTEQ    TokenType = 207
	LT       TokenType = 208
	LTEQ     TokenType = 209
	GT       TokenType = 210
	GTEQ     TokenType = 211
	UNKNOWN  TokenType = 999
)

// NewLexer create a new lexer
func NewLexer(source string) *Lexer {
	lexer := &Lexer{
		Source:  source + "\n",
		CurChar: 0,
		CurPos:  -1,
	}
	NextChar(lexer)
	return lexer
}

func NextChar(l *Lexer) {
	l.CurPos += 1
	if l.CurPos >= len(l.Source) {
		l.CurChar = 0
	} else {
		l.CurChar = l.Source[l.CurPos]
	}

	return
}

func Peek(l *Lexer) byte {
	if (l.CurPos + 1) >= len(l.Source) {
		return 0
	}
	return l.Source[l.CurPos+1]
}

func Abort(l *Lexer, message string) {
	log.Fatalln("Lexing error. " + message)
}

func SkipWhitespace(l *Lexer) {
	for string(l.CurChar) == " " || string(l.CurChar) == "\t" || string(l.CurChar) == "\r" {
		NextChar(l)
	}
}

func SkipComment(l *Lexer) {
	if string(l.CurChar) == "#" {
		for string(l.CurChar) != "\n" {
			NextChar(l)
		}
	}
}

func GetToken(l *Lexer) Token {
	SkipWhitespace(l)
	SkipComment(l)
	var token Token

	current := string(l.CurChar)
	if current == "+" {
		token = Token{TokenText: current, TokenType: PLUS}
	} else if current == "-" {
		token = Token{TokenText: current, TokenType: MINUS}
	} else if current == "*" {
		token = Token{TokenText: current, TokenType: ASTERISK}
	} else if current == "/" {
		token = Token{TokenText: current, TokenType: SLASH}
	} else if current == "=" {
		if string(Peek(l)) == "=" {
			lastChar := current
			NextChar(l)
			token = Token{lastChar + string(l.CurChar), EQEQ}
		} else {
			token = Token{current, EQ}
		}
	} else if current == ">" {
		if string(Peek(l)) == "=" {
			lastChar := current
			NextChar(l)
			token = Token{lastChar + string(l.CurChar), GTEQ}
		} else {
			token = Token{current, GT}
		}
	} else if current == "<" {
		if string(Peek(l)) == "=" {
			lastChar := current
			NextChar(l)
			token = Token{lastChar + string(l.CurChar), LTEQ}
		} else {
			token = Token{current, LT}
		}
	} else if current == "!" {
		if string(Peek(l)) == "=" {
			lastChar := current
			NextChar(l)
			token = Token{lastChar + string(l.CurChar), NOTEQ}
		} else {
			Abort(l, "Expected !=, got !"+string(Peek(l)))
		}
	} else if current == "\"" {
		NextChar(l)
		startPos := l.CurPos

		for string(l.CurChar) != "\"" {
			current := string(l.CurChar)
			if current == "\r" || current == "\n" || current == "\t" || current == "\\" || current == "%" {
				Abort(l, "Illegal character in string.")
			}
			NextChar(l)
		}

		tokenText := l.Source[startPos:l.CurPos]
		token = Token{TokenText: tokenText, TokenType: STRING}
	} else if unicode.IsDigit(rune(l.CurChar)) {
		startPos := l.CurPos
		for unicode.IsDigit(rune(Peek(l))) {
			NextChar(l)
		}

		if string(Peek(l)) == "." {
			NextChar(l)

			if !unicode.IsDigit(rune(Peek(l))) {
				Abort(l, "Illegal character in number")
			}

			for unicode.IsDigit(rune(Peek(l))) {
				NextChar(l)
			}
		}

		tokenNumber := l.Source[startPos : l.CurPos+1]
		token = Token{TokenText: tokenNumber, TokenType: NUMBER}
	} else if unicode.IsLetter(rune(l.CurChar)) {
		startPos := l.CurPos

		for unicode.IsLetter(rune(Peek(l))) || unicode.IsDigit(rune(Peek(l))) {
			NextChar(l)
		}

		tokenText := l.Source[startPos : l.CurPos+1]

		keyword := isKeyword(tokenText)

		if keyword == UNKNOWN {
			token = Token{TokenText: tokenText, TokenType: IDENT}
		} else {
			token = Token{TokenText: tokenText, TokenType: keyword}
		}

	} else if current == "\n" {
		token = Token{TokenText: current, TokenType: NEWLINE}
	} else if current == string(byte(0)) {
		token = Token{TokenText: current, TokenType: EOF}
	} else {
		Abort(l, "Unknown token: "+current)
	}

	NextChar(l)
	return token
}

func isKeyword(word string) TokenType {
	keywords := map[string]TokenType{
		"LABEL":    LABEL,
		"GOTO":     GOTO,
		"PRINT":    PRINT,
		"INPUT":    INPUT,
		"LET":      LET,
		"IF":       IF,
		"THEN":     THEN,
		"ENDIF":    ENDIF,
		"WHILE":    WHILE,
		"REPEAT":   REPEAT,
		"ENDWHILE": ENDWHILE,
	}

	if tokenType, found := keywords[word]; found {
		return tokenType
	}

	return UNKNOWN
}

// Util function

func TokenTypeName(tokenType TokenType) string {
	switch tokenType {
	case EOF:
		return "EOF"
	case NEWLINE:
		return "NEWLINE"
	case NUMBER:
		return "NUMBER"
	case IDENT:
		return "IDENT"
	case STRING:
		return "STRING"
	case LABEL:
		return "LABEL"
	case GOTO:
		return "GOTO"
	case PRINT:
		return "PRINT"
	case INPUT:
		return "INPUT"
	case LET:
		return "LET"
	case IF:
		return "IF"
	case THEN:
		return "THEN"
	case ENDIF:
		return "ENDIF"
	case WHILE:
		return "WHILE"
	case REPEAT:
		return "REPEAT"
	case ENDWHILE:
		return "ENDWHILE"
	case EQ:
		return "EQ"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case ASTERISK:
		return "ASTERISK"
	case SLASH:
		return "SLASH"
	case EQEQ:
		return "EQEQ"
	case NOTEQ:
		return "NOTEQ"
	case LT:
		return "LT"
	case LTEQ:
		return "LTEQ"
	case GT:
		return "GT"
	case GTEQ:
		return "GTEQ"
	case UNKNOWN:
		return "UNKNOWN"
	default:
		return "Invalid TokenType"
	}
}

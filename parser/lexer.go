package parser

import (
	"fmt"
	"regexp"
)

type regexPattern struct {
	regex   *regexp.Regexp
	handler regexHandler
}

type regexHandler func(l *lexer, regex *regexp.Regexp)

func defaultHandler(tokenType TokenType, value string) regexHandler {
	return func(l *lexer, regex *regexp.Regexp) {
		l.advance(len(value))
		token := NewToken(tokenType, value)
		l.push(token)
	}
}

func skipHandler(l *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(l.remainder())
	l.advance(match[1])
}

func operationHandler(l *lexer, regex *regexp.Regexp) {
	match := regex.FindString(l.remainder())
	token := NewToken(OPERATION, match)
	l.push(token)
	l.advance(len(match))
}

func numberHandler(l *lexer, regex *regexp.Regexp) {
	match := regex.FindString(l.remainder())
	token := NewToken(CONSTANT, match)
	l.push(token)
	l.advance(len(match))
}

type lexer struct {
	Tokens   []Token
	input    string
	pos      int
	patterns []regexPattern
}

func Lex(input string) []Token {
	l := newLexer(input)

	for !l.at_eof() {
		matched := false

		for _, pattern := range l.patterns {
			loc := pattern.regex.FindStringIndex(l.remainder())
			if loc != nil && loc[0] == 0 {
				pattern.handler(l, pattern.regex)
				matched = true
				break
			}
		}
		if !matched {
			panic(fmt.Sprintf("Unrecognise token %v", l.remainder()))
		}
	}

	l.push(NewToken(EOF, "EOF"))
	return l.Tokens
}

func (l *lexer) advance(n int) {
	l.pos += n
}

func (l *lexer) remainder() string {
	return l.input[l.pos:]
}

func (l *lexer) push(token Token) {
	l.Tokens = append(l.Tokens, token)
}

func (l *lexer) at_eof() bool {
	return l.pos >= len(l.input)
}

func newLexer(input string) *lexer {
	return &lexer{
		Tokens: make([]Token, 0),
		input:  input,
		pos:    0,
		patterns: []regexPattern{
			{regexp.MustCompile(`\r?\n`), skipHandler},
			{regexp.MustCompile(`\s?[, ]\s`), skipHandler},
			{regexp.MustCompile(`\(`), defaultHandler(OPEN_PAREN, "(")},
			{regexp.MustCompile(`\)`), defaultHandler(CLOSE_PAREN, ")")},
			{regexp.MustCompile(`[a-zA-Z]+[0-9]*`), operationHandler},
			{regexp.MustCompile(`[-]?[0-9]+(\.[0-9]+)?`), numberHandler},
		},
	}
}

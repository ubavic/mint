package parser

import (
	"bufio"
	"io"
	"slices"
	"unicode"
)

type TokenType uint

const (
	Identifier TokenType = iota
	LeftBrace
	RightBrace
	Text
	EOF
)

type Token struct {
	Type    TokenType
	Content string
	Line    int
	Column  int
}

type Tokenizer struct {
	input      *bufio.Reader
	line       int
	column     int
	lastLine   int
	lastColumn int
}

func NewTokenizer(input *bufio.Reader) Tokenizer {
	return Tokenizer{
		input:  input,
		line:   1,
		column: 1,
	}
}

func (tokenizer *Tokenizer) Tokenize() []Token {
	tokens := []Token{}
	var newTokens []Token

	for {
		startLine := tokenizer.line
		startCol := tokenizer.column

		r, _, err := tokenizer.readRune()
		if err != nil {
			if err == io.EOF {
				tokens = append(tokens, Token{Type: EOF, Line: startLine, Column: startCol})
				return tokens
			}

			panic(err)
		}

		switch r {
		case '{':
			newTokens = []Token{{Type: LeftBrace, Content: "{", Line: startLine, Column: startCol}}
		case '}':
			newTokens = []Token{{Type: RightBrace, Content: "}", Line: startLine, Column: startCol}}
		case '@':
			newTokens = tokenizer.tokenizeIdentifier("", startLine, startCol)
		default:
			tokenizer.unreadRune()
			newTokens = tokenizer.tokenizeText("", startLine, startCol)
		}

		tokens = append(tokens, newTokens...)
	}

}

func (tokenizer *Tokenizer) tokenizeText(start string, startLine, startCol int) []Token {
	text := start

	for {
		r, _, err := tokenizer.readRune()

		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		if slices.Contains([]rune("{}"), r) {
			tokenizer.unreadRune()
			break
		} else if r == '@' {
			atLine := tokenizer.lastLine
			atCol := tokenizer.lastColumn

			nextRune, _, err := tokenizer.readRune()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					panic(err)
				}
			}

			if slices.Contains([]rune("{}@"), nextRune) {
				r = nextRune
			} else {
				identifier := tokenizer.tokenizeIdentifier(string(nextRune), atLine, atCol)
				return append([]Token{{Type: Text, Content: text, Line: startLine, Column: startCol}}, identifier...)
			}
		}

		text += string(r)
	}

	return []Token{
		{Type: Text, Content: text, Line: startLine, Column: startCol},
	}
}

// Tokenize identifier or a escaped sequence: `@@`, `@{`, `@}`
func (tokenizer *Tokenizer) tokenizeIdentifier(start string, startLine, startCol int) []Token {
	identifier := start
	firstPass := start == ""

	for {
		r, _, err := tokenizer.readRune()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		if slices.Contains([]rune("{} @"), r) {
			if firstPass {
				return tokenizer.tokenizeText(string(r), tokenizer.lastLine, tokenizer.lastColumn)
			}

			err := tokenizer.unreadRune()
			if err != nil {
				panic(err)
			}

			break
		}

		identifier += string(r)
		firstPass = false
	}

	return []Token{
		{Type: Identifier, Content: identifier, Line: startLine, Column: startCol},
	}
}

func EqualStreams(a, b []Token) bool {
	if a == nil {
		return b == nil
	}

	if b == nil {
		return a == nil
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}

		if a[i].Line != b[i].Line || a[i].Column != b[i].Column {
			return false
		}
	}

	return true
}

func (tokenizer *Tokenizer) readRune() (rune, int, error) {
	r, size, err := tokenizer.input.ReadRune()
	if err != nil {
		return r, size, err
	}

	tokenizer.lastLine = tokenizer.line
	tokenizer.lastColumn = tokenizer.column

	if r == '\n' {
		tokenizer.line++
		tokenizer.column = 1
	} else {
		tokenizer.column++
	}

	return r, size, nil
}

func (tokenizer *Tokenizer) unreadRune() error {
	err := tokenizer.input.UnreadRune()
	if err != nil {
		return err
	}

	tokenizer.line = tokenizer.lastLine
	tokenizer.column = tokenizer.lastColumn

	return nil
}

func (t Token) ContainsWhitespaceOnly() bool {
	for _, r := range t.Content {
		if !unicode.IsSpace(r) {
			return false
		}
	}

	return true
}

func (t Token) String() string {
	switch t.Type {
	case Identifier:
		return "\x1b[91m" + t.Content + "\x1b[0m"
	case Text:
		return "\x1b[93m\"" + t.Content + "\"\x1b[0m"
	case LeftBrace, RightBrace:
		return "\x1b[95m" + t.Content + "\x1b[0m"
	case EOF:
		return "\x1b[96mEOF\x1b[0m"
	default:
		return t.Content
	}
}

func (t TokenType) String() string {
	switch t {
	case Identifier:
		return "\x1b[91mIdentifier\x1b[0m"
	case Text:
		return "\x1b[93mText\x1b[0m"
	case LeftBrace:
		return "\x1b[95mLeftBrace\x1b[0m"
	case RightBrace:
		return "\x1b[95mRightBrace\x1b[0m"
	case EOF:
		return "\x1b[96mEOF\x1b[0m"
	default:
		return "Unknown"
	}
}

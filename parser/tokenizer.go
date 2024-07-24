package parser

import (
	"bufio"
	"io"
	"slices"
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
}

type Tokenizer struct {
	input *bufio.Reader
}

func NewTokenizer(input *bufio.Reader) Tokenizer {
	return Tokenizer{
		input: input,
	}
}

func (tokenizer *Tokenizer) Tokenize() []Token {
	tokens := []Token{}

	for {
		r, _, err := tokenizer.input.ReadRune()
		if err != nil {
			if err == io.EOF {
				tokens = append(tokens, Token{Type: EOF})
				return tokens
			}

			panic(err)
		}

		var token Token

		switch r {
		case '{':
			token = Token{Type: LeftBrace, Content: "{"}
		case '}':
			token = Token{Type: RightBrace, Content: "}"}
		case '@':
			identifier := tokenizer.tokenizeIdentifier()
			token = Token{Type: Identifier, Content: identifier}
		default:
			tokenizer.input.UnreadRune()
			text := tokenizer.tokenizeText()
			token = Token{Type: Text, Content: text}
		}

		tokens = append(tokens, token)
	}

}

func (tokenizer *Tokenizer) tokenizeText() string {
	text := ""

	for {
		r, _, err := tokenizer.input.ReadRune()

		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		if slices.Contains([]rune("{}@"), r) {
			tokenizer.input.UnreadRune()
			break
		}

		text += string(r)
	}

	return text
}

func (tokenizer *Tokenizer) tokenizeIdentifier() string {
	identifier := ""

	for {
		r, _, err := tokenizer.input.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		if slices.Contains([]rune("{} @"), r) {
			err := tokenizer.input.UnreadRune()
			if err != nil {
				panic(err)
			}

			break
		}

		identifier += string(r)
	}

	return identifier
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

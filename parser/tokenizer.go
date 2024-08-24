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
	var newTokens []Token

	for {
		r, _, err := tokenizer.input.ReadRune()
		if err != nil {
			if err == io.EOF {
				tokens = append(tokens, Token{Type: EOF})
				return tokens
			}

			panic(err)
		}

		switch r {
		case '{':
			newTokens = []Token{{Type: LeftBrace, Content: "{"}}
		case '}':
			newTokens = []Token{{Type: RightBrace, Content: "}"}}
		case '@':
			newTokens = tokenizer.tokenizeIdentifier("")
		default:
			tokenizer.input.UnreadRune()
			newTokens = tokenizer.tokenizeText("")
		}

		tokens = append(tokens, newTokens...)
	}

}

func (tokenizer *Tokenizer) tokenizeText(start string) []Token {
	text := start

	for {
		r, _, err := tokenizer.input.ReadRune()

		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		if slices.Contains([]rune("{}"), r) {
			tokenizer.input.UnreadRune()
			break
		} else if r == '@' {

			nextRune, _, err := tokenizer.input.ReadRune()
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
				identifier := tokenizer.tokenizeIdentifier(string(nextRune))

				return append([]Token{{Type: Text, Content: text}}, identifier...)
			}
		}

		text += string(r)
	}

	return []Token{
		{Type: Text, Content: text},
	}
}

// Tokenize identifier or a escaped sequence: `@@`, `@{`, `@}`
func (tokenizer *Tokenizer) tokenizeIdentifier(start string) []Token {
	identifier := start
	firstPass := true

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
			if firstPass {
				return tokenizer.tokenizeText(string(r))
			}

			err := tokenizer.input.UnreadRune()
			if err != nil {
				panic(err)
			}

			break
		}

		identifier += string(r)
		firstPass = false
	}

	return []Token{
		Token{Type: Identifier, Content: identifier},
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
	}

	return true
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

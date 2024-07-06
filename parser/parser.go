package parser

type TokenType uint

const (
	Identifier TokenType = iota
	LeftBrace
	RightBrace
	Text
)

type Token struct {
	Type    TokenType
	Content string
}

func Tokenize(in string) []Token {
	tokens := []Token{}

	for _, r := range in {
		var token Token

		switch r {
		case '{':
			token = Token{Type: LeftBrace, Content: "{"}
		case '}':
			token = Token{Type: RightBrace, Content: "}"}

		}

		tokens = append(tokens, token)
	}

	return tokens
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
		return "\x1b[93m" + t.Content + "\x1b[0m"
	case LeftBrace, RightBrace:
		return "\x1b[95m" + t.Content + "\x1b[0m"
	default:
		return t.Content
	}
}

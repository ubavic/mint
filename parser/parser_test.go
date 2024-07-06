package parser_test

import (
	"testing"

	"github.com/ubavic/mint/parser"
)

func TestTokenizer(t *testing.T) {
	result := parser.Tokenize("")
	if result == nil {
		t.Error("Result must not be nil")
	}

	result = parser.Tokenize("{}")
	expects := []parser.Token{
		{Type: parser.LeftBrace, Content: "{"},
		{Type: parser.RightBrace, Content: "}"},
	}

	if !parser.EqualStreams(result, expects) {
		t.Error("Streams are not equal")
	}
}

func Test_EqualStreams(t *testing.T) {
	if !parser.EqualStreams(nil, nil) {
		t.Error("Streams should be equal")
	}

	if !parser.EqualStreams([]parser.Token{}, []parser.Token{}) {
		t.Error("Streams should be equal")
	}

	if parser.EqualStreams([]parser.Token{}, nil) {
		t.Error("Streams should be different")
	}

	stream1 := []parser.Token{
		{Type: parser.Identifier, Content: ""},
	}

	stream2 := []parser.Token{
		{Type: parser.Identifier, Content: ""},
	}

	if !parser.EqualStreams(stream1, stream2) {
		t.Error("Streams should be equal")
	}
}

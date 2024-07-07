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

	result = parser.Tokenize("{hello world}")
	expects = []parser.Token{
		{Type: parser.LeftBrace, Content: "{"},
		{Type: parser.Text, Content: "hello world"},
		{Type: parser.RightBrace, Content: "}"},
	}

	if !parser.EqualStreams(result, expects) {
		t.Errorf("Streams are not equal. Expected %v got %v", expects, result)
	}

	result = parser.Tokenize("{ }")
	expects = []parser.Token{
		{Type: parser.LeftBrace, Content: "{"},
		{Type: parser.Text, Content: " "},
		{Type: parser.RightBrace, Content: "}"},
	}

	if !parser.EqualStreams(result, expects) {
		t.Errorf("Streams are not equal. Expected %v got %v", expects, result)
	}

	result = parser.Tokenize(" { }")
	expects = []parser.Token{
		{Type: parser.Text, Content: " "},
		{Type: parser.LeftBrace, Content: "{"},
		{Type: parser.Text, Content: " "},
		{Type: parser.RightBrace, Content: "}"},
	}

	if !parser.EqualStreams(result, expects) {
		t.Errorf("Streams are not equal. Expected %v got %v", expects, result)
	}

	result = parser.Tokenize(" ")
	expects = []parser.Token{
		{Type: parser.Text, Content: " "},
	}

	if !parser.EqualStreams(result, expects) {
		t.Errorf("Streams are not equal. Expected %v got %v", expects, result)
	}

	result = parser.Tokenize("helloWorld}")
	expects = []parser.Token{
		{Type: parser.Text, Content: "helloWorld"},
		{Type: parser.RightBrace, Content: "}"},
	}

	if !parser.EqualStreams(result, expects) {
		t.Errorf("Streams are not equal. Expected %v got %v", expects, result)
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

package parser_test

import (
	"bufio"
	"fmt"
	"strings"
	"testing"

	"github.com/ubavic/mint/parser"
)

func TestTokenizer(t *testing.T) {
	testCases := []struct {
		input          string
		expectedResult []parser.Token
	}{
		{
			input: "",
			expectedResult: []parser.Token{
				{Type: parser.EOF, Content: "", Line: 1, Column: 1},
			},
		},
		{
			input: "{}",
			expectedResult: []parser.Token{
				{Type: parser.LeftBrace, Content: "{", Line: 1, Column: 1},
				{Type: parser.RightBrace, Content: "}", Line: 1, Column: 2},
				{Type: parser.EOF, Content: "", Line: 1, Column: 3},
			},
		},
		{
			input: "{hello world}",
			expectedResult: []parser.Token{
				{Type: parser.LeftBrace, Content: "{", Line: 1, Column: 1},
				{Type: parser.Text, Content: "hello world", Line: 1, Column: 2},
				{Type: parser.RightBrace, Content: "}", Line: 1, Column: 13},
				{Type: parser.EOF, Content: "", Line: 1, Column: 14},
			},
		},
		{
			input: "{ }",
			expectedResult: []parser.Token{
				{Type: parser.LeftBrace, Content: "{", Line: 1, Column: 1},
				{Type: parser.Text, Content: " ", Line: 1, Column: 2},
				{Type: parser.RightBrace, Content: "}", Line: 1, Column: 3},
				{Type: parser.EOF, Content: "", Line: 1, Column: 4},
			},
		},
		{
			input: " { }",
			expectedResult: []parser.Token{
				{Type: parser.Text, Content: " ", Line: 1, Column: 1},
				{Type: parser.LeftBrace, Content: "{", Line: 1, Column: 2},
				{Type: parser.Text, Content: " ", Line: 1, Column: 3},
				{Type: parser.RightBrace, Content: "}", Line: 1, Column: 4},
				{Type: parser.EOF, Content: "", Line: 1, Column: 5},
			},
		},
		{
			input: " ",
			expectedResult: []parser.Token{
				{Type: parser.Text, Content: " ", Line: 1, Column: 1},
				{Type: parser.EOF, Content: "", Line: 1, Column: 2},
			},
		},
		{
			input: "helloWorld}",
			expectedResult: []parser.Token{
				{Type: parser.Text, Content: "helloWorld", Line: 1, Column: 1},
				{Type: parser.RightBrace, Content: "}", Line: 1, Column: 11},
				{Type: parser.EOF, Content: "", Line: 1, Column: 12},
			},
		},
		{
			input: "@p",
			expectedResult: []parser.Token{
				{Type: parser.Identifier, Content: "p", Line: 1, Column: 1},
				{Type: parser.EOF, Content: "", Line: 1, Column: 3},
			},
		},
		{
			input: "@p @a",
			expectedResult: []parser.Token{
				{Type: parser.Identifier, Content: "p", Line: 1, Column: 1},
				{Type: parser.Text, Content: " ", Line: 1, Column: 3},
				{Type: parser.Identifier, Content: "a", Line: 1, Column: 4},
				{Type: parser.EOF, Content: "", Line: 1, Column: 6},
			},
		},
		{
			input: "@{@@@}",
			expectedResult: []parser.Token{
				{Type: parser.Text, Content: "{@}", Line: 1, Column: 2},
				{Type: parser.EOF, Content: "", Line: 1, Column: 7},
			},
		},
		{
			input: "@p{hello there}\n\n@p{how are you?}",
			expectedResult: []parser.Token{
				{Type: parser.Identifier, Content: "p", Line: 1, Column: 1},
				{Type: parser.LeftBrace, Content: "{", Line: 1, Column: 3},
				{Type: parser.Text, Content: "hello there", Line: 1, Column: 4},
				{Type: parser.RightBrace, Content: "}", Line: 1, Column: 15},
				{Type: parser.Text, Content: "\n\n", Line: 1, Column: 16},
				{Type: parser.Identifier, Content: "p", Line: 3, Column: 1},
				{Type: parser.LeftBrace, Content: "{", Line: 3, Column: 3},
				{Type: parser.Text, Content: "how are you?", Line: 3, Column: 4},
				{Type: parser.RightBrace, Content: "}", Line: 3, Column: 16},
				{Type: parser.EOF, Content: "", Line: 3, Column: 17},
			},
		},
		{
			input: "@p#paragraph",
			expectedResult: []parser.Token{
				{Type: parser.Identifier, Content: "p", Line: 1, Column: 1},
				{Type: parser.CommandId, Content: "paragraph", Line: 1, Column: 3},
				{Type: parser.EOF, Content: "", Line: 1, Column: 13},
			},
		},
		{
			input: "@p#paragraph-1",
			expectedResult: []parser.Token{
				{Type: parser.Identifier, Content: "p", Line: 1, Column: 1},
				{Type: parser.CommandId, Content: "paragraph-1", Line: 1, Column: 3},
				{Type: parser.EOF, Content: "", Line: 1, Column: 15},
			},
		},
		{
			input: "@p{#just-text}",
			expectedResult: []parser.Token{
				{Type: parser.Identifier, Content: "p", Line: 1, Column: 1},
				{Type: parser.LeftBrace, Content: "{", Line: 1, Column: 3},
				{Type: parser.Text, Content: "#just-text", Line: 1, Column: 4},
				{Type: parser.RightBrace, Content: "}", Line: 1, Column: 14},
				{Type: parser.EOF, Content: "", Line: 1, Column: 15},
			},
		},
	}

	for i, testCase := range testCases {
		t.Run(
			fmt.Sprintf("TestTokenizer%d", i),
			func(t *testing.T) {
				reader := bufio.NewReader(strings.NewReader(testCase.input))
				tokenizer := parser.NewTokenizer(reader)

				result := tokenizer.Tokenize()
				if !parser.EqualStreams(result, testCase.expectedResult) {
					t.Errorf("Streams are not equal. Expected \n%v \ngot \n%v", testCase.expectedResult, result)
				}
			},
		)
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

func TestTokenizerVerbatimArgument(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("{>>hello world {}<<}"))
	tokenizer := parser.NewTokenizer(reader)

	result := tokenizer.Tokenize()
	if len(result) != 4 {
		t.Fatalf("expected 4 tokens, got %d", len(result))
	}

	if result[0].Type != parser.LeftBrace {
		t.Fatalf("expected first token to be LeftBrace, got %v", result[0].Type)
	}

	if result[1].Type != parser.Text || result[1].Content != "hello world {}" {
		t.Fatalf("expected verbatim text %q, got token %#v", "hello world {}", result[1])
	}

	if result[2].Type != parser.RightBrace {
		t.Fatalf("expected third token to be RightBrace, got %v", result[2].Type)
	}

	if result[3].Type != parser.EOF {
		t.Fatalf("expected EOF token, got %v", result[3].Type)
	}
}

func TestTokenizerVerbatimMarkersMustBeAdjacent(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("{ >> hello << }"))
	tokenizer := parser.NewTokenizer(reader)

	result := tokenizer.Tokenize()
	if len(result) != 4 {
		t.Fatalf("expected 4 tokens, got %d", len(result))
	}

	if result[1].Type != parser.Text || result[1].Content != " >> hello << " {
		t.Fatalf("expected normal text token %q, got token %#v", " >> hello << ", result[1])
	}
}

package parser_test

import (
	"fmt"
	"testing"

	"github.com/ubavic/mint/parser"
)

func TestTokenizer(t *testing.T) {

	testCases := []struct {
		input          string
		expectedResult []parser.Token
	}{
		{
			input:          "",
			expectedResult: []parser.Token{},
		},
		{
			input: "{}",
			expectedResult: []parser.Token{
				{Type: parser.LeftBrace, Content: "{"},
				{Type: parser.RightBrace, Content: "}"},
			},
		},
		{
			input: "{hello world}",
			expectedResult: []parser.Token{
				{Type: parser.LeftBrace, Content: "{"},
				{Type: parser.Text, Content: "hello world"},
				{Type: parser.RightBrace, Content: "}"},
			},
		},
		{
			input: "{ }",
			expectedResult: []parser.Token{
				{Type: parser.LeftBrace, Content: "{"},
				{Type: parser.Text, Content: " "},
				{Type: parser.RightBrace, Content: "}"},
			},
		},
		{
			input: " { }",
			expectedResult: []parser.Token{
				{Type: parser.Text, Content: " "},
				{Type: parser.LeftBrace, Content: "{"},
				{Type: parser.Text, Content: " "},
				{Type: parser.RightBrace, Content: "}"},
			},
		},
		{
			input: " ",
			expectedResult: []parser.Token{
				{Type: parser.Text, Content: " "},
			},
		},
		{
			input: "helloWorld}",
			expectedResult: []parser.Token{
				{Type: parser.Text, Content: "helloWorld"},
				{Type: parser.RightBrace, Content: "}"},
			},
		},
	}

	for i, testCase := range testCases {
		t.Run(
			fmt.Sprintf("TestTokenizer%d", i),
			func(t *testing.T) {
				result := parser.Tokenize(testCase.input)
				if !parser.EqualStreams(result, testCase.expectedResult) {
					t.Errorf("Streams are not equal. Expected %v got %v", testCase.expectedResult, result)
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

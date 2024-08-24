package schema_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ubavic/mint/parser"
	"github.com/ubavic/mint/schema"
)

func TestSchemaValidator(t *testing.T) {
	testCases := []struct {
		Commands            []schema.Command
		Tokens              []parser.Token
		Groups              []schema.Group
		AllowedRootChildren string
		ExpectedError       error
	}{
		{
			Commands: []schema.Command{
				{Command: "c1", Arguments: 0},
				{Command: "c2", Arguments: 0},
			},
			Tokens: []parser.Token{
				{Type: parser.Identifier, Content: "c1"},
				{Type: parser.Identifier, Content: "c2"},
				{Type: parser.EOF, Content: ""},
			},
			ExpectedError: nil,
		},
		{
			Commands: []schema.Command{
				{Command: "c1", Arguments: 0},
			},
			Tokens: []parser.Token{
				{Type: parser.Identifier, Content: "c2"},
				{Type: parser.EOF, Content: ""},
			},
			ExpectedError: schema.ErrCommandNotFound,
		},
		{
			Commands: []schema.Command{
				{Command: "c1", Arguments: 2},
			},
			Tokens: []parser.Token{
				{Type: parser.Identifier, Content: "c1"},
				{Type: parser.LeftBrace, Content: "{"},
				{Type: parser.RightBrace, Content: "}"},
				{Type: parser.LeftBrace, Content: "{"},
				{Type: parser.RightBrace, Content: "}"},
				{Type: parser.EOF, Content: ""},
			},
			ExpectedError: nil,
		},
		{
			Commands: []schema.Command{
				{Command: "c1", Arguments: 3},
			},
			Tokens: []parser.Token{
				{Type: parser.Identifier, Content: "c1"},
				{Type: parser.LeftBrace, Content: "{"},
				{Type: parser.RightBrace, Content: "}"},
				{Type: parser.LeftBrace, Content: "{"},
				{Type: parser.RightBrace, Content: "}"},
				{Type: parser.EOF, Content: ""},
			},
			ExpectedError: schema.ErrCommandInvalidArguments,
		},
		{
			Commands: []schema.Command{
				{Command: "c1", Arguments: 0},
			},
			Tokens: []parser.Token{
				{Type: parser.Identifier, Content: "c1"},
				{Type: parser.EOF, Content: ""},
			},
			AllowedRootChildren: "G1",
			Groups: []schema.Group{
				{
					Name:     "G1",
					Commands: []string{"c1"},
				},
			},
			ExpectedError: nil,
		},
		{
			Commands: []schema.Command{
				{Command: "c1", Arguments: 0},
			},
			Tokens: []parser.Token{
				{Type: parser.Identifier, Content: "c1"},
				{Type: parser.EOF, Content: ""},
			},
			AllowedRootChildren: "G1",
			Groups: []schema.Group{
				{
					Name:     "G1",
					Commands: []string{},
				},
			},
			ExpectedError: schema.ErrCommandNotFound,
		},
	}

	for i, testCase := range testCases {
		t.Run(
			fmt.Sprintf("Test_schema_validator_%d", i),
			func(t *testing.T) {
				sc := schema.Schema{
					Source: schema.Source{
						AllowedRootChildren: testCase.AllowedRootChildren,
						Commands:            testCase.Commands,
						Groups:              testCase.Groups,
					},
				}

				np := parser.NewParser(testCase.Tokens, &sc)
				_, err := np.Parse()

				if err == nil && testCase.ExpectedError != nil {
					t.Fatalf("Expected an error \"%v\", but got no error", testCase.ExpectedError)
				} else if err != nil && testCase.ExpectedError == nil {
					t.Fatalf("Expected no error but got an error: \"%v\"", err)
				} else if err != nil && testCase.ExpectedError != nil && !errors.Is(err, testCase.ExpectedError) {
					t.Fatalf("Expected to find the error \"%v\" in the error \"%v\"", testCase.ExpectedError, err)
				}
			},
		)
	}

}

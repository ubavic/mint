package parser_test

import (
	"encoding/json"
	"fmt"
	"slices"
	"testing"

	"github.com/ubavic/mint/parser"
)

func Test_Parser(t *testing.T) {

	testCases := []struct {
		input          []parser.Token
		expectedResult parser.Element
	}{
		{
			input: []parser.Token{},
			expectedResult: &parser.Block{
				Nodes: []parser.Element{},
			},
		},
		{
			input: []parser.Token{
				{Type: parser.Identifier, Content: "p"},
				{Type: parser.EOF, Content: ""},
			},
			expectedResult: &parser.Block{
				Nodes: []parser.Element{
					&parser.Command{
						Name:      "p",
						Arguments: []parser.Element{},
					},
				},
			},
		},
		{
			input: []parser.Token{
				{Type: parser.Identifier, Content: "a1"},
				{Type: parser.Identifier, Content: "a2"},
				{Type: parser.EOF, Content: ""},
			},
			expectedResult: &parser.Block{
				Nodes: []parser.Element{
					&parser.Command{
						Name:      "a1",
						Arguments: []parser.Element{},
					},
					&parser.Command{
						Name:      "a2",
						Arguments: []parser.Element{},
					},
				},
			},
		},
		{
			input: []parser.Token{
				{Type: parser.Identifier, Content: "hello"},
				{Type: parser.Text, Content: "world"},
				{Type: parser.Identifier, Content: "bye"},
				{Type: parser.EOF, Content: ""},
			},
			expectedResult: &parser.Block{
				Nodes: []parser.Element{
					&parser.Command{
						Name:      "hello",
						Arguments: []parser.Element{},
					},
					&parser.TextContent{
						TextContent: "world",
					},
					&parser.Command{
						Name:      "bye",
						Arguments: []parser.Element{},
					},
				},
			},
		},
		{
			input: []parser.Token{
				{Type: parser.Identifier, Content: "cmd"},
				{Type: parser.LeftBrace, Content: "{"},
				{Type: parser.Text, Content: "arg"},
				{Type: parser.RightBrace, Content: "}"},
				{Type: parser.EOF, Content: ""},
			},
			expectedResult: &parser.Block{
				Nodes: []parser.Element{
					&parser.Command{
						Name: "cmd",
						Arguments: []parser.Element{
							&parser.Block{
								Nodes: []parser.Element{
									&parser.TextContent{
										TextContent: "arg",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: []parser.Token{
				{Type: parser.Identifier, Content: "cmd"},
				{Type: parser.LeftBrace, Content: "{"},
				{Type: parser.RightBrace, Content: "}"},
				{Type: parser.EOF, Content: ""},
			},
			expectedResult: &parser.Block{
				Nodes: []parser.Element{
					&parser.Command{
						Name: "cmd",
						Arguments: []parser.Element{
							&parser.Block{
								Nodes: []parser.Element{},
							},
						},
					},
				},
			},
		},
		{
			input: []parser.Token{
				{Type: parser.Identifier, Content: "cmd"},
				{Type: parser.LeftBrace, Content: "{"},
				{Type: parser.Text, Content: "arg1"},
				{Type: parser.RightBrace, Content: "}"},
				{Type: parser.LeftBrace, Content: "{"},
				{Type: parser.Text, Content: "arg2"},
				{Type: parser.RightBrace, Content: "}"},
				{Type: parser.EOF, Content: ""},
			},
			expectedResult: &parser.Block{
				Nodes: []parser.Element{
					&parser.Command{
						Name: "cmd",
						Arguments: []parser.Element{
							&parser.Block{
								Nodes: []parser.Element{
									&parser.TextContent{
										TextContent: "arg1",
									},
								},
							},
							&parser.Block{
								Nodes: []parser.Element{
									&parser.TextContent{
										TextContent: "arg2",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: []parser.Token{
				{Type: parser.Identifier, Content: "cmd1"},
				{Type: parser.LeftBrace, Content: "{"},
				{Type: parser.Identifier, Content: "cmd2"},
				{Type: parser.LeftBrace, Content: "{"},
				{Type: parser.RightBrace, Content: "}"},
				{Type: parser.RightBrace, Content: "}"},
				{Type: parser.LeftBrace, Content: "{"},
				{Type: parser.Identifier, Content: "cmd3"},
				{Type: parser.RightBrace, Content: "}"},
				{Type: parser.EOF, Content: ""},
			},
			expectedResult: &parser.Block{
				Nodes: []parser.Element{
					&parser.Command{
						Name: "cmd1",
						Arguments: []parser.Element{
							&parser.Block{
								Nodes: []parser.Element{
									&parser.Command{
										Name: "cmd2",
										Arguments: []parser.Element{
											&parser.Block{
												Nodes: []parser.Element{},
											},
										},
									},
								},
							},
							&parser.Block{
								Nodes: []parser.Element{
									&parser.Command{
										Name:      "cmd3",
										Arguments: []parser.Element{},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: []parser.Token{
				{Type: parser.Identifier, Content: "cmd"},
				{Type: parser.Text, Content: "       "},
				{Type: parser.LeftBrace, Content: "{"},
				{Type: parser.RightBrace, Content: "}"},
				{Type: parser.EOF, Content: ""},
			},
			expectedResult: &parser.Block{
				Nodes: []parser.Element{
					&parser.Command{
						Name: "cmd",
						Arguments: []parser.Element{
							&parser.Block{
								Nodes: []parser.Element{},
							},
						},
					},
				},
			},
		},
		{
			input: []parser.Token{
				{Type: parser.Identifier, Content: "cmd"},
				{Type: parser.Text, Content: "       "},
				{Type: parser.LeftBrace, Content: "{"},
				{Type: parser.RightBrace, Content: "}"},
				{Type: parser.Text, Content: "       \n\n"},
				{Type: parser.LeftBrace, Content: "{"},
				{Type: parser.RightBrace, Content: "}"},
				{Type: parser.EOF, Content: ""},
			},
			expectedResult: &parser.Block{
				Nodes: []parser.Element{
					&parser.Command{
						Name: "cmd",
						Arguments: []parser.Element{
							&parser.Block{
								Nodes: []parser.Element{},
							},
							&parser.Block{
								Nodes: []parser.Element{},
							},
						},
					},
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Run(
			fmt.Sprintf("Test_Parser_%d\n", i),
			func(t *testing.T) {
				t.Log(testCase.input)

				np := parser.NewParser(testCase.input)

				result, err := np.Parse()
				if err != nil {
					t.Fatalf("Expected no error, got \"%s\"", err)
				}

				resultJson, _ := json.Marshal(result)
				expectedResultJson, _ := json.Marshal(testCase.expectedResult)

				if !slices.Equal(resultJson, expectedResultJson) {
					t.Errorf("Results are not equal. Expected \n%v\ngot:\n%v\n", string(expectedResultJson), string(resultJson))
				}
			})
	}
}

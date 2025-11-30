package schema_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ubavic/mint/schema"
)

func TestValidSchemaCheck(t *testing.T) {
	schemas := []schema.Schema{
		{
			Mint: "v0.1",
		},
		{
			Mint: "v0.2",
			Source: schema.Source{
				Commands: []schema.Command{},
			},
		},
		{
			Mint: "v0.3.1",
			Source: schema.Source{
				Commands: []schema.Command{
					{
						Command: "paragraph",
					},
					{
						Command: "title",
					},
				},
				Groups: []schema.Group{
					{
						Name:     "RootElements",
						Commands: []string{"paragraph", "title"},
					},
				},
				AllowedRootCommands: "RootElements",
			},
		},
	}

	for i, s := range schemas {
		t.Run(
			fmt.Sprintf("Test_Valid_Schema_Check_%d", i),
			func(t *testing.T) {
				err := s.Check()
				if err != nil {
					t.Fatalf("expected no error but got \"%v\"", err)
				}
			},
		)
	}
}

func TestInvalidSchemaCheck(t *testing.T) {
	testCases := []struct {
		schema        schema.Schema
		expectedError error
	}{
		{
			schema: schema.Schema{
				Mint: "Lorem ipsum",
			},
			expectedError: schema.ErrInvalidVersionFormat,
		},
		{
			schema: schema.Schema{
				Mint:    "v0.1",
				Version: "Lorem ipsum",
			},
			expectedError: schema.ErrInvalidVersionFormat,
		},
		{
			schema: schema.Schema{
				Mint: "v0.1",
				Source: schema.Source{
					Commands: []schema.Command{
						{
							Command: "paragraph",
						},
						{
							Command: "paragraph",
						},
					},
				},
			},
			expectedError: schema.ErrMultipleComandDefinitions,
		},
		{
			schema: schema.Schema{
				Mint: "v0.1",
				Source: schema.Source{
					Commands: []schema.Command{
						{},
					},
				},
			},
			expectedError: schema.ErrEmptyCommandName,
		},
		{
			schema: schema.Schema{
				Mint: "v0.1",
				Source: schema.Source{
					Commands: []schema.Command{
						{
							Command: "lorem ipsum",
						},
					},
				},
			},
			expectedError: schema.ErrInvalidCommandName,
		},
		{
			schema: schema.Schema{
				Mint: "v0.1",
				Source: schema.Source{
					Groups: []schema.Group{
						{},
					},
				},
			},
			expectedError: schema.ErrEmptyGroupName,
		},
		{
			schema: schema.Schema{
				Mint: "v0.1",
				Source: schema.Source{
					Groups: []schema.Group{
						{
							Name: "RootElements",
						},
						{
							Name: "RootElements",
						},
					},
				},
			},
			expectedError: schema.ErrMultipleGroupDefinitions,
		},
		{
			schema: schema.Schema{
				Mint: "v0.1",
				Source: schema.Source{
					Groups: []schema.Group{
						{
							Name: "Root Elements",
						},
					},
				},
			},
			expectedError: schema.ErrInvalidGroupName,
		},
		{
			schema: schema.Schema{
				Mint: "v0.1",
				Source: schema.Source{
					Groups: []schema.Group{
						{
							Name:     "RootElements",
							Commands: []string{"paragraph"},
						},
					},
				},
			},
			expectedError: schema.ErrUnknownCommand,
		},
		{
			schema: schema.Schema{
				Mint: "v0.1",
				Source: schema.Source{
					AllowedRootCommands: "RootElements",
				},
			},
			expectedError: schema.ErrUnknownGroup,
		},
	}

	for i, testCase := range testCases {
		t.Run(
			fmt.Sprintf("Test_Valid_Schema_Check_%d", i),
			func(t *testing.T) {
				err := testCase.schema.Check()
				if !errors.Is(err, testCase.expectedError) {
					t.Fatalf("expected error \"%v\" but got error \"%v\"", testCase.expectedError, err)
				}
			},
		)
	}
}

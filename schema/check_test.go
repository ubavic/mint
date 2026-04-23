package schema_test

import (
	"errors"
	"fmt"
	"strings"
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

func TestSchemaCheckWithTargetsAndArguments(t *testing.T) {
	validSchema := schema.Schema{
		Mint: "v0.3.1",
		Source: schema.Source{
			Commands: []schema.Command{
				{
					Command: "paragraph",
					Arguments: []schema.Argument{
						{
							Name:         "text",
							Type:         schema.ArgumentTypeString,
							Pattern:      "^[a-z]+$",
							DefaultValue: "hello",
						},
						{
							Name:         "count",
							Type:         schema.ArgumentTypeNumber,
							DefaultValue: "3",
						},
					},
				},
				{
					Command: "enabled",
					Arguments: []schema.Argument{
						{
							Name:         "value",
							Type:         schema.ArgumentTypeBoolean,
							DefaultValue: "true",
						},
					},
				},
				{
					Command: "ref",
					Arguments: []schema.Argument{
						{
							Name: "node",
							Type: schema.ArgumentTypeReference,
						},
					},
				},
				{
					Command: "anything",
					Arguments: []schema.Argument{
						{
							Name: "value",
							Type: schema.ArgumentTypeAny,
						},
					},
				},
			},
			Groups: []schema.Group{
				{
					Name:     "RootElements",
					Commands: []string{"paragraph", "enabled", "ref", "anything"},
				},
			},
			AllowedRootCommands: "RootElements",
		},
		Targets: []schema.Target{
			{
				Name: "html",
				Commands: []struct {
					Command    string `yaml:"command"`
					Expression string `yaml:"expression"`
				}{
					{
						Command:    "paragraph",
						Expression: "<p>$1 ($2)</p>",
					},
					{
						Command:    "enabled",
						Expression: `<flag enabled="$1" />`,
					},
					{
						Command:    "ref",
						Expression: "<ref>$1</ref>",
					},
					{
						Command:    "anything",
						Expression: "<any>$1</any>",
					},
				},
			},
		},
	}

	if err := validSchema.Check(); err != nil {
		t.Fatalf("expected no error but got %v", err)
	}
}

func TestInvalidSchemaCheckForArgumentsAndTargets(t *testing.T) {
	testCases := []struct {
		name            string
		schema          schema.Schema
		expectedError   error
		expectedMessage string
	}{
		{
			name: "empty argument name",
			schema: schema.Schema{
				Mint: "v0.1",
				Source: schema.Source{
					Commands: []schema.Command{
						{
							Command: "paragraph",
							Arguments: []schema.Argument{
								{},
							},
						},
					},
				},
			},
			expectedError: schema.ErrEmptyArgumentName,
		},
		{
			name: "invalid argument name",
			schema: schema.Schema{
				Mint: "v0.1",
				Source: schema.Source{
					Commands: []schema.Command{
						{
							Command: "paragraph",
							Arguments: []schema.Argument{
								{
									Name: "invalid name",
									Type: schema.ArgumentTypeString,
								},
							},
						},
					},
				},
			},
			expectedError: schema.ErrInvalidArgumentName,
		},
		{
			name: "invalid argument type",
			schema: schema.Schema{
				Mint: "v0.1",
				Source: schema.Source{
					Commands: []schema.Command{
						{
							Command: "paragraph",
							Arguments: []schema.Argument{
								{
									Name: "value",
									Type: schema.ArgumentType("json"),
								},
							},
						},
					},
				},
			},
			expectedError: schema.ErrInvalidArgumentType,
		},
		{
			name: "duplicate target name",
			schema: schema.Schema{
				Mint: "v0.1",
				Source: schema.Source{
					Commands: []schema.Command{
						{Command: "paragraph"},
					},
				},
				Targets: []schema.Target{
					{Name: "html"},
					{Name: "html"},
				},
			},
			expectedError: schema.ErrMultipleTargetDefinitions,
		},
		{
			name: "empty target name",
			schema: schema.Schema{
				Mint: "v0.1",
				Source: schema.Source{
					Commands: []schema.Command{
						{Command: "paragraph"},
					},
				},
				Targets: []schema.Target{
					{},
				},
			},
			expectedError: schema.ErrEmptyTargetName,
		},
		{
			name: "invalid target name",
			schema: schema.Schema{
				Mint: "v0.1",
				Source: schema.Source{
					Commands: []schema.Command{
						{Command: "paragraph"},
					},
				},
				Targets: []schema.Target{
					{Name: "html target"},
				},
			},
			expectedError: schema.ErrInvalidTargetName,
		},
		{
			name: "unknown command in target",
			schema: schema.Schema{
				Mint: "v0.1",
				Source: schema.Source{
					Commands: []schema.Command{
						{Command: "paragraph"},
					},
				},
				Targets: []schema.Target{
					{
						Name: "html",
						Commands: []struct {
							Command    string `yaml:"command"`
							Expression string `yaml:"expression"`
						}{
							{
								Command:    "title",
								Expression: "<h1>$1</h1>",
							},
						},
					},
				},
			},
			expectedError: schema.ErrUnknownCommand,
		},
		{
			name: "too many target arguments",
			schema: schema.Schema{
				Mint: "v0.1",
				Source: schema.Source{
					Commands: []schema.Command{
						{
							Command: "paragraph",
							Arguments: []schema.Argument{
								{
									Name: "text",
									Type: schema.ArgumentTypeString,
								},
							},
						},
					},
				},
				Targets: []schema.Target{
					{
						Name: "html",
						Commands: []struct {
							Command    string `yaml:"command"`
							Expression string `yaml:"expression"`
						}{
							{
								Command:    "paragraph",
								Expression: "<p>$1 $2</p>",
							},
						},
					},
				},
			},
			expectedMessage: "command paragraph in target html has more arguments than the command definition",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.schema.Check()

			if testCase.expectedMessage != "" {
				if err == nil {
					t.Fatal("expected an error but got nil")
				}

				if got, want := err.Error(), testCase.expectedMessage; got != want {
					t.Fatalf("expected error %q but got %q", want, got)
				}

				return
			}

			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("expected error %v but got %v", testCase.expectedError, err)
			}
		})
	}
}

func TestArgumentCheckValidationErrors(t *testing.T) {
	testCases := []struct {
		name          string
		argument      schema.Argument
		expectedError string
	}{
		{
			name: "minimum not supported for string",
			argument: schema.Argument{
				Name: "value",
				Type: schema.ArgumentTypeString,
				Min:  1,
			},
			expectedError: "minimum value is not supported for string type",
		},
		{
			name: "maximum not supported for boolean",
			argument: schema.Argument{
				Name: "value",
				Type: schema.ArgumentTypeBoolean,
				Max:  1,
			},
			expectedError: "maximum value is not supported for boolean type",
		},
		{
			name: "pattern not supported for number",
			argument: schema.Argument{
				Name:    "value",
				Type:    schema.ArgumentTypeNumber,
				Pattern: ".*",
			},
			expectedError: "pattern is not supported for number type",
		},
		{
			name: "invalid regex pattern",
			argument: schema.Argument{
				Name:    "value",
				Type:    schema.ArgumentTypeString,
				Pattern: "(",
			},
			expectedError: "invalid argument pattern",
		},
		{
			name: "default string does not match pattern",
			argument: schema.Argument{
				Name:         "value",
				Type:         schema.ArgumentTypeString,
				Pattern:      "^[a-z]+$",
				DefaultValue: "123",
			},
			expectedError: "default value 123 does not match pattern ^[a-z]+$",
		},
		{
			name: "invalid numeric default",
			argument: schema.Argument{
				Name:         "value",
				Type:         schema.ArgumentTypeNumber,
				DefaultValue: "NaNnope",
			},
			expectedError: "invalid default value NaNnope for number type",
		},
		{
			name: "invalid boolean default",
			argument: schema.Argument{
				Name:         "value",
				Type:         schema.ArgumentTypeBoolean,
				DefaultValue: "yes",
			},
			expectedError: "default value yes is not a boolean",
		},
		{
			name: "reference default not supported",
			argument: schema.Argument{
				Name:         "value",
				Type:         schema.ArgumentTypeReference,
				DefaultValue: "node-1",
			},
			expectedError: "default value is not supported for reference type",
		},
		{
			name: "any default not supported",
			argument: schema.Argument{
				Name:         "value",
				Type:         schema.ArgumentTypeAny,
				DefaultValue: "anything",
			},
			expectedError: "default value is not supported for any type",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.argument.Check()
			if err == nil {
				t.Fatal("expected an error but got nil")
			}

			if got := err.Error(); !strings.HasPrefix(got, testCase.expectedError) {
				t.Fatalf("expected error starting with %q but got %q", testCase.expectedError, got)
			}
		})
	}
}

func TestSchemaCheckCompatibility(t *testing.T) {
	testCases := []struct {
		name            string
		schema          schema.Schema
		compilerVersion string
		expectedError   error
		expectedMessage string
	}{
		{
			name: "pre 1.0 exact minor match",
			schema: schema.Schema{
				Mint: "v0.3.5",
			},
			compilerVersion: "0.3.0",
		},
		{
			name: "pre 1.0 different minor is rejected",
			schema: schema.Schema{
				Mint: "v0.4.0",
			},
			compilerVersion: "0.3.1",
			expectedError:   schema.ErrIncompatibleMintVersion,
			expectedMessage: "must match the compiler major and minor version exactly",
		},
		{
			name: "major version must match after 1.0",
			schema: schema.Schema{
				Mint: "v2.0.0",
			},
			compilerVersion: "1.4.0",
			expectedError:   schema.ErrIncompatibleMintVersion,
			expectedMessage: "supported schemas must use major version 1",
		},
		{
			name: "after 1.0 older minor versions are supported",
			schema: schema.Schema{
				Mint: "v1.2.7",
			},
			compilerVersion: "1.4.0",
		},
		{
			name: "after 1.0 newer minor versions are rejected",
			schema: schema.Schema{
				Mint: "v1.5.0",
			},
			compilerVersion: "1.4.2",
			expectedError:   schema.ErrIncompatibleMintVersion,
			expectedMessage: "must have minor version <= 4",
		},
		{
			name: "compiler version can include v prefix",
			schema: schema.Schema{
				Mint: "v1.4.0",
			},
			compilerVersion: "v1.4.2",
		},
		{
			name: "invalid compiler version format",
			schema: schema.Schema{
				Mint: "v1.4.0",
			},
			compilerVersion: "version-one",
			expectedError:   schema.ErrInvalidVersionFormat,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.schema.CheckCompatibility(testCase.compilerVersion)

			if testCase.expectedError == nil {
				if err != nil {
					t.Fatalf("expected no error but got %v", err)
				}

				return
			}

			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("expected error %v but got %v", testCase.expectedError, err)
			}

			if testCase.expectedMessage != "" && !strings.Contains(err.Error(), testCase.expectedMessage) {
				t.Fatalf("expected error to contain %q but got %q", testCase.expectedMessage, err.Error())
			}
		})
	}
}

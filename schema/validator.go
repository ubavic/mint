package schema

import (
	"errors"
	"fmt"
	"slices"

	"github.com/ubavic/mint/parser"
)

var ErrCommandNotAllowed = errors.New("command not allowed")
var ErrCommandNotFound = errors.New("command not found")
var ErrCommandInvalidArguments = errors.New("command has invalid arguments")
var ErrGroupNotFound = errors.New("group not found")

func (s *Schema) Validate(document parser.Element) error {
	return s.validate(document, nil)
}

func (s *Schema) validate(document parser.Element, parent *string) error {
	var parentAllowedCommands []string
	var err error

	if parent == nil {
		if s.Source.AllowedRootChildren != "" {
			parentAllowedCommands, err = s.getGroup(s.Source.AllowedRootChildren)
			if err != nil {
				return fmt.Errorf("%w: parent allowed commands", err)
			}
		}
	} else {
		command, err := s.getCommand(*parent)
		if err != nil {
			return fmt.Errorf("%w: command %s", err, command.Command)
		}

		parentAllowedCommands, err = s.getGroup(command.Command)
		if err != nil {
			return fmt.Errorf("%w: command %s", err, command.Command)
		}
	}

	if parentAllowedCommands == nil {
		return nil
	}

	for _, el := range document.Content() {
		if command, ok := el.(*parser.Command); ok {
			if !slices.Contains(parentAllowedCommands, command.Name) {
				return fmt.Errorf("%w: command %s is not in list %v", ErrCommandNotFound, command.Name, parentAllowedCommands)
			}
		}
	}

	return nil
}

func (s *Schema) ValidateSingleCommand(name string, args int) error {
	for _, command := range s.Source.Commands {
		if command.Command == name {
			if command.Arguments == args {
				return nil
			} else {
				return fmt.Errorf("%w: command %s requires %d arguments, but %d is given", ErrCommandInvalidArguments, name, command.Arguments, args)
			}
		}
	}

	return fmt.Errorf("%w: command %s is not found in the schema", ErrCommandNotFound, name)
}

func (s *Schema) getCommand(commandName string) (*Command, error) {
	for _, command := range s.Source.Commands {
		if command.Command == commandName {
			return &command, nil
		}
	}

	return nil, ErrCommandNotFound
}

func (s *Schema) getGroup(groupName string) ([]string, error) {
	for _, group := range s.Source.Groups {
		if group.Name == groupName {
			return group.Commands, nil
		}
	}

	return nil, ErrGroupNotFound
}

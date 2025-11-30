package schema

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
)

var ErrMultipleComandDefinitions = errors.New("command defined multiple times")
var ErrEmptyCommandName = errors.New("command name is empty")
var ErrInvalidCommandName = errors.New("command name is invalid")
var ErrUnknownCommand = errors.New("unknown command")

var ErrMultipleGroupDefinitions = errors.New("command defined multiple times")
var ErrEmptyGroupName = errors.New("command name is empty")
var ErrInvalidGroupName = errors.New("command name is invalid")
var ErrUnknownGroup = errors.New("unknown group")

var ErrInvalidVersionFormat = errors.New("invalid version format")

func (s Schema) Check() error {

	nameValidator := regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]*$")
	versionValidator := regexp.MustCompile(`^v\d+\.\d+(\.\d+)?$`)

	if !versionValidator.MatchString(s.Mint) {
		return fmt.Errorf("%w: %s", ErrInvalidVersionFormat, s.Mint)
	}

	if s.Version != "" {
		if !versionValidator.MatchString(s.Version) {
			return fmt.Errorf("%w: %s", ErrInvalidVersionFormat, s.Version)
		}
	}

	commandNames := make([]string, 0, len(s.Source.Commands))
	for _, cmd := range s.Source.Commands {
		if cmd.Command == "" {
			return ErrEmptyCommandName
		}

		if !nameValidator.MatchString(cmd.Command) {
			return fmt.Errorf("%w: %s", ErrInvalidCommandName, cmd.Command)
		}

		if slices.Contains(commandNames, cmd.Command) {
			return fmt.Errorf("%w: %s", ErrMultipleComandDefinitions, cmd.Command)
		} else {
			commandNames = append(commandNames, cmd.Command)
		}
	}

	groupNames := make([]string, 0, len(s.Source.Commands))
	for _, group := range s.Source.Groups {
		if group.Name == "" {
			return ErrEmptyGroupName
		}

		if !nameValidator.MatchString(group.Name) {
			return fmt.Errorf("%w: %s", ErrInvalidGroupName, group.Name)
		}

		if slices.Contains(groupNames, group.Name) {
			return fmt.Errorf("%w: %s", ErrMultipleGroupDefinitions, group.Name)
		} else {
			groupNames = append(commandNames, group.Name)
		}

		for _, cmd := range group.Commands {
			if !slices.Contains(commandNames, cmd) {
				return fmt.Errorf("%w %s in group %s", ErrUnknownCommand, cmd, group.Name)
			}
		}
	}

	if s.Source.AllowedRootCommands != "" {
		if !slices.Contains(groupNames, s.Source.AllowedRootCommands) {
			return fmt.Errorf("%w %s as allowed root comands", ErrUnknownGroup, s.Source.AllowedRootCommands)
		}
	}

	return nil
}

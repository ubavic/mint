package schema

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var ErrMultipleComandDefinitions = errors.New("command defined multiple times")
var ErrEmptyCommandName = errors.New("command name is empty")
var ErrInvalidCommandName = errors.New("command name is invalid")
var ErrUnknownCommand = errors.New("unknown command")

var ErrEmptyArgumentName = errors.New("argument name is empty")
var ErrInvalidArgumentName = errors.New("argument name is invalid")
var ErrInvalidArgumentType = errors.New("argument type is invalid")

var ErrMultipleGroupDefinitions = errors.New("command defined multiple times")
var ErrEmptyGroupName = errors.New("command name is empty")
var ErrInvalidGroupName = errors.New("command name is invalid")
var ErrUnknownGroup = errors.New("unknown group")

var ErrEmptyTargetName = errors.New("target name is empty")
var ErrInvalidTargetName = errors.New("target name is invalid")
var ErrMultipleTargetDefinitions = errors.New("target defined multiple times")
var ErrUnknownTarget = errors.New("unknown target")

var ErrInvalidVersionFormat = errors.New("invalid version format")
var ErrIncompatibleMintVersion = errors.New("schema mint version is not supported by this compiler")

var nameValidator = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]*$")
var schemaVersionValidator = regexp.MustCompile(`^v\d+\.\d+(\.\d+)?$`)
var compilerVersionValidator = regexp.MustCompile(`^v?\d+\.\d+(\.\d+)?$`)

type mintVersion struct {
	Major int
	Minor int
	Patch int
}

func (s Schema) Check() error {
	if !schemaVersionValidator.MatchString(s.Mint) {
		return fmt.Errorf("%w: %s", ErrInvalidVersionFormat, s.Mint)
	}

	if s.Version != "" {
		if !schemaVersionValidator.MatchString(s.Version) {
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

		for _, arg := range cmd.Arguments {
			if err := arg.Check(); err != nil {
				return fmt.Errorf("command %s has invalid argument %s: %w", cmd.Command, arg.Name, err)
			}
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

	targetNames := make([]string, 0, len(s.Targets))
	for _, target := range s.Targets {
		if slices.Contains(targetNames, target.Name) {
			return fmt.Errorf("%w: %s", ErrMultipleTargetDefinitions, target.Name)
		} else {
			targetNames = append(targetNames, target.Name)
		}

		if target.Name == "" {
			return ErrEmptyTargetName
		}

		if !nameValidator.MatchString(target.Name) {
			return fmt.Errorf("%w: %s", ErrInvalidTargetName, target.Name)
		}

		argumentMatcher := regexp.MustCompile(`\$(\d+)`)

		for _, cmd := range target.Commands {
			matches := argumentMatcher.FindAllStringSubmatch(cmd.Expression, -1)
			commandIndex := slices.Index(commandNames, cmd.Command)

			if commandIndex == -1 {
				return fmt.Errorf("%w %s in target %s", ErrUnknownCommand, cmd.Command, target.Name)
			}

			if len(matches) > len(s.Source.Commands[commandIndex].Arguments) {
				return fmt.Errorf("command %s in target %s has more arguments than the command definition", cmd.Command, target.Name)
			}
		}
	}

	return nil
}

func (s Schema) CheckCompatibility(compilerVersion string) error {
	schemaVersion, err := parseMintVersion(s.Mint, true)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidVersionFormat, s.Mint)
	}

	compilerSemver, err := parseMintVersion(compilerVersion, false)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidVersionFormat, compilerVersion)
	}

	compilerDisplayVersion := schemaVersionString(compilerSemver)

	if schemaVersion.Major != compilerSemver.Major {
		return fmt.Errorf("%w: schema %s is incompatible with compiler %s; supported schemas must use major version %d", ErrIncompatibleMintVersion, s.Mint, compilerDisplayVersion, compilerSemver.Major)
	}

	if compilerSemver.Major == 0 {
		if schemaVersion.Minor != compilerSemver.Minor {
			return fmt.Errorf("%w: schema %s is incompatible with compiler %s; while Mint is below v1.0.0, supported schemas must match the compiler major and minor version exactly (expected v0.%d.x)", ErrIncompatibleMintVersion, s.Mint, compilerDisplayVersion, compilerSemver.Minor)
		}

		return nil
	}

	if schemaVersion.Minor > compilerSemver.Minor {
		return fmt.Errorf("%w: schema %s is incompatible with compiler %s; supported schemas for major version %d must have minor version <= %d", ErrIncompatibleMintVersion, s.Mint, compilerDisplayVersion, compilerSemver.Major, compilerSemver.Minor)
	}

	return nil
}

func (argument *Argument) Check() error {
	if argument.Name == "" {
		return ErrEmptyArgumentName
	}

	if !nameValidator.MatchString(argument.Name) {
		return fmt.Errorf("%w: %s", ErrInvalidArgumentName, argument.Name)
	}

	if argument.Type != ArgumentTypeNumber {
		if argument.Min != 0 {
			return fmt.Errorf("minimum value is not supported for %s type", argument.Type)
		}
		if argument.Max != 0 {
			return fmt.Errorf("maximum value is not supported for %s type", argument.Type)
		}
	}

	if argument.Type != ArgumentTypeString {
		if argument.Pattern != "" {
			return fmt.Errorf("pattern is not supported for %s type", argument.Type)
		}
	}

	switch argument.Type {
	case ArgumentTypeString:
		if argument.Pattern != "" {
			pattern, err := regexp.Compile(argument.Pattern)
			if err != nil {
				return fmt.Errorf("invalid argument pattern: %w", err)
			}

			if argument.DefaultValue != "" {
				if !pattern.MatchString(argument.DefaultValue) {
					return fmt.Errorf("default value %s does not match pattern %s", argument.DefaultValue, argument.Pattern)
				}
			}
		}
	case ArgumentTypeNumber:
		if argument.DefaultValue != "" {
			defaultValueFloat, err := strconv.ParseFloat(argument.DefaultValue, 64)
			if err != nil {
				return fmt.Errorf("invalid default value %s for number type: %w", argument.DefaultValue, err)
			}
			argument.defaultValueFloat = float64(defaultValueFloat)
		}
	case ArgumentTypeBoolean:
		if argument.DefaultValue != "" {
			if argument.DefaultValue != "true" && argument.DefaultValue != "false" {
				return fmt.Errorf("default value %s is not a boolean", argument.DefaultValue)
			}
			argument.defaultValueBool = argument.DefaultValue == "true"
		}
	case ArgumentTypeReference:
		if argument.DefaultValue != "" {
			return fmt.Errorf("default value is not supported for reference type")
		}

	case ArgumentTypeAny, "":
		if argument.DefaultValue != "" {
			return fmt.Errorf("default value is not supported for any type")
		}
	default:
		return fmt.Errorf("%w: %s", ErrInvalidArgumentType, argument.Type)
	}

	return nil
}

func parseMintVersion(version string, requirePrefix bool) (mintVersion, error) {
	validator := compilerVersionValidator
	if requirePrefix {
		validator = schemaVersionValidator
	}

	if !validator.MatchString(version) {
		return mintVersion{}, ErrInvalidVersionFormat
	}

	parts := strings.Split(strings.TrimPrefix(version, "v"), ".")
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return mintVersion{}, fmt.Errorf("%w: %s", ErrInvalidVersionFormat, version)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return mintVersion{}, fmt.Errorf("%w: %s", ErrInvalidVersionFormat, version)
	}

	patch := 0
	if len(parts) == 3 {
		patch, err = strconv.Atoi(parts[2])
		if err != nil {
			return mintVersion{}, fmt.Errorf("%w: %s", ErrInvalidVersionFormat, version)
		}
	}

	return mintVersion{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

func schemaVersionString(version mintVersion) string {
	return fmt.Sprintf("v%d.%d.%d", version.Major, version.Minor, version.Patch)
}

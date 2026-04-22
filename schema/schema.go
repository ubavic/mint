package schema

type Schema struct {
	Mint    string   `yaml:"mint"`
	Name    string   `yaml:"name"`
	Author  string   `yaml:"author"`
	Version string   `yaml:"version"`
	Source  Source   `yaml:"source"`
	Targets []Target `yaml:"targets"`
}

type Command struct {
	Command     string     `yaml:"command"`
	Arguments   []Argument `yaml:"arguments"`
	Description string     `yaml:"description"`
}

type ArgumentType string

const (
	ArgumentTypeString    ArgumentType = "string"
	ArgumentTypeNumber    ArgumentType = "number"
	ArgumentTypeBoolean   ArgumentType = "boolean"
	ArgumentTypeReference ArgumentType = "reference"
	ArgumentTypeAny       ArgumentType = "any"
)

type Argument struct {
	Name              string
	Description       string
	Type              ArgumentType
	DefaultValue      string  `yaml:"defaultValue,omitempty"`
	defaultValueFloat float64 `yaml:"defaultValueFloat,omitempty"`
	defaultValueBool  bool    `yaml:"defaultValueBool,omitempty"`
	Min               float64 `yaml:"min,omitempty"`
	Max               float64 `yaml:"max,omitempty"`
	Pattern           string  `yaml:"pattern,omitempty"`
}

type Target struct {
	Name      string `yaml:"name"`
	Extension string `yaml:"extension"`
	Commands  []struct {
		Command    string `yaml:"command"`
		Expression string `yaml:"expression"`
	} `yaml:"commands"`
}

type Source struct {
	AllowedRootCommands string    `yaml:"allowedRootChildren,omitempty"`
	Commands            []Command `yaml:"commands"`
	Groups              []Group   `yaml:"groups"`
}

type Group struct {
	Name     string   `yaml:"name"`
	Commands []string `yaml:"commands"`
}

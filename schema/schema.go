package schema

type Schema struct {
	Mint    string   `yaml:"mint"`
	Name    string   `yaml:"name"`
	Version string   `yaml:"version"`
	Source  Source   `yaml:"source"`
	Targets []Target `yaml:"targets"`
}

type Command struct {
	Command     string `yaml:"command"`
	Arguments   int    `yaml:"arguments"`
	Description string `yaml:"description"`
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
	AllowedRootCommands string    `yaml:"allowedRootChildren"`
	Commands            []Command `yaml:"commands"`
	Groups              []Group   `yaml:"groups"`
}

type Group struct {
	Name     string   `yaml:"name"`
	Commands []string `yaml:"commands"`
}

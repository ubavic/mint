package parser

import (
	"encoding/json"
	"fmt"
)

type Element interface {
	fmt.Stringer
	Content() []Element
}

type Block struct {
	Nodes []Element
}

func (doc *Block) Content() []Element {
	return doc.Nodes
}

func (doc *Block) String() string {
	result := ""

	for _, node := range doc.Nodes {
		result += node.String() + "\n"
	}

	return result
}

func (block *Block) Json() []byte {
	json, err := json.Marshal(block)
	if err != nil {
		panic(err.Error())
	}

	return json
}

type Command struct {
	Name      string
	Arguments []Element
}

func (com *Command) Content() []Element {
	return com.Arguments
}

func (cmd *Command) String() string {
	result := "@" + cmd.Name + "\n"

	for _, arg := range cmd.Arguments {
		result += "  " + arg.String() + "\n"
	}

	return result
}

type TextContent struct {
	TextContent string
}

func (tc *TextContent) Content() []Element {
	return []Element{tc}
}

func (tc *TextContent) String() string {
	return tc.TextContent
}

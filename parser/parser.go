package parser

import (
	"errors"
	"fmt"
)

type Parser struct {
	tokens          []Token
	currentPosition int
	validator       Validator
}

func NewParser(tokens []Token, validator Validator) Parser {
	if validator == nil {
		panic("Validator must not be nil.")
	}

	parser := Parser{
		tokens:          tokens,
		currentPosition: 0,
		validator:       validator,
	}

	return parser
}

func (p *Parser) Parse() (*Block, error) {
	document, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	if p.currentToken().Type != EOF {
		return nil, errors.New("didn't reach end of file")
	}

	err = p.validator.Validate(document)
	if err != nil {
		return nil, fmt.Errorf("parsing error: %w", err)
	}

	return document, nil
}

func (p *Parser) parseBlock() (*Block, error) {
	block := Block{
		Nodes: []Element{},
	}

	for {
		currentToken := p.currentToken()

		switch currentToken.Type {
		case EOF:
			return &block, nil
		case Identifier:
			command := Command{
				Name: currentToken.Content,
			}

			args, err := p.parseArguments()
			if err != nil {
				return nil, err
			}

			err = p.validator.ValidateSingleCommand(command.Name, len(args))
			if err != nil {
				return nil, err
			}

			command.Arguments = args
			block.Nodes = append(block.Nodes, &command)
		case Text:
			tc := TextContent{TextContent: currentToken.Content}
			block.Nodes = append(block.Nodes, &tc)
		case LeftBrace:
			return nil, fmt.Errorf("%w %s", ErrUnexpectedToken, currentToken.String())
		case RightBrace:
			return &block, nil
		}

		p.next()
	}

}

func (p *Parser) parseArguments() ([]Element, error) {
	arguments := []Element{}
	start := true

	for {
		currentToken := p.currentToken()

		switch currentToken.Type {
		case LeftBrace:
			element, err := p.parseArgument()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, element)
		case Text:
			if currentToken.ContainsWhitespaceOnly() {
				p.next()
			} else {
				p.unwind()
				return arguments, nil
			}
		case RightBrace, EOF:
			p.unwind()
			return arguments, nil
		case Identifier:
			if !start {
				p.unwind()
				return arguments, nil
			}
			p.next()
		}

		start = false
	}
}

func (p *Parser) parseArgument() (Element, error) {
	err := p.parseToken(LeftBrace)
	if err != nil {
		return nil, err
	}

	el, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	err = p.parseToken(RightBrace)
	if err != nil {
		return nil, err
	}

	return el, nil

}

func (p *Parser) parseToken(tt TokenType) error {
	if p.currentToken().Type == tt {
		p.next()
		return nil
	} else {
		return fmt.Errorf("%w: expected %s, got %s", ErrUnexpectedToken, tt.String(), p.peek().Type.String())
	}
}

func (p *Parser) next() {
	p.currentPosition += 1
}

func (p *Parser) unwind() {
	p.currentPosition -= 1
}

func (p *Parser) peek() Token {
	if p.currentPosition+1 >= len(p.tokens) {
		return Token{Type: EOF}
	}

	return p.tokens[p.currentPosition+1]
}

func (p *Parser) currentToken() Token {
	if p.currentPosition >= len(p.tokens) {
		return Token{Type: EOF}
	}

	return p.tokens[p.currentPosition]
}

package parser

type Validator interface {
	Validate(document Element) error
	ValidateSingleCommand(name string, args int) error
}

type OptimisticValidator struct{}

func (b *OptimisticValidator) Validate(document Element) error {
	return nil
}

func (s *OptimisticValidator) ValidateSingleCommand(name string, args int) error {
	return nil
}

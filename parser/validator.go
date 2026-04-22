package parser

type Validator interface {
	Validate(document Element) error
	ValidateCommand(name string, args []Element) error
}

type OptimisticValidator struct{}

func (b *OptimisticValidator) Validate(document Element) error {
	return nil
}

func (s *OptimisticValidator) ValidateCommand(name string, args []Element) error {
	return nil
}

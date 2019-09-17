package usecases

//go:generate mockery -name=CodeGenerator -inpkg -case=underscore -testonly

// CodeGenerator ...
type CodeGenerator interface {
	GenerateCode() (string, error)
}

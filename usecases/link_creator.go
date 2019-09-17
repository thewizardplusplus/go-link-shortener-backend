package usecases

// CodeGenerator ...
type CodeGenerator interface {
	GenerateCode() (string, error)
}

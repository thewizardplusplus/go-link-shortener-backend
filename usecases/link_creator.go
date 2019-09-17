package usecases

//go:generate mockery -name=CodeGenerator -inpkg -case=underscore -testonly

// CodeGenerator ...
type CodeGenerator interface {
	GenerateCode() (string, error)
}

// LinkCreator ...
type LinkCreator struct {
	LinkGetter    LinkGetter
	LinkSetter    LinkSetter
	CodeGenerator CodeGenerator
}

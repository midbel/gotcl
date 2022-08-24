package interp

type Namespace struct {
	Name     string
	Parent   *Namespace
	Children []*Namespace
	Exported []string

	CommandSet
	*Env
}

func Prepare(name string) *Namespace {
	return &Namespace{
		Name:       name,
		CommandSet: EmptySet(),
		Env:        Environ(),
	}
}

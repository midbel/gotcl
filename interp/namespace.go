package interp

type Namespace struct {
	Name     string
	Parent   *Namespace
	Children []*Namespace
	Exported []string

	CommandSet
	*Env
}

func Global() *Namespace {
	return &Namespace{
		CommandSet: DefaultSet(),
		Env:        Environ(),
	}
}

func Prepare(name string) *Namespace {
	return &Namespace{
		Name:       name,
		CommandSet: EmptySet(),
		Env:        Environ(),
	}
}

func (ns *Namespace) Root() bool {
	return ns.Parent == nil
}

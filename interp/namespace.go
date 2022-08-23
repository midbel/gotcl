package interp

type Namespace struct {
	Name     string
	Parent   *Namespace
	Children []*Namespace

	CommandSet
	*Env
}

package interp

import (
	"fmt"
	"sort"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/gotcl/stdlib"
)

type Namespace struct {
	Name     string
	Parent   *Namespace
	Children []*Namespace
	Exported []string

	CommandSet
	env env.Environment
}

func Global() *Namespace {
	root := create("", DefaultSet())
	tcl := create("tcl", EmptySet())
	tcl.Parent = root

	mop := Mathop()
	mop.Parent = tcl

	tcl.Children = append(tcl.Children, mop)
	root.Children = append(root.Children, tcl)

	return root
}

func Mathop() *Namespace {
	return create("mathop", MathopSet())
}

func Prepare(name string) *Namespace {
	return create(name, EmptySet())
}

func (ns *Namespace) Root() bool {
	return ns.Parent == nil
}

func (ns *Namespace) Command(names []string) (stdlib.Executer, error) {
	if names[0] == "" && !ns.Root() {
		return nil, fmt.Errorf("relative search from non root namespace")
	}
	if names[0] == "" {
		names = names[1:]
	}
	if len(names) == 1 {
		exec, err := ns.Lookup(names[0])
		if err != nil && !ns.Root() {
			return ns.Parent.Command(names)
		}
		return exec, err
	}
	sub, _, err := ns.getNS(names[0])
	if err != nil {
		return nil, err
	}
	return sub.Command(names[1:])
}

func (ns *Namespace) Lookup(name string) (stdlib.Executer, error) {
	exec, err := ns.CommandSet.Lookup(name)
	if err == nil {
		return exec, err
	}
	if ns.Root() {
		return nil, fmt.Errorf("%s: undefined proc", name)
	}
	return ns.Parent.Lookup(name)
}

func (ns *Namespace) RegisterProc(name, args, body string) error {
	return ns.CommandSet.RegisterProc(name, args, body)
}

func (ns *Namespace) Get(names []string) (*Namespace, error) {
	sub, i, err := ns.getNS(names[0])
	if err == nil {
		if len(names) == 1 {
			return sub, nil
		}
		return sub.Get(names[1:])
	}
	var (
		curr = Prepare(names[0])
		tmp  = append([]*Namespace{curr}, ns.Children[i:]...)
	)
	curr.Parent = ns
	ns.Children = append(ns.Children[:i], tmp...)
	return curr, nil
}

func (ns *Namespace) getNS(name string) (*Namespace, int, error) {
	i := sort.Search(len(ns.Children), func(i int) bool {
		return name >= ns.Children[i].Name
	})
	if i < len(ns.Children) && ns.Children[i].Name == name {
		return ns.Children[i], i, nil
	}
	return nil, i, fmt.Errorf("%s: namespace not defined", name)
}

func create(name string, set CommandSet) *Namespace {
	return &Namespace{
		Name:       name,
		CommandSet: set,
		env:        env.EmptyEnv(),
	}
}

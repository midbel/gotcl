package interp

import (
	"errors"
	"fmt"
	"os/exec"
	"sort"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/gotcl/stdlib"
	"github.com/midbel/slices"
)

var ErrUndefined = errors.New("undefined name")

type Namespace struct {
	Name     string
	parent   *Namespace
	children []*Namespace

	env *env.Env
	CommandSet
	unknown stdlib.CommandFunc
}

func EmptyNS(name string) *Namespace {
	return createNS(name, make(CommandSet))
}

func GlobalNS() *Namespace {
	ns := createNS("", DefaultSet())
	ns.unknown = func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
		var (
			name   = slices.Fst(args).String()
			values []string
		)
		for _, a := range slices.Rest(args) {
			values = append(values, a.String())
		}
		res, err := exec.Command(name, values...).Output()
		return env.Str(string(res)), err
	}

	var (
		mathfunc = createNS("mathfunc", MathfuncSet())
		mathop   = createNS("mathop", MathfuncSet())
		tcl      = emptyNS("tcl")
	)
	tcl.RegisterNS(mathfunc)
	tcl.RegisterNS(mathop)
	ns.RegisterNS(tcl)
	ns.RegisterNS(UtilNS())
	return ns
}

func UtilNS() *Namespace {
	ns := createNS("util", UtilSet())
	ns.env.Define("version", env.Str("1.12.189"))
	return ns
}

func emptyNS(name string) *Namespace {
	return createNS(name, make(CommandSet))
}

func createNS(name string, set CommandSet) *Namespace {
	return &Namespace{
		Name:       name,
		CommandSet: set,
		env:        env.EmptyEnv(),
	}
}

func (n *Namespace) GetName() string {
	return n.Name
}

func (n *Namespace) Resolve(v string) (env.Value, error) {
	return n.env.Resolve(v)
}

func (n *Namespace) RegisterNS(ns *Namespace) error {
	ns.parent = n
	x := sort.Search(len(n.children), func(i int) bool {
		return ns.Name >= n.children[i].Name
	})
	if x < len(n.children) && n.children[x].Name == ns.Name {
		return fmt.Errorf("namespace %s already exists", ns.Name)
	}
	tmp := append([]*Namespace{ns}, n.children[x:]...)
	n.children = append(n.children[:x], tmp...)
	return nil
}

func (n *Namespace) LookupNS(name []string) (*Namespace, error) {
	if len(name) == 0 {
		return n, nil
	}
	name, err := n.validNS(name)
	if err != nil {
		return nil, err
	}
	ns, err := n.lookupNS(name[0])
	if err == nil {
		if len(name) == 1 {
			return ns, nil
		}
		return ns.LookupNS(name[1:])
	}
	return nil, err
}

func (n *Namespace) RegisterExec(name []string, exec stdlib.Executer) error {
	name, err := n.validNS(name)
	if err != nil {
		return err
	}
	if len(name) == 1 {
		n.CommandSet[name[0]] = exec
		return nil
	}
	ns, err := n.lookupNS(name[0])
	if err == nil {
		return ns.RegisterExec(name[1:], exec)
	}
	return err
}

func (n *Namespace) LookupExec(name []string) (stdlib.Executer, error) {
	name, err := n.validNS(name)
	if err != nil {
		return nil, err
	}
	if len(name) > 1 && len(n.children) == 0 {
		return nil, fmt.Errorf("executer (lookup) %s (%s): %w", name[0], n.FQN(), ErrUndefined)
	}
	if len(name) == 1 {
		exec, ok := n.CommandSet[name[0]]
		if ok {
			return exec, nil
		}
		if !n.Root() {
			return n.parent.LookupExec(name)
		}
		return nil, fmt.Errorf("executer (lookup) %s (%s): %w", name[0], n.FQN(), ErrUndefined)
	}
	ns, err := n.lookupNS(name[0])
	if err == nil {
		return ns.LookupExec(name[1:])
	}
	return nil, err
}

func (n *Namespace) Root() bool {
	return n.parent == nil
}

func (n *Namespace) Parent() *Namespace {
	return n.parent
}

func (n *Namespace) Children() []*Namespace {
	return n.children
}

func (n *Namespace) FQN() string {
	if n.Root() {
		return "::"
	}
	if n.parent.Root() {
		return "::" + n.Name
	}
	return n.parent.FQN() + "::" + n.Name
}

func (n *Namespace) lookupNS(name string) (*Namespace, error) {
	x := sort.Search(len(n.children), func(i int) bool {
		return name >= n.children[i].Name
	})
	if x < len(n.children) && n.children[x].Name == name {
		return n.children[x], nil
	}
	return nil, fmt.Errorf("namespace %s (%s): %w", name, n.FQN(), ErrUndefined)
}

func (n *Namespace) validNS(name []string) ([]string, error) {
	if len(name) > 0 && name[0] == "" {
		if !n.Root() {
			return nil, fmt.Errorf("namespace %s: invalid name", name)
		}
		name = name[1:]
	}
	return name, nil
}

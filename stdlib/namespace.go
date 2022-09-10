package stdlib

import (
	"fmt"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/slices"
)

type NamespaceHandler interface {
	Interpreter
	RegisterNS(string, string) error
	UnregisterNS(string) error
	CurrentNS() string
	ParentNS(string) (string, error)
	ChildrenNS(string) ([]string, error)
	HasNS(string) bool
	DefineVar(string, env.Value)
	// DefineUnknown()
}

type namespaceHandleFunc func(NamespaceHandler, []env.Value) (env.Value, error)

func RunVariable() Executer {
	return Builtin{
		Name:  "variable",
		Safe:  true,
		Arity: 1,
		Run:   wrapNamespaceFunc(runVariable),
	}
}

func MakeNamespace() Executer {
	e := Ensemble{
		Name: "namespace",
		Safe: true,
		List: []Executer{
			Builtin{
				Name:  "eval",
				Arity: 2,
				Run:   wrapNamespaceFunc(namespaceCreate),
			},
			Builtin{
				Name:  "delete",
				Arity: 1,
				Run:   wrapNamespaceFunc(namespaceDelete),
			},
			Builtin{
				Name: "current",
				Run:  wrapNamespaceFunc(namespaceCurrent),
			},
			Builtin{
				Name:  "parent",
				Arity: 1,
				Run:   wrapNamespaceFunc(namespaceParent),
			},
			Builtin{
				Name:  "children",
				Arity: 1,
				Run:   wrapNamespaceFunc(namespaceCurrent),
			},
			Builtin{
				Name:  "exists",
				Arity: 1,
				Run:   wrapNamespaceFunc(namespaceExists),
			},
			Builtin{
				Name: "unknown",
				Run:  wrapNamespaceFunc(namespaceUnknown),
			},
			Builtin{
				Name: "export",
			},
			Builtin{
				Name: "import",
			},
			Builtin{
				Name: "forget",
			},
		},
	}
	return sortEnsembleCommands(e)
}

func runVariable(i NamespaceHandler, args []env.Value) (env.Value, error) {
	i.DefineVar(slices.Fst(args).String(), slices.Lst(args))
	return nil, nil
}

func namespaceCurrent(i NamespaceHandler, args []env.Value) (env.Value, error) {
	curr := i.CurrentNS()
	return env.Str(curr), nil
}

func namespaceParent(i NamespaceHandler, args []env.Value) (env.Value, error) {
	parent, err := i.ParentNS(slices.Fst(args).String())
	if err != nil {
		return nil, err
	}
	return env.Str(parent), nil
}

func namespaceChildren(i NamespaceHandler, args []env.Value) (env.Value, error) {
	list, err := i.ChildrenNS(slices.Fst(args).String())
	if err != nil {
		return nil, err
	}
	return env.ListFromStrings(list), nil
}

func namespaceExists(i NamespaceHandler, args []env.Value) (env.Value, error) {
	has := i.HasNS(slices.Fst(args).String())
	return env.Bool(has), nil
}

func namespaceUnknown(i NamespaceHandler, args []env.Value) (env.Value, error) {
	return nil, nil
}

func namespaceCreate(i NamespaceHandler, args []env.Value) (env.Value, error) {
	err := i.RegisterNS(slices.Fst(args).String(), slices.Snd(args).String())
	return env.EmptyStr(), err
}

func namespaceDelete(i NamespaceHandler, args []env.Value) (env.Value, error) {
	err := i.UnregisterNS(slices.Fst(args).String())
	return env.EmptyStr(), err
}

func wrapNamespaceFunc(do namespaceHandleFunc) CommandFunc {
	return func(i Interpreter, args []env.Value) (env.Value, error) {
		is, ok := i.(NamespaceHandler)
		if !ok {
			return nil, fmt.Errorf("interpreter can not handle interpreter children")
		}
		return do(is, args)
	}
}

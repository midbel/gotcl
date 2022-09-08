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
}

type namespaceHandleFunc func(NamespaceHandler, []env.Value) (env.Value, error)

func MakeNamespace() Executer {
	e := Ensemble{
		Name: "namespace",
		List: []Executer{
			Builtin{
				Name:  "eval",
				Arity: 2,
				Safe:  true,
				Run:   wrapNamespaceFunc(namespaceCreate),
			},
			Builtin{
				Name:  "delete",
				Arity: 1,
				Safe:  true,
				Run:   wrapNamespaceFunc(namespaceDelete),
			},
		},
	}
	return sortEnsembleCommands(e)
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

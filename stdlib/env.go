package stdlib

import (
	"github.com/midbel/gotcl/env"
	"github.com/midbel/slices"
)

func RunSet() Executer {
	return Builtin{
		Name:  "set",
		Arity: 2,
		Safe:  true,
		Run:   runSet,
	}
}

func RunUnset() Executer {
	return Builtin{
		Name:  "unset",
		Arity: 1,
		Safe:  true,
		Options: []Option{
			{
				Name:  "nocomplain",
				Flag:  true,
				Value: env.False(),
				Check: CheckBool,
			},
		},
		Run: runUnset,
	}
}

func runSet(i Interpreter, args []env.Value) (env.Value, error) {
	i.Define(slices.Fst(args).String(), slices.Snd(args))
	return slices.Snd(args), nil
}

func runUnset(i Interpreter, args []env.Value) (env.Value, error) {
	i.Delete(slices.Fst(args).String())
	return nil, nil
}

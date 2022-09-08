package stdlib

import (
	"github.com/midbel/gotcl/env"
	"github.com/midbel/slices"
)

func RunList() Executer {
	return Builtin{
		Name:  "list",
		Arity: 1,
		Safe:  true,
		Run: runList,
	}
}

func RunListLen() Executer {
	return Builtin{
		Name:  "llength",
		Arity: 1,
		Safe:  true,
		Run: runLlength,
	}
}

func runList(i Interpreter, args []env.Value) (env.Value, error) {
	return slices.Fst(args).ToList()
}

func runLlength(i Interpreter, args []env.Value) (env.Value, error) {
	list, err := slices.Fst(args).ToList()
	if err != nil {
		return nil, err
	}
	n, ok := list.(interface{ Len() int })
	if !ok {
		return env.Int(0), nil
	}
	return env.Int(int64(n.Len())), nil
}

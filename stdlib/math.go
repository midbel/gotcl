package stdlib

import (
	"github.com/midbel/gotcl/env"
	"github.com/midbel/slices"
)

func RunIncr() Executer {
	return Builtin{
		Name:  "incr",
		Arity: 1,
		Safe:  true,
		Run: runIncr,
	}
}

func runIncr(i Interpreter, args []env.Value) (env.Value, error) {
	v, err := i.Resolve(slices.Fst(args).String())
	if err != nil {
		return nil, err
	}
	n, err := env.ToInt(v)
	if err != nil {
		return nil, err
	}
	res := env.Int(int64(n) + 1)
	i.Define(slices.Fst(args).String(), res)
	return res, nil
}

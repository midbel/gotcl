package stdlib

import (
	"strings"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/slices"
)

func RunUplevel() Executer {
	return Builtin{
		Name:     "uplevel",
		Arity:    1,
		Variadic: true,
		Safe:     false,
		Run: func(i Interpreter, args []env.Value) (env.Value, error) {
			var (
				level int
				abs   bool
			)
			if len(args) > 1 {
				x, a, err := env.ToLevel(slices.Fst(args))
				if err != nil {
					return nil, err
				}
				level, abs, args = x, a, slices.Rest(args)
			} else {
				level++
			}
			if !abs {
				level = i.Depth() - level
			}
			return i.Level(strings.NewReader(slices.Fst(args).String()), level)
		},
	}
}

func RunUpvar() Executer {
	return Builtin{
		Name:     "upvar",
		Arity:    2,
		Variadic: true,
		Safe:     false,
		Run: func(i Interpreter, args []env.Value) (env.Value, error) {
			var level int
			if n := len(args) % 2; n == 0 {
				level++
			} else {
				x, err := env.ToInt(slices.Fst(args))
				if err != nil {
					return nil, err
				}
				level = x
				args = slices.Rest(args)
			}
			for j := 0; j < len(args); j += 2 {
				var (
					src = slices.At(args, j)
					dst = slices.At(args, j+1)
				)
				if err := i.LinkVar(src.String(), dst.String(), level); err != nil {
					return nil, err
				}
			}
			return env.EmptyStr(), nil
		},
	}
}

func RunEval() Executer {
	return Builtin{
		Name:     "eval",
		Help:     "eval given script",
		Variadic: true,
		Safe:     false,
		Run: func(i Interpreter, args []env.Value) (env.Value, error) {
			tmp := env.ListFrom(args...)
			return i.Execute(strings.NewReader(tmp.String()))
		},
	}
}

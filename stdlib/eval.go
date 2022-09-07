package stdlib

import (
	"fmt"
	"io"
	"strings"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/slices"
)

type DeferHandler interface {
	Interpreter
	RegisterDefer(string) error
}

func RunDefer() Executer {
	return Builtin{
		Name:  "defer",
		Arity: 1,
		Safe:  true,
		Run: func(i Interpreter, args []env.Value) (env.Value, error) {
			h, ok := i.(DeferHandler)
			if !ok {
				return nil, fmt.Errorf("interpreter can not register defer call")
			}
			return env.EmptyStr(), h.RegisterDefer(slices.Fst(args).String())
		},
	}
}

func RunProc() Executer {
	return Builtin{
		Name:  "proc",
		Arity: 3,
		Safe:  true,
		Run: func(i Interpreter, args []env.Value) (env.Value, error) {
			h, ok := i.(interface {
				RegisterProc(string, string, string) error
			})
			if !ok {
				return nil, fmt.Errorf("interpreter can not register defer call")
			}
			var (
				name = slices.Fst(args).String()
				list = slices.Snd(args).String()
				body = slices.Lst(args).String()
			)
			return nil, h.RegisterProc(name, body, list)
		},
	}
}

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
			n, ok := i.(interface {
				ExecuteLevel(io.Reader, int, bool) (env.Value, error)
			})
			if !ok {
				return nil, fmt.Errorf("interpreter can not execute script in a parent level")
			}
			return n.ExecuteLevel(strings.NewReader(slices.Fst(args).String()), level, abs)
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
			k, ok := i.(interface {
				LinkVar(string, string, int) error
			})
			if !ok {
				return nil, fmt.Errorf("interpreter can not create link between variables")
			}
			for j := 0; j < len(args); j += 2 {
				var (
					src = slices.At(args, j)
					dst = slices.At(args, j+1)
				)
				if err := k.LinkVar(src.String(), dst.String(), level); err != nil {
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

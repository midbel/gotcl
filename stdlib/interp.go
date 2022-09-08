package stdlib

import (
	"fmt"
	"strings"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/slices"
)

type InterpHandler interface {
	Interpreter
	RegisterInterpreter([]string, bool) (string, error)
	UnregisterInterpreter([]string) error
	LookupInterpreter(paths []string) (InterpHandler, error)
	InterpretersList() []string
	IsSafe() bool
}

type interpHandleFunc func(InterpHandler, []env.Value) (env.Value, error)

func MakeInterp() Executer {
	e := Ensemble{
		Name: "interp",
		List: []Executer{
			Builtin{
				Name:  "create",
				Arity: 1,
				Options: []Option{
					{
						Name:  "safe",
						Flag:  true,
						Value: env.False(),
						Check: CheckBool,
					},
				},
				Run: wrapInterpFunc(interpCreate),
			},
			Builtin{
				Name:  "delete",
				Arity: 1,
				Run:   wrapInterpFunc(interpDelete),
			},
			Builtin{
				Name:  "issafe",
				Arity: 1,
				Run:   wrapInterpFunc(interpIssafe),
			},
			Builtin{
				Name:     "eval",
				Variadic: true,
				Run:      wrapInterpFunc(interpEval),
			},
			Builtin{
				Name:  "children",
				Arity: 1,
				Run:   wrapInterpFunc(interpChildren),
			},
		},
	}
	return sortEnsembleCommands(e)
}

func interpCreate(i InterpHandler, args []env.Value) (env.Value, error) {
	paths, err := env.ToStringList(slices.Fst(args))
	if err != nil {
		return nil, err
	}
	safe, err := i.Resolve("safe")
	if err != nil {
		return nil, err
	}
	val, err := i.RegisterInterpreter(paths, env.ToBool(safe))
	return env.Str(val), err
}

func interpDelete(i InterpHandler, args []env.Value) (env.Value, error) {
	paths, err := env.ToStringList(slices.Fst(args))
	if err != nil {
		return nil, err
	}
	return nil, i.UnregisterInterpreter(paths)
}

func interpIssafe(i InterpHandler, args []env.Value) (env.Value, error) {
	paths, err := env.ToStringList(slices.Fst(args))
	if err != nil {
		return nil, err
	}
	i, err = i.LookupInterpreter(paths)
	if err != nil {
		return nil, err
	}
	return env.Bool(i.IsSafe()), nil
}

func interpEval(i InterpHandler, args []env.Value) (env.Value, error) {
	paths, err := env.ToStringList(slices.Fst(args))
	if err != nil {
		return nil, err
	}
	i, err = i.LookupInterpreter(paths)
	if err != nil {
		return nil, err
	}
	return i.Execute(strings.NewReader(slices.Snd(args).String()))
}

func interpChildren(i InterpHandler, args []env.Value) (env.Value, error) {
	paths, err := env.ToStringList(slices.Fst(args))
	if err != nil {
		return nil, err
	}
	i, err = i.LookupInterpreter(paths)
	if err != nil {
		return nil, err
	}
	list := i.InterpretersList()
	return env.ListFromStrings(list), nil
}

func wrapInterpFunc(do interpHandleFunc) CommandFunc {
	return func(i Interpreter, args []env.Value) (env.Value, error) {
		is, ok := i.(InterpHandler)
		if !ok {
			return nil, fmt.Errorf("interpreter can not handle interpreter children")
		}
		return do(is, args)
	}
}

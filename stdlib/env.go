package stdlib

import (
	"fmt"
	"strings"

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

func RunRename() Executer {
	return Builtin{
		Name:  "rename",
		Arity: 2,
		Run:   runRename,
	}
}

func RunAppend() Executer {
	return Builtin{
		Name:     "append",
		Arity:    1,
		Variadic: true,
		Run:      runAppend,
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

func runRename(i Interpreter, args []env.Value) (env.Value, error) {
	r, ok := i.(interface{ Rename(string, string) error })
	if !ok {
		return nil, fmt.Errorf("interpreter can not rename commands/procedure")
	}
	return nil, r.Rename(slices.Fst(args).String(), slices.Snd(args).String())
}

func runAppend(i Interpreter, args []env.Value) (env.Value, error) {
	val, err := i.Resolve(slices.Fst(args).String())
	if err != nil {
		return nil, err
	}
	list := []string{val.String()}
	for _, a := range slices.Rest(args) {
		list = append(list, a.String())
	}
	val = env.Str(strings.Join(list, ""))
	i.Define(slices.Fst(args).String(), val)
	return val, nil
}

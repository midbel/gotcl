package stdlib

import (
  "fmt"

  "github.com/midbel/gotcl/env"
  "github.com/midbel/slices"
)

func RunTypeOf() Executer {
	return Builtin{
		Name:  "typeof",
		Arity: 1,
		Safe:  true,
		Run: func(i Interpreter, args []env.Value) (env.Value, error) {
			typ := fmt.Sprintf("%T", slices.Fst(args))
			return env.Str(typ), nil
		},
	}
}

func RunDefer() Executer {
	return Builtin{
		Name:  "defer",
		Arity: 1,
		Safe:  true,
		Run: func(i Interpreter, args []env.Value) (env.Value, error) {
			var (
				name = fmt.Sprintf("defer%d", i.Count())
				body = slices.Fst(args).String()
			)
			exec, _ := createProcedure(name, body, "")
			i.registerDefer(exec)
			return env.EmptyStr(), nil
		},
	}
}

func RunHelp() Executer {
	return Builtin{
		Name:  "help",
		Help:  "retrieve help of given builtin command",

		Arity: 1,
		Safe:  true,
		Run: func(i Interpreter, args []env.Value) (env.Value, error) {
			help, err := i.GetHelp(slices.Fst(args).String())
			if err != nil {
				return nil, err
			}
			return env.Str(help), nil
		},
	}
}

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

func RunHelp() Executer {
	return Builtin{
		Name: "help",
		Help: "retrieve help of given builtin command",

		Arity: 1,
		Safe:  true,
		Run: func(i Interpreter, args []env.Value) (env.Value, error) {
			h, ok := i.(interface{ GetHelp(string) (string, error) })
			if !ok {
				return nil, fmt.Errorf("interpreter can not extract help from command")
			}
			help, err := h.GetHelp(slices.Fst(args).String())
			if err != nil {
				return nil, err
			}
			return env.Str(help), nil
		},
	}
}

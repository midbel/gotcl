package stdlib

import (
	"fmt"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/slices"
)

type PrintHandler interface {
	Print(string, string, bool) error
}

func RunPuts() Executer {
	return Builtin{
		Name:  "puts",
		Help:  "print a message to given channel (default to stdout)",
		Arity: 1,
		Safe:  true,
		Options: []Option{
			{
				Name:  "nonewline",
				Flag:  true,
				Value: env.False(),
				Check: CheckBool,
			},
			{
				Name:     "channel",
				Value:    env.Str("stdout"),
				Required: true,
				Check:    CheckString,
			},
		},
		Run: func(i Interpreter, args []env.Value) (env.Value, error) {
			ch, err := i.Resolve("channel")
			if err != nil {
				return nil, err
			}
			ph, ok := i.(PrintHandler)
			if !ok {
				return nil, fmt.Errorf("interpreter can not print message to channel")
			}
			nonl, _ := i.Resolve("nonewline")
			err = ph.Print(ch.String(), slices.Fst(args).String(), !env.ToBool(nonl))
			return env.EmptyStr(), err
		},
	}
}

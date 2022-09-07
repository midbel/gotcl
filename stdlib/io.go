package stdlib

import (
	"fmt"
	"io"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/gotcl/stdlib/argparse"
	"github.com/midbel/slices"
)

func RunPuts() Executer {
	return Builtin{
		Name:  "puts",
		Help:  "print a message to given channel (default to stdout)",
		Arity: 1,
		Safe:  true,
		Options: []argparse.Option{
			{
				Name:  "nonewline",
				Flag:  true,
				Value: env.False(),
				Check: argparse.CheckBool,
			},
			{
				Name:     "channel",
				Value:    env.Str("stdout"),
				Required: true,
				Check:    argparse.CheckString,
			},
		},
		Run: func(i Interpreter, args []env.Value) (env.Value, error) {
			str, err := i.Resolve("channel")
			if err != nil {
				return nil, err
			}
			var ch io.Writer
			switch str.String() {
			case "stdout":
				ch = i.Out
			case "stderr":
				ch = i.Err
			default:
				return nil, nil
			}
			fmt.Fprintln(ch, slices.Fst(args))
			return env.EmptyStr(), nil
		},
	}
}

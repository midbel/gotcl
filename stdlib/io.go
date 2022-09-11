package stdlib

import (
	"fmt"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/slices"
)

type ChannelHandler interface {
	Open(string, string) (string, error)
	Close(string) error
	Eof(string) (bool, error)

	Seek(string, int, int) (int64, error)
	Tell(string) (int64, error)
	Gets(string) (string, error)
	Read(string, int) (string, error)
}

type channelFunc func(ChannelHandler, []env.Value) (env.Value, error)

func wrapChannelFunc(do channelFunc) CommandFunc {
	return func(i Interpreter, args []env.Value) (env.Value, error) {
		ch, ok := i.(ChannelHandler)
		if !ok {
			return nil, fmt.Errorf("interpreter can not handle files")
		}
		return do(ch, args)
	}
}

type PrintHandler interface {
	Print(string, string, bool) error
}

func MakeChan() Executer {
	e := Ensemble{
		Name: "chan",
		List: []Executer{},
	}
	return sortEnsembleCommands(e)
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
		Run: runPuts,
	}
}

func RunOpen() Executer {
	return Builtin{
		Name: "",
		Safe: true,
		Run:  wrapChannelFunc(runOpen),
	}
}

func RunClose() Executer {
	return Builtin{
		Name: "",
		Safe: true,
		Run:  wrapChannelFunc(runClose),
	}
}

func RunEof() Executer {
	return Builtin{
		Name: "eof",
		Safe: true,
		Run:  wrapChannelFunc(runEof),
	}
}

func RunSeek() Executer {
	return Builtin{
		Name: "seek",
		Safe: true,
		Run:  wrapChannelFunc(runSeek),
	}
}

func RunTell() Executer {
	return Builtin{
		Name: "tell",
		Safe: true,
		Run:  wrapChannelFunc(runTell),
	}
}

func RunGets() Executer {
	return Builtin{
		Name: "gets",
		Safe: true,
		Run:  wrapChannelFunc(runGets),
	}
}

func RunRead() Executer {
	return Builtin{
		Name: "read",
		Safe: true,
		Run:  wrapChannelFunc(runRead),
	}
}

func runPuts(i Interpreter, args []env.Value) (env.Value, error) {
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
}

func runOpen(ch ChannelHandler, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runClose(ch ChannelHandler, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runEof(ch ChannelHandler, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runSeek(ch ChannelHandler, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runTell(ch ChannelHandler, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runGets(ch ChannelHandler, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runRead(ch ChannelHandler, args []env.Value) (env.Value, error) {
	return nil, nil
}

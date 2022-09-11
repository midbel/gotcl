package stdlib

import (
	"fmt"
	"strings"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/slices"
)

type PrintHandler interface {
	Print(string, string) error
	Println(string, string) error
}

type ChannelHandler interface {
	Interpreter

	Open(string, string) (string, error)
	Close(string) error
	Eof(string) (bool, error)

	Seek(string, int, int) (int64, error)
	Tell(string) (int64, error)
	Gets(string) (string, error)
	Read(string, int) (string, error)

	PrintHandler
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
		Run: wrapChannelFunc(runPuts),
	}
}

func RunOpen() Executer {
	return Builtin{
		Name:     "open",
		Safe:     true,
		Arity:    1,
		Variadic: true,
		Run:      wrapChannelFunc(runOpen),
	}
}

func RunClose() Executer {
	return Builtin{
		Name:  "close",
		Safe:  true,
		Arity: 1,
		Run:   wrapChannelFunc(runClose),
	}
}

func RunEof() Executer {
	return Builtin{
		Name:  "eof",
		Safe:  true,
		Arity: 1,
		Run:   wrapChannelFunc(runEof),
	}
}

func RunSeek() Executer {
	return Builtin{
		Name:     "seek",
		Safe:     true,
		Arity:    2,
		Variadic: true,
		Run:      wrapChannelFunc(runSeek),
	}
}

func RunTell() Executer {
	return Builtin{
		Name:  "tell",
		Safe:  true,
		Arity: 1,
		Run:   wrapChannelFunc(runTell),
	}
}

func RunGets() Executer {
	return Builtin{
		Name:  "gets",
		Safe:  true,
		Arity: 1,
		Run:   wrapChannelFunc(runGets),
	}
}

func RunRead() Executer {
	return Builtin{
		Name:     "read",
		Safe:     true,
		Arity:    1,
		Variadic: true,
		Options: []Option{
			{
				Name:  "nonewline",
				Flag:  true,
				Value: env.False(),
				Check: CheckBool,
			},
		},
		Run: wrapChannelFunc(runRead),
	}
}

func runPuts(ch ChannelHandler, args []env.Value) (env.Value, error) {
	var (
		nonl, _   = ch.Resolve("nonewline")
		file, err = ch.Resolve("channel")
	)
	if err != nil {
		return nil, err
	}
	ch.Print(file.String(), slices.Fst(args).String())
	if !env.ToBool(nonl) {
		ch.Println(file.String(), "")
	}
	return env.EmptyStr(), nil
}

func runOpen(ch ChannelHandler, args []env.Value) (env.Value, error) {
	var mode string
	if v := slices.Snd(args); v != nil {
		mode = v.String()
	}
	file, err := ch.Open(slices.Fst(args).String(), mode)
	return env.Str(file), err
}

func runClose(ch ChannelHandler, args []env.Value) (env.Value, error) {
	return nil, ch.Close(slices.Fst(args).String())
}

func runEof(ch ChannelHandler, args []env.Value) (env.Value, error) {
	ok, err := ch.Eof(slices.Fst(args).String())
	return env.Bool(ok), err
}

func runSeek(ch ChannelHandler, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runTell(ch ChannelHandler, args []env.Value) (env.Value, error) {
	tell, err := ch.Tell(slices.Fst(args).String())
	return env.Int(tell), err
}

func runGets(ch ChannelHandler, args []env.Value) (env.Value, error) {
	str, err := ch.Gets(slices.Fst(args).String())
	return env.Str(str), err
}

func runRead(ch ChannelHandler, args []env.Value) (env.Value, error) {
	var (
		size    int
		err     error
		nonl, _ = ch.Resolve("nonewline")
	)
	size, err = env.ToInt(slices.Lst(args))
	if err != nil {
		return nil, err
	}
	str, err := ch.Read(slices.Fst(args).String(), size)
	if err == nil && !env.ToBool(nonl) {
		str = strings.TrimSpace(str)
	}
	return env.Str(str), err
}

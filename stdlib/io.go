package stdlib

import (
	"fmt"
	"io"
	"strings"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/gotcl/glob"
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

	Channels() []string

	Seek(string, int, int) (int64, error)
	Tell(string) (int64, error)
	Gets(string) (string, error)
	Read(string, int) (string, error)

	Copy(string, string, int) (int64, error)

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
		Safe: true,
		List: []Executer{
			Builtin{
				Name:  "close",
				Arity: 1,
				Run:   wrapChannelFunc(chanClose),
			},
			Builtin{
				Name:  "copy",
				Arity: 2,
				Options: []Option{
					{
						Name:  "size",
						Value: env.Zero(),
						Check: CheckNumber,
					},
				},
				Run: wrapChannelFunc(chanCopy),
			},
			Builtin{
				Name:  "eof",
				Arity: 1,
				Run:   wrapChannelFunc(chanEof),
			},
			Builtin{
				Name:     "gets",
				Arity:    1,
				Variadic: true,
				Run:      wrapChannelFunc(chanGets),
			},
			Builtin{
				Name:     "names",
				Variadic: true,
				Run:      wrapChannelFunc(chanNames),
			},
			Builtin{
				Name:  "puts",
				Arity: 1,
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
				Run: wrapChannelFunc(chanPuts),
			},
			Builtin{
				Name:     "read",
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
				Run: wrapChannelFunc(chanRead),
			},
			Builtin{
				Name:     "seek",
				Arity:    2,
				Variadic: true,
				Run:      wrapChannelFunc(chanSeek),
			},
			Builtin{
				Name:  "tell",
				Arity: 1,
				Run:   wrapChannelFunc(chanTell),
			},
		},
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
		Run: wrapChannelFunc(chanPuts),
	}
}

func RunOpen() Executer {
	return Builtin{
		Name:     "open",
		Safe:     true,
		Arity:    1,
		Variadic: true,
		Run:      wrapChannelFunc(chanOpen),
	}
}

func RunClose() Executer {
	return Builtin{
		Name:  "close",
		Safe:  true,
		Arity: 1,
		Run:   wrapChannelFunc(chanClose),
	}
}

func RunEof() Executer {
	return Builtin{
		Name:  "eof",
		Safe:  true,
		Arity: 1,
		Run:   wrapChannelFunc(chanEof),
	}
}

func RunSeek() Executer {
	return Builtin{
		Name:     "seek",
		Safe:     true,
		Arity:    2,
		Variadic: true,
		Run:      wrapChannelFunc(chanSeek),
	}
}

func RunTell() Executer {
	return Builtin{
		Name:  "tell",
		Safe:  true,
		Arity: 1,
		Run:   wrapChannelFunc(chanTell),
	}
}

func RunGets() Executer {
	return Builtin{
		Name:     "gets",
		Safe:     true,
		Arity:    1,
		Variadic: true,
		Run:      wrapChannelFunc(chanGets),
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
		Run: wrapChannelFunc(chanRead),
	}
}

func chanCopy(ch ChannelHandler, args []env.Value) (env.Value, error) {
	var (
		size, _ = ch.Resolve("size")
		z, _    = env.ToInt(size)
		in      = slices.Fst(args).String()
		out     = slices.Snd(args).String()
	)
	n, err := ch.Copy(in, out, z)
	return env.Int(n), err
}

func chanNames(ch ChannelHandler, args []env.Value) (env.Value, error) {
	var pat string
	if v := slices.Fst(args); v != nil {
		pat = v.String()
	}
	list := glob.Filter(ch.Channels(), pat)
	return env.ListFromStrings(list), nil
}

func chanPuts(ch ChannelHandler, args []env.Value) (env.Value, error) {
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

func chanOpen(ch ChannelHandler, args []env.Value) (env.Value, error) {
	var mode string
	if v := slices.Snd(args); v != nil {
		mode = v.String()
	}
	file, err := ch.Open(slices.Fst(args).String(), mode)
	return env.Str(file), err
}

func chanClose(ch ChannelHandler, args []env.Value) (env.Value, error) {
	return nil, ch.Close(slices.Fst(args).String())
}

func chanEof(ch ChannelHandler, args []env.Value) (env.Value, error) {
	ok, err := ch.Eof(slices.Fst(args).String())
	return env.Bool(ok), err
}

func chanSeek(ch ChannelHandler, args []env.Value) (env.Value, error) {
	var offset int
	if x := slices.At(args, 2); x != nil {
		o, err := env.ToInt(x)
		if err != nil {
			return nil, err
		}
		offset = o
	}
	var whence int
	switch w := slices.Snd(args).String(); w {
	case "start":
		whence = io.SeekStart
	case "end":
		whence = io.SeekEnd
	case "current":
		whence = io.SeekCurrent
	default:
		return nil, fmt.Errorf("%s: unknown value for whence", w)
	}
	tell, err := ch.Seek(slices.Fst(args).String(), offset, whence)
	return env.Int(tell), err
}

func chanTell(ch ChannelHandler, args []env.Value) (env.Value, error) {
	tell, err := ch.Tell(slices.Fst(args).String())
	return env.Int(tell), err
}

func chanGets(ch ChannelHandler, args []env.Value) (env.Value, error) {
	str, err := ch.Gets(slices.Fst(args).String())
	if err != nil {
		return nil, err
	}
	res := env.Str(str)
	if v := slices.Snd(args); v != nil {
		ch.Define(v.String(), res)
	}
	return res, nil
}

func chanRead(ch ChannelHandler, args []env.Value) (env.Value, error) {
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

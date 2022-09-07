package stdlib

import (
	"fmt"
	"os"
	"sort"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/slices"
)

type CommandHandler interface {
	Depth() int
	Count() int
	Commands(string) []string
}

type commandHandlerFunc func(CommandHandler, []env.Value) (env.Value, error)

func wrapCommandHandler(do commandHandlerFunc) CommandFunc {
	return func(i Interpreter, args []env.Value) (env.Value, error) {
		ch, ok := i.(CommandHandler)
		if !ok {
			return nil, fmt.Errorf("interperter can not introspect Commands")
		}
		return do(ch, args)
	}
}

type ProcHandler interface {
	ProcList(string) []string
	ProcBody(string) (string, error)
	ProcArgs(string) ([]string, error)
	ProcDefault(string, string) (string, bool, error)
}

type procHandlerFunc func(ProcHandler, []env.Value) (env.Value, error)

func wrapProcHandler(do procHandlerFunc) CommandFunc {
	return func(i Interpreter, args []env.Value) (env.Value, error) {
		ph, ok := i.(ProcHandler)
		if !ok {
			return nil, fmt.Errorf("interperter can not introspect procedures")
		}
		return do(ph, args)
	}
}

type VariableHandler interface {
	Globals(string) []string
	Locals(string) []string
	Variables(string) []string
}

type variableHandlerFunc func(VariableHandler, []env.Value) (env.Value, error)

func wrapVariableHandler(do variableHandlerFunc) CommandFunc {
	return func(i Interpreter, args []env.Value) (env.Value, error) {
		vh, ok := i.(VariableHandler)
		if !ok {
			return nil, fmt.Errorf("interperter can not introspect variables")
		}
		return do(vh, args)
	}
}

func MakeInfo() Executer {
	e := Ensemble{
		Name: "info",
		List: []Executer{
			Builtin{
				Name: "commands",
				Run: wrapCommandHandler(func(ch CommandHandler, args []env.Value) (env.Value, error) {
					pat := slices.Fst(args)
					if pat == nil {
						pat = env.EmptyStr()
					}
					return env.ListFromStrings(ch.Commands(pat.String())), nil
				}),
			},
			Builtin{
				Name: "complete",
				Run: func(i Interpreter, args []env.Value) (env.Value, error) {
					return nil, nil
				},
			},
			Builtin{
				Name: "version",
				Run: func(i Interpreter, args []env.Value) (env.Value, error) {
					return nil, nil
				},
			},
			Builtin{
				Name: "hostname",
				Run: func(i Interpreter, args []env.Value) (env.Value, error) {
					str, err := os.Hostname()
					return env.Str(str), err
				},
			},
			Builtin{
				Name: "nameofexecutable",
				Run: func(i Interpreter, args []env.Value) (env.Value, error) {
					str, err := os.Executable()
					return env.Str(str), err
				},
			},
			Builtin{
				Name: "count",
				Run: wrapCommandHandler(func(ch CommandHandler, args []env.Value) (env.Value, error) {
					c := ch.Count()
					return env.Int(int64(c)), nil
				}),
			},
			Builtin{
				Name: "depth",
				Run: wrapCommandHandler(func(ch CommandHandler, args []env.Value) (env.Value, error) {
					d := ch.Depth()
					return env.Int(int64(d)), nil
				}),
			},
			Builtin{
				Name:     "procs",
				Variadic: true,
				Run: wrapProcHandler(func(ph ProcHandler, args []env.Value) (env.Value, error) {
					pat := slices.Fst(args)
					if pat == nil {
						pat = env.EmptyStr()
					}
					list := ph.ProcList(pat.String())
					return env.ListFromStrings(list), nil
				}),
			},
			Builtin{
				Name:  "args",
				Arity: 1,
				Run: wrapProcHandler(func(ph ProcHandler, args []env.Value) (env.Value, error) {
					list, err := ph.ProcArgs(slices.Fst(args).String())
					if err != nil {
						return nil, err
					}
					return env.ListFromStrings(list), nil
				}),
			},
			Builtin{
				Name:  "body",
				Arity: 1,
				Run: wrapProcHandler(func(ph ProcHandler, args []env.Value) (env.Value, error) {
					body, err := ph.ProcBody(slices.Fst(args).String())
					if err != nil {
						return nil, err
					}
					return env.Str(body), nil
				}),
			},
			Builtin{
				Name:  "default",
				Arity: 2,
				Run: wrapProcHandler(func(ph ProcHandler, args []env.Value) (env.Value, error) {
					arg, set, err := ph.ProcDefault(slices.Fst(args).String(), slices.Snd(args).String())
					if err != nil {
						return nil, err
					}
					if !set {
						return nil, fmt.Errorf("%s has not arg named %s", slices.Fst(args), slices.Snd(args))
					}
					return env.Str(arg), nil
				}),
			},
			Builtin{
				Name: "globals",
				Run: func(i Interpreter, args []env.Value) (env.Value, error) {
					return nil, nil
				},
			},
			Builtin{
				Name: "locals",
				Run: func(i Interpreter, args []env.Value) (env.Value, error) {
					return nil, nil
				},
			},
			Builtin{
				Name: "vars",
				Run: func(i Interpreter, args []env.Value) (env.Value, error) {
					return nil, nil
				},
			},
		},
	}
	sort.Slice(e.List, func(i, j int) bool {
		return e.List[i].GetName() < e.List[j].GetName()
	})
	return e
}

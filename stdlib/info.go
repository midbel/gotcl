package stdlib

import (
	"fmt"
	"os"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/slices"
)

type CommandHandler interface {
	Depth() int
	Count() int
	Commands(string) []string
	CurrentFrame(int) (string, []string, error)
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
				Name:  "complete",
				Arity: 1,
				Run:   infoComplete,
			},
			Builtin{
				Name: "tclversion",
				Run:  infoVersion,
			},
			Builtin{
				Name: "hostname",
				Run:  infoHostname,
			},
			Builtin{
				Name: "nameofexecutable",
				Run:  infoExecutable,
			},
			Builtin{
				Name: "commands",
				Run:  wrapCommandHandler(infoCommands),
			},
			Builtin{
				Name: "cmdcount",
				Run:  wrapCommandHandler(infoCommandCount),
			},
			Builtin{
				Name:     "level",
				Variadic: true,
				Run:      wrapCommandHandler(infoCommandLevel),
			},
			Builtin{
				Name:     "procs",
				Variadic: true,
				Run:      wrapProcHandler(infoProcedures),
			},
			Builtin{
				Name:  "args",
				Arity: 1,
				Run:   wrapProcHandler(infoProcedureArgs),
			},
			Builtin{
				Name:  "body",
				Arity: 1,
				Run:   wrapProcHandler(infoProcedureBody),
			},
			Builtin{
				Name:  "default",
				Arity: 2,
				Run:   wrapProcHandler(infoProcedureDefault),
			},
			Builtin{
				Name:  "exists",
				Arity: 1,
				Run:   infoExists,
			},
			Builtin{
				Name:     "globals",
				Variadic: true,
				Run:      wrapVariableHandler(infoGlobals),
			},
			Builtin{
				Name:     "locals",
				Variadic: true,
				Run:      wrapVariableHandler(infoLocals),
			},
			Builtin{
				Name:     "vars",
				Variadic: true,
				Run:      wrapVariableHandler(infoVariables),
			},
		},
	}
	return sortEnsembleCommands(e)
}

func infoVersion(i Interpreter, args []env.Value) (env.Value, error) {
	return env.Str(i.Version()), nil
}

func infoComplete(i Interpreter, args []env.Value) (env.Value, error) {
	ok := i.IsComplete(slices.Fst(args).String())
	return env.Bool(ok), nil
}

func infoHostname(i Interpreter, args []env.Value) (env.Value, error) {
	str, err := os.Hostname()
	return env.Str(str), err
}

func infoExecutable(i Interpreter, args []env.Value) (env.Value, error) {
	str, err := os.Executable()
	return env.Str(str), err
}

func infoExists(i Interpreter, args []env.Value) (env.Value, error) {
	_, err := i.Resolve(slices.Fst(args).String())
	if err == nil {
		return env.True(), nil
	}
	return env.False(), nil
}

func infoCommands(ch CommandHandler, args []env.Value) (env.Value, error) {
	pat := slices.Fst(args)
	if pat == nil {
		pat = env.EmptyStr()
	}
	return env.ListFromStrings(ch.Commands(pat.String())), nil
}

func infoCommandCount(ch CommandHandler, args []env.Value) (env.Value, error) {
	c := ch.Count()
	return env.Int(int64(c)), nil
}

func infoCommandLevel(ch CommandHandler, args []env.Value) (env.Value, error) {
	if len(args) == 0 {
		d := ch.Depth()
		return env.Int(int64(d)), nil
	}
	n, err := env.ToInt(slices.Fst(args))
	if err != nil {
		return nil, err
	}
	cmd, params, err := ch.CurrentFrame(int(n))
	if err != nil {
		return nil, err
	}
	return env.ListFromStrings(slices.Prepend(cmd, params)), nil
}

func infoProcedures(ph ProcHandler, args []env.Value) (env.Value, error) {
	pat := slices.Fst(args)
	if pat == nil {
		pat = env.EmptyStr()
	}
	list := ph.ProcList(pat.String())
	return env.ListFromStrings(list), nil
}

func infoProcedureArgs(ph ProcHandler, args []env.Value) (env.Value, error) {
	list, err := ph.ProcArgs(slices.Fst(args).String())
	if err != nil {
		return nil, err
	}
	return env.ListFromStrings(list), nil
}

func infoProcedureBody(ph ProcHandler, args []env.Value) (env.Value, error) {
	body, err := ph.ProcBody(slices.Fst(args).String())
	if err != nil {
		return nil, err
	}
	return env.Str(body), nil
}

func infoProcedureDefault(ph ProcHandler, args []env.Value) (env.Value, error) {
	arg, set, err := ph.ProcDefault(slices.Fst(args).String(), slices.Snd(args).String())
	if err != nil {
		return nil, err
	}
	if !set {
		return nil, fmt.Errorf("%s has not arg named %s", slices.Fst(args), slices.Snd(args))
	}
	return env.Str(arg), nil
}

func infoGlobals(vh VariableHandler, args []env.Value) (env.Value, error) {
	pat := slices.Fst(args)
	if pat == nil {
		pat = env.EmptyStr()
	}
	list := vh.Globals(pat.String())
	return env.ListFromStrings(list), nil
}

func infoLocals(vh VariableHandler, args []env.Value) (env.Value, error) {
	pat := slices.Fst(args)
	if pat == nil {
		pat = env.EmptyStr()
	}
	list := vh.Locals(pat.String())
	return env.ListFromStrings(list), nil
}

func infoVariables(vh VariableHandler, args []env.Value) (env.Value, error) {
	pat := slices.Fst(args)
	if pat == nil {
		pat = env.EmptyStr()
	}
	var (
		gs = vh.Globals(pat.String())
		ls = vh.Locals(pat.String())
	)
	return env.ListFromStrings(append(gs, ls...)), nil
}

package stdlib

import (
	"os"
	"sort"

	"github.com/midbel/gotcl/env"
)

type CommandHandler interface {
	Depth() int
	Count() int
	List(string) []string
}

type ProcHandler interface {
	Interpreter
	ProcList(string) []string
	ProcBody(string) (string, error)
	ProcArgs(string) ([]string, error)
	ProcDefault(string, string) (string, bool, error)
	RegisterProc(string, string, string) error
}

type VariableHandler interface {
	Globals(string) []string
	Locals(string) []string
	Variables(string) []string
}

func MakeInfo() Executer {
	e := Ensemble{
		Name: "info",
		List: []Executer{
			Builtin{
				Name: "commands",
				Run: func(i Interpreter, args []env.Value) (env.Value, error) {
					return nil, nil
				},
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
				Run: func(i Interpreter, args []env.Value) (env.Value, error) {
					return nil, nil
				},
			},
			Builtin{
				Name: "depth",
				Run: func(i Interpreter, args []env.Value) (env.Value, error) {
					return nil, nil
				},
			},
			Builtin{
				Name: "procs",
				Run: func(i Interpreter, args []env.Value) (env.Value, error) {
					return nil, nil
				},
			},
			Builtin{
				Name: "args",
				Run: func(i Interpreter, args []env.Value) (env.Value, error) {
					return nil, nil
				},
			},
			Builtin{
				Name: "body",
				Run: func(i Interpreter, args []env.Value) (env.Value, error) {
					return nil, nil
				},
			},
			Builtin{
				Name: "default",
				Run: func(i Interpreter, args []env.Value) (env.Value, error) {
					return nil, nil
				},
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

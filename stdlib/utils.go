package stdlib

import (
	"fmt"
	"os"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/slices"
)

func RunTypeOf() Executer {
	return Builtin{
		Name:  "typeof",
		Arity: 1,
		Safe:  true,
		Run:   runTypeof,
	}
}

func RunHelp() Executer {
	return Builtin{
		Name:  "help",
		Help:  "retrieve help of given builtin command",
		Arity: 1,
		Safe:  true,
		Options: []Option{
			{
				Name:  "usage",
				Flag:  true,
				Value: env.False(),
				Check: CheckBool,
			},
		},
		Run: runHelp,
	}
}

func RunChdir() Executer {
	return Builtin{
		Name:     "cd",
		Variadic: true,
		Run:      runChdir,
	}
}

func RunPid() Executer {
	return Builtin{
		Name:     "pid",
		Variadic: true,
		Run:      runPid,
	}
}

func RunPwd() Executer {
	return Builtin{
		Name:     "pwd",
		Variadic: true,
		Run:      runPwd,
	}
}

func runChdir(i Interpreter, args []env.Value) (env.Value, error) {
	dir := slices.Fst(args)
	if dir == nil {
		d, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		dir = env.Str(d)
	}
	return nil, os.Chdir(dir.String())
}

func runPid(i Interpreter, args []env.Value) (env.Value, error) {
	pid := os.Getpid()
	return env.Int(int64(pid)), nil
}

func runPwd(i Interpreter, args []env.Value) (env.Value, error) {
	var (
		pwd, err = os.Getwd()
		res      env.Value
	)
	if err == nil {
		res = env.Str(pwd)
	}
	return res, nil
}

func runTypeof(i Interpreter, args []env.Value) (env.Value, error) {
	typ := fmt.Sprintf("%T", slices.Fst(args))
	return env.Str(typ), nil
}

func runHelp(i Interpreter, args []env.Value) (env.Value, error) {
	h, ok := i.(interface{ GetHelp(string) (string, error) })
	if !ok {
		return nil, fmt.Errorf("interpreter can not extract help from command")
	}
	help, err := h.GetHelp(slices.Fst(args).String())
	if err != nil {
		return nil, err
	}
	return env.Str(help), nil
}

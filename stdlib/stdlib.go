package stdlib

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/slices"
)

var (
	ErrArgument    = errors.New("wrong number of arguments given")
	ErrUsage       = errors.New("bad usage")
	ErrUnknown     = errors.New("unknown command")
	ErrImplemented = errors.New("command not implemented")
	ErrExit        = errors.New("exit")
	ErrIndex       = errors.New("index out of range")
)

type Executer interface {
	Execute(Interpreter, []string) (string, error)
}

type CommandFunc func(Interpreter, []string) (string, error)

type Interpreter interface {
	Version() string

	env.Environment
	Do(string, func(string) (string, error)) (string, error)
	Execute(io.Reader) (string, error)
	ExecuteUp(io.Reader, int) (string, error)

	Split(string) ([]string, error)

	Print(string, string) error
	Println(string, string) error

	RegisterProc(string, string, string) error
	UnregisterProc(string) error
	Rename(string, string)

	RegisterNS(string, string) error
	UnregisterNS(string) error
	RegisterVar(string, string) error
	ResolveVar(string) (string, error)
	CurrentNS() string
	ParentNS(string) (string, error)
	ChildrenNS(string, string) ([]string, error)
	ExistsNS(string) bool
}

func makeEnsemble(name string, set map[string]CommandFunc) CommandFunc {
	return func(i Interpreter, args []string) (string, error) {
		exec, ok := set[slices.Fst(args)]
		if !ok {
			return "", fmt.Errorf("%s %s: %w", name, slices.Fst(args), ErrUnknown)
		}
		if exec == nil {
			return "", fmt.Errorf("%s: %w", slices.Fst(args), ErrImplemented)
		}
		return exec(i, slices.Rest(args))
	}
}

func executeBool(i Interpreter, str string) (bool, error) {
	if str == "" {
		return false, nil
	}
	res, err := i.Execute(strings.NewReader(str))
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(res)
}

type setter interface {
	Set(string) error
}

type argFunc func(*flag.FlagSet) (int, bool)

func parseArgs(name string, args []string, fn argFunc) ([]string, error) {
	if fn == nil {
		fn = func(_ *flag.FlagSet) (int, bool) { return 0, false }
	}
	var (
		set     = flag.NewFlagSet(name, flag.ContinueOnError)
		num, eq = fn(set)
	)
	err := set.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, ErrUsage)
	}
	args = set.Args()
	if (eq && len(args) != num) || (!eq && len(args) < num) {
		err = fmt.Errorf("%s: %w", name, ErrArgument)
	}
	return args, err
}

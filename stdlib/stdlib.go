package stdlib

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/gotcl/stdlib/argparse"
	"github.com/midbel/slices"
)

var ErrArgument = errors.New("wrong number of argument given")

type Interpreter interface {
	Execute(io.Reader) (env.Value, error)
	IsSafe() bool

	Resolve(string) (env.Value, error)
	Define(string, env.Value)
	Delete(string)
}

type CommandFunc func(Interpreter, []env.Value) (env.Value, error)

type Executer interface {
	GetName() string
	IsSafe() bool
	Execute(Interpreter, []env.Value) (env.Value, error)
}

type Ensemble struct {
	Name  string
	Usage string
	Help  string
	Safe  bool
	List  []Executer
}

func (e Ensemble) IsSafe() bool {
	return e.Safe
}

func (e Ensemble) GetName() string {
	return e.Name
}

func (e Ensemble) Execute(i Interpreter, args []env.Value) (env.Value, error) {
	name := slices.Fst(args).String()
	x := sort.Search(len(e.List), func(i int) bool {
		return e.List[i].GetName() >= name
	})
	if x >= len(e.List) || e.List[x].GetName() != name {
		return nil, fmt.Errorf("%s %s: command not defined", e.Name, name)
	}
	return e.List[x].Execute(i, slices.Rest(args))
}

type Builtin struct {
	Name     string
	Usage    string
	Help     string
	Safe     bool
	Arity    int
	Variadic bool
	Run      CommandFunc
	Options  []argparse.Option
}

func (b Builtin) IsSafe() bool {
	return b.Safe
}

func (b Builtin) GetName() string {
	return b.Name
}

func (b Builtin) Execute(i Interpreter, args []env.Value) (env.Value, error) {
	if b.Run == nil {
		return nil, fmt.Errorf("%s: command is not runnable", b.Name)
	}
	args, err := b.parseOptions(i, args)
	if err != nil {
		return nil, err
	}
	if err := b.parseArgs(args); err != nil {
		return nil, err
	}
	return b.Run(i, args)
}

func (b Builtin) parseArgs(args []env.Value) error {
	if n := len(args); n != b.Arity {
		if !b.Variadic || (b.Variadic && n < b.Arity) {
			return fmt.Errorf("%s: %w: want %d, got %d", b.Name, ErrArgument, b.Arity, n)
		}
	}
	return nil
}

func (b Builtin) parseOptions(i Interpreter, args []env.Value) ([]env.Value, error) {
	if len(b.Options) == 0 {
		return args, nil
	}
	sort.Slice(b.Options, func(i, j int) bool {
		return b.Options[i].Name < b.Options[j].Name
	})
	for _, o := range b.Options {
		if o.Value == nil {
			continue
		}
		i.Define(o.Name, o.Value)
	}
	var j int
	for ; j < len(args) && j < len(b.Options); j++ {
		str := args[j].String()
		if str == "--" || !strings.HasPrefix(str, "-") {
			break
		}
		str = strings.TrimPrefix(str, "-")
		x, err := argparse.IsSet(b.Options, str)
		if err != nil {
			return nil, err
		}
		if b.Options[x].Flag {
			i.Define(b.Options[x].Name, env.True())
			continue
		}
		if check := b.Options[x].Check; check != nil {
			if err := check(args[j+1]); err != nil {
				return nil, err
			}
			b.Options[j].Value = args[j+1]
			i.Define(str, b.Options[j].Value)
		}
		j++
	}
	if err := argparse.IsValid(b.Options); err != nil {
		return nil, err
	}
	return args[j:], nil
}

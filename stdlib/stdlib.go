package stdlib

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/slices"
)

var (
	ErrArgument = errors.New("wrong number of argument given")
	ErrExit     = errors.New("exit")
	ErrReturn   = errors.New("return")
	ErrBreak    = errors.New("break")
	ErrContinue = errors.New("continue")
)

type Interpreter interface {
	Execute(io.Reader) (env.Value, error)
	IsComplete(string) bool
	IsSafe() bool
	Version() string

	Resolve(string) (env.Value, error)
	Define(string, env.Value)
	Delete(string)
}

const (
	ErrorOk int = iota
	ErrorErr
	ErrorRet
	ErrorBreak
	ErrorContinue
)

type Error struct {
	Err   error
	Code  int
	Level int
}

func ErrorWithCode(msg string, code int) error {
	return Error{
		Err:  errors.New(msg),
		Code: code,
	}
}

func ErrorFromError(err error) error {
	return Error{
		Err:  err,
		Code: ErrorErr,
	}
}

func DefaultError(msg string) error {
	return ErrorWithCode(msg, ErrorErr)
}

func (e Error) Error() string {
	return e.Err.Error()
}

func (e Error) Unwrap() error {
	return e.Err
}

type CommandFunc func(Interpreter, []env.Value) (env.Value, error)

type Executer interface {
	GetName() string
	IsSafe() bool
	Scoped() bool
	Execute(Interpreter, []env.Value) (env.Value, error)
}

type Ensemble struct {
	Name  string
	Usage string
	Help  string
	Safe  bool
	List  []Executer
}

func (_ Ensemble) Scoped() bool {
	return false
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

func sortEnsembleCommands(e Ensemble) Ensemble {
	sort.Slice(e.List, func(i, j int) bool {
		return e.List[i].GetName() < e.List[j].GetName()
	})
	return e
}

type Builtin struct {
	Name     string
	Usage    string
	Help     string
	Safe     bool
	Arity    int
	Variadic bool
	Scope    bool
	Run      CommandFunc
	Options  []Option
}

func (b Builtin) Scoped() bool {
	return b.Scope
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
		if _, ok := args[j].(env.String); !ok || str == "--" || !strings.HasPrefix(str, "-") {
			if str == "--" {
				j++
			}
			break
		}
		str = strings.TrimPrefix(str, "-")
		x, err := IsSet(b.Options, str)
		if err != nil {
			return nil, err
		}
		if b.Options[x].Flag {
			i.Define(b.Options[x].Name, env.True())
			continue
		}
		if check := b.Options[x].Check; check != nil {
			val, err := check(args[j+1])
			if err != nil {
				return nil, err
			}
			b.Options[j].Value = val
			i.Define(str, b.Options[j].Value)
		}
		j++
	}
	if err := IsValid(b.Options); err != nil {
		return nil, err
	}
	return args[j:], nil
}

var ErrType = errors.New("invalid type given")

type CheckFunc func(env.Value) (env.Value, error)

type Option struct {
	env.Value
	Name     string
	Help     string
	Flag     bool
	Required bool
	Check    CheckFunc
}

func CheckBool(v env.Value) (env.Value, error) {
	_, ok := v.(env.Boolean)
	if !ok {
		return nil, ErrType
	}
	return v, nil
}

func CheckNumber(v env.Value) (env.Value, error) {
	_, ok := v.(env.Number)
	if !ok {
		return nil, ErrType
	}
	return v, nil
}

func CheckString(v env.Value) (env.Value, error) {
	_, ok := v.(env.String)
	if !ok {
		return nil, ErrType
	}
	return v, nil
}

func OneOf(cs ...CheckFunc) CheckFunc {
	return func(v env.Value) (env.Value, error) {
		var (
			val env.Value
			err error
		)
		for i := range cs {
			if val, err = cs[i](v); err == nil {
				break
			}
		}
		return val, err
	}
}

func CombineCheck(cs ...CheckFunc) CheckFunc {
	return func(v env.Value) (env.Value, error) {
		var (
			val env.Value
			err error
		)
		for i := range cs {
			if val, err = cs[i](v); err != nil {
				return nil, err
			}
		}
		return val, nil
	}
}

func IsValid(list []Option) error {
	ok := slices.Every(list, func(o Option) bool {
		if !o.Required {
			return true
		}
		return o.Required && o.Value != nil
	})
	if !ok {
		return fmt.Errorf("required options are not provided!")
	}
	return nil
}

func IsSet(list []Option, name string) (int, error) {
	x := sort.Search(len(list), func(i int) bool {
		return list[i].Name >= name
	})
	if x < len(list) && list[x].Name == name {
		return x, nil
	}
	return 0, fmt.Errorf("%s: option not supported", name)
}

func testScript(i Interpreter, v env.Value) (bool, error) {
	v, err := i.Execute(strings.NewReader(v.String()))
	if err != nil {
		return false, err
	}
	return env.ToBool(v), nil
}

func hasError(es ...error) error {
	for i := range es {
		if es[i] != nil {
			return es[i]
		}
	}
	return nil
}

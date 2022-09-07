package interp

import (
	"errors"
	"fmt"
	"strings"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/gotcl/stdlib"
	"github.com/midbel/gotcl/word"
	"github.com/midbel/slices"
)

var ErrSyntax = errors.New("syntax error")

type CommandSet map[string]stdlib.Executer

func EmptySet() CommandSet {
	return make(CommandSet)
}

func DefaultSet() CommandSet {
	set := EmptySet()
	set.registerCmd("puts", stdlib.RunPuts())
	set.registerCmd("set", stdlib.RunSet())
	set.registerCmd("unset", stdlib.RunUnset())
	set.registerCmd("list", stdlib.RunList())
	set.registerCmd("llength", stdlib.RunListLen())
	set.registerCmd("proc", stdlib.RunProc())
	set.registerCmd("string", stdlib.MakeString())
	set.registerCmd("interp", stdlib.MakeInterp())
	set.registerCmd("eval", stdlib.RunEval())
	set.registerCmd("upvar", stdlib.RunUpvar())
	set.registerCmd("uplevel", stdlib.RunUplevel())
	set.registerCmd("incr", stdlib.RunIncr())
	set.registerCmd("namespace", stdlib.MakeNamespace())
	set.registerCmd("parray", stdlib.PrintArray())
	set.registerCmd("array", stdlib.MakeArray())
	return set
}

func UtilSet() CommandSet {
	set := EmptySet()
	set.registerCmd("defer", stdlib.RunDefer())
	set.registerCmd("typeof", stdlib.RunTypeOf())
	set.registerCmd("help", stdlib.RunHelp())
	return set
}

func (cs CommandSet) registerCmd(name string, exec stdlib.Executer) {
	cs[name] = exec
}

type procedure struct {
	Name     string
	Body     string
	Args     []argument
	Variadic bool
}

func createProcedure(name, body, args string) (stdlib.Executer, error) {
	p := procedure{
		Name: name,
		Body: strings.TrimSpace(body),
	}
	args = strings.TrimSpace(args)
	if len(args) != 0 {
		as, err := parseArguments(args)
		if err != nil {
			return nil, err
		}
		p.Args = as
		if a := slices.Lst(p.Args); a.Name == "args" {
			p.Variadic = true
		}
	}
	return p, nil
}

func (p procedure) GetName() string {
	return p.Name
}

func (_ procedure) IsSafe() bool {
	return true
}

func (p procedure) Execute(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
	for j, a := range p.Args {
		if j < len(args) {
			a.Default = args[j]
		}
		i.Define(a.Name, a.Default)
	}
	return i.Execute(strings.NewReader(p.Body))
}

type argument struct {
	Name    string
	Default env.Value
}

func createArg(n string, v env.Value) argument {
	return argument{
		Name:    n,
		Default: v,
	}
}

func parseArguments(str string) ([]argument, error) {
	argWithDefault := func(str string) ([]string, error) {
		scan, err := word.Scan(strings.NewReader(str))
		if err != nil {
			return nil, err
		}
		scan.KeepBlanks(false)
		var words []string
		for {
			w := scan.Scan()
			if w.Type == word.EOF {
				break
			}
			if w.Type != word.Literal && w.Type != word.Block {
				return nil, ErrSyntax
			}
			words = append(words, w.Literal)
		}
		if len(words) != 2 {
			return nil, ErrSyntax
		}
		return words, nil
	}
	var list []argument

	scan, err := word.Scan(strings.NewReader(str))
	if err != nil {
		return nil, err
	}
	scan.KeepBlanks(false)
	var (
		seen  = make(map[string]struct{})
		dummy = struct{}{}
	)
	for {
		w := scan.Scan()
		if w.Type == word.EOF {
			break
		}
		if w.Type == word.Illegal {
			return nil, ErrSyntax
		}
		var a argument
		switch w.Type {
		case word.Literal:
			a = createArg(w.Literal, nil)
		case word.Block:
			ws, err := argWithDefault(w.Literal)
			if err != nil {
				return nil, err
			}
			a = createArg(ws[0], env.Str(ws[1]))
		default:
			return nil, ErrSyntax
		}
		if _, ok := seen[a.Name]; ok {
			return nil, fmt.Errorf("%s: duplicate argument", a.Name)
		}
		seen[a.Name] = dummy
		list = append(list, a)
	}
	return list, nil
}

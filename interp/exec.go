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
	set.registerCmd("proc", stdlib.RunProc())
	set.registerCmd("string", stdlib.MakeString())
	set.registerCmd("interp", stdlib.MakeInterp())
	set.registerCmd("eval", stdlib.RunEval())
	set.registerCmd("upvar", stdlib.RunUpvar())
	set.registerCmd("uplevel", stdlib.RunUplevel())
	set.registerCmd("incr", stdlib.RunIncr())
	set.registerCmd("decr", stdlib.RunDecr())
	set.registerCmd("namespace", stdlib.MakeNamespace())
	set.registerCmd("variable", stdlib.RunVariable())
	set.registerCmd("parray", stdlib.PrintArray())
	set.registerCmd("array", stdlib.MakeArray())
	set.registerCmd("info", stdlib.MakeInfo())
	set.registerCmd("clock", stdlib.MakeClock())
	set.registerCmd("append", stdlib.RunAppend())
	set.registerCmd("rename", stdlib.RunRename())
	set.registerCmd("global", stdlib.RunGlobal())
	set.registerCmd("time", stdlib.RunTime())
	set.registerCmd("exit", stdlib.RunExit())
	set.registerCmd("cd", stdlib.RunChdir())
	set.registerCmd("pid", stdlib.RunPid())
	set.registerCmd("pwd", stdlib.RunPwd())
	set.registerCmd("try", stdlib.RunTry())
	set.registerCmd("throw", stdlib.RunThrow())
	set.registerCmd("error", stdlib.RunError())
	set.registerCmd("catch", stdlib.RunCatch())
	set.registerCmd("if", stdlib.RunIf())
	set.registerCmd("switch", stdlib.RunSwitch())
	set.registerCmd("for", stdlib.RunFor())
	set.registerCmd("while", stdlib.RunWhile())
	set.registerCmd("break", stdlib.RunBreak())
	set.registerCmd("continue", stdlib.RunContinue())
	set.registerCmd("expr", stdlib.RunExpr())
	set.registerCmd("unknown", stdlib.RunUnknown())
	set.registerCmd("return", stdlib.RunReturn())
	set.registerCmd("open", stdlib.RunOpen())
	set.registerCmd("close", stdlib.RunClose())
	set.registerCmd("eof", stdlib.RunEof())
	set.registerCmd("seek", stdlib.RunSeek())
	set.registerCmd("tell", stdlib.RunTell())
	set.registerCmd("gets", stdlib.RunGets())
	set.registerCmd("read", stdlib.RunRead())
	set.registerCmd("chan", stdlib.MakeChan())
	set.registerCmd("file", stdlib.MakeFile())
	set.registerCmd("list", stdlib.RunList())
	set.registerCmd("split", stdlib.RunSplit())
	set.registerCmd("llength", stdlib.RunLLength())
	set.registerCmd("lset", stdlib.RunLSet())
	set.registerCmd("lsort", stdlib.RunLSort())
	set.registerCmd("lsearch", stdlib.RunLSearch())
	set.registerCmd("lreverse", stdlib.RunLReverse())
	set.registerCmd("lreplace", stdlib.RunLReplace())
	set.registerCmd("lrepeat", stdlib.RunLRepeat())
	set.registerCmd("lindex", stdlib.RunLIndex())
	set.registerCmd("lmap", stdlib.RunLMap())
	set.registerCmd("lrange", stdlib.RunLRange())
	set.registerCmd("lassign", stdlib.RunLAssign())
	set.registerCmd("lappend", stdlib.RunLAppend())
	set.registerCmd("linsert", stdlib.RunLInsert())
	return set
}

func UtilSet() CommandSet {
	set := EmptySet()
	set.registerCmd("defer", stdlib.RunDefer())
	set.registerCmd("typeof", stdlib.RunTypeOf())
	set.registerCmd("help", stdlib.RunHelp())
	return set
}

func MathfuncSet() CommandSet {
	set := EmptySet()
	set.registerCmd("abs", stdlib.RunAbs())
	set.registerCmd("acos", stdlib.RunAcos())
	set.registerCmd("asin", stdlib.RunAsin())
	set.registerCmd("atan", stdlib.RunAtan())
	set.registerCmd("atan2", stdlib.RunAtan2())
	set.registerCmd("cos", stdlib.RunCos())
	set.registerCmd("cosh", stdlib.RunCosh())
	set.registerCmd("sin", stdlib.RunSin())
	set.registerCmd("sinh", stdlib.RunSinh())
	set.registerCmd("tan", stdlib.RunTan())
	set.registerCmd("tanh", stdlib.RunTanh())
	set.registerCmd("hypot", stdlib.RunHypot())
	set.registerCmd("bool", stdlib.RunBool())
	set.registerCmd("double", stdlib.RunDouble())
	set.registerCmd("entier", stdlib.RunEntier())
	set.registerCmd("ceil", stdlib.RunCeil())
	set.registerCmd("floor", stdlib.RunFloor())
	set.registerCmd("round", stdlib.RunRound())
	set.registerCmd("fmod", stdlib.RunFmod())
	set.registerCmd("int", stdlib.RunInt())
	set.registerCmd("exp", stdlib.RunExp())
	set.registerCmd("log", stdlib.RunLog())
	set.registerCmd("log10", stdlib.RunLog10())
	set.registerCmd("max", stdlib.RunMax())
	set.registerCmd("min", stdlib.RunMin())
	set.registerCmd("pow", stdlib.RunRaise())
	set.registerCmd("rand", stdlib.RunRand())
	set.registerCmd("srand", stdlib.RunSrand())
	set.registerCmd("isqrt", stdlib.RunIsqrt())
	set.registerCmd("sqrt", stdlib.RunSqrt())
	set.registerCmd("wide", stdlib.RunWide())
	set.registerCmd("rad", stdlib.RunRadian())
	set.registerCmd("deg", stdlib.RunDegree())
	return set
}

func MathopSet() CommandSet {
	set := EmptySet()
	set.registerCmd("+", stdlib.RunAdd())
	set.registerCmd("-", stdlib.RunSub())
	set.registerCmd("*", stdlib.RunMul())
	set.registerCmd("/", stdlib.RunDiv())
	set.registerCmd("%", stdlib.RunMod())
	set.registerCmd("**", stdlib.RunPow())
	set.registerCmd("==", stdlib.RunEq())
	set.registerCmd("!=", stdlib.RunNe())
	set.registerCmd("<", stdlib.RunLt())
	set.registerCmd("<=", stdlib.RunLe())
	set.registerCmd(">", stdlib.RunGt())
	set.registerCmd(">=", stdlib.RunGe())
	set.registerCmd("!", stdlib.RunNot())
	return set
}

func (cs CommandSet) Rename(prev, next string) error {
	e, ok := cs[prev]
	if !ok {
		return fmt.Errorf("%s: command not defined", prev)
	}
	delete(cs, prev)
	if next != "" {
		cs[next] = e
	}
	return nil
}

func (cs CommandSet) Procedures() []string {
	var list []string
	for k, e := range cs {
		if _, ok := e.(procedure); !ok {
			continue
		}
		list = append(list, k)
	}
	return list
}

func (cs CommandSet) Args(proc string) ([]string, error) {
	e, ok := cs[proc]
	if !ok {
		return nil, fmt.Errorf("%s: undefined", proc)
	}
	p, ok := e.(procedure)
	if !ok {
		return nil, fmt.Errorf("%s is not a procedure", proc)
	}
	var list []string
	for _, a := range p.Args {
		list = append(list, a.Name)
	}
	return list, nil
}

func (cs CommandSet) Body(proc string) (string, error) {
	e, ok := cs[proc]
	if !ok {
		return "", fmt.Errorf("%s: undefined", proc)
	}
	p, ok := e.(procedure)
	if !ok {
		return "", fmt.Errorf("%s is not a procedure", proc)
	}
	return p.Body, nil
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

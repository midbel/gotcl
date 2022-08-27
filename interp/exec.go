package interp

import (
	"errors"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/midbel/gotcl/stdlib"
	"github.com/midbel/gotcl/glob"
)

var ErrLookup = errors.New("procedure not defined")

const MaxArity = 255

type CommandSet map[string]stdlib.Executer

func EmptySet() CommandSet {
	return make(CommandSet)
}

func MathfuncSet() CommandSet {
	set := EmptySet()

	set.registerCmd("abs", stdlib.RunAbs)
	set.registerCmd("acos", stdlib.RunAcos)
	set.registerCmd("asin", stdlib.RunAsin)
	set.registerCmd("atan", stdlib.RunAtan)
	set.registerCmd("atan2", stdlib.RunAtan2)
	set.registerCmd("cos", stdlib.RunCos)
	set.registerCmd("cosh", stdlib.RunCosh)
	set.registerCmd("sin", stdlib.RunSin)
	set.registerCmd("sinh", stdlib.RunSinh)
	set.registerCmd("tan", stdlib.RunTan)
	set.registerCmd("tanh", stdlib.RunTanh)
	set.registerCmd("hypot", stdlib.RunHypot)
	set.registerCmd("bool", stdlib.RunBool)
	set.registerCmd("double", stdlib.RunDouble)
	set.registerCmd("entier", stdlib.RunEntier)
	set.registerCmd("ceil", stdlib.RunCeil)
	set.registerCmd("floor", stdlib.RunFloor)
	set.registerCmd("round", stdlib.RunRound)
	set.registerCmd("fmod", stdlib.RunFmod)
	set.registerCmd("int", stdlib.RunInt)
	set.registerCmd("exp", stdlib.RunExp)
	set.registerCmd("log", stdlib.RunLog)
	set.registerCmd("log10", stdlib.RunLog10)
	set.registerCmd("max", stdlib.RunMax)
	set.registerCmd("min", stdlib.RunMin)
	set.registerCmd("pow", stdlib.RunRaise)
	set.registerCmd("rand", stdlib.RunRand)
	set.registerCmd("srand", stdlib.RunSrand)
	set.registerCmd("isqrt", stdlib.RunIsqrt)
	set.registerCmd("sqrt", stdlib.RunSqrt)

	return set
}

func MathopSet() CommandSet {
	set := EmptySet()
	set.registerCmd("!", stdlib.RunNot)
	set.registerCmd("+", stdlib.RunAdd)
	set.registerCmd("-", stdlib.RunSub)
	set.registerCmd("*", stdlib.RunMul)
	set.registerCmd("**", stdlib.RunPow)
	set.registerCmd("/", stdlib.RunDiv)
	set.registerCmd("%", stdlib.RunMod)
	set.registerCmd("==", stdlib.RunEq)
	set.registerCmd("!=", stdlib.RunNe)
	set.registerCmd("<", stdlib.RunLt)
	set.registerCmd("<=", stdlib.RunLe)
	set.registerCmd(">", stdlib.RunGt)
	set.registerCmd(">=", stdlib.RunGe)
	return set
}

func DefaultSet() CommandSet {
	set := EmptySet()

	set.registerCmd("set", stdlib.RunSet)
	set.registerCmd("unset", stdlib.RunUnset)
	set.registerCmd("global", stdlib.RunGlobal)
	set.registerCmd("upvar", stdlib.RunUpvar)
	set.registerCmd("uplevel", stdlib.RunUpLevel)
	set.registerCmd("append", stdlib.RunAppend)
	set.registerCmd("rename", stdlib.RunRename)
	set.registerCmd("info", stdlib.RunInfos())
	set.registerCmd("string", stdlib.RunString())
	set.registerCmd("clock", stdlib.RunClock())
	set.registerCmd("namespace", stdlib.RunNamespace())
	set.registerCmd("variable", stdlib.RunVariable)
	set.registerCmd("file", stdlib.RunFile())
	set.registerCmd("exit", stdlib.RunExit)
	set.registerCmd("time", stdlib.RunTime)
	set.registerCmd("proc", stdlib.RunProc)
	set.registerCmd("expr", stdlib.RunExpr)
	set.registerCmd("incr", stdlib.RunIncr)
	set.registerCmd("decr", stdlib.RunDecr)
	set.registerCmd("puts", stdlib.RunPuts)
	set.registerCmd("cd", stdlib.RunChdir)
	set.registerCmd("pid", stdlib.RunPid)
	set.registerCmd("pwd", stdlib.RunPwd)
	set.registerCmd("for", stdlib.RunFor)
	set.registerCmd("while", stdlib.RunWhile)
	set.registerCmd("if", stdlib.RunIf)
	set.registerCmd("switch", stdlib.RunSwitch)
	set.registerCmd("break", stdlib.RunBreak)
	set.registerCmd("continue", stdlib.RunContinue)
	set.registerCmd("open", stdlib.RunOpen)
	set.registerCmd("close", stdlib.RunClose)
	set.registerCmd("eof", stdlib.RunEof)
	set.registerCmd("seek", stdlib.RunSeek)
	set.registerCmd("tell", stdlib.RunTell)
	set.registerCmd("gets", stdlib.RunGets)
	set.registerCmd("read", stdlib.RunRead)
	set.registerCmd("unknown", stdlib.RunUnknown)
	set.registerCmd("try", stdlib.RunTry)
	set.registerCmd("throw", stdlib.RunThrow)
	set.registerCmd("catch", stdlib.RunCatch)
	set.registerCmd("error", stdlib.RunError)

	return set
}

func (cs CommandSet) Lookup(name string) (stdlib.Executer, error) {
	exec, ok := cs[name]
	if !ok {
		return nil, fmt.Errorf("%s: undefined", name)
	}
	return exec, nil
}

func (cs CommandSet) CmdList(pat string) []string {
	return cs.getList(pat, func(e stdlib.Executer) bool {
		_, ok := e.(procedure)
		return !ok
	})
}

func (cs CommandSet) ProcList(pat string) []string {
	return cs.getList(pat, func(e stdlib.Executer) bool {
		_, ok := e.(procedure)
		return ok
	})
}

func (cs CommandSet) ProcArgs(proc string) ([]string, error) {
	p, err := cs.find(proc)
	if err != nil {
		return nil, err
	}
	var args []string
	for _, a := range p.Args {
		args = append(args, a.Name)
	}
	return args, nil
}

func (cs CommandSet) ProcBody(proc string) (string, error) {
	p, err := cs.find(proc)
	return p.Body, err
}

func (cs CommandSet) ProcDefault(proc, arg string) (string, bool, error) {
	p, err := cs.find(proc)
	if err != nil {
		return "", false, err
	}
	args := make([]argument, len(p.Args))
	copy(args, p.Args)
	sort.Slice(args, func(i, j int) bool {
		return args[i].Name < args[j].Name
	})
	x := sort.Search(len(args), func(i int) bool {
		return args[i].Name >= arg
	})
	if x < len(args) && args[x].Name == arg {
		return args[x].Default, args[x].Default != "", nil
	}
	return "", false, fmt.Errorf("%s: argument not defined", arg)
}

func (cs CommandSet) Rename(prev, next string) {
	cmd, ok := cs[prev]
	if !ok {
		return
	}
	delete(cs, prev)
	if next == "" {
		return
	}
	cs[next] = cmd
}

func (cs CommandSet) UnregisterProc(name string) error {
	delete(cs, name)
	return nil
}

func (cs CommandSet) RegisterProc(name, args, body string) error {
	p, err := createProcedure(args, body)
	if err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	cs.registerFunc(name, p)
	return nil
}

func (cs CommandSet) registerCmd(name string, cmd stdlib.CommandFunc) {
	cs.registerFunc(name, createExecuter(cmd))
}

func (cs CommandSet) registerFunc(name string, cmd stdlib.Executer) {
	cs[name] = cmd
}

func (cs CommandSet) find(proc string) (procedure, error) {
	var (
		p  procedure
		ok bool
	)
	e, ok := cs[proc]
	if !ok {
		return p, fmt.Errorf("%s: not defined", proc)
	}
	if p, ok = e.(procedure); !ok {
		return p, fmt.Errorf("%s: not defined with proc command", proc)
	}
	return p, nil
}

func (cs CommandSet) getList(pat string, keep func(e stdlib.Executer) bool) []string {
	list := make([]string, 0, len(cs))
	for k, e := range cs {
		if !keep(e) {
			continue
		}
		list = append(list, k)
	}
	return glob.Filter(list, pat)
}

type argument struct {
	Name    string
	Default string
}

type executer struct {
	stdlib.CommandFunc
}

func createExecuter(cmd stdlib.CommandFunc) stdlib.Executer {
	return executer{
		CommandFunc: cmd,
	}
}

func (e executer) Execute(i stdlib.Interpreter, args []string) (string, error) {
	return e.CommandFunc(i, args)
}

type procedure struct {
	origin   string
	Body     string
	variadic bool
	Args     []argument
}

func createProcedure(args, body string) (procedure, error) {
	proc := procedure{
		Body:     body,
		variadic: false,
	}
	args = strings.TrimSpace(args)
	if len(args) > 0 {
		for {
			var a argument
			args, a.Name, a.Default = splitArg(args)
			proc.Args = append(proc.Args, a)
			if len(proc.Args) >= MaxArity {
				return proc, fmt.Errorf("too many arguments given (max: %d)!", MaxArity)
			}
			if a.Name == "" && a.Default == "" {
				return proc, fmt.Errorf("syntax error")
			}
			if args == "" {
				break
			}
		}
	}
	return proc, nil
}

func (p procedure) Execute(i stdlib.Interpreter, args []string) (string, error) {
	for j, a := range p.Args {
		if j < len(args) {
			a.Default = args[j]
		}
		i.Define(a.Name, a.Default)
	}
	if len(args) > len(p.Args) && !p.variadic {
		return "", stdlib.ErrArgument
	}
	return i.Execute(strings.NewReader(p.Body))
}

func splitArg(str string) (string, string, string) {
	var n, d string
	if strings.HasPrefix(str, "{") {
		tmp, rest, ok := strings.Cut(str[1:], "}")
		if !ok {
			return "", "", ""
		}
		parts := strings.SplitN(tmp, " ", 2)
		n, d = parts[0], parts[1]
		str = strings.TrimSpace(rest)
	} else {
		tmp, rest, _ := strings.Cut(str, " ")
		n = tmp
		str = strings.TrimSpace(rest)
	}
	return str, n, d
}

func unknownDefault(_ stdlib.Interpreter, args []string) (string, error) {
	res, err := exec.Command(args[0], args[1:]...).Output()
	return string(res), err
}

// type frame struct {
// 	*Namespace
// 	env env.Environment
// }

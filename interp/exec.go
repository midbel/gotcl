package interp

import (
	"fmt"
	"sort"
	"strings"

	"github.com/midbel/gotcl/stdlib"
)

type CommandSet map[string]stdlib.Executer

func EmptySet() CommandSet {
	return make(CommandSet)
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
	set.registerCmd("::tcl::mathop::!", stdlib.RunNot)
	set.registerCmd("::tcl::mathop::+", stdlib.RunAdd)
	set.registerCmd("::tcl::mathop::-", stdlib.RunSub)
	set.registerCmd("::tcl::mathop::*", stdlib.RunMul)
	set.registerCmd("::tcl::mathop::**", stdlib.RunPow)
	set.registerCmd("::tcl::mathop::/", stdlib.RunDiv)
	set.registerCmd("::tcl::mathop::%", stdlib.RunMod)
	set.registerCmd("::tcl::mathop::==", stdlib.RunEq)
	set.registerCmd("::tcl::mathop::!=", stdlib.RunNe)
	set.registerCmd("::tcl::mathop::<", stdlib.RunLt)
	set.registerCmd("::tcl::mathop::<=", stdlib.RunLe)
	set.registerCmd("::tcl::mathop::>", stdlib.RunGt)
	set.registerCmd("::tcl::mathop::>=", stdlib.RunGe)
	set.registerCmd("open", stdlib.RunOpen)
	set.registerCmd("close", stdlib.RunClose)
	set.registerCmd("eof", stdlib.RunEof)
	set.registerCmd("seek", stdlib.RunSeek)
	set.registerCmd("tell", stdlib.RunTell)
	set.registerCmd("gets", stdlib.RunGets)
	set.registerCmd("read", stdlib.RunRead)

	return set
}

func (cs CommandSet) Lookup(name string) (stdlib.Executer, error) {
	exec, ok := cs[name]
	if !ok {
		return nil, fmt.Errorf("%s: undefined", name)
	}
	return exec, nil
}

func (cs CommandSet) List(pat string) []string {
	list := make([]string, 0, len(cs))
	for k, e := range cs {
		if _, ok := e.(procedure); !ok {
			continue
		}
		list = append(list, k)
	}
	return list
}

func (cs CommandSet) Args(proc string) ([]string, error) {
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

func (cs CommandSet) Body(proc string) (string, error) {
	p, err := cs.find(proc)
	return p.Body, err
}

func (cs CommandSet) Default(proc, arg string) (string, bool, error) {
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

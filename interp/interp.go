package interp

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/gotcl/glob"
	"github.com/midbel/gotcl/stdlib"
	"github.com/midbel/slices"
)

const Version = "0.1.2"

const (
	argc     = "argc"
	argv     = "argv"
	arg0     = "arg0"
	tclcmd   = "tcl_command"
	tclver   = "tcl_version"
	tcldepth = "tcl_depth"
)

type Interp struct {
	Env      env.Environment
	Commands map[string]stdlib.Executer
	Files    map[string]*os.File

	Echo  bool
	Depth int
	Count int
}

func New() stdlib.Interpreter {
	i := &Interp{
		Commands: make(map[string]stdlib.Executer),
		Env:      env.EmptyEnv(),
		Files:    make(map[string]*os.File),
	}
	i.registerCmd("set", stdlib.RunSet)
	i.registerCmd("unset", stdlib.RunUnset)
	i.registerCmd("global", stdlib.RunGlobal)
	i.registerCmd("upvar", stdlib.RunUpvar)
	i.registerCmd("append", stdlib.RunAppend)
	i.registerCmd("rename", stdlib.RunRename)
	i.registerCmd("info", stdlib.RunInfos())
	i.registerCmd("string", stdlib.RunString())
	i.registerCmd("clock", stdlib.RunClock())
	i.registerCmd("file", stdlib.RunFile())
	i.registerCmd("exit", stdlib.RunExit)
	i.registerCmd("proc", stdlib.RunProc)
	i.registerCmd("expr", stdlib.RunExpr)
	i.registerCmd("incr", stdlib.RunIncr)
	i.registerCmd("decr", stdlib.RunDecr)
	i.registerCmd("puts", stdlib.RunPuts)
	i.registerCmd("cd", stdlib.RunChdir)
	i.registerCmd("pid", stdlib.RunPid)
	i.registerCmd("pwd", stdlib.RunPwd)
	i.registerCmd("for", stdlib.RunFor)
	i.registerCmd("while", stdlib.RunWhile)
	i.registerCmd("if", stdlib.RunIf)
	i.registerCmd("switch", stdlib.RunSwitch)
	i.registerCmd("break", stdlib.RunBreak)
	i.registerCmd("continue", stdlib.RunContinue)
	i.registerCmd("::tcl::mathop::!", stdlib.RunNot)
	i.registerCmd("::tcl::mathop::+", stdlib.RunAdd)
	i.registerCmd("::tcl::mathop::-", stdlib.RunSub)
	i.registerCmd("::tcl::mathop::*", stdlib.RunMul)
	i.registerCmd("::tcl::mathop::**", stdlib.RunPow)
	i.registerCmd("::tcl::mathop::/", stdlib.RunDiv)
	i.registerCmd("::tcl::mathop::%", stdlib.RunMod)
	i.registerCmd("::tcl::mathop::==", stdlib.RunEq)
	i.registerCmd("::tcl::mathop::!=", stdlib.RunNe)
	i.registerCmd("::tcl::mathop::<", stdlib.RunLt)
	i.registerCmd("::tcl::mathop::<=", stdlib.RunLe)
	i.registerCmd("::tcl::mathop::>", stdlib.RunGt)
	i.registerCmd("::tcl::mathop::>=", stdlib.RunGe)

	return i
}

func (i *Interp) CmdDepth() int {
	return i.Depth
}

func (i *Interp) CmdCount() int {
	return i.Count
}

func (i *Interp) Procedures(pat string) []string {
	var list []string
	for n, e := range i.Commands {
		if _, ok := e.(procedure); ok {
			list = append(list, n)
		}
	}
	return glob.Filter(list, pat)
}

func (i *Interp) ProcBody(name string) (string, error) {
	e, ok := i.Commands[name]
	if !ok {
		return "", fmt.Errorf("%s: not defined", name)
	}
	p, ok := e.(procedure)
	if !ok {
		return "", fmt.Errorf("%s: not defined with proc command", name)
	}
	return p.Body, nil
}

func (i *Interp) ProcArgs(name string) ([]string, error) {
	e, ok := i.Commands[name]
	if !ok {
		return nil, fmt.Errorf("%s: not defined", name)
	}
	p, ok := e.(procedure)
	if !ok {
		return nil, fmt.Errorf("%s: not defined with proc command", name)
	}
	var args []string
	for _, a := range p.Args {
		args = append(args, a.Name)
	}
	return args, nil
}

func (i *Interp) ProcDefault(name string, arg string) (string, bool, error) {
	return "", false, nil
}

func (i *Interp) Version() string {
	return Version
}

func (i *Interp) Out(str string) {
	fmt.Fprintln(os.Stdout, str)
}

func (i *Interp) Err(str string) {
	fmt.Fprintln(os.Stderr, str)
}

func (i *Interp) Define(name, value string) error {
	if isSpecial(name) {
		return env.ErrForbidden
	}
	return i.Env.Define(name, value)
}

func (i *Interp) Resolve(name string) (string, error) {
	switch name {
	case argc:
		n := len(os.Args) - 1
		return strconv.Itoa(n), nil
	case argv:
		return strings.Join(os.Args[1:], " "), nil
	case arg0:
		return os.Args[0], nil
	case tclcmd:
		return strconv.Itoa(i.Count), nil
	case tcldepth:
		return strconv.Itoa(i.Depth), nil
	case tclver:
		return i.Version(), nil
	default:
		return i.Env.Resolve(name)
	}
}

func (i *Interp) Delete(name string) error {
	if isSpecial(name) {
		return env.ErrForbidden
	}
	return i.Env.Delete(name)
}

func (i *Interp) Exists(name string) bool {
	if isSpecial(name) {
		return true
	}
	return i.Env.Exists(name)
}

func (i *Interp) IsSet(name string) bool {
	if isSpecial(name) {
		return true
	}
	return i.Env.IsSet(name)
}

func (i *Interp) Do(name string, do func(string) (string, error)) (string, error) {
	res, err := i.Resolve(name)
	if err != nil {
		return res, err
	}
	res, err = do(res)
	if err == nil {
		err = i.Define(name, res)
	}
	return res, err
}

func (i *Interp) RegisterFunc(name, args, body string) error {
	p, err := createProcedure(args, body)
	if err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	i.registerFunc(name, p)
	return nil
}

func (i *Interp) registerCmd(name string, cmd stdlib.CommandFunc) {
	i.registerFunc(name, createExecuter(cmd))
}

func (i *Interp) registerFunc(name string, cmd stdlib.Executer) {
	i.Commands[name] = cmd
}

func (i *Interp) UnregisterFunc(name string) {
	delete(i.Commands, name)
}

func (i *Interp) RenameFunc(prev, next string) {
	cmd, ok := i.Commands[prev]
	if !ok {
		return
	}
	delete(i.Commands, prev)
	if next == "" {
		return
	}
	i.Commands[next] = cmd
}

func (i *Interp) Sub() stdlib.Interpreter {
	s := *i
	s.Depth++
	s.Commands = make(map[string]stdlib.Executer)
	for k, cmd := range i.Commands {
		s.Commands[k] = cmd
	}
	s.Env = env.EnclosedEnv(i.Env)
	return &s
}

func (i *Interp) Split(str string) ([]string, error) {
	list, err := scan(str)
	if err == nil {
		list = slices.Filter(list, func(v string) bool { return v != "" })
	}
	return list, err
}

func (i *Interp) Execute(r io.Reader) (string, error) {
	defer func() {
		i.Depth--
	}()
	i.Depth++
	return i.execute(r)
}

func (i *Interp) execute(r io.Reader) (string, error) {
	b, err := Build(r)
	if err != nil {
		return "", err
	}
	var last string
	for {
		c, err := b.Next(i)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return last, err
		}
		last, err = i.executeCmd(c)
		i.Count++
		if err != nil {
			switch {
			case errors.Is(err, stdlib.ErrExit):
			case errors.Is(err, stdlib.ErrReturn):
				err = nil
			default:
			}
			return last, err
		}
	}
	return last, nil
}

func (i *Interp) executeCmd(c *Command) (string, error) {
	exec, ok := i.Commands[c.Cmd]
	if !ok {
		return i.executeExt(c)
	}
	if i.Echo {
		i.Out(fmt.Sprintf("execute: %s", c.Cmd))
	}
	res, err := exec.Execute(i, c.Args)
	if err != nil {
		err = fmt.Errorf("%s: %w", c.Cmd, err)
	}
	return res, err
}

func (i *Interp) executeExt(c *Command) (string, error) {
	res, err := exec.Command(c.Cmd, c.Args...).Output()
	return string(res), err
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
	i = i.Sub()
	for j, a := range p.Args {
		if j < len(args) {
			a.Default = args[j]
		}
		i.Define(a.Name, a.Default)
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

func isSpecial(name string) bool {
	switch name {
	default:
		return false
	case argc:
	case argv:
	case arg0:
	case tclcmd:
	case tclver:
	case tcldepth:
	}
	return true
}

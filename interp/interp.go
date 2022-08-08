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
	"github.com/midbel/gotcl/stdlib"
)

const Version = "0.1.2"

type Interp struct {
	Env      env.Environment
	Commands map[string]stdlib.CommandFunc
	Files    map[string]*os.File

	Echo  bool
	Depth int
	Count int
}

func New() stdlib.Interpreter {
	i := &Interp{
		Commands: make(map[string]stdlib.CommandFunc),
		Env:      env.EmptyEnv(),
		Files:    make(map[string]*os.File),
	}
	i.registerFunc("set", stdlib.RunSet)
	i.registerFunc("unset", stdlib.RunUnset)
	i.registerFunc("append", stdlib.RunAppend)
	i.registerFunc("rename", stdlib.RunRename)
	i.registerFunc("infos", stdlib.RunInfos())
	i.registerFunc("string", stdlib.RunString())
	i.registerFunc("clock", stdlib.RunClock())
	i.registerFunc("file", stdlib.RunFile())
	i.registerFunc("exit", stdlib.RunExit)
	i.registerFunc("proc", stdlib.RunProc)
	i.registerFunc("expr", stdlib.RunExpr)
	i.registerFunc("incr", stdlib.RunIncr)
	i.registerFunc("decr", stdlib.RunDecr)
	i.registerFunc("puts", stdlib.RunPuts)
	i.registerFunc("cd", stdlib.RunChdir)
	i.registerFunc("pid", stdlib.RunPid)
	i.registerFunc("pwd", stdlib.RunPwd)
	i.registerFunc("for", stdlib.RunFor)
	i.registerFunc("while", stdlib.RunWhile)
	i.registerFunc("if", stdlib.RunIf)
	i.registerFunc("switch", stdlib.RunSwitch)
	i.registerFunc("break", stdlib.RunBreak)
	i.registerFunc("continue", stdlib.RunContinue)
	i.registerFunc("::tcl::mathop::!", stdlib.RunNot)
	i.registerFunc("::tcl::mathop::+", stdlib.RunAdd)
	i.registerFunc("::tcl::mathop::-", stdlib.RunSub)
	i.registerFunc("::tcl::mathop::*", stdlib.RunMul)
	i.registerFunc("::tcl::mathop::**", stdlib.RunPow)
	i.registerFunc("::tcl::mathop::/", stdlib.RunDiv)
	i.registerFunc("::tcl::mathop::%", stdlib.RunMod)
	i.registerFunc("::tcl::mathop::==", stdlib.RunEq)
	i.registerFunc("::tcl::mathop::!=", stdlib.RunNe)
	i.registerFunc("::tcl::mathop::<", stdlib.RunLt)
	i.registerFunc("::tcl::mathop::<=", stdlib.RunLe)
	i.registerFunc("::tcl::mathop::>", stdlib.RunGt)
	i.registerFunc("::tcl::mathop::>=", stdlib.RunGe)

	return i
}

func (i *Interp) Version() string {
	return Version
}

func (i *Interp) Open(file string) error {
	return nil
}

func (i *Interp) Close(file string) error {
	return nil
}

func (i *Interp) Read(file string, size int) (string, error) {
	return "", nil
}

func (i *Interp) Seek(file string, offset int, whence int) error {
	return nil
}

func (i *Interp) Out(str string) {
	fmt.Fprintln(os.Stdout, str)
}

func (i *Interp) Err(str string) {
	fmt.Fprintln(os.Stderr, str)
}

func (i *Interp) Define(name, value string) error {
	return i.Env.Define(name, value)
}

func (i *Interp) Resolve(name string) (string, error) {
	switch name {
	case "argc":
		n := len(os.Args) - 1
		return strconv.Itoa(n), nil
	case "argv":
		return strings.Join(os.Args[1:], " "), nil
	case "arg0":
		return os.Args[0], nil
	case "tcl_command":
		return strconv.Itoa(i.Count), nil
	case "tcl_depth":
		return strconv.Itoa(i.Depth), nil
	}
	return i.Env.Resolve(name)
}

func (i *Interp) Delete(name string) error {
	return i.Env.Delete(name)
}

func (i *Interp) Exists(name string) bool {
	return i.Env.Exists(name)
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
	cmd, err := makeCommand(args, body)
	if err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	i.registerFunc(name, cmd)
	return nil
}

func (i *Interp) registerFunc(name string, cmd stdlib.CommandFunc) {
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
	s.Env = env.EnclosedEnv(i)
	return &s
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
	res, err := exec(i, c.Args)
	if err != nil {
		err = fmt.Errorf("%s: %w", c.Cmd, err)
	}
	return res, err
}

func (i *Interp) executeExt(c *Command) (string, error) {
	res, err := exec.Command(c.Cmd, c.Args...).Output()
	return string(res), err
}

func makeCommand(args, body string) (stdlib.CommandFunc, error) {
	type arg struct {
		Name    string
		Default string
	}
	var list []arg
	args = strings.TrimSpace(args)
	for {
		var a arg
		args, a.Name, a.Default = splitArg(args)
		list = append(list, a)
		if a.Name == "" && a.Default == "" {
			return nil, fmt.Errorf("syntax error")
		}
		if args == "" {
			break
		}
	}
	return func(i stdlib.Interpreter, args []string) (string, error) {
		i = i.Sub()
		for j, a := range list {
			if j < len(args) {
				a.Default = args[j]
			}
			i.Define(a.Name, a.Default)
		}
		return i.Execute(strings.NewReader(body))
	}, nil
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

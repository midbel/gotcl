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
	"github.com/midbel/slices"
)

const Version = "0.1.2"

const (
	argc     = "argc"
	argv     = "argv"
	arg0     = "argv0"
	tclcmd   = "tcl_command"
	tclver   = "tcl_version"
	tcldepth = "tcl_depth"
)

type Interp struct {
	*Env
	*FileSet
	CommandSet

	Echo  bool
	Count int
}

func New() stdlib.Interpreter {
	i := &Interp{
		CommandSet: DefaultSet(),
		FileSet:    Stdio(),
		Env:        Environ(),
	}
	return i
}

func (i *Interp) RegisterNS(name, script string) error {
	return nil
}

func (i *Interp) UnregisterNS(name string) error {
	return nil
}

func (i *Interp) CmdDepth() int {
	return i.Env.Depth()
}

func (i *Interp) CmdCount() int {
	return i.Count
}

func (i *Interp) Version() string {
	return Version
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
		return strconv.Itoa(i.CmdDepth()), nil
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

func (i *Interp) Link(dst string, src string) error {
	if isSpecial(src) {
		return env.ErrForbidden
	}
	return i.Env.Link(dst, src)
}

func (i *Interp) LinkAt(dst string, src string, level int) error {
	if isSpecial(src) {
		return env.ErrForbidden
	}
	return i.Env.LinkAt(dst, src, level)
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

func (i *Interp) Split(str string) ([]string, error) {
	list, err := scan(str)
	if err == nil {
		list = slices.Filter(list, func(v string) bool { return v != "" })
	}
	return list, err
}

func (i *Interp) Execute(r io.Reader) (string, error) {
	return i.execute(r)
}

func (i *Interp) ExecuteUp(r io.Reader, level int) (string, error) {
	var (
		err error
		old = i.Env
	)
	defer func() {
		i.Env = old
	}()
	i.Env, err = i.Env.Sub(level)
	if err != nil {
		return "", err
	}
	return i.Execute(r)
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
	exec, err := i.Lookup(c.Cmd)
	if err != nil {
		return i.executeExt(c)
	}
	if _, ok := exec.(procedure); ok {
		i.Env.Append()
		defer i.Env.Pop()
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

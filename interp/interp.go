package interp

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
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
	arg0     = "argv0"
	tclcmd   = "tcl_command"
	tclver   = "tcl_version"
	tcldepth = "tcl_depth"
)

type Interp struct {
	*FileSet
	*Env

	root *Namespace
	*Namespace

	Count int
}

func New() stdlib.Interpreter {
	i := &Interp{
		FileSet: Stdio(),
		Env:     Environ(),
		root:    Global(),
	}
	i.Namespace = i.root
	return i
}

func (i *Interp) Globals(pat string) []string {
	var (
		list = []string{argc, argv, arg0, tclcmd, tclver, tcldepth}
		root = i.Env.Root()
	)
	if a, ok := root.(interface{ All() []string }); ok {
		list = append(list, a.All()...)
	}
	sort.Strings(list)
	return glob.Filter(list, pat)
}

func (i *Interp) Locals(pat string) []string {
	a, ok := i.Env.Current().(interface{ All() []string })
	if !ok {
		return nil
	}
	list := a.All()
	sort.Strings(list)
	return glob.Filter(list, pat)
}

func (i *Interp) Variables(pat string) []string {
	var list []string
	list = append(list, i.Globals("")...)
	list = append(list, i.Locals("")...)
	if !i.Namespace.Root() {
		a, ok := i.Namespace.env.(interface{ All() []string })
		if ok {
			list = append(list, a.All()...)
		}
	}
	sort.Strings(list)
	return glob.Filter(list, pat)
}

func (i *Interp) Valid(cmd string) (bool, error) {
	b, err := Build(strings.NewReader(cmd))
	if err != nil {
		return false, err
	}
	_, err = b.Next(i)
	return err == nil, err
}

func (i *Interp) ResolveVar(name string) (string, error) {
	if i.Namespace.Root() {
		return "", fmt.Errorf("variable can not be resolved in global namespace")
	}
	return i.Namespace.env.Resolve(name)
}

func (i *Interp) RegisterVar(name, value string) error {
	if i.Namespace.Root() {
		return fmt.Errorf("variable can not be defined in global namespace")
	}
	i.Namespace.env.Define(name, value)
	return nil
}

func (i *Interp) ExistsNS(name string) bool {
	if name == "" {
		return true
	}
	_, err := i.lookupNS(name)
	return err == nil
}

func (i *Interp) CurrentNS() string {
	return i.Namespace.FQN()
}

func (i *Interp) ParentNS(name string) (string, error) {
	var (
		ns  *Namespace
		err error
	)
	if name == "" {
		ns = i.Namespace
	} else {
		ns, err = i.lookupNS(name)
		if err != nil {
			return "", err
		}
	}
	if ns.Root() {
		return "", nil
	}
	return ns.Parent.FQN(), nil
}

func (i *Interp) ChildrenNS(name, pat string) ([]string, error) {
	var (
		ns  *Namespace
		err error
	)
	if name == "" {
		ns = i.Namespace
	} else {
		ns, err = i.lookupNS(name)
		if err != nil {
			return nil, err
		}
	}
	var list []string
	for _, c := range ns.Children {
		list = append(list, c.FQN())
	}
	return glob.Filter(list, pat), nil
}

func (i *Interp) RegisterNS(name, script string) error {
	names := strings.Split(name, "::")
	if len(names) == 0 {
		return fmt.Errorf("invalid name for namespace")
	}
	ns, err := i.Namespace.GetOrCreate(names)
	if err != nil {
		return err
	}
	old := i.Namespace
	defer func() {
		i.Namespace = old
	}()
	i.Namespace = ns
	_, err = i.Execute(strings.NewReader(script))
	return err
}

func (i *Interp) UnregisterNS(name string) error {
	if name == "" {
		return nil
	}
	ns, err := i.lookupNS(name)
	if err != nil {
		return err
	}
	if ns.Root() {
		return fmt.Errorf("global namespace can not be deleted")
	}
	ns.Parent.Delete(ns.Name)
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
		return i.resolve(name)
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
	var (
		names = strings.Split(c.Cmd, "::")
		exec  stdlib.Executer
		root  bool
		err   error
	)
	if len(names) > 0 && names[0] == "" {
		exec, err = i.root.Command(names[1:])
		root = true
	} else {
		exec, err = i.Command(names)
	}
	if err != nil {
		var cerr *CmdError
		if errors.As(err, &cerr) && cerr.Unknown != nil {
			args := []string{c.Cmd}
			return cerr.Unknown.Execute(i, append(args, c.Args...))
		}
		return "", err
	}
	if _, ok := exec.(procedure); ok {
		var sub *Namespace
		if ns := names[:len(names)-1]; root {
			sub, err = i.root.Get(ns)
		} else {
			sub, err = i.Namespace.Get(ns)
		}
		if err != nil {
			return "", err
		}
		i.Env.Append()
		defer i.Env.Pop()

		old := i.Namespace
		defer func() {
			i.Namespace = old
		}()
		i.Namespace = sub
	}

	res, err := exec.Execute(i, c.Args)
	if err != nil {
		err = fmt.Errorf("%s: %w", c.Cmd, err)
	}
	return res, err
}

func (i *Interp) resolve(name string) (string, error) {
	v, err := i.Env.Resolve(name)
	if err == nil {
		return v, err
	}
	names := strings.Split(name, "::")
	if len(names) == 1 {
		return i.Namespace.env.Resolve(name)
	}
	name = names[len(names)-1]
	ns, err := i.lookupNS(strings.Join(names[:len(names)-1], "::"))
	if err != nil {
		return "", err
	}
	return ns.env.Resolve(name)
}

func (i *Interp) lookupNS(name string) (*Namespace, error) {
	names := strings.Split(name, "::")
	if len(names) == 0 {
		return nil, fmt.Errorf("invalid namespace name")
	}
	if len(names) > 0 && names[0] == "" {
		return i.root.Get(names[1:])
	}
	return i.Namespace.Get(names)
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

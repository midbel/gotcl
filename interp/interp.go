package interp

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/gotcl/glob"
	"github.com/midbel/gotcl/stdlib"
	"github.com/midbel/slices"
)

const Version = "0.0.1"

type Interpreter struct {
	last env.Value
	err  error

	count  int
	safe   bool
	frames []*Frame

	*Fileset

	name     string
	parent   *Interpreter
	children []*Interpreter
}

// type Interactive struct {
// 	*Interpreter
// 	history []*Command
// }
//
// func Interact() *Interactive {
// 	i := Interactive{
// 		Interpreter: defaultInterpreter("", true),
// 	}
// 	return &i
// }

func Interpret() *Interpreter {
	return defaultInterpreter("", true)
}

func defaultInterpreter(name string, safe bool) *Interpreter {
	i := Interpreter{
		safe:    safe,
		name:    name,
		Fileset: Stdio(),
	}
	i.pushDefault(GlobalNS())
	return &i
}

func (i *Interpreter) Version() string {
	return Version
}

func (i *Interpreter) Globals(pat string) []string {
	list := i.rootFrame().Names()
	return glob.Filter(list, pat)
}

func (i *Interpreter) Locals(pat string) []string {
	list := i.currentFrame().Names()
	return glob.Filter(list, pat)
}

func (i *Interpreter) RegisterNS(name, body string) error {
	ns := emptyNS(name)
	if err := i.currentNS().RegisterNS(ns); err != nil {
		return err
	}
	i.pushDefault(ns)
	defer i.pop()

	_, err := i.Execute(strings.NewReader(body))
	return err
}

func (i *Interpreter) UnregisterNS(name string) error {
	return nil
}

func (i *Interpreter) DefineVar(name string, v env.Value) {

}

func (i *Interpreter) CurrentNS() string {
	return i.currentNS().GetName()
}

func (i *Interpreter) ParentNS(n string) (string, error) {
	var (
		name    = strings.Split(n, "::")
		ns, err = i.rootNS().LookupNS(name)
	)
	if err != nil {
		return "", err
	}
	ns = ns.Parent()
	if ns == nil {
		return "", nil
	}
	return ns.FQN(), nil
}

func (i *Interpreter) ChildreNS(n string) ([]string, error) {
	var (
		name    = strings.Split(n, "::")
		ns, err = i.rootNS().LookupNS(name)
	)
	if err != nil {
		return nil, err
	}
	var list []string
	for _, i := range ns.Children() {
		list = append(list, i.GetName())
	}
	return list, nil
}

func (i *Interpreter) HasNS(n string) bool {
	var (
		name   = strings.Split(n, "::")
		_, err = i.rootNS().LookupNS(name)
	)
	return err == nil
}

func (i *Interpreter) LinkVar(src, dst string, level int) error {
	if i.Depth() <= 1 {
		return fmt.Errorf("can not link variables in global level")
	}
	depth := i.Depth() - 1
	if depth < level {
		return fmt.Errorf("can not link variables in level %d", level)
	}
	depth -= level
	i.currentFrame().Define(dst, env.NewLink(src, depth))
	return nil
}

func (i *Interpreter) Root() bool {
	return i.parent == nil
}

func (i *Interpreter) GetName() string {
	return i.name
}

func (i *Interpreter) IsSafe() bool {
	return i.safe
}

func (i *Interpreter) LookupInterpreter(name []string) (*Interpreter, error) {
	if len(name) == 0 {
		return i, nil
	}
	x := sort.Search(len(i.children), func(j int) bool {
		return i.children[j].name >= name[0]
	})
	if x < len(i.children) && i.children[x].name == name[0] {
		if len(name) == 1 {
			return i.children[x], nil
		}
		return i.children[x].LookupInterpreter((name[1:]))
	}
	return nil, fmt.Errorf("%s: interpreter not registered", name[0])
}

func (i *Interpreter) RegisterInterpreter(name []string, safe bool) (string, error) {
	p, err := i.LookupInterpreter(name[:len(name)-1])
	if err != nil {
		return "", err
	}
	s := defaultInterpreter(name[len(name)-1], safe)
	s.parent = p

	x := sort.Search(len(p.children), func(i int) bool {
		return p.children[i].name >= s.name
	})
	if x < len(p.children) && p.children[x].name == s.name {
		return "", fmt.Errorf("%s: interpreter already registered", s.name)
	}
	tmp := append([]*Interpreter{s}, p.children[x:]...)
	p.children = append(p.children, tmp...)
	return s.name, nil
}

func (i *Interpreter) UnregisterInterpreter(name []string) error {
	p, err := i.LookupInterpreter(name[:len(name)-1])
	if err != nil {
		return err
	}
	if n := len(p.children); n == 0 {
		return nil
	} else if n == 1 {
		p.children = p.children[:0]
	}
	n := name[len(name)-1]
	x := sort.Search(len(p.children), func(i int) bool {
		return p.children[i].name >= n
	})
	if x < len(p.children) && p.children[x].name == n {
		return fmt.Errorf("%s: interpreter not registered", n)
	}
	p.children = append(p.children[:x], p.children[x+1:]...)
	return nil
}

func (i *Interpreter) InterpretersList() []string {
	var list []string
	for _, c := range i.children {
		list = append(list, c.name)
	}
	return list
}

func (i *Interpreter) GetHelp(name string) (string, error) {
	exec, err := i.currentNS().LookupExec([]string{name})
	if err != nil {
		return "", err
	}
	var help string
	switch e := exec.(type) {
	case stdlib.Builtin:
		help = e.Help
	case stdlib.Ensemble:
		help = e.Help
	default:
		return "", fmt.Errorf("%s: can not retrieve help", name)
	}
	return help, nil
}

func (i *Interpreter) RegisterDefer(body string) error {
	name := fmt.Sprintf("defer%d", i.Count())
	exec, err := createProcedure(name, body, "")
	if err == nil {
		i.registerDefer(exec)
	}
	return err
}

func (i *Interpreter) RegisterProc(name, body, args string) error {
	exec, err := createProcedure(name, body, args)
	if err == nil {
		i.currentNS().RegisterExec([]string{name}, exec)
	}
	return err
}

func (i *Interpreter) Define(n string, v env.Value) {
	tmp, err := i.currentFrame().Resolve(n)
	if err == nil {
		k, ok := tmp.(env.Link)
		if ok {
			i.frames[k.At()].env.Define(k.String(), v)
			return
		}
	}
	i.currentFrame().Define(n, v)
}

func (i *Interpreter) Delete(n string) {
	v, err := i.currentFrame().Resolve(n)
	if err == nil {
		k, ok := v.(env.Link)
		if ok {
			i.frames[k.At()].env.Delete(k.String())
		}
	}
	i.currentFrame().Delete(n)
}

func (i *Interpreter) Rename(prev, next string) error {
	return i.currentNS().Rename(prev, next)
}

func (i *Interpreter) Resolve(n string) (env.Value, error) {
	name := strings.Split(n, "::")
	if len(name) == 1 {
		v, err := i.currentFrame().Resolve(n)
		if err != nil {
			return nil, err
		}
		if k, ok := v.(env.Link); ok {
			v, err = i.frames[k.At()].env.Resolve(k.String())
		}
		return v, err
	}
	var (
		ps  = slices.Slice(name)
		vs  = slices.Lst(name)
		ns  *Namespace
		err error
	)
	if ps[0] == "" {
		ns, err = i.rootNS().LookupNS(ps[1:])
	} else {
		ns, err = i.currentNS().LookupNS(ps)
	}
	if err != nil {
		return nil, err
	}
	v, err := ns.Resolve(vs)
	if err != nil {
		return nil, err
	}
	if k, ok := v.(env.Link); ok {
		v, err = i.frames[k.At()].env.Resolve(k.String())
	}
	return v, err
}

func (i *Interpreter) Commands(pat string) []string {
	var list []string
	for k := range i.currentNS().CommandSet {
		list = append(list, k)
	}
	return glob.Filter(list, pat)
}

func (i *Interpreter) ProcList(pat string) []string {
	list := i.currentNS().Procedures()
	return glob.Filter(list, pat)
}

func (i *Interpreter) ProcArgs(proc string) ([]string, error) {
	return i.currentNS().Args(proc)
}

func (i *Interpreter) ProcBody(proc string) (string, error) {
	return i.currentNS().Body(proc)
}

func (i *Interpreter) ProcDefault(proc, arg string) (string, bool, error) {
	return "", false, nil
}

func (i *Interpreter) Depth() int {
	return len(i.frames)
}

func (i *Interpreter) Count() int {
	return i.count
}

func (i *Interpreter) ExecuteLevel(r io.Reader, level int, abs bool) (env.Value, error) {
	if !abs {
		level = i.Depth() - level
	}
	old := append([]*Frame{}, i.frames...)
	defer func() {
		i.frames = old
	}()
	i.frames = i.frames[:level]
	return i.Execute(r)
}

func (i Interpreter) IsComplete(cmd string) bool {
	return false
}

func (i *Interpreter) CurrentFrame(level int) (string, []string, error) {
	var (
		f  = i.currentFrame()
		as []string
	)
	for i := range f.args {
		as = append(as, f.args[i].String())
	}
	return f.cmd.GetName(), as, nil
}

func (i *Interpreter) Execute(r io.Reader) (env.Value, error) {
	if i.currentNS().Root() && i.count == 0 {
		defer i.executeDefer()
	}
	p, err := New(r)
	if err != nil {
		return nil, err
	}
	for i.err != nil {
		c, err := p.Parse(i)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		i.last, i.err = i.execute(c)
	}
	return i.last, i.err
}

func (i *Interpreter) execute(c *Command) (env.Value, error) {
	var (
		parts = strings.Split(c.Name.String(), "::")
		ns    *Namespace
		err   error
	)
	if slices.Fst(parts) == "" {
		parts = slices.Rest(parts)
	}
	if n := len(parts); n > 1 {
		ns, err = i.currentNS().LookupNS(slices.Slice(parts))
	} else {
		ns = i.currentNS()
	}
	if err != nil {
		return nil, err
	}
	exec, err := ns.LookupExec(slices.Take(parts, len(parts)-1))
	if err != nil {
		if ns.unknown != nil {
			return ns.unknown(i, slices.Prepend(c.Name, c.Args))
		}
		return nil, err
	}
	if !i.isSafe(exec) {
		return nil, fmt.Errorf("command %s: can not be execute in unsafe interpreter", c.Name.String())
	}
	if _, ok := exec.(procedure); ok {
		i.pushDefault(ns)
		defer i.executeDefer()
	}
	f := i.currentFrame()
	f.cmd = exec
	f.args = c.Args

	defer func() {
		i.count++
	}()
	return exec.Execute(i, c.Args)
}

func (i *Interpreter) isSafe(exec stdlib.Executer) bool {
	if i.Root() {
		return true
	}
	return !i.safe || (i.safe && exec.IsSafe())
}

func (i *Interpreter) pushDefault(ns *Namespace) {
	i.push(ns, nil, nil)
}

func (i *Interpreter) push(ns *Namespace, e stdlib.Executer, as []env.Value) {
	f := &Frame{
		env:  env.EmptyEnv(),
		ns:   ns,
		cmd:  e,
		args: as,
	}
	i.frames = append(i.frames, f)
}

func (i *Interpreter) pop() {
	n := len(i.frames)
	if n == 1 {
		return
	}
	i.frames = i.frames[:n-1]
}

func (i *Interpreter) rootNS() *Namespace {
	f := slices.Fst(i.frames)
	return f.ns
}

func (i *Interpreter) currentNS() *Namespace {
	f := slices.Lst(i.frames)
	return f.ns
}

func (i *Interpreter) rootFrame() *Frame {
	return slices.Fst(i.frames)
}

func (i *Interpreter) currentFrame() *Frame {
	return slices.Lst(i.frames)
}

func (i *Interpreter) registerDefer(exec stdlib.Executer) {
	curr := i.currentFrame()
	curr.deferred = append(curr.deferred, exec)
}

func (i *Interpreter) executeDefer() {
	defer i.pop()
	var (
		last = i.last
		list = slices.Lst(i.frames).deferred
	)
	for j := len(list) - 1; j >= 0; j-- {
		list[j].Execute(i, nil)
	}
	i.last = last
}

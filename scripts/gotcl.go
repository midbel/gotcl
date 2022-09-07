package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/gotcl/stdlib"
	"github.com/midbel/gotcl/word"
	"github.com/midbel/slices"
)

var (
	ErrArgument  = errors.New("wrong number of argument given")
	ErrCast      = errors.New("type can not be casted")
	ErrUndefined = errors.New("undefined name")
	ErrSyntax    = errors.New("syntax error")
)

type Command struct {
	Name env.Value
	Args []env.Value
}

type Parser struct {
	scan *word.Scanner
	curr word.Word
	peek word.Word
}

func New(r io.Reader) (*Parser, error) {
	scan, err := word.Scan(r)
	if err != nil {
		return nil, err
	}
	p := Parser{
		scan: scan,
	}
	p.next()
	p.next()
	return &p, nil
}

func (p *Parser) Parse(i *Interpreter) (*Command, error) {
	p.skipEmptyLines()
	if p.done() {
		return nil, io.EOF
	}
	var (
		c   Command
		err error
	)
	c.Name, err = p.parse(i)
	if err != nil {
		return nil, err
	}
	for !p.done() && !p.curr.IsEOL() {
		arg, err := p.parse(i)
		if err != nil {
			return nil, err
		}
		c.Args = append(c.Args, arg)
	}
	p.next()
	return &c, nil
}

func (p *Parser) parse(i *Interpreter) (env.Value, error) {
	p.skipBlank()
	var vs []env.Value
	for !p.isEnd() {
		if p.curr.Type == word.Illegal {
			return nil, ErrSyntax
		}
		v, err := substitute(p.curr, i)
		if err != nil {
			return nil, err
		}
		vs = append(vs, v)
		p.next()
	}
	if p.isBlank() {
		p.next()
	}
	return list2str(vs), nil
}

func (p *Parser) next() {
	p.curr = p.peek
	p.peek = p.scan.Scan()
}

func (p *Parser) done() bool {
	return p.curr.Type == word.EOF
}

func (p *Parser) skipEnd() {
	for p.isEnd() {
		p.next()
	}
}

func (p *Parser) skipBlank() {
	for p.isBlank() && !p.done() {
		p.next()
	}
}

func (p *Parser) skipEmptyLines() {
	for p.curr.IsEOL() && !p.done() {
		p.next()
	}
}

func (p *Parser) isEnd() bool {
	return p.curr.IsEOL() || p.isBlank()
}

func (p *Parser) isBlank() bool {
	return p.curr.Type == word.Blank
}

func list2str(list []env.Value) env.Value {
	if len(list) == 1 {
		return list[0]
	}
	var str strings.Builder
	for i := range list {
		str.WriteString(list[i].String())
	}
	return env.Str(str.String())
}

func substitute(curr word.Word, i *Interpreter) (env.Value, error) {
	split := func(str string, i *Interpreter) (env.Value, error) {
		scan, err := word.Scan(strings.NewReader(str))
		if err != nil {
			return nil, err
		}
		var list []env.Value
		for {
			w := scan.Split()
			if w.Type == word.EOF {
				break
			}
			val, err := substitute(w, i)
			if err != nil {
				return nil, err
			}
			list = append(list, val)
		}
		return list2str(list), nil
	}
	var (
		val env.Value
		err error
	)
	switch curr.Type {
	case word.Literal, word.Block:
		val = env.Str(curr.Literal)
	case word.Variable:
		val, err = i.Resolve(curr.Literal)
	case word.Quote:
		val, err = split(curr.Literal, i)
	case word.Script:
		val, err = i.Execute(strings.NewReader(curr.Literal))
	default:
		err = fmt.Errorf("%s: %w", curr, ErrSyntax)
	}
	return val, err
}

func RunProc() stdlib.Executer {
	return stdlib.Builtin{
		Name:  "proc",
		Arity: 3,
		Safe:  true,
		Run: func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
			var (
				name = slices.Fst(args).String()
				list = slices.Snd(args).String()
				body = slices.Lst(args).String()
			)
			exec, err := createProcedure(name, body, list)
			if err == nil {
				i.RegisterProc(name, exec)
			}
			return nil, err
		},
	}
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
	set.registerCmd("proc", RunProc())
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

type Namespace struct {
	Name     string
	parent   *Namespace
	children []*Namespace

	env *env.Env
	CommandSet
	unknown stdlib.CommandFunc
}

func EmptyNS(name string) *Namespace {
	return createNS(name, make(CommandSet))
}

func GlobalNS() *Namespace {
	ns := createNS("", DefaultSet())
	ns.RegisterNS(UtilNS())
	ns.unknown = func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
		var (
			name   = slices.Fst(args).String()
			values []string
		)
		for _, a := range slices.Rest(args) {
			values = append(values, a.String())
		}
		res, err := exec.Command(name, values...).Output()
		return env.Str(string(res)), err
	}
	return ns
}

func UtilNS() *Namespace {
	ns := createNS("util", UtilSet())
	ns.env.Define("version", env.Str("1.12.189"))
	return ns
}

func emptyNS(name string) *Namespace {
	return createNS(name, make(CommandSet))
}

func createNS(name string, set CommandSet) *Namespace {
	return &Namespace{
		Name:       name,
		CommandSet: set,
		env:        env.EmptyEnv(),
	}
}

func (n *Namespace) GetName() string {
	return n.Name
}

func (n *Namespace) Resolve(v string) (env.Value, error) {
	return n.env.Resolve(v)
}

func (n *Namespace) RegisterNS(ns *Namespace) error {
	ns.parent = n
	x := sort.Search(len(n.children), func(i int) bool {
		return ns.Name >= n.children[i].Name
	})
	if x < len(n.children) && n.children[x].Name == ns.Name {
		return fmt.Errorf("namespace %s already exists", ns.Name)
	}
	tmp := append([]*Namespace{ns}, n.children[x:]...)
	n.children = append(n.children[:x], tmp...)
	return nil
}

func (n *Namespace) LookupNS(name []string) (*Namespace, error) {
	if len(name) == 0 {
		return n, nil
	}
	name, err := n.validNS(name)
	if err != nil {
		return nil, err
	}
	ns, err := n.lookupNS(name[0])
	if err == nil {
		if len(name) == 1 {
			return ns, nil
		}
		return ns.LookupNS(name[1:])
	}
	return nil, err
}

func (n *Namespace) RegisterExec(name []string, exec stdlib.Executer) error {
	name, err := n.validNS(name)
	if err != nil {
		return err
	}
	if len(name) == 1 {
		n.CommandSet[name[0]] = exec
		return nil
	}
	ns, err := n.lookupNS(name[0])
	if err == nil {
		return ns.RegisterExec(name[1:], exec)
	}
	return err
}

func (n *Namespace) LookupExec(name []string) (stdlib.Executer, error) {
	name, err := n.validNS(name)
	if err != nil {
		return nil, err
	}
	if len(name) > 1 && len(n.children) == 0 {
		return nil, fmt.Errorf("executer (lookup) %s (%s): %w", name[0], n.FQN(), ErrUndefined)
	}
	if len(name) == 1 {
		exec, ok := n.CommandSet[name[0]]
		if ok {
			return exec, nil
		}
		if !n.Root() {
			return n.parent.LookupExec(name)
		}
		return nil, fmt.Errorf("executer (lookup) %s (%s): %w", name[0], n.FQN(), ErrUndefined)
	}
	ns, err := n.lookupNS(name[0])
	if err == nil {
		return ns.LookupExec(name[1:])
	}
	return nil, err
}

func (n *Namespace) Root() bool {
	return n.parent == nil
}

func (n *Namespace) FQN() string {
	if n.Root() {
		return "::"
	}
	if n.parent.Root() {
		return "::" + n.Name
	}
	return n.parent.FQN() + "::" + n.Name
}

func (n *Namespace) lookupNS(name string) (*Namespace, error) {
	x := sort.Search(len(n.children), func(i int) bool {
		return name >= n.children[i].Name
	})
	if x < len(n.children) && n.children[x].Name == name {
		return n.children[x], nil
	}
	return nil, fmt.Errorf("namespace %s (%s): %w", name, n.FQN(), ErrUndefined)
}

func (n *Namespace) validNS(name []string) ([]string, error) {
	if len(name) > 0 && name[0] == "" {
		if !n.Root() {
			return nil, fmt.Errorf("namespace %s: invalid name", name)
		}
		name = name[1:]
	}
	return name, nil
}

type Frame struct {
	env      *env.Env
	ns       *Namespace
	deferred []stdlib.Executer
}

func (f *Frame) Define(n string, v env.Value) {
	f.env.Define(n, v)
}

func (f *Frame) Delete(n string) {
	f.env.Delete(n)
}

func (f *Frame) Resolve(n string) (env.Value, error) {
	v, err := f.env.Resolve(n)
	if err == nil {
		return v, err
	}
	return f.ns.Resolve(n)
}

type Interpreter struct {
	last   env.Value
	count  int
	safe   bool
	frames []*Frame

	Out io.Writer
	Err io.Writer

	name     string
	parent   *Interpreter
	children []*Interpreter
}

func Interpret() *Interpreter {
	return defaultInterpreter("", true)
}

func defaultInterpreter(name string, safe bool) *Interpreter {
	i := Interpreter{
		Out:  os.Stdout,
		Err:  os.Stderr,
		safe: safe,
		name: name,
	}
	i.push(GlobalNS())
	return &i
}

func (i *Interpreter) RegisterNS(name, body string) error {
	ns := emptyNS(name)
	if err := i.currentNS().RegisterNS(ns); err != nil {
		return err
	}
	i.push(ns)
	defer i.pop()

	// body = strings.TrimSpace(body)
	_, err := i.Execute(strings.NewReader(body))
	return err
}

func (i *Interpreter) UnregisterNS(name string) error {
	return nil
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

func (i *Interpreter) RegisterProc(name string, exec stdlib.Executer) {
	i.currentNS().RegisterExec([]string{name}, exec)
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

func (i *Interpreter) Depth() int {
	return len(i.frames)
}

func (i *Interpreter) Count() int {
	return i.count
}

func (i *Interpreter) Level(r io.Reader, level int) (env.Value, error) {
	old := append([]*Frame{}, i.frames...)
	defer func() {
		i.frames = old
	}()
	i.frames = i.frames[:level]
	return i.Execute(r)
}

func (i *Interpreter) Execute(r io.Reader) (env.Value, error) {
	if i.currentNS().Root() && i.count == 0 {
		defer i.executeDefer()
	}
	p, err := New(r)
	if err != nil {
		return nil, err
	}
	for {
		c, err := p.Parse(i)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		i.last, err = i.execute(c)
		if err != nil {
			return nil, err
		}
		// if i.last != nil {
		// 	fmt.Fprintln(i.Out, ">> ", i.last)
		// }
	}
	return i.last, nil
}

func (i *Interpreter) execute(c *Command) (env.Value, error) {
	var (
		parts = strings.Split(c.Name.String(), "::")
		ns    *Namespace
		err   error
	)
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
		i.push(ns)
		defer i.executeDefer()
	}
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

func (i *Interpreter) push(ns *Namespace) {
	f := &Frame{
		env: env.EmptyEnv(),
		ns:  ns,
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

func main() {
	flag.Parse()

	r, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer r.Close()

	i := Interpret()
	v, err := i.Execute(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
	if v != nil {
		fmt.Fprintf(os.Stdout, "%[1]s (%[1]T)", v, v)
		fmt.Fprintln(os.Stdout)
	}
}

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/midbel/gotcl/word"
	"github.com/midbel/slices"
)

var (
	ErrCast      = errors.New("type can not be casted")
	ErrUndefined = errors.New("undefined name")
	ErrSyntax    = errors.New("syntax error")
)

type Value interface {
	fmt.Stringer

	ToList() (Value, error)
	ToArray() (Value, error)
	ToNumber() (Value, error)
	ToString() (Value, error)
	ToBoolean() (Value, error)
}

func isTrue(v Value) bool {
	v, err := v.ToBoolean()
	if err != nil {
		return false
	}
	b, ok := v.(Boolean)
	if !ok {
		return ok
	}
	return b.value
}

type Array struct {
	values map[string]Value
}

func ZipArr(keys []string, values []Value) Value {
	return EmptyArr()
}

func EmptyArr() Value {
	return Array{
		values: make(map[string]Value),
	}
}

func (a Array) String() string {
	var str strings.Builder
	for k, v := range a.values {
		str.WriteString(k)
		str.WriteString(" ")
		str.WriteString(v.String())
	}
	return str.String()
}

func (a Array) ToList() (Value, error) {
	return nil, nil
}

func (a Array) ToArray() (Value, error) {
	return a, nil
}

func (a Array) ToNumber() (Value, error) {
	return nil, ErrCast
}

func (a Array) ToString() (Value, error) {
	return Str(a.String()), nil
}

func (a Array) ToBoolean() (Value, error) {
	return Bool(len(a.values) != 0), nil
}

type List struct {
	values []Value
}

func ListFrom(vs ...Value) Value {
	if len(vs) == 0 {
		return EmptyList()
	}
	var i List
	i.values = append(i.values, vs...)
	return i
}

func EmptyList() Value {
	var i List
	return i
}

func (i List) String() string {
	var list []string
	for _, v := range i.values {
		list = append(list, v.String())
	}
	return strings.Join(list, " ")
}

func (i List) Len() int {
	return len(i.values)
}

func (i List) ToList() (Value, error) {
	return i, nil
}

func (i List) ToArray() (Value, error) {
	if len(i.values)%2 != 0 {
		return nil, ErrCast
	}
	var (
		ks []string
		vs []Value
	)
	for j := 0; j < len(i.values); j += 2 {
		ks = append(ks, i.values[j].String())
		vs = append(vs, i.values[j+1])
	}
	return ZipArr(ks, vs), nil
}

func (i List) ToNumber() (Value, error) {
	return nil, ErrCast
}

func (i List) ToString() (Value, error) {
	return Str(i.String()), nil
}

func (i List) ToBoolean() (Value, error) {
	return nil, ErrCast
}

type String struct {
	value string
}

func Str(str string) Value {
	return String{value: str}
}

func EmptyStr() Value {
	return Str("")
}

func (s String) String() string {
	return s.value
}

func (s String) ToList() (Value, error) {
	return split(s.value)
}

func (s String) ToArray() (Value, error) {
	list, err := s.ToList()
	if err != nil {
		return nil, err
	}
	return list.ToArray()
}

func (s String) ToNumber() (Value, error) {
	n, err := strconv.ParseFloat(s.value, 64)
	if err != nil {
		return nil, err
	}
	return Float(n), nil
}

func (s String) ToString() (Value, error) {
	return s, nil
}

func (s String) ToBoolean() (Value, error) {
	return Bool(s.value != ""), nil
}

type Boolean struct {
	value bool
}

func False() Value {
	return Bool(false)
}

func True() Value {
	return Bool(true)
}

func Bool(b bool) Value {
	return Boolean{value: b}
}

func (b Boolean) String() string {
	if b.value {
		return "1"
	}
	return "0"
}

func (b Boolean) ToList() (Value, error) {
	return ListFrom(b), nil
}

func (b Boolean) ToArray() (Value, error) {
	return nil, ErrCast
}

func (b Boolean) ToNumber() (Value, error) {
	if !b.value {
		return Zero(), nil
	}
	return Float(1), nil
}

func (b Boolean) ToString() (Value, error) {
	return Str(b.String()), nil
}

func (b Boolean) ToBoolean() (Value, error) {
	return b, nil
}

type Number struct {
	value float64
}

func Float(f float64) Value {
	return Number{value: f}
}

func Int(i int64) Value {
	return Float(float64(i))
}

func Zero() Value {
	return Float(0)
}

func (n Number) String() string {
	return strconv.FormatFloat(n.value, 'g', -1, 64)
}

func (n Number) ToList() (Value, error) {
	return ListFrom(n), nil
}

func (n Number) ToArray() (Value, error) {
	return nil, ErrCast
}

func (n Number) ToNumber() (Value, error) {
	return n, nil
}

func (n Number) ToString() (Value, error) {
	str := strconv.FormatFloat(n.value, 'g', -1, 64)
	return Str(str), nil
}

func (n Number) ToBoolean() (Value, error) {
	return Bool(int(n.value) == 0), nil
}

type Env struct {
	values map[string]Value
}

func EmptyEnv() *Env {
	return &Env{
		values: make(map[string]Value),
	}
}

func (e *Env) Delete(n string) {
	delete(e.values, n)
}

func (e *Env) Define(n string, v Value) {
	e.values[n] = v
}

func (e *Env) Resolve(n string) (Value, error) {
	v, ok := e.values[n]
	if !ok {
		return nil, fmt.Errorf("%s: %w", n, ErrUndefined)
	}
	return v, nil
}

type Command struct {
	Name Value
	Args []Value
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

func (p *Parser) parse(i *Interpreter) (Value, error) {
	p.skipBlank()
	var vs []Value
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
	for p.isBlank() {
		p.next()
	}
}

func (p *Parser) isEnd() bool {
	return p.curr.IsEOL() || p.isBlank()
}

func (p *Parser) isBlank() bool {
	return p.curr.Type == word.Blank
}

func list2str(list []Value) Value {
	if len(list) == 1 {
		return list[0]
	}
	var str strings.Builder
	for i := range list {
		str.WriteString(list[i].String())
	}
	return Str(str.String())
}

func split(str string) (Value, error) {
	str = strings.TrimSpace(str)
	scan, err := word.Scan(strings.NewReader(str))
	if err != nil {
		return nil, err
	}
	var list List
	for {
		w := scan.Scan()
		if w.Type == word.EOF {
			break
		}
		if w.Type == word.Blank {
			continue
		}
		switch w.Type {
		case word.Literal:
		case word.Block:
			w.Literal = fmt.Sprintf("{%s}", w.Literal)
		case word.Variable:
			w.Literal = fmt.Sprintf("$%s", w.Literal)
		case word.Script:
			w.Literal = fmt.Sprintf("[%s]", w.Literal)
		case word.Quote:
			w.Literal = fmt.Sprintf("\"%s\"", w.Literal)
		default:
			return nil, fmt.Errorf("%s: %w", w, ErrSyntax)
		}
		list.values = append(list.values, Str(w.Literal))
	}
	return list, nil
}

func substitute(curr word.Word, i *Interpreter) (Value, error) {
	split := func(str string, i *Interpreter) (Value, error) {
		scan, err := word.Scan(strings.NewReader(str))
		if err != nil {
			return nil, err
		}
		var list []Value
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
		val Value
		err error
	)
	switch curr.Type {
	case word.Literal, word.Block:
		val = Str(curr.Literal)
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

type CommandFunc func(*Interpreter, []Value) (Value, error)

func RunTypeOf(i *Interpreter, args []Value) (Value, error) {
	typ := fmt.Sprintf("%T", slices.Fst(args))
	return Str(typ), nil
}

func RunDefer(i *Interpreter, args []Value) (Value, error) {
	var (
		name = fmt.Sprintf("defer%d", i.Count())
		body = slices.Fst(args).String()
	)
	exec, _ := createProcedure(name, body, "")
	i.registerDefer(exec)
	return Str(""), nil
}

func RunProc(i *Interpreter, args []Value) (Value, error) {
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
}

func RunSet(i *Interpreter, args []Value) (Value, error) {
	i.Define(slices.Fst(args).String(), slices.Snd(args))
	return slices.Snd(args), nil
}

func RunUnset(i *Interpreter, args []Value) (Value, error) {
	i.Delete(slices.Fst(args).String())
	return nil, nil
}

func RunPuts(i *Interpreter, args []Value) (Value, error) {
	fmt.Fprintln(i.Out, slices.Fst(args))
	return nil, nil
}

func RunList(i *Interpreter, args []Value) (Value, error) {
	return slices.Fst(args).ToList()
}

func RunListLen(i *Interpreter, args []Value) (Value, error) {
	list, err := slices.Fst(args).ToList()
	if err != nil {
		return nil, err
	}
	n, ok := list.(interface{ Len() int })
	if !ok {
		return Int(0), nil
	}
	return Int(int64(n.Len())), nil
}

type Executer interface {
	Execute(*Interpreter, []Value) (Value, error)
}

type argument struct {
	Name    string
	Default string
}

type procedure struct {
	Name string
	Body string
	Args []argument
}

func createProcedure(name, body, args string) (Executer, error) {
	p := procedure{
		Name: name,
		Body: strings.TrimSpace(body),
	}
	return p, nil
}

func (p procedure) Execute(i *Interpreter, args []Value) (Value, error) {
	return i.Execute(strings.NewReader(p.Body))
}

type cmdExecuter struct {
	fn CommandFunc
}

func fromCommandFunc(fn CommandFunc) Executer {
	return cmdExecuter{fn: fn}
}

func (c cmdExecuter) Execute(i *Interpreter, args []Value) (Value, error) {
	return c.fn(i, args)
}

type CommandSet map[string]Executer

func EmptySet() CommandSet {
	return make(CommandSet)
}

func DefaultSet() CommandSet {
	set := EmptySet()
	set.registerCmd("puts", fromCommandFunc(RunPuts))
	set.registerCmd("set", fromCommandFunc(RunSet))
	set.registerCmd("unset", fromCommandFunc(RunUnset))
	set.registerCmd("list", fromCommandFunc(RunList))
	set.registerCmd("llength", fromCommandFunc(RunListLen))
	set.registerCmd("proc", fromCommandFunc(RunProc))
	return set
}

func UtilSet() CommandSet {
	set := EmptySet()
	set.registerCmd("defer", fromCommandFunc(RunDefer))
	set.registerCmd("typeof", fromCommandFunc(RunTypeOf))
	return set
}

func (cs CommandSet) registerCmd(name string, exec Executer) {
	cs[name] = exec
}

type Namespace struct {
	Name     string
	parent   *Namespace
	children []*Namespace

	env *Env
	CommandSet
	unknown Executer
}

func EmptyNS(name string) *Namespace {
	return createNS(name, make(CommandSet))
}

func GlobalNS() *Namespace {
	ns := createNS("", DefaultSet())
	ns.RegisterNS(UtilNS())
	return ns
}

func UtilNS() *Namespace {
	ns := createNS("util", UtilSet())
	return ns
}

func createNS(name string, set CommandSet) *Namespace {
	return &Namespace{
		Name:       name,
		CommandSet: set,
		env:        EmptyEnv(),
	}
}

func (n *Namespace) Resolve(v string) (Value, error) {
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

func (n *Namespace) RegisterExec(name []string, exec Executer) error {
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

func (n *Namespace) LookupExec(name []string) (Executer, error) {
	name, err := n.validNS(name)
	if err != nil {
		return nil, err
	}
	if len(name) > 1 && len(n.children) == 0 {
		return nil, fmt.Errorf("executer (lookup) %s: %w", name[0], ErrUndefined)
	}
	if len(name) == 1 {
		exec, ok := n.CommandSet[name[0]]
		if ok {
			return exec, nil
		}
		if !n.Root() {
			return n.parent.LookupExec(name)
		}
		return nil, fmt.Errorf("executer (lookup) %s: %w", name[0], ErrUndefined)
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
	return nil, fmt.Errorf("namespace %s: %w", name, ErrUndefined)
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
	env      *Env
	ns       *Namespace
	deferred []Executer
}

func (f *Frame) Define(n string, v Value) {
	f.env.Define(n, v)
}

func (f *Frame) Delete(n string) {
	f.env.Delete(n)
}

func (f *Frame) Resolve(n string) (Value, error) {
	v, err := f.env.Resolve(n)
	if err == nil {
		return v, err
	}
	return f.ns.Resolve(n)
}

type Interpreter struct {
	last   Value
	count  int
	frames []*Frame

	Out io.Writer
	Err io.Writer

	name     string
	parent   *Interpreter
	children []*Interpreter
}

func Interpret() *Interpreter {
	i := Interpreter{
		Out: os.Stdout,
		Err: os.Stderr,
	}
	i.push(GlobalNS())
	return &i
}

func (i *Interpreter) RegisterProc(name string, exec Executer) {
	i.currentNS().RegisterExec([]string{name}, exec)
}

func (i *Interpreter) Define(n string, v Value) {
	i.currentFrame().Define(n, v)
}

func (i *Interpreter) Delete(n string) {
	i.currentFrame().Delete(n)
}

func (i *Interpreter) Resolve(n string) (Value, error) {
	return i.currentFrame().Resolve(n)
}

func (i *Interpreter) Depth() int {
	return len(i.frames)
}

func (i *Interpreter) Count() int {
	return i.count
}

func (i *Interpreter) Execute(r io.Reader) (Value, error) {
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
	}
	return i.last, nil
}

func (i *Interpreter) execute(c *Command) (Value, error) {
	var (
		parts = strings.Split(c.Name.String(), "::")
		ns    *Namespace
		err   error
	)
	if n := len(parts); n > 1 {
		ns, err = i.currentNS().LookupNS(parts[:len(parts)-1])
	} else {
		ns = i.currentNS()
	}
	if err != nil {
		return nil, err
	}
	exec, err := ns.LookupExec(parts[len(parts)-1:])
	if err != nil {
		return nil, err
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

func (i *Interpreter) push(ns *Namespace) {
	f := &Frame{
		env: EmptyEnv(),
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

func (i *Interpreter) currentFrame() *Frame {
	return slices.Lst(i.frames)
}

func (i *Interpreter) registerDefer(exec Executer) {
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

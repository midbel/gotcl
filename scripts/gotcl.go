package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
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
	ToNumber() (Value, error)
	ToString() (Value, error)
	ToBoolean() (Value, error)
}

type List struct {
	values []Value
}

func EmptyList() Value {
	return List{}
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
	return splitString(s.value)
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

type CommandFunc func(*Interpreter, []Value) (Value, error)

func RunTypeOf(i *Interpreter, args []Value) (Value, error) {
	typ := fmt.Sprintf("%T", slices.Fst(args))
	return Str(typ), nil
}

func RunDefer(i *Interpreter, args []Value) (Value, error) {
	i.registerDefer(slices.Fst(args).String())
	return nil, nil
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
	fmt.Println(slices.Fst(args))
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

type Frame struct {
	env      *Env
	deferred []string
}

func Prepare() *Frame {
	return &Frame{
		env: EmptyEnv(),
	}
}

type Interpreter struct {
	last   Value
	frames []*Frame
}

func Interpret() Interpreter {
	return Interpreter{}
}

func (i *Interpreter) push() {
	i.frames = append(i.frames, Prepare())
}

func (i *Interpreter) registerDefer(script string) {
	x := len(i.frames)
	i.frames[x-1].deferred = append(i.frames[x-1].deferred, script)
}

func (i *Interpreter) deferred() {
	defer i.pop()
	var (
		x = len(i.frames)
		a = i.last
	)
	for _, str := range i.frames[x-1].deferred {
		i.Execute(strings.NewReader(str))
	}
	i.last = a
}

func (i *Interpreter) pop() {
	n := len(i.frames)
	if n == 1 {
		return
	}
	i.frames = i.frames[:n-1]
}

func (i *Interpreter) Define(n string, v Value) {
	x := len(i.frames)
	i.frames[x-1].env.Define(n, v)
}

func (i *Interpreter) Delete(n string) {
	x := len(i.frames)
	i.frames[x-1].env.Delete(n)
}

func (i *Interpreter) Resolve(n string) (Value, error) {
	x := len(i.frames)
	v, err := i.frames[x-1].env.Resolve(n)
	if err != nil && x > 0 {
		x--
		if x > 0 {
			return i.frames[x-1].env.Resolve(n)
		}
		return nil, err
	}
	return v, err
}

func (i *Interpreter) Execute(r io.Reader) (Value, error) {
	i.push()
	defer i.deferred()

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
	var exec CommandFunc
	switch name := c.Name.String(); name {
	case "puts":
		exec = RunPuts
	case "set":
		exec = RunSet
	case "unset":
		exec = RunUnset
	case "list":
		exec = RunList
	case "llength":
		exec = RunListLen
	case "typeof":
		exec = RunTypeOf
	case "defer":
		exec = RunDefer
	default:
		return nil, fmt.Errorf("command %s: %w", name, ErrUndefined)
	}
	return exec(i, c.Args)
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

func split(str string, i *Interpreter) (Value, error) {
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

func splitString(str string) (Value, error) {
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
		list.values = append(list.values, Str(w.Literal))
	}
	return list, nil
}

func substitute(curr word.Word, i *Interpreter) (Value, error) {
	var (
		val Value
		err error
	)
	switch curr.Type {
	case word.Literal:
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

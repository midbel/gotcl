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

	"github.com/midbel/gotcl/word"
	"github.com/midbel/slices"
)

var (
	ErrArgument  = errors.New("wrong number of argument given")
	ErrCast      = errors.New("type can not be casted")
	ErrUndefined = errors.New("undefined name")
	ErrSyntax    = errors.New("syntax error")
	ErrType      = errors.New("wrong type given")
)

type Namer interface {
	GetName() string
}

type NamedTree[T any] struct {
	Name     string
	Parent   *T
	Children []T
}

func (n NamedTree[T]) Root() bool {
	return n.Parent == nil
}

func (n NamedTree[T]) GetName() string {
	return n.Name
}

type Value interface {
	fmt.Stringer

	ToList() (Value, error)
	ToArray() (Value, error)
	ToNumber() (Value, error)
	ToString() (Value, error)
	ToBoolean() (Value, error)
}

func asStringList(v Value) ([]string, error) {
	v, err := v.ToList()
	if err != nil {
		return nil, err
	}
	x, ok := v.(List)
	if !ok {
		return nil, nil
	}
	var list []string
	for i := range x.values {
		list = append(list, x.values[i].String())
	}
	return list, nil
}

func asLevel(v Value) (int, bool, error) {
	var (
		str = v.String()
		abs bool
	)
	if strings.HasPrefix(str, "#") {
		abs = true
		str = strings.TrimPrefix(str, "#")
	}
	lvl, err := strconv.Atoi(str)
	return lvl, abs, err
}

func asInt(v Value) (int, error) {
	n, err := v.ToNumber()
	if err != nil {
		return 0, err
	}
	x, ok := n.(Number)
	if !ok {
		return 0, nil
	}
	return int(x.value), nil
}

func asFloat(v Value) (float64, error) {
	n, err := v.ToNumber()
	if err != nil {
		return 0, err
	}
	x, ok := n.(Number)
	if !ok {
		return 0, nil
	}
	return x.value, nil
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

func (a Array) Get(n string) Value {
	return a.values[n]
}

func (a Array) Set(n string, v Value) {
	a.values[n] = v
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

func ListFromStrings(vs []string) Value {
	var list []Value
	for i := range vs {
		list = append(list, Str(vs[i]))
	}
	return ListFrom(list...)
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

type Link struct {
	Value
	level int
}

func createLink(name string, level int) Value {
	return Link{
		Value: Str(name),
		level: level,
	}
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

func scan(str string) ([]string, error) {
	s, err := word.Scan(strings.NewReader(str))
	if err != nil {
		return nil, err
	}
	var list []string
	for {
		w := s.Scan()
		if w.Type == word.Illegal {
			return nil, fmt.Errorf("illegal token")
		}
		if w.Type == word.EOF {
			break
		}
		if w.Literal != "" {
			list = append(list, w.Literal)
		}
	}
	return list, nil
}

type CommandFunc func(*Interpreter, []Value) (Value, error)

type Executer interface {
	Execute(*Interpreter, []Value) (Value, error)
	IsSafe() bool
}

type option struct {
	Value
	Name     string
	Flag     bool
	Required bool
	Check    func(Value) error
}

func checkBool(v Value) error {
	_, ok := v.(Boolean)
	if !ok {
		return ErrType
	}
	return nil
}

func checkNumber(v Value) error {
	_, ok := v.(Number)
	if !ok {
		return ErrType
	}
	return nil
}

func checkString(v Value) error {
	_, ok := v.(String)
	if !ok {
		return ErrType
	}
	return nil
}

func checkChannel(v Value) error {
	switch v.String() {
	case "stdout":
	case "stderr":
	default:
		return fmt.Errorf("%s: unknown channel id", v.String())
	}
	return nil
}

func combineCheck(cs ...func(Value) error) func(Value) error {
	return func(v Value) error {
		for i := range cs {
			if err := cs[i](v); err != nil {
				return err
			}
		}
		return nil
	}
}

type Ensemble struct {
	Name  string
	Usage string
	Help  string
	Safe  bool
	List  []Executer
}

func MakeNamespace() Executer {
	e := Ensemble{
		Name: "namespace",
		List: []Executer{
			Builtin{
				Name:  "eval",
				Arity: 2,
				Safe:  true,
				Run: func(i *Interpreter, args []Value) (Value, error) {
					err := i.RegisterNS(slices.Fst(args).String(), slices.Snd(args).String())
					return EmptyStr(), err
				},
			},
			Builtin{
				Name:  "delete",
				Arity: 1,
				Safe:  true,
				Run: func(i *Interpreter, args []Value) (Value, error) {
					err := i.UnregisterNS(slices.Fst(args).String())
					return EmptyStr(), err
				},
			},
		},
	}
	sort.Slice(e.List, func(i, j int) bool {
		return getName(e.List[i]) < getName(e.List[j])
	})
	return e
}

func MakeArray() Executer {
	e := Ensemble{
		Name: "array",
		List: []Executer{
			Builtin{
				Name:  "set",
				Arity: 2,
				Run: func(i *Interpreter, args []Value) (Value, error) {
					arr, err := i.Resolve(slices.Fst(args).String())
					if err != nil {
						arr = EmptyArr()
					}
					list, err := scan(slices.Snd(args).String())
					if err != nil {
						return nil, err
					}
					if len(list)%2 != 0 {
						return nil, fmt.Errorf("invalid length")
					}
					s := arr.(Array)
					for i := 0; i < len(list); i += 2 {
						s.Set(list[i], Str(list[i+1]))
					}
					i.Define(slices.Fst(args).String(), s)
					return nil, nil
				},
			},
			Builtin{
				Name:  "get",
				Arity: 1,
				Run: func(i *Interpreter, args []Value) (Value, error) {
					arr, err := i.Resolve(slices.Fst(args).String())
					if err != nil {
						return nil, err
					}
					arr, err = arr.ToArray()
					if err != nil {
						return nil, err
					}
					var (
						g  = arr.(Array)
						vs []Value
					)
					for k, v := range g.values {
						vs = append(vs, ListFrom(Str(k), v))
					}
					return ListFrom(vs...), nil
				},
			},
			Builtin{
				Name:  "names",
				Arity: 1,
				Run: func(i *Interpreter, args []Value) (Value, error) {
					return nil, nil
				},
			},
			Builtin{
				Name:  "size",
				Arity: 1,
				Run: func(i *Interpreter, args []Value) (Value, error) {
					return nil, nil
				},
			},
		},
	}
	sort.Slice(e.List, func(i, j int) bool {
		return getName(e.List[i]) < getName(e.List[j])
	})
	return e
}

func MakeInterp() Executer {
	e := Ensemble{
		Name: "interp",
		List: []Executer{
			Builtin{
				Name:  "create",
				Arity: 1,
				Options: []option{
					{
						Name:  "safe",
						Flag:  true,
						Value: False(),
						Check: checkBool,
					},
				},
				Run: func(i *Interpreter, args []Value) (Value, error) {
					paths, err := asStringList(slices.Fst(args))
					if err != nil {
						return nil, err
					}
					safe, err := i.Resolve("safe")
					if err != nil {
						return nil, err
					}
					val, err := i.RegisterInterpreter(paths, isTrue(safe))
					return Str(val), err
				},
			},
			Builtin{
				Name:  "delete",
				Arity: 1,
				Run: func(i *Interpreter, args []Value) (Value, error) {
					paths, err := asStringList(slices.Fst(args))
					if err != nil {
						return nil, err
					}
					return nil, i.UnregisterInterpreter(paths)
				},
			},
			Builtin{
				Name:  "issafe",
				Arity: 1,
				Run: func(i *Interpreter, args []Value) (Value, error) {
					paths, err := asStringList(slices.Fst(args))
					if err != nil {
						return nil, err
					}
					i, err = i.LookupInterpreter(paths)
					if err != nil {
						return nil, err
					}
					return Bool(i.IsSafe()), nil
				},
			},
			Builtin{
				Name:     "eval",
				Variadic: true,
				Run: func(i *Interpreter, args []Value) (Value, error) {
					paths, err := asStringList(slices.Fst(args))
					if err != nil {
						return nil, err
					}
					i, err = i.LookupInterpreter(paths)
					if err != nil {
						return nil, err
					}
					return i.Execute(strings.NewReader(slices.Snd(args).String()))
				},
			},
			Builtin{
				Name:  "children",
				Arity: 1,
				Run: func(i *Interpreter, args []Value) (Value, error) {
					paths, err := asStringList(slices.Fst(args))
					if err != nil {
						return nil, err
					}
					i, err = i.LookupInterpreter(paths)
					if err != nil {
						return nil, err
					}
					list := i.InterpretersList()
					return ListFromStrings(list), nil
				},
			},
		},
	}
	sort.Slice(e.List, func(i, j int) bool {
		return getName(e.List[i]) < getName(e.List[j])
	})
	return e
}

func MakeString() Executer {
	e := Ensemble{
		Name: "string",
		List: []Executer{
			Builtin{
				Name:  "tolower",
				Arity: 1,
				Run: func(i *Interpreter, args []Value) (Value, error) {
					return withString(slices.Fst(args), strings.ToLower)
				},
			},
			Builtin{
				Name:  "toupper",
				Arity: 1,
				Run: func(i *Interpreter, args []Value) (Value, error) {
					return withString(slices.Fst(args), strings.ToUpper)
				},
			},
			Builtin{
				Name:  "length",
				Arity: 1,
				Run: func(i *Interpreter, args []Value) (Value, error) {
					return withString(slices.Fst(args), func(s string) string {
						return strconv.Itoa(len(s))
					})
				},
			},
			Builtin{
				Name:  "repeat",
				Arity: 2,
				Run: func(i *Interpreter, args []Value) (Value, error) {
					c, err := asInt(slices.Snd(args))
					if err != nil {
						return nil, err
					}
					return withString(slices.Fst(args), func(s string) string {
						return strings.Repeat(s, c)
					})
				},
			},
		},
	}
	sort.Slice(e.List, func(i, j int) bool {
		return getName(e.List[i]) < getName(e.List[j])
	})
	return e
}

func withString(v Value, do func(str string) string) (Value, error) {
	str, err := v.ToString()
	if err != nil {
		return nil, err
	}
	return Str(do(str.String())), nil
}

func (e Ensemble) IsSafe() bool {
	return e.Safe
}

func (e Ensemble) GetName() string {
	return e.Name
}

func (e Ensemble) Execute(i *Interpreter, args []Value) (Value, error) {
	name := slices.Fst(args).String()
	x := sort.Search(len(e.List), func(i int) bool {
		return getName(e.List[i]) >= name
	})
	if x >= len(e.List) || getName(e.List[x]) != name {
		return nil, fmt.Errorf("%s %s: command not defined", e.Name, name)
	}
	return e.List[x].Execute(i, slices.Rest(args))
}

func getName(e Executer) string {
	switch e := e.(type) {
	case Builtin:
		return e.Name
	case Ensemble:
		return e.Name
	case procedure:
		return e.Name
	default:
		return ""
	}
}

type Builtin struct {
	Name     string
	Usage    string
	Help     string
	Safe     bool
	Arity    int
	Variadic bool
	Run      CommandFunc
	Options  []option
}

func (b Builtin) IsSafe() bool {
	return b.Safe
}

func (b Builtin) GetName() string {
	return b.Name
}

func (b Builtin) Execute(i *Interpreter, args []Value) (Value, error) {
	if b.Run == nil {
		return nil, fmt.Errorf("%s: command is not runnable", b.Name)
	}
	args, err := b.parseOptions(i, args)
	if err != nil {
		return nil, err
	}
	if err := b.parseArgs(args); err != nil {
		return nil, err
	}
	return b.Run(i, args)
}

func (b Builtin) parseArgs(args []Value) error {
	if n := len(args); n != b.Arity {
		if !b.Variadic || (b.Variadic && n < b.Arity) {
			return fmt.Errorf("%s: %w: want %d, got %d", b.Name, ErrArgument, b.Arity, n)
		}
	}
	return nil
}

func (b Builtin) parseOptions(i *Interpreter, args []Value) ([]Value, error) {
	if len(b.Options) == 0 {
		return args, nil
	}
	sort.Slice(b.Options, func(i, j int) bool {
		return b.Options[i].Name < b.Options[j].Name
	})
	for _, o := range b.Options {
		if o.Value == nil {
			continue
		}
		i.Define(o.Name, o.Value)
	}
	var j int
	for ; j < len(args) && j < len(b.Options); j++ {
		str := args[j].String()
		if str == "--" || !strings.HasPrefix(str, "-") {
			break
		}
		str = strings.TrimPrefix(str, "-")
		x, err := isSet(b.Options, str)
		if err != nil {
			return nil, err
		}
		if b.Options[x].Flag {
			i.Define(b.Options[x].Name, True())
			continue
		}
		if check := b.Options[x].Check; check != nil {
			if err := check(args[j+1]); err != nil {
				return nil, err
			}
			b.Options[j].Value = args[j+1]
			i.Define(str, b.Options[j].Value)
		}
		j++
	}
	if err := isValid(b.Options); err != nil {
		return nil, err
	}
	return args[j:], nil
}

func isValid(list []option) error {
	ok := slices.Every(list, func(o option) bool {
		if !o.Required {
			return true
		}
		return o.Required && o.Value != nil
	})
	if !ok {
		return fmt.Errorf("required options are not provided!")
	}
	return nil
}

func isSet(list []option, name string) (int, error) {
	x := sort.Search(len(list), func(i int) bool {
		return list[i].Name >= name
	})
	if x < len(list) && list[x].Name == name {
		return x, nil
	}
	return 0, fmt.Errorf("%s: option not supported", name)
}

func RunTypeOf() Executer {
	return Builtin{
		Name:  "typeof",
		Arity: 1,
		Safe:  true,
		Run: func(i *Interpreter, args []Value) (Value, error) {
			typ := fmt.Sprintf("%T", slices.Fst(args))
			return Str(typ), nil
		},
	}
}

func RunDefer() Executer {
	return Builtin{
		Name:  "defer",
		Arity: 1,
		Safe:  true,
		Run: func(i *Interpreter, args []Value) (Value, error) {
			var (
				name = fmt.Sprintf("defer%d", i.Count())
				body = slices.Fst(args).String()
			)
			exec, _ := createProcedure(name, body, "")
			i.registerDefer(exec)
			return Str(""), nil
		},
	}
}

func RunHelp() Executer {
	return Builtin{
		Name:  "help",
		Help:  "retrieve help of given builtin command",
		Arity: 1,
		Safe:  true,
		Run: func(i *Interpreter, args []Value) (Value, error) {
			help, err := i.GetHelp(slices.Fst(args).String())
			if err != nil {
				return nil, err
			}
			return Str(help), nil
		},
	}
}

func RunProc() Executer {
	return Builtin{
		Name:  "proc",
		Arity: 3,
		Safe:  true,
		Run: func(i *Interpreter, args []Value) (Value, error) {
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

func RunUplevel() Executer {
	return Builtin{
		Name:     "uplevel",
		Arity:    1,
		Variadic: true,
		Safe:     false,
		Run: func(i *Interpreter, args []Value) (Value, error) {
			var (
				level int
				abs   bool
			)
			if len(args) > 1 {
				x, a, err := asLevel(slices.Fst(args))
				if err != nil {
					return nil, err
				}
				level, abs, args = x, a, slices.Rest(args)
			} else {
				level++
			}
			if !abs {
				level = i.Depth() - level
			}
			return i.Level(strings.NewReader(slices.Fst(args).String()), level)
		},
	}
}

func RunUpvar() Executer {
	return Builtin{
		Name:     "upvar",
		Arity:    2,
		Variadic: true,
		Safe:     false,
		Run: func(i *Interpreter, args []Value) (Value, error) {
			var level int
			if n := len(args) % 2; n == 0 {
				level++
			} else {
				x, err := asInt(slices.Fst(args))
				if err != nil {
					return nil, err
				}
				level = x
				args = slices.Rest(args)
			}
			for j := 0; j < len(args); j += 2 {
				var (
					src = slices.At(args, j)
					dst = slices.At(args, j+1)
				)
				if err := i.LinkVar(src.String(), dst.String(), level); err != nil {
					return nil, err
				}
			}
			return EmptyStr(), nil
		},
	}
}

func RunSet() Executer {
	return Builtin{
		Name:  "set",
		Arity: 2,
		Safe:  true,
		Run: func(i *Interpreter, args []Value) (Value, error) {
			i.Define(slices.Fst(args).String(), slices.Snd(args))
			return slices.Snd(args), nil
		},
	}
}

func RunUnset() Executer {
	return Builtin{
		Name:  "unset",
		Arity: 1,
		Safe:  true,
		Options: []option{
			{
				Name:  "nocomplain",
				Flag:  true,
				Value: False(),
				Check: checkBool,
			},
		},
		Run: func(i *Interpreter, args []Value) (Value, error) {
			i.Delete(slices.Fst(args).String())
			return nil, nil
		},
	}
}

func RunIncr() Executer {
	return Builtin{
		Name:  "incr",
		Arity: 1,
		Safe:  true,
		Run: func(i *Interpreter, args []Value) (Value, error) {
			v, err := i.Resolve(slices.Fst(args).String())
			if err != nil {
				return nil, err
			}
			n, err := asInt(v)
			if err != nil {
				return nil, err
			}
			res := Int(int64(n) + 1)
			i.Define(slices.Fst(args).String(), res)
			return res, nil
		},
	}
}

func RunEval() Executer {
	return Builtin{
		Name:     "eval",
		Help:     "eval given script",
		Variadic: true,
		Safe:     false,
		Run: func(i *Interpreter, args []Value) (Value, error) {
			tmp := ListFrom(args...)
			return i.Execute(strings.NewReader(tmp.String()))
		},
	}
}

func RunPrintArray() Executer {
	return Builtin{
		Name:  "parray",
		Arity: 1,
		Safe:  true,
		Run: func(i *Interpreter, args []Value) (Value, error) {
			arr, err := i.Resolve(slices.Fst(args).String())
			if err != nil {
				return nil, err
			}
			arr, err = arr.ToArray()
			if err != nil {
				return nil, err
			}
			vs := arr.(Array)
			for k, v := range vs.values {
				fmt.Fprintf(i.Out, "%s(%s) = %s", slices.Fst(args), k, v)
				fmt.Fprintln(i.Out)
			}
			return nil, nil
		},
	}
}

func RunPuts() Executer {
	return Builtin{
		Name:  "puts",
		Help:  "print a message to given channel (default to stdout)",
		Arity: 1,
		Safe:  true,
		Options: []option{
			{
				Name:  "nonewline",
				Flag:  true,
				Value: False(),
				Check: checkBool,
			},
			{
				Name:     "channel",
				Value:    Str("stdout"),
				Required: true,
				Check:    combineCheck(checkString, checkChannel),
			},
		},
		Run: func(i *Interpreter, args []Value) (Value, error) {
			str, err := i.Resolve("channel")
			if err != nil {
				return nil, err
			}
			var ch io.Writer
			switch str.String() {
			case "stdout":
				ch = i.Out
			case "stderr":
				ch = i.Err
			default:
				return nil, nil
			}
			fmt.Fprintln(ch, slices.Fst(args))
			return Str(""), nil
		},
	}
}

func RunList() Executer {
	return Builtin{
		Name:  "list",
		Arity: 1,
		Safe:  true,
		Run: func(i *Interpreter, args []Value) (Value, error) {
			return slices.Fst(args).ToList()
		},
	}
}

func RunListLen() Executer {
	return Builtin{
		Name:  "llength",
		Arity: 1,
		Safe:  true,
		Run: func(i *Interpreter, args []Value) (Value, error) {
			list, err := slices.Fst(args).ToList()
			if err != nil {
				return nil, err
			}
			n, ok := list.(interface{ Len() int })
			if !ok {
				return Int(0), nil
			}
			return Int(int64(n.Len())), nil
		},
	}
}

type argument struct {
	Name    string
	Default Value
}

func createArg(n string, v Value) argument {
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
			a = createArg(ws[0], Str(ws[1]))
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

func createProcedure(name, body, args string) (Executer, error) {
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

func (p procedure) Execute(i *Interpreter, args []Value) (Value, error) {
	for j, a := range p.Args {
		if j < len(args) {
			a.Default = args[j]
		}
		i.Define(a.Name, a.Default)
	}
	return i.Execute(strings.NewReader(p.Body))
}

type CommandSet map[string]Executer

func EmptySet() CommandSet {
	return make(CommandSet)
}

func DefaultSet() CommandSet {
	set := EmptySet()
	set.registerCmd("puts", RunPuts())
	set.registerCmd("set", RunSet())
	set.registerCmd("unset", RunUnset())
	set.registerCmd("list", RunList())
	set.registerCmd("llength", RunListLen())
	set.registerCmd("proc", RunProc())
	set.registerCmd("string", MakeString())
	set.registerCmd("interp", MakeInterp())
	set.registerCmd("eval", RunEval())
	set.registerCmd("upvar", RunUpvar())
	set.registerCmd("uplevel", RunUplevel())
	set.registerCmd("incr", RunIncr())
	set.registerCmd("namespace", MakeNamespace())
	set.registerCmd("parray", RunPrintArray())
	set.registerCmd("array", MakeArray())
	return set
}

func UtilSet() CommandSet {
	set := EmptySet()
	set.registerCmd("defer", RunDefer())
	set.registerCmd("typeof", RunTypeOf())
	set.registerCmd("help", RunHelp())
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
	unknown CommandFunc
}

func EmptyNS(name string) *Namespace {
	return createNS(name, make(CommandSet))
}

func GlobalNS() *Namespace {
	ns := createNS("", DefaultSet())
	ns.RegisterNS(UtilNS())
	ns.unknown = func(i *Interpreter, args []Value) (Value, error) {
		var (
			name   = slices.Fst(args).String()
			values []string
		)
		for _, a := range slices.Rest(args) {
			values = append(values, a.String())
		}
		res, err := exec.Command(name, values...).Output()
		return Str(string(res)), err
	}
	return ns
}

func UtilNS() *Namespace {
	ns := createNS("util", UtilSet())
	ns.env.Define("version", Str("1.12.189"))
	return ns
}

func emptyNS(name string) *Namespace {
	return createNS(name, make(CommandSet))
}

func createNS(name string, set CommandSet) *Namespace {
	return &Namespace{
		Name:       name,
		CommandSet: set,
		env:        EmptyEnv(),
	}
}

func (n *Namespace) GetName() string {
	return n.Name
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
	i.currentFrame().Define(dst, createLink(src, depth))
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
	case Builtin:
		help = e.Help
	case Ensemble:
		help = e.Help
	default:
		return "", fmt.Errorf("%s: can not retrieve help", name)
	}
	return help, nil
}

func (i *Interpreter) RegisterProc(name string, exec Executer) {
	i.currentNS().RegisterExec([]string{name}, exec)
}

func (i *Interpreter) Define(n string, v Value) {
	tmp, err := i.currentFrame().Resolve(n)
	if err == nil {
		k, ok := tmp.(Link)
		if ok {
			i.frames[k.level].env.Define(k.Value.String(), v)
			return
		}
	}
	i.currentFrame().Define(n, v)
}

func (i *Interpreter) Delete(n string) {
	v, err := i.currentFrame().Resolve(n)
	if err == nil {
		k, ok := v.(Link)
		if ok {
			i.frames[k.level].env.Delete(k.Value.String())
		}
	}
	i.currentFrame().Delete(n)
}

func (i *Interpreter) Resolve(n string) (Value, error) {
	name := strings.Split(n, "::")
	if len(name) == 1 {
		v, err := i.currentFrame().Resolve(n)
		if err != nil {
			return nil, err
		}
		if k, ok := v.(Link); ok {
			v, err = i.frames[k.level].env.Resolve(k.Value.String())
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
	if k, ok := v.(Link); ok {
		v, err = i.frames[k.level].env.Resolve(k.Value.String())
	}
	return v, err
}

func (i *Interpreter) Depth() int {
	return len(i.frames)
}

func (i *Interpreter) Count() int {
	return i.count
}

func (i *Interpreter) Level(r io.Reader, level int) (Value, error) {
	old := append([]*Frame{}, i.frames...)
	defer func() {
		i.frames = old
	}()
	i.frames = i.frames[:level]
	return i.Execute(r)
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
		// if i.last != nil {
		// 	fmt.Fprintln(i.Out, ">> ", i.last)
		// }
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

func (i *Interpreter) isSafe(exec Executer) bool {
	if i.Root() {
		return true
	}
	return !i.safe || (i.safe && exec.IsSafe())
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

func (i *Interpreter) rootFrame() *Frame {
	return slices.Fst(i.frames)
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

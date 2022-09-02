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
	Name string
	Parent *T
	Children []T
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

type Executer interface {
	Execute(*Interpreter, []Value) (Value, error)
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
	Name string
	Usage string
	Help string
	Safe bool
	List []Executer
}

func MakeInterp() Executer {
	return Ensemble{
		Name: "interp",
		List: []Executer{
			Builtin{
				Name: "create",
				Run: func(i *Interpreter, args []Value) (Value, error) {
					return nil, nil
				},
			},
			Builtin{
				Name: "delete",
				Run: func(i *Interpreter, args []Value) (Value, error) {
					return nil, nil
				},
			},
			Builtin{
				Name: "issafe",
				Run: func(i *Interpreter, args []Value) (Value, error) {
					return Bool(i.safe), nil
				},
			},
			Builtin{
				Name: "eval",
				Run: func(i *Interpreter, args []Value) (Value, error) {
					return nil, nil
				},
			},
			Builtin{
				Name: "children",
				Run: func(i *Interpreter, args []Value) (Value, error) {
					list := i.InterpretersList()
					return ListFromStrings(list), nil
				},
			},
		},
	}
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

func (b Builtin) GetName() string {
	return b.Name
}

func (b Builtin) Execute(i *Interpreter, args []Value) (Value, error) {
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
			return fmt.Errorf("%w: want %d, got %d", ErrArgument, b.Arity, n)
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
		Safe: true,
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
		Safe: true,
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
		Safe: true,
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
		Safe: true,
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

func RunSet() Executer {
	return Builtin{
		Name:  "set",
		Arity: 2,
		Safe: true,
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
		Safe: true,
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

func RunPuts() Executer {
	return Builtin{
		Name:  "puts",
		Help:  "print a message to given channel (default to stdout)",
		Arity: 1,
		Safe: true,
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
		Safe: true,
		Run: func(i *Interpreter, args []Value) (Value, error) {
			return slices.Fst(args).ToList()
		},
	}
}

func RunListLen() Executer {
	return Builtin{
		Name:  "llength",
		Arity: 1,
		Safe: true,
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
	i := Interpreter{
		Out: os.Stdout,
		Err: os.Stderr,
	}
	i.push(GlobalNS())
	return &i
}

func (i *Interpreter) Root() bool {
	return i.parent == nil
}

func (i *Interpreter) GetName() string {
	return i.name
}

func (i *Interpreter) LookupInterpreter(name []string) (*Interpreter, error) {
	return nil, nil
}

func (i *Interpreter) RegisterInterpreter(name []string) (string, error) {
	return "", nil
}

func (i *Interpreter) UnregisterInterpreter(name []string) (string, error) {
	return "", nil
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
	i.currentFrame().Define(n, v)
}

func (i *Interpreter) Delete(n string) {
	i.currentFrame().Delete(n)
}

func (i *Interpreter) Resolve(n string) (Value, error) {
	name := strings.Split(n, "::")
	if len(name) == 1 {
		return i.currentFrame().Resolve(n)
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
	return ns.Resolve(vs)
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
	var safe bool
	switch e := exec.(type) {
	case procedure:
		i.push(ns)
		defer i.executeDefer()
		safe = true
	case Builtin:
		safe = i.isSafe(e.Safe)
	case Ensemble:
		safe = i.isSafe(e.Safe)
	default:
	}
	if !safe {
		return nil, fmt.Errorf("command: can not be execute in unsafe interpreter")
	}
	defer func() {
		i.count++
	}()
	return exec.Execute(i, c.Args)
}

func (i *Interpreter) isSafe(safe bool) bool {
	if i.Root() || !i.safe {
		return true
	}
	return i.safe && safe
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

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
	"github.com/midbel/gotcl/stdlib/argparse"
	"github.com/midbel/gotcl/word"
	"github.com/midbel/slices"
)

var (
	ErrArgument  = errors.New("wrong number of argument given")
	ErrCast      = errors.New("type can not be casted")
	ErrUndefined = errors.New("undefined name")
	ErrSyntax    = errors.New("syntax error")
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

func checkChannel(v env.Value) error {
	switch v.String() {
	case "stdout":
	case "stderr":
	default:
		return fmt.Errorf("%s: unknown channel id", v.String())
	}
	return nil
}

func MakeNamespace() stdlib.Executer {
	e := stdlib.Ensemble{
		Name: "namespace",
		List: []stdlib.Executer{
			stdlib.Builtin{
				Name:  "eval",
				Arity: 2,
				Safe:  true,
				Run: func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
					err := i.RegisterNS(slices.Fst(args).String(), slices.Snd(args).String())
					return env.EmptyStr(), err
				},
			},
			stdlib.Builtin{
				Name:  "delete",
				Arity: 1,
				Safe:  true,
				Run: func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
					err := i.UnregisterNS(slices.Fst(args).String())
					return env.EmptyStr(), err
				},
			},
		},
	}
	sort.Slice(e.List, func(i, j int) bool {
		return getName(e.List[i]) < getName(e.List[j])
	})
	return e
}

func MakeArray() stdlib.Executer {
	e := stdlib.Ensemble{
		Name: "array",
		List: []stdlib.Executer{
			stdlib.Builtin{
				Name:  "set",
				Arity: 2,
				Run: func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
					arr, err := i.Resolve(slices.Fst(args).String())
					if err != nil {
						arr = env.EmptyArr()
					}
					list, err := scan(slices.Snd(args).String())
					if err != nil {
						return nil, err
					}
					if len(list)%2 != 0 {
						return nil, fmt.Errorf("invalid length")
					}
					s := arr.(env.Array)
					for i := 0; i < len(list); i += 2 {
						s.Set(list[i], env.Str(list[i+1]))
					}
					i.Define(slices.Fst(args).String(), s)
					return nil, nil
				},
			},
			stdlib.Builtin{
				Name:  "get",
				Arity: 1,
				Run: func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
					arr, err := i.Resolve(slices.Fst(args).String())
					if err != nil {
						return nil, err
					}
					arr, err = arr.ToArray()
					if err != nil {
						return nil, err
					}
					var (
						g  = arr.(env.Array)
						vs []env.Value
					)
					for k, v := range g.values {
						vs = append(vs, env.ListFrom(env.Str(k), v))
					}
					return env.ListFrom(vs...), nil
				},
			},
			stdlib.Builtin{
				Name:  "names",
				Arity: 1,
				Run: func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
					arr, err := i.Resolve(slices.Fst(args).String())
					if err != nil {
						return nil, err
					}
					arr, err = arr.ToArray()
					if err != nil {
						return nil, err
					}
					var (
						g  = arr.(env.Array)
						vs []string
					)
					for k := range g.values {
						vs = append(vs, k)
					}
					return env.ListFromStrings(vs), nil
				},
			},
			stdlib.Builtin{
				Name:  "size",
				Arity: 1,
				Run: func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
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

func MakeString() stdlib.Executer {
	e := stdlib.Ensemble{
		Name: "string",
		List: []stdlib.Executer{
			stdlib.Builtin{
				Name:  "tolower",
				Arity: 1,
				Run: func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
					return withString(slices.Fst(args), strings.ToLower)
				},
			},
			stdlib.Builtin{
				Name:  "toupper",
				Arity: 1,
				Run: func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
					return withString(slices.Fst(args), strings.ToUpper)
				},
			},
			stdlib.Builtin{
				Name:  "length",
				Arity: 1,
				Run: func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
					return withString(slices.Fst(args), func(s string) string {
						return strconv.Itoa(len(s))
					})
				},
			},
			stdlib.Builtin{
				Name:  "repeat",
				Arity: 2,
				Run: func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
					c, err := env.ToInt(slices.Snd(args))
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

func withString(v env.Value, do func(str string) string) (env.Value, error) {
	str, err := v.ToString()
	if err != nil {
		return nil, err
	}
	return env.Str(do(str.String())), nil
}

func (e stdlib.Ensemble) IsSafe() bool {
	return e.Safe
}

func (e stdlib.Ensemble) GetName() string {
	return e.Name
}

func (e stdlib.Ensemble) Execute(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
	name := slices.Fst(args).String()
	x := sort.Search(len(e.List), func(i int) bool {
		return getName(e.List[i]) >= name
	})
	if x >= len(e.List) || getName(e.List[x]) != name {
		return nil, fmt.Errorf("%s %s: command not defined", e.Name, name)
	}
	return e.List[x].Execute(i, slices.Rest(args))
}

func getName(e stdlib.Executer) string {
	switch e := e.(type) {
	case stdlib.Builtin:
		return e.Name
	case stdlib.Ensemble:
		return e.Name
	case procedure:
		return e.Name
	default:
		return ""
	}
}

func RunTypeOf() stdlib.Executer {
	return stdlib.Builtin{
		Name:  "typeof",
		Arity: 1,
		Safe:  true,
		Run: func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
			typ := fmt.Sprintf("%T", slices.Fst(args))
			return env.Str(typ), nil
		},
	}
}

func RunDefer() stdlib.Executer {
	return stdlib.Builtin{
		Name:  "defer",
		Arity: 1,
		Safe:  true,
		Run: func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
			var (
				name = fmt.Sprintf("defer%d", i.Count())
				body = slices.Fst(args).String()
			)
			exec, _ := createProcedure(name, body, "")
			i.registerDefer(exec)
			return env.EmptyStr(), nil
		},
	}
}

func RunHelp() stdlib.Executer {
	return stdlib.Builtin{
		Name:  "help",
		Help:  "retrieve help of given builtin command",
		Arity: 1,
		Safe:  true,
		Run: func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
			help, err := i.GetHelp(slices.Fst(args).String())
			if err != nil {
				return nil, err
			}
			return env.Str(help), nil
		},
	}
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

func RunUplevel() stdlib.Executer {
	return stdlib.Builtin{
		Name:     "uplevel",
		Arity:    1,
		Variadic: true,
		Safe:     false,
		Run: func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
			var (
				level int
				abs   bool
			)
			if len(args) > 1 {
				x, a, err := env.ToLevel(slices.Fst(args))
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

func RunUpvar() stdlib.Executer {
	return stdlib.Builtin{
		Name:     "upvar",
		Arity:    2,
		Variadic: true,
		Safe:     false,
		Run: func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
			var level int
			if n := len(args) % 2; n == 0 {
				level++
			} else {
				x, err := env.ToInt(slices.Fst(args))
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
			return env.EmptyStr(), nil
		},
	}
}

func RunSet() stdlib.Executer {
	return stdlib.Builtin{
		Name:  "set",
		Arity: 2,
		Safe:  true,
		Run: func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
			i.Define(slices.Fst(args).String(), slices.Snd(args))
			return slices.Snd(args), nil
		},
	}
}

func RunUnset() stdlib.Executer {
	return stdlib.Builtin{
		Name:  "unset",
		Arity: 1,
		Safe:  true,
		Options: []argparse.Option{
			{
				Name:  "nocomplain",
				Flag:  true,
				Value: env.False(),
				Check: argparse.CheckBool,
			},
		},
		Run: func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
			i.Delete(slices.Fst(args).String())
			return nil, nil
		},
	}
}

func RunIncr() stdlib.Executer {
	return stdlib.Builtin{
		Name:  "incr",
		Arity: 1,
		Safe:  true,
		Run: func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
			v, err := i.Resolve(slices.Fst(args).String())
			if err != nil {
				return nil, err
			}
			n, err := env.ToInt(v)
			if err != nil {
				return nil, err
			}
			res := env.Int(int64(n) + 1)
			i.Define(slices.Fst(args).String(), res)
			return res, nil
		},
	}
}

func RunEval() stdlib.Executer {
	return stdlib.Builtin{
		Name:     "eval",
		Help:     "eval given script",
		Variadic: true,
		Safe:     false,
		Run: func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
			tmp := env.ListFrom(args...)
			return i.Execute(strings.NewReader(tmp.String()))
		},
	}
}

func RunPrintArray() stdlib.Executer {
	return stdlib.Builtin{
		Name:  "parray",
		Arity: 1,
		Safe:  true,
		Run: func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
			arr, err := i.Resolve(slices.Fst(args).String())
			if err != nil {
				return nil, err
			}
			arr, err = arr.ToArray()
			if err != nil {
				return nil, err
			}
			vs := arr.(env.Array)
			for k, v := range vs.values {
				fmt.Fprintf(i.Out, "%s(%s) = %s", slices.Fst(args), k, v)
				fmt.Fprintln(i.Out)
			}
			return nil, nil
		},
	}
}

func RunPuts() stdlib.Executer {
	return stdlib.Builtin{
		Name:  "puts",
		Help:  "print a message to given channel (default to stdout)",
		Arity: 1,
		Safe:  true,
		Options: []argparse.Option{
			{
				Name:  "nonewline",
				Flag:  true,
				Value: env.False(),
				Check: argparse.CheckBool,
			},
			{
				Name:     "channel",
				Value:    env.Str("stdout"),
				Required: true,
				Check:    argparse.CombineCheck(argparse.CheckString, checkChannel),
			},
		},
		Run: func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
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
			return env.EmptyStr(), nil
		},
	}
}

func RunList() stdlib.Executer {
	return stdlib.Builtin{
		Name:  "list",
		Arity: 1,
		Safe:  true,
		Run: func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
			return slices.Fst(args).ToList()
		},
	}
}

func RunListLen() stdlib.Executer {
	return stdlib.Builtin{
		Name:  "llength",
		Arity: 1,
		Safe:  true,
		Run: func(i stdlib.Interpreter, args []env.Value) (env.Value, error) {
			list, err := slices.Fst(args).ToList()
			if err != nil {
				return nil, err
			}
			n, ok := list.(interface{ Len() int })
			if !ok {
				return env.Int(0), nil
			}
			return env.Int(int64(n.Len())), nil
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

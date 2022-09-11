package stdlib

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/gotcl/expr"
	"github.com/midbel/gotcl/expr/types"
	"github.com/midbel/gotcl/glob"
	"github.com/midbel/slices"
)

type DeferHandler interface {
	Interpreter
	RegisterDefer(string) error
}

type LinkHandler interface {
	LinkVar(string, string, int) error
}

func RunDefer() Executer {
	return Builtin{
		Name:  "defer",
		Arity: 1,
		Safe:  true,
		Run:   runDefer,
	}
}

func RunProc() Executer {
	return Builtin{
		Name:  "proc",
		Arity: 3,
		Safe:  true,
		Run:   runProc,
	}
}

func RunUplevel() Executer {
	return Builtin{
		Name:     "uplevel",
		Arity:    1,
		Variadic: true,
		Safe:     false,
		Run:      runUplevel,
	}
}

func RunUpvar() Executer {
	return Builtin{
		Name:     "upvar",
		Arity:    2,
		Variadic: true,
		Safe:     false,
		Run:      runUpvar,
	}
}

func RunGlobal() Executer {
	return Builtin{
		Name:     "global",
		Arity:    2,
		Variadic: true,
		Safe:     false,
		Run:      runGlobal,
	}
}

func RunEval() Executer {
	return Builtin{
		Name:     "eval",
		Help:     "eval given script",
		Variadic: true,
		Safe:     false,
		Run:      runEval,
	}
}

func RunTime() Executer {
	return Builtin{
		Name:  "time",
		Arity: 1,
		Run:   runTime,
	}
}

func RunExit() Executer {
	return Builtin{
		Name:     "exit",
		Arity:    0,
		Variadic: true,
		Run:      runExit,
	}
}

func RunReturn() Executer {
	return Builtin{
		Name: "return",
		Safe: true,
		Options: []Option{
			{
				Name:  "code",
				Value: env.Zero(),
				Check: CheckNumber,
			},
		},
		Run: runReturn,
	}
}

func RunUnknown() Executer {
	return Builtin{
		Name: "unknown",
		// Safe: true,
		// Run: runUnknown,
	}
}

func RunExpr() Executer {
	return Builtin{
		Name:     "expr",
		Variadic: true,
		Safe:     true,
		Run:      runExpr,
	}
}

func RunContinue() Executer {
	return Builtin{
		Name: "continue",
		Safe: true,
		Run:  runContinue,
	}
}

func RunBreak() Executer {
	return Builtin{
		Name: "break",
		Safe: true,
		Run:  runBreak,
	}
}

func RunFor() Executer {
	return Builtin{
		Name:  "for",
		Arity: 4,
		Safe:  true,
		Run:   runFor,
	}
}

func RunWhile() Executer {
	return Builtin{
		Name:  "while",
		Arity: 2,
		Safe:  true,
		Run:   runFor,
	}
}

func RunSwitch() Executer {
	return Builtin{
		Name:  "switch",
		Arity: 2,
		Safe:  true,
		Options: []Option{
			{
				Name:  "exact",
				Flag:  true,
				Value: env.False(),
				Check: CheckBool,
			},
			{
				Name:  "glob",
				Flag:  true,
				Value: env.False(),
				Check: CheckBool,
			},
			{
				Name:  "nocase",
				Flag:  true,
				Value: env.False(),
				Check: CheckBool,
			},
		},
		Run: runSwitch,
	}
}

func RunIf() Executer {
	return Builtin{
		Name:     "if",
		Arity:    2,
		Variadic: true,
		Safe:     true,
		Run:      runIf,
	}
}

func RunError() Executer {
	return Builtin{
		Name: "error",
	}
}

func RunCatch() Executer {
	return Builtin{
		Name: "catch",
	}
}

func RunThrow() Executer {
	return Builtin{
		Name: "throw",
	}
}

func RunTry() Executer {
	return Builtin{
		Name: "try",
	}
}

func runIf(i Interpreter, args []env.Value) (env.Value, error) {
	for len(args) > 0 {
		b, err := testScript(i, slices.Fst(args))
		if err != nil {
			return nil, err
		}
		if next := slices.Snd(args).String(); next == "then" {
			args = slices.Take(args, 2)
		} else {
			args = slices.Take(args, 1)
		}
		if b {
			return i.Execute(strings.NewReader(slices.Fst(args).String()))
		}
		args = slices.Rest(args)
		if kw := slices.Fst(args).String(); kw == "else" {
			break
		} else if kw == "elseif" {
			args = slices.Rest(args)
		}
	}
	return i.Execute(strings.NewReader(slices.Lst(args).String()))
}

func runSwitch(i Interpreter, args []env.Value) (env.Value, error) {
	list, err := env.ToStringList(slices.Lst(args))
	if err != nil {
		return nil, err
	}
	if len(list)%2 != 0 {
		return nil, fmt.Errorf("invalid argument")
	}
	var (
		nocase, _ = i.Resolve("nocase")
		exact, _  = i.Resolve("exact")
		match, _  = i.Resolve("glob")
		input     = slices.Fst(args).String()
	)
	if env.ToBool(nocase) {
		input = strings.ToLower(input)
	}
	var alt string
	for j := 0; j < len(list); j += 2 {
		if list[j] == "default" {
			if j != len(list)-2 {
				return nil, fmt.Errorf("syntax error! default must be the last pattern")
			}
			alt = list[j+1]
			break
		}
		pat := list[j]
		if env.ToBool(nocase) {
			pat = strings.ToLower(pat)
		}
		switch {
		default:
		case env.ToBool(exact) && pat == input:
			return i.Execute(strings.NewReader(list[j+1]))
		case env.ToBool(match) && glob.Match(input, pat):
			return i.Execute(strings.NewReader(list[j+1]))
		}
	}
	if alt != "" {
		return i.Execute(strings.NewReader(alt))
	}
	return nil, nil
}

func runFor(i Interpreter, args []env.Value) (env.Value, error) {
	_, err := i.Execute(strings.NewReader(slices.Fst(args).String()))
	if err != nil {
		return nil, err
	}
	return runLoop(i, slices.Snd(args), slices.Lst(args), slices.At(args, 2))
}

func runWhile(i Interpreter, args []env.Value) (env.Value, error) {
	return runLoop(i, slices.Fst(args), slices.Lst(args), nil)
}

func runReturn(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, ErrReturn
}

func runBreak(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, ErrBreak
}

func runContinue(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, ErrContinue
}

func runExpr(i Interpreter, args []env.Value) (env.Value, error) {
	var str strings.Builder
	for i := range args {
		str.WriteString(args[i].String())
	}
	p, err := expr.Parse(str.String())
	if err != nil {
		return nil, err
	}
	expr, err := p.Parse()
	if err != nil {
		return nil, err
	}
	res, err := expr.Eval(i)
	if err != nil {
		return nil, err
	}
	var val env.Value
	switch res.(type) {
	case types.Boolean:
		b, _ := types.AsBool(res)
		val = env.Bool(b)
	case types.Integer:
		i, _ := types.AsInt(res)
		val = env.Int(i)
	case types.Real:
		f, _ := types.AsFloat(res)
		val = env.Float(f)
	default:
		return env.Str(res.String()), nil
	}
	return val, nil
}

func runTime(i Interpreter, args []env.Value) (env.Value, error) {
	var (
		now    = time.Now()
		_, err = i.Execute(strings.NewReader(slices.Fst(args).String()))
	)
	return env.Str(time.Since(now).String()), err
}

func runExit(i Interpreter, args []env.Value) (env.Value, error) {
	var res env.Value
	if len(args) == 0 {
		res = env.Zero()
	} else {
		x, err := slices.Fst(args).ToNumber()
		if err != nil {
			return nil, err
		}
		res = x
	}
	return res, ErrExit
}

func runEval(i Interpreter, args []env.Value) (env.Value, error) {
	tmp := env.ListFrom(args...)
	return i.Execute(strings.NewReader(tmp.String()))
}

func runDefer(i Interpreter, args []env.Value) (env.Value, error) {
	h, ok := i.(DeferHandler)
	if !ok {
		return nil, fmt.Errorf("interpreter can not register defer call")
	}
	return env.EmptyStr(), h.RegisterDefer(slices.Fst(args).String())
}

func runProc(i Interpreter, args []env.Value) (env.Value, error) {
	h, ok := i.(interface {
		RegisterProc(string, string, string) error
	})
	if !ok {
		return nil, fmt.Errorf("interpreter can not register defer call")
	}
	var (
		name = slices.Fst(args).String()
		list = slices.Snd(args).String()
		body = slices.Lst(args).String()
	)
	return nil, h.RegisterProc(name, body, list)
}

func runUplevel(i Interpreter, args []env.Value) (env.Value, error) {
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
	n, ok := i.(interface {
		ExecuteLevel(io.Reader, int, bool) (env.Value, error)
	})
	if !ok {
		return nil, fmt.Errorf("interpreter can not execute script in a parent level")
	}
	return n.ExecuteLevel(strings.NewReader(slices.Fst(args).String()), level, abs)
}

func runGlobal(i Interpreter, args []env.Value) (env.Value, error) {
	k, ok := i.(LinkHandler)
	if !ok {
		return nil, fmt.Errorf("interpreter can not create link between variables")
	}
	return env.EmptyStr(), linkVars(k, args, 0)
}

func runUpvar(i Interpreter, args []env.Value) (env.Value, error) {
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
	k, ok := i.(LinkHandler)
	if !ok {
		return nil, fmt.Errorf("interpreter can not create link between variables")
	}
	return env.EmptyStr(), linkVars(k, args, level)
}

func linkVars(k LinkHandler, args []env.Value, level int) error {
	for j := 0; j < len(args); j += 2 {
		var (
			src = slices.At(args, j)
			dst = slices.At(args, j+1)
		)
		if err := k.LinkVar(src.String(), dst.String(), level); err != nil {
			return err
		}
	}
	return nil
}

func runLoop(i Interpreter, cdt, body, next env.Value) (env.Value, error) {
	var res env.Value
	for {
		b, err := testScript(i, cdt)
		if err != nil {
			return nil, err
		}
		if !b {
			break
		}
		res, err = i.Execute(strings.NewReader(body.String()))
		if err != nil && !errors.Is(err, ErrContinue) {
			if errors.Is(err, ErrBreak) {
				err = nil
			}
			return nil, err
		}
		if next == nil {
			continue
		}
		_, err = i.Execute(strings.NewReader(next.String()))
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

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
		Safe:  true,
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
			{
				Name:  "level",
				Value: env.Int(1),
				Check: CheckNumber,
			},
		},
		Run: runReturn,
	}
}

func RunUnknown() Executer {
	return Builtin{
		Name:     "unknown",
		Arity:    1,
		Variadic: true,
		Run:      runUnknown,
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

func RunForeach() Executer {
	return Builtin{
		Name:     "foreach",
		Arity:    3,
		Variadic: true,
		Safe:     true,
		Run:      runForeach,
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
		Name:     "error",
		Safe:     true,
		Arity:    1,
		Variadic: true,
		Run:      runError,
	}
}

func RunCatch() Executer {
	return Builtin{
		Name:     "catch",
		Safe:     true,
		Arity:    1,
		Variadic: true,
		Run:      runCatch,
	}
}

func RunThrow() Executer {
	return Builtin{
		Name:  "throw",
		Safe:  true,
		Arity: 2,
		Run:   runThrow,
	}
}

func RunTry() Executer {
	return Builtin{
		Name:     "try",
		Safe:     true,
		Arity:    1,
		Variadic: true,
		Run:      runTry,
	}
}

func runTry(i Interpreter, args []env.Value) (env.Value, error) {
	var finally string
	if v := slices.At(args, len(args)-2); v != nil {
		if v.String() == "finally" {
			finally = slices.Lst(args).String()
		}
		args = slices.Take(args, len(args)-2)
	}
	var (
		res, err = i.Execute(strings.NewReader(slices.Fst(args).String()))
		errtry   error
		errfin   error
	)
	if err != nil {
		e, ok := err.(Error)
		if !ok {
			return nil, err
		}
		args = slices.Rest(args)
		if len(args)%4 != 0 {
			return nil, fmt.Errorf("syntax error")
		}
		for j := 0; j < len(args); j += 4 {
			if v := slices.At(args, j); v.String() != "on" {
				return nil, fmt.Errorf("syntax error")
			}
			c, err := env.ToInt(slices.At(args, j+1))
			if err != nil {
				return nil, err
			}
			if e.Code == c {
				names, err := env.ToStringList(slices.At(args, j+2))
				if err != nil {
					return nil, err
				}
				switch len(names) {
				case 1:
					i.Define(names[0], env.Int(int64(e.Code)))
				case 2:
					i.Define(names[0], env.Int(int64(e.Code)))
					i.Define(names[1], env.Str(err.Error()))
				default:
				}
				script := slices.At(args, j+3)
				_, errtry = i.Execute(strings.NewReader(script.String()))
				break
			}
		}
	}
	if finally != "" {
		_, errfin = i.Execute(strings.NewReader(finally))
	}
	return res, hasError(err, errtry, errfin)
}

func runThrow(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runCatch(i Interpreter, args []env.Value) (env.Value, error) {
	var (
		res, err = i.Execute(strings.NewReader(slices.Fst(args).String()))
		name     = slices.Snd(args)
		code     int64
	)
	if err != nil {
		code = int64(ErrorErr)
		if e, ok := err.(Error); ok {
			code = int64(e.Code)
		}
		res = env.Str(err.Error())
	}
	i.Define(name.String(), res)
	return env.Int(code), nil
}

func runError(i Interpreter, args []env.Value) (env.Value, error) {
	code := ErrorErr
	if v := slices.Snd(args); v != nil {
		c, err := env.ToInt(v)
		if err != nil {
			return nil, ErrorFromError(err)
		}
		code = c
	}
	return nil, ErrorWithCode(slices.Fst(args).String(), code)
}

func runUnknown(i Interpreter, args []env.Value) (env.Value, error) {
	uh, ok := i.(interface {
		RegisterUnknown(string, []env.Value) error
	})
	if !ok {
		return nil, fmt.Errorf("interpreter can not register unknown handler")
	}
	err := uh.RegisterUnknown(slices.Fst(args).String(), slices.Rest(args))
	return nil, ErrorFromError(err)
}

func runIf(i Interpreter, args []env.Value) (env.Value, error) {
	var alt string
	if v := slices.At(args, len(args)-2); v != nil {
		if v.String() == "else" {
			alt = slices.Lst(args).String()
			args = slices.Take(args, len(args)-2)
		}
	}
	for len(args) > 0 {
		b, err := testScript(i, slices.Fst(args))
		if err != nil {
			return nil, ErrorFromError(err)
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
		if kw := slices.Fst(args); kw != nil && kw.String() == "elseif" {
			args = slices.Rest(args)
		}
	}
	if alt != "" {
		return i.Execute(strings.NewReader(alt))
	}
	return env.EmptyStr(), nil
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

func runForeach(i Interpreter, args []env.Value) (env.Value, error) {
	if len(args) % 3 != 0 {
		return nil, fmt.Errorf("wrong number of arguments given")
	}
	var res env.Value
	for j := 0; j < len(args); j+=3 {
		list, err := slices.At(args, j+1).ToList()
		if err != nil {
			return nil, err
		}
		for _, a := range list.(env.List).Values() {
			i.Define(slices.At(args, j).String(), a)
			res, err = i.Execute(strings.NewReader(slices.At(args, j+2).String()))
			if err != nil && !errors.Is(err, ErrContinue) {
				if errors.Is(err, ErrBreak) {
					break
				}
				return nil, err
			}
		}
	}
	return res, nil
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

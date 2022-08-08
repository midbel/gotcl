package stdlib

import (
	"errors"
	"flag"
	"strings"

	"github.com/midbel/gotcl/expr"
	"github.com/midbel/slices"
)

var (
	ErrReturn   = errors.New("return")
	ErrBreak    = errors.New("break")
	ErrContinue = errors.New("continue")
)

func RunFor(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("for", args, func(_ *flag.FlagSet) (int, bool) {
		return 4, true
	})
	if err != nil {
		return "", err
	}
	_, err = i.Execute(strings.NewReader(slices.Fst(args)))
	if err != nil {
		return "", err
	}
	return runLoop(i, slices.Snd(args), slices.Lst(args), slices.At(args, 2))
}

func RunWhile(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("while", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, true
	})
	if err != nil {
		return "", err
	}
	return runLoop(i, slices.Fst(args), slices.Lst(args), "")
}

func RunIf(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("if", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, false
	})
	if err != nil {
		return "", err
	}
	for len(args) > 0 {
		b, err := executeBool(i, slices.Fst(args))
		if err != nil {
			return "", err
		}
		if next := slices.Snd(args); next == "then" {
			args = slices.Take(args, 2)
		} else {
			args = slices.Take(args, 1)
		}
		if b {
			return i.Execute(strings.NewReader(slices.Fst(args)))
		}
		args = slices.Rest(args)
		if kw := slices.Fst(args); kw == "else" {
			break
		} else if kw == "elseif" {
			args = slices.Rest(args)
		}
	}
	return i.Execute(strings.NewReader(slices.Lst(args)))
}

func RunSwitch(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunBreak(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("break", args, nil)
	if err != nil {
		return "", err
	}
	return "", ErrBreak
}

func RunContinue(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("continue", args, nil)
	if err != nil {
		return "", err
	}
	return "", ErrContinue
}

func RunReturn(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("return", args, nil)
	if err != nil {
		return "", err
	}
	return "", ErrReturn
}

func RunExpr(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("expr", args, nil)
	if err != nil {
		return "", err
	}
	var str strings.Builder
	for i := range args {
		str.WriteString(args[i])
	}
	p, err := expr.Parse(str.String())
	if err != nil {
		return "", err
	}
	expr, err := p.Parse()
	if err != nil {
		return "", err
	}
	res, err := expr.Eval(i)
	if err != nil {
		return "", err
	}
	return res.String(), nil
}

func RunProc(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("proc", args, func(_ *flag.FlagSet) (int, bool) {
		return 3, true
	})
	if err != nil {
		return "", err
	}
	err = i.RegisterFunc(slices.Fst(args), slices.Snd(args), slices.Lst(args))
	return "", err
}

func runLoop(i Interpreter, test, next, body string) (string, error) {
	var res string
	for {
		b, err := executeBool(i, test)
		if err != nil || !b {
			return res, err
		}
		res, err = i.Execute(strings.NewReader(body))
		if err != nil && !errors.Is(err, ErrContinue) {
			if errors.Is(err, ErrBreak) {
				err = nil
			}
			return res, err
		}
		if next == "" {
			continue
		}
		_, err = i.Execute(strings.NewReader(next))
		if err != nil {
			return "", err
		}
	}
	return res, nil
}

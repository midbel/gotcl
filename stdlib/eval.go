package stdlib

import (
	"errors"
	"flag"
	"strconv"
	"strings"

	"github.com/midbel/gotcl/expr"
	"github.com/midbel/gotcl/glob"
	"github.com/midbel/slices"
)

var (
	ErrReturn   = errors.New("return")
	ErrBreak    = errors.New("break")
	ErrContinue = errors.New("continue")
)

func RunTry(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunCatch(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunError(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

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
	args, err := parseArgs("switch", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, true
	})
	if err != nil {
		return "", err
	}
	list, err := i.Split(slices.Snd(args))
	if err != nil {
		return "", err
	}
	if len(list)%2 != 0 {
		return "", ErrArgument
	}
	var alt string
	for j := 0; j < len(list); j += 2 {
		if list[j] == "default" {
			alt = list[j+1]
			continue
		}
		if glob.Match(slices.Fst(args), list[j]) {
			return i.Execute(strings.NewReader(list[j+1]))
		}
	}
	if alt != "" {
		return i.Execute(strings.NewReader(alt))
	}
	return "", nil
}

func RunBreak(i Interpreter, args []string) (string, error) {
	_, err := parseArgs("break", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, true
	})
	if err != nil {
		return "", err
	}
	return "", ErrBreak
}

func RunContinue(i Interpreter, args []string) (string, error) {
	_, err := parseArgs("continue", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, true
	})
	if err != nil {
		return "", err
	}
	return "", ErrContinue
}

func RunReturn(i Interpreter, args []string) (string, error) {
	_, err := parseArgs("return", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, true
	})
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

func RunUpLevel(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("uplevel", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, false
	})
	if err != nil {
		return "", err
	}
	var (
		level  = 1
		script = slices.Fst(args)
	)
	if len(args) == 2 {
		level, err = strconv.Atoi(slices.Fst(args))
		if err != nil {
			return "", err
		}
		script = slices.Snd(args)
	}
	return i.ExecuteUp(strings.NewReader(script), level)
}

func RunProc(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("proc", args, func(_ *flag.FlagSet) (int, bool) {
		return 3, true
	})
	if err != nil {
		return "", err
	}
	err = i.RegisterProc(slices.Fst(args), slices.Snd(args), slices.Lst(args))
	return "", err
}

func RunUnknown(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
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

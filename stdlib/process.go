package stdlib

import (
	"flag"
	"strings"
	"time"

	"github.com/midbel/slices"
)

func RunTime(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("time", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	now := time.Now()
	_, err = i.Execute(strings.NewReader(slices.Fst(args)))
	if err != nil {
		return "", err
	}
	return time.Since(now).String(), nil
}

func RunExit(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("exit", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, false
	})
	if err != nil {
		return "", err
	}
	var res string
	switch len(args) {
	case 0:
		res = "0"
	case 1:
		res = slices.Fst(args)
	default:
	}
	return res, ErrExit
}

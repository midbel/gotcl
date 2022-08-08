package stdlib

import (
	"flag"

	"github.com/midbel/slices"
)

func RunExit(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("exit", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, false
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

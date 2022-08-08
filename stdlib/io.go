package stdlib

import (
	"flag"

	"github.com/midbel/slices"
)

func RunPuts(i Interpreter, args []string) (string, error) {
	var nonl bool
	args, err := parseArgs("puts", args, func(set *flag.FlagSet) (int, bool) {
		set.BoolVar(&nonl, "nonewline", nonl, "nonewline")
		return 1, false
	})
	if err != nil {
		return "", err
	}
	var (
		msg = slices.Fst(args)
	  print func(string) = i.Out
	)
	if len(args) == 2 {
		msg = slices.Lst(args)
		if dst := slices.Fst(args); dst == "stdout" {
			print = i.Out
		} else if dst == "stderr" {
			print = i.Err
		} else {
			// TBD
		}
	}
	print(msg)
	return "", nil
}

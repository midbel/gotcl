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
	i.Out(slices.Fst(args))
	return "", nil
}

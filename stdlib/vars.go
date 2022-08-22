package stdlib

import (
	"flag"

	"github.com/midbel/slices"
)

func RunUpvar(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunGlobal(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunSet(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("set", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, false
	})
	if err != nil {
		return "", err
	}
	value := slices.At(args, 1)
	if value == "" {
		return i.Resolve(slices.Fst(args))
	}
	return value, i.Define(slices.Fst(args), value)
}

func RunUnset(i Interpreter, args []string) (string, error) {
	var nocomplain bool
	args, err := parseArgs("unset", args, func(set *flag.FlagSet) (int, bool) {
		set.BoolVar(&nocomplain, "nocomplain", nocomplain, "nocomplain")
		return 0, false
	})
	if err != nil {
		return "", err
	}
	for _, n := range args {
		err := i.Delete(n)
		if err != nil && !nocomplain {
			return "", err
		}
	}
	return "", nil
}

func RunRename(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("rename", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, true
	})
	if err != nil {
		return "", err
	}
	i.RenameFunc(slices.Fst(args), slices.Snd(args))
	return "", nil
}

func RunAppend(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("append", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, false
	})
	if err != nil {
		return "", err
	}
	for _, n := range args {
		i.Delete(n)
	}
	return i.Do(slices.Fst(args), func(str string) (string, error) {
		return str + slices.Snd(args), nil
	})
}

package stdlib

import (
	"flag"
	"strconv"

	"github.com/midbel/slices"
)

func RunUpvar(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("upvar", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, false
	})
	if err != nil {
		return "", err
	}
	var level int
	if len(args) > 1 {
		level, err = strconv.Atoi(slices.Fst(args))
		if err != nil {
			return "", nil
		}
	}
	args = slices.Rest(args)
	if len(args)%2 != 0 {
		return "", ErrArgument
	}
	for j := 0; j < len(args); j += 2 {
		if err := i.LinkAt(args[j], args[j+1], level); err != nil {
			return "", err
		}
	}
	return "", nil
}

func RunGlobal(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("global", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, false
	})
	if err != nil {
		return "", err
	}
	for _, a := range args {
		if err := i.Link(a); err != nil {
			return "", err
		}
	}
	return "", nil
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

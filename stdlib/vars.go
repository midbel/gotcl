package stdlib

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/midbel/slices"
)

type Linker interface {
	Link(string, string) error
	LinkAt(string, string, int) error
}

func RunUpvar(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("upvar", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, false
	})
	if err != nil {
		return "", err
	}
	lk, ok := i.(Linker)
	if !ok {
		return "", fmt.Errorf("interpreter can not create link between variable")
	}
	var level int64
	if len(args)%2 == 0 {
		level++
	} else {
		level, err = strconv.ParseInt(slices.Fst(args), 0, 64)
		if err != nil {
			return "", err
		}
		args = slices.Rest(args)
	}
	for i := 0; i < len(args); i += 2 {
		if err := lk.LinkAt(args[i], args[i+1], int(level)); err != nil {
			return "", err
		}
	}
	return "", nil
}

func RunGlobal(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("global", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, false
	})
	if err != nil {
		return "", err
	}
	lk, ok := i.(Linker)
	if !ok {
		return "", fmt.Errorf("interpreter can not create link between variable")
	}
	for i := 0; i < len(args); i += 2 {
		if err := lk.Link(args[i], args[i+1]); err != nil {
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
	i.Rename(slices.Fst(args), slices.Snd(args))
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

package stdlib

import (
	"flag"
	"fmt"
	"io"
	"strconv"

	"github.com/midbel/gotcl/stdlib/conv"
	"github.com/midbel/slices"
)

type FileManager interface {
	Open(string, string) (string, error)
	Close(string) error
	Eof(string) (bool, error)

	Seek(string, int, int) (int64, error)
	Tell(string) (int64, error)
	Gets(string) (string, error)
}

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
		msg   = slices.Fst(args)
		file  string
		print func(string, string) error = i.Println
	)
	if nonl {
		print = i.Print
	}
	if len(args) == 2 {
		msg = slices.Lst(args)
		file = slices.Fst(args)
	}
	return "", print(file, msg)
}

func RunOpen(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("open", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, false
	})
	if err != nil {
		return "", err
	}
	return withFile(i, func(fm FileManager) (string, error) {
		return fm.Open(slices.Fst(args), slices.Snd(args))
	})
}

func RunClose(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("close", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return withFile(i, func(fm FileManager) (string, error) {
		return "", fm.Close(slices.Fst(args))
	})
}

func RunEof(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("eof", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return withFile(i, func(fm FileManager) (string, error) {
		ok, err := fm.Eof(slices.Fst(args))
		return conv.Bool(ok), err
	})
}

func RunSeek(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("seek", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, false
	})
	if err != nil {
		return "", err
	}
	offset, err := strconv.Atoi(slices.Snd(args))
	var whence int
	switch slices.Lst(args) {
	case "start", "":
		whence = io.SeekStart
	case "current":
		whence = io.SeekCurrent
	case "end":
		whence = io.SeekEnd
	default:
		return "", fmt.Errorf("%s: unknown origin given")
	}
	return withFile(i, func(fm FileManager) (string, error) {
		tell, err := fm.Seek(slices.Fst(args), offset, whence)
		return strconv.FormatInt(tell, 10), err
	})
}

func RunTell(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("tell", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return withFile(i, func(fm FileManager) (string, error) {
		tell, err := fm.Tell(slices.Fst(args))
		return strconv.FormatInt(tell, 10), err
	})
}

func RunGets(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("gets", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return withFile(i, func(fm FileManager) (string, error) {
		return fm.Gets(slices.Fst(args))
	})
}

func withFile(i Interpreter, fn func(FileManager) (string, error)) (string, error) {
	fm, ok := i.(FileManager)
	if !ok {
		return "", fmt.Errorf("")
	}
	return fn(fm)
}

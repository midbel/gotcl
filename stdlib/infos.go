package stdlib

import (
	"flag"
	"os"
	"strconv"

	"github.com/midbel/slices"
)

func RunInfos() CommandFunc {
	set := map[string]CommandFunc{
		"exists":           runExists,
		"hostname":         runHost,
		"tclversion":       runVersion,
		"nameofexecutable": runExecutable,
	}
	return makeEnsemble("infos", set)
}

func runExecutable(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("nameofexecutable", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, true
	})
	if err != nil {
		return "", err
	}
	exec, err := os.Executable()
	if err != nil {
		return "", err
	}
	return exec, nil
}

func runHost(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("hostname", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, true
	})
	if err != nil {
		return "", err
	}
	host, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return host, nil
}

func runExists(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("exists", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	ok := i.Exists(slices.Fst(args))
	return strconv.FormatBool(ok), nil
}

func runVersion(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("version", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, true
	})
	if err != nil {
		return "", err
	}
	return i.Version(), nil
}

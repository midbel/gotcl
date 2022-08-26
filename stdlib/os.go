package stdlib

import (
	"flag"
	"os"
	"strconv"

	"github.com/midbel/slices"
)

func RunChdir(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("cd", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, false
	})
	if err != nil {
		return "", err
	}
	var dir string
	switch len(args) {
	case 0:
		dir, err = os.UserHomeDir()
	case 1:
		dir = slices.Fst(args)
	default:
	}
	if err != nil {
		return "", err
	}
	return "", os.Chdir(dir)
}

func RunPid(i Interpreter, args []string) (string, error) {
	_, err := parseArgs("pid", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, true
	})
	if err != nil {
		return "", err
	}
	pid := os.Getpid()
	return strconv.Itoa(pid), nil
}

func RunPwd(i Interpreter, args []string) (string, error) {
	_, err := parseArgs("pwd", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, true
	})
	if err != nil {
		return "", err
	}
	return os.Getwd()
}

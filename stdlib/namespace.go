package stdlib

import (
	"flag"
	"strings"

	"github.com/midbel/slices"
)

func RunNamespace() CommandFunc {
	set := map[string]CommandFunc{
		"eval":     runEvalNS,
		"current":  runCurrentNS,
		"parent":   runParentNS,
		"children": runChildrenNS,
		"delete":   runDeleteNS,
		"exists":   runExistNS,
		"export":   runExportNS,
		"import":   runImportNS,
		"forget":   runForgetNS,
		"unknown":  runUnknownNS,
	}
	return makeEnsemble("namespace", set)
}

func RunVariable(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func runUnknownNS(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func runEvalNS(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("eval", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, true
	})
	if err != nil {
		return "", err
	}
	return "", i.RegisterNS(slices.Fst(args), slices.Snd(args))
}

func runDeleteNS(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("delete", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return "", i.UnregisterNS(slices.Fst(args))
}

func runExistNS(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("exists", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return "", err
}

func runCurrentNS(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("current", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, true
	})
	if err != nil {
		return "", err
	}
	return i.CurrentNS(), err
}

func runParentNS(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("parent", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, true
	})
	if err != nil {
		return "", err
	}
	return i.ParentNS(), err
}

func runChildrenNS(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("children", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, true
	})
	if err != nil {
		return "", err
	}
	return strings.Join(i.ChildrenNS(), " "), err
}

func runExportNS(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("export", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, false
	})
	if err != nil {
		return "", err
	}
	return "", err
}

func runImportNS(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("import", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, false
	})
	if err != nil {
		return "", err
	}
	return "", err
}

func runForgetNS(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("forget", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, false
	})
	if err != nil {
		return "", err
	}
	return "", err
}

package stdlib

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/midbel/gotcl/stdlib/conv"
	"github.com/midbel/slices"
)

type VariableIntrospecter interface {
	Globals(string) []string
	Locals(string) []string
	Variables(string) []string
}

type CommandIntrospecter interface {
	CmdDepth() int
	CmdCount() int
	CmdList(string) []string
}

type ProcIntrospecter interface {
	ProcList(string) []string
	ProcBody(string) (string, error)
	ProcArgs(string) ([]string, error)
	ProcDefault(string, string) (string, bool, error)
}

func RunInfos() CommandFunc {
	set := map[string]CommandFunc{
		"exists":           runExists,
		"hostname":         runHost,
		"tclversion":       runVersion,
		"nameofexecutable": runExecutable,
		"args":             runProcArgs,
		"body":             runProcBody,
		"cmdcount":         runCmdCount,
		"commands":         runCommands,
		"default":          runProcDefaultArg,
		"globals":          runGlobals,
		"level":            runCmdDepth,
		"locals":           runLocals,
		"procs":            runProcs,
		"vars":             runVars,
		"complete":         runComplete,
	}
	return makeEnsemble("infos", set)
}

func runComplete(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("complete", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	ok, err := i.Valid(slices.Fst(args))
	return conv.Bool(ok), err
}

func runCommands(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("commands", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, false
	})
	if err != nil {
		return "", err
	}
	return introspectCmd(i, func(ci CommandIntrospecter) (string, error) {
		list := ci.CmdList(slices.Fst(args))
		return strings.Join(list, " "), nil
	})
}

func runCmdCount(i Interpreter, args []string) (string, error) {
	_, err := parseArgs("cmdcount", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, true
	})
	if err != nil {
		return "", err
	}
	return introspectCmd(i, func(ci CommandIntrospecter) (string, error) {
		n := strconv.Itoa(ci.CmdCount())
		return n, nil
	})
}

func runCmdDepth(i Interpreter, args []string) (string, error) {
	_, err := parseArgs("level", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, false
	})
	if err != nil {
		return "", err
	}
	return introspectCmd(i, func(ci CommandIntrospecter) (string, error) {
		n := strconv.Itoa(ci.CmdDepth())
		return n, nil
	})
}

func runProcs(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("procs", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, false
	})
	if err != nil {
		return "", err
	}
	return introspectProc(i, func(pi ProcIntrospecter) (string, error) {
		list := pi.ProcList(slices.Fst(args))
		return strings.Join(list, " "), nil
	})
}

func runProcArgs(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("args", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return introspectProc(i, func(pi ProcIntrospecter) (string, error) {
		args, err := pi.ProcArgs(slices.Fst(args))
		if err != nil {
			return "", err
		}
		return strings.Join(args, " "), nil
	})
}

func runProcBody(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("body", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return introspectProc(i, func(pi ProcIntrospecter) (string, error) {
		return pi.ProcBody(slices.Fst(args))
	})
}

func runProcDefaultArg(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("default", args, func(_ *flag.FlagSet) (int, bool) {
		return 3, true
	})
	if err != nil {
		return "", err
	}
	return introspectProc(i, func(pi ProcIntrospecter) (string, error) {
		val, ok, err := pi.ProcDefault(slices.Fst(args), slices.Snd(args))
		if err != nil {
			return "", err
		}
		if !ok {
			return conv.False(), nil
		}
		return conv.True(), i.Define(slices.Lst(args), val)
	})
}

func runGlobals(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("globals", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, false
	})
	if err != nil {
		return "", err
	}
	return introspectVars(i, func(vi VariableIntrospecter) (string, error) {
		args = vi.Globals(slices.Fst(args))
		return strings.Join(args, " "), nil
	})
}

func runLocals(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("locals", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, false
	})
	if err != nil {
		return "", err
	}
	return introspectVars(i, func(vi VariableIntrospecter) (string, error) {
		args = vi.Locals(slices.Fst(args))
		return strings.Join(args, " "), nil
	})
}

func runVars(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("vars", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, false
	})
	if err != nil {
		return "", err
	}
	return introspectVars(i, func(vi VariableIntrospecter) (string, error) {
		args = vi.Variables(slices.Fst(args))
		return strings.Join(args, " "), nil
	})
}

func runExecutable(i Interpreter, args []string) (string, error) {
	_, err := parseArgs("nameofexecutable", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, true
	})
	if err != nil {
		return "", err
	}
	return os.Executable()
}

func runHost(i Interpreter, args []string) (string, error) {
	_, err := parseArgs("hostname", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, true
	})
	if err != nil {
		return "", err
	}
	return os.Hostname()
}

func runExists(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("exists", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	ok := i.Exists(slices.Fst(args))
	return conv.Bool(ok), nil
}

func runVersion(i Interpreter, args []string) (string, error) {
	_, err := parseArgs("version", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, true
	})
	if err != nil {
		return "", err
	}
	return i.Version(), nil
}

func introspectProc(i Interpreter, do func(pi ProcIntrospecter) (string, error)) (string, error) {
	pi, ok := i.(ProcIntrospecter)
	if !ok {
		return "", fmt.Errorf("interpreter can not check procedure(s)")
	}
	return do(pi)
}

func introspectCmd(i Interpreter, do func(pi CommandIntrospecter) (string, error)) (string, error) {
	pi, ok := i.(CommandIntrospecter)
	if !ok {
		return "", fmt.Errorf("interpreter can not check command(s)")
	}
	return do(pi)
}

func introspectVars(i Interpreter, do func(vi VariableIntrospecter) (string, error)) (string, error) {
	vi, ok := i.(VariableIntrospecter)
	if !ok {
		return "", fmt.Errorf("interpreter can not check variable(s)")
	}
	return do(vi)
}

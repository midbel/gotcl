package stdlib

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/midbel/slices"
)

type CommandIntrospecter interface {
	CmdDepth() int
	CmdCount() int
}

type ProcIntrospecter interface {
	Procedures(string) []string
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
	}
	return makeEnsemble("infos", set)
}

func runCommands(i Interpreter, args []string) (string, error) {
	return "", nil
}

func runCmdCount(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("cmdcount", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, true
	})
	if err != nil {
		return "", err
	}
	return introspectCmd(i, func(ci CommandIntrospecter) (string, error) {
		fmt.Println("cmdcount", ci.CmdCount())
		n := strconv.Itoa(ci.CmdCount())
		return n, nil
	})
}

func runCmdDepth(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("level", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, false
	})
	if err != nil {
		return "", err
	}
	return introspectCmd(i, func(ci CommandIntrospecter) (string, error) {
		fmt.Println("cmddepth")
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
		for _, p := range pi.Procedures(slices.Fst(args)) {
			i.Out(p)
		}
		return "", nil
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
		for _, a := range args {
			i.Out(a)
		}
		return "", nil
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
		body, err := pi.ProcBody(slices.Fst(args))
		if err != nil {
			return "", err
		}
		i.Out(body)
		return "", nil
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
			return "0", nil
		}
		return "1", i.Define(slices.Lst(args), val)
	})
}

func runGlobals(i Interpreter, args []string) (string, error) {
	return "", nil
}

func runLocals(i Interpreter, args []string) (string, error) {
	return "", nil
}

func runVars(i Interpreter, args []string) (string, error) {
	return "", nil
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

func introspectProc(i Interpreter, do func(pi ProcIntrospecter) (string, error)) (string, error) {
	pi, ok := i.(ProcIntrospecter)
	if !ok {
		return "", fmt.Errorf("interpreter can not check proc")
	}
	return do(pi)
}

func introspectCmd(i Interpreter, do func(pi CommandIntrospecter) (string, error)) (string, error) {
	pi, ok := i.(CommandIntrospecter)
	if !ok {
		return "", fmt.Errorf("interpreter can not check command")
	}
	return do(pi)
}

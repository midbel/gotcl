package stdlib

import (
	"flag"
	"os"
	"path/filepath"
	"strconv"

	"github.com/midbel/slices"
)

func RunFile() CommandFunc {
	set := map[string]CommandFunc{
		"attributes": runAttributes,
		"channels":   runChannels,
		"copy":       runCopy,
		"delete":     runDelete,
		"dirname":    runDirname,
		"executable": runIsExecutable,
		"exists":     runFileExists,
		"extension":  runExtension,
		"isdir":      runIsDir,
		"isfile":     runIsFile,
		"join":       runJoin,
		"move":       runMove,
		"mkdir":      runMkdir,
		"normalize":  runNormalize,
		"readable":   runReadable,
		"size":       runSize,
		"stat":       runFileStat,
		"type":       runFileType,
		"writable":   runWritable,
		"rootname":   runRootname,
		"owned":      runOwned,
	}
	return makeEnsemble("file", set)
}

func runAttributes(i Interpreter, args []string) (string, error) {
	return "", nil
}

func runChannels(i Interpreter, args []string) (string, error) {
	return "", nil
}

func runCopy(i Interpreter, args []string) (string, error) {
	return "", nil
}

func runDelete(i Interpreter, args []string) (string, error) {
	return "", os.RemoveAll(slices.Fst(args))
}

func runDirname(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("dirname", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return filepath.Dir(slices.Fst(args)), nil
}

func runIsExecutable(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("executable", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return "", nil
}

func runFileExists(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("exists", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	if _, err = os.Stat(slices.Fst(args)); err == nil {
		return "1", nil
	}
	return "0", nil
}

func runExtension(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("extension", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return filepath.Ext(slices.Fst(args)), nil
}

func runIsDir(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("isdir", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	fi, err := os.Stat(slices.Fst(args))
	if err != nil {
		return "0", err
	}
	return strconv.FormatBool(fi.IsDir()), nil
}

func runIsFile(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("isfile", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	fi, err := os.Stat(slices.Fst(args))
	if err != nil {
		return "0", err
	}
	return strconv.FormatBool(fi.Mode().IsRegular()), nil
}

func runJoin(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("join", args, nil)
	if err != nil {
		return "", err
	}
	return filepath.Join(args...), nil
}

func runMkdir(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("mkdir", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return "", os.MkdirAll(slices.Fst(args), 0755)
}

func runMove(i Interpreter, args []string) (string, error) {
	return "", nil
}

func runNormalize(i Interpreter, args []string) (string, error) {
	return "", nil
}

func runOwned(i Interpreter, args []string) (string, error) {
	return "", nil
}

func runReadable(i Interpreter, args []string) (string, error) {
	return "", nil
}

func runRootname(i Interpreter, args []string) (string, error) {
	return "", nil
}

func runSize(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("size", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	fi, err := os.Stat(slices.Fst(args))
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(fi.Size(), 10), nil
}

func runFileStat(i Interpreter, args []string) (string, error) {
	return "", nil
}

func runFileType(i Interpreter, args []string) (string, error) {
	return "", nil
}

func runWritable(i Interpreter, args []string) (string, error) {
	return "", nil
}

package stdlib

import (
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/midbel/gotcl/stdlib/conv"
	"github.com/midbel/gotcl/stdlib/ioutil"
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
		"link":       runLink,
		"rename":     runMove,
		"mkdir":      runMkdir,
		"mtime":      runModTime,
		"normalize":  runNormalize,
		"readable":   runReadable,
		"readlink":   runReadLink,
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
	return "", ErrImplemented
}

func runChannels(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func runCopy(i Interpreter, args []string) (string, error) {
	var force bool
	args, err := parseArgs("copy", args, func(set *flag.FlagSet) (int, bool) {
		set.BoolVar(&force, "force", force, "force")
		return 1, false
	})
	if err != nil {
		return "", err
	}
	return "", ioutil.Copy(slices.Slice(args), slices.Lst(args), force)
}

func runMove(i Interpreter, args []string) (string, error) {
	var force bool
	_, err := parseArgs("move", args, func(set *flag.FlagSet) (int, bool) {
		set.BoolVar(&force, "force", force, "force")
		return 1, false
	})
	if err != nil {
		return "", err
	}
	return "", ErrImplemented
}

func runDelete(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("delete", args, func(set *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return "", os.RemoveAll(slices.Fst(args))
}

func runMkdir(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("mkdir", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err == nil {
		err = os.MkdirAll(slices.Fst(args), 0755)
	}
	return "", err
}

func runLink(i Interpreter, args []string) (string, error) {
	var (
		symlink  bool
		hardlink bool
	)
	args, err := parseArgs("link", args, func(set *flag.FlagSet) (int, bool) {
		set.BoolVar(&symlink, "symbolic", symlink, "symbolic")
		set.BoolVar(&hardlink, "hard", hardlink, "hard")
		return 1, false
	})
	if err != nil {
		return "", err
	}
	if len(args) == 1 {
		return runReadLink(i, args)
	}
	if hardlink {
		err = os.Link(slices.Snd(args), slices.Fst(args))
	} else {
		err = os.Symlink(slices.Snd(args), slices.Fst(args))
	}
	return slices.Snd(args), err
}

func runReadLink(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("readlink", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return os.Readlink(slices.Fst(args))
}

func runRootname(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("rootname", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	file := filepath.Base(slices.Fst(args))
	return strings.TrimSuffix(file, filepath.Ext(file)), nil
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

func runDirname(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("dirname", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return filepath.Dir(slices.Fst(args)), nil
}

func runFileExists(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("exists", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	_, err = os.Stat(slices.Fst(args))
	return conv.Bool(err == nil), nil
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
		return conv.False(), err
	}
	return conv.Bool(fi.IsDir()), nil
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
		return conv.False(), err
	}
	return conv.Bool(fi.Mode().IsRegular()), nil
}

func runIsExecutable(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("executable", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	ok := ioutil.Executable(slices.Fst(args), os.Getuid())
	return conv.Bool(ok), nil
}

func runOwned(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("owned", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	ok := ioutil.Owned(slices.Fst(args), os.Getuid())
	return conv.Bool(ok), nil
}

func runReadable(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("readable", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	ok := ioutil.Readable(slices.Fst(args), os.Getuid())
	return conv.Bool(ok), nil
}

func runWritable(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("writable", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	ok := ioutil.Writable(slices.Fst(args), os.Getuid())
	return conv.Bool(ok), nil
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

func runModTime(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("mtime", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	fi, err := os.Stat(slices.Fst(args))
	if err != nil {
		return "", err
	}
	mod := fi.ModTime()
	return strconv.FormatInt(mod.Unix(), 10), nil
}

func runFileType(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("type", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return ioutil.Type(slices.Fst(args))
}

func runFileStat(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func runJoin(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("join", args, nil)
	if err != nil {
		return "", err
	}
	return filepath.Join(args...), nil
}

func runNormalize(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

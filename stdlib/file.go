package stdlib

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/gotcl/stdlib/ioutil"
	"github.com/midbel/slices"
)

func MakeFile() Executer {
	e := Ensemble{
		Name: "file",
		Safe: true,
		List: []Executer{
			Builtin{
				Name: "attributes",
				Run:  fileAttributes,
			},
			Builtin{
				Name: "channels",
				Run:  fileChannels,
			},
			Builtin{
				Name:     "copy",
				Variadic: true,
				Arity:    2,
				Options: []Option{
					{
						Name:  "force",
						Flag:  true,
						Value: env.False(),
						Check: CheckBool,
					},
				},
				Run: fileCopy,
			},
			Builtin{
				Name:     "rename",
				Variadic: true,
				Arity:    2,
				Options: []Option{
					{
						Name:  "force",
						Flag:  true,
						Value: env.False(),
						Check: CheckBool,
					},
				},
				Run: fileMove,
			},
			Builtin{
				Name:  "delete",
				Arity: 1,
				Run:   fileDelete,
			},
			Builtin{
				Name:  "mkdir",
				Arity: 1,
				Run:   fileMkdir,
			},
			Builtin{
				Name:     "link",
				Arity:    1,
				Variadic: true,
				Options: []Option{
					{
						Name:  "symbolic",
						Flag:  true,
						Value: env.False(),
						Check: CheckBool,
					},
					{
						Name:  "hard",
						Flag:  true,
						Value: env.False(),
						Check: CheckBool,
					},
				},
				Run: fileLink,
			},
			Builtin{
				Name:  "executable",
				Arity: 1,
				Run:   fileIsExecutable,
			},
			Builtin{
				Name:  "exists",
				Arity: 1,
				Run:   fileFileExists,
			},
			Builtin{
				Name:  "dirname",
				Arity: 1,
				Run:   fileDirname,
			},
			Builtin{
				Name:  "rootname",
				Arity: 1,
				Run:   fileRootname,
			},
			Builtin{
				Name:  "extension",
				Arity: 1,
				Run:   fileExtension,
			},
			Builtin{
				Name:     "join",
				Arity:    1,
				Variadic: true,
				Run:      fileJoin,
			},
			Builtin{
				Name:  "normalize",
				Arity: 1,
				Run:   fileNormalize,
			},
			Builtin{
				Name:  "isdir",
				Arity: 1,
				Run:   fileIsDir,
			},
			Builtin{
				Name:  "isfile",
				Arity: 1,
				Run:   fileIsFile,
			},
			Builtin{
				Name:  "mtime",
				Arity: 1,
				Run:   fileModTime,
			},
			Builtin{
				Name:  "readable",
				Arity: 1,
				Run:   fileReadable,
			},
			Builtin{
				Name:  "writable",
				Arity: 1,
				Run:   fileWritable,
			},
			Builtin{
				Name:  "owned",
				Arity: 1,
				Run:   fileOwned,
			},
			Builtin{
				Name:  "readlink",
				Arity: 1,
				Run:   fileReadLink,
			},
			Builtin{
				Name:  "size",
				Arity: 1,
				Run:   fileSize,
			},
			Builtin{
				Name:  "stat",
				Arity: 1,
				Run:   fileStat,
			},
			Builtin{
				Name:  "type",
				Arity: 1,
				Run:   fileType,
			},
		},
	}
	return sortEnsembleCommands(e)
}

func fileAttributes(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func fileChannels(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func fileCopy(i Interpreter, args []env.Value) (env.Value, error) {
	var (
		force, _ = i.Resolve("force")
		list     []string
	)
	for _, a := range slices.Slice(args) {
		vs, err := env.ToStringList(a)
		if err != nil {
			return nil, err
		}
		list = append(list, vs...)
	}
	return nil, ioutil.Copy(list, slices.Lst(args).String(), env.ToBool(force))
}

func fileMove(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func fileDelete(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, os.RemoveAll(slices.Fst(args).String())
}

func fileMkdir(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, os.MkdirAll(slices.Fst(args).String(), 0755)
}

func fileLink(i Interpreter, args []env.Value) (env.Value, error) {
	if len(args) == 1 {
		return fileReadLink(i, args)
	}
	var (
		hard, _ = i.Resolve("hard")
		err     error
	)
	if env.ToBool(hard) {
		err = os.Link(slices.Snd(args).String(), slices.Fst(args).String())
	} else {
		err = os.Symlink(slices.Snd(args).String(), slices.Fst(args).String())
	}
	return slices.Snd(args), err
}

func fileReadLink(i Interpreter, args []env.Value) (env.Value, error) {
	str, err := os.Readlink(slices.Fst(args).String())
	return env.Str(str), err
}

func fileRootname(i Interpreter, args []env.Value) (env.Value, error) {
	file := filepath.Base(slices.Fst(args).String())
	file = strings.TrimSuffix(file, filepath.Ext(file))
	return env.Str(file), nil
}

func fileExtension(i Interpreter, args []env.Value) (env.Value, error) {
	ext := filepath.Ext(slices.Fst(args).String())
	return env.Str(ext), nil
}

func fileDirname(i Interpreter, args []env.Value) (env.Value, error) {
	dir := filepath.Dir(slices.Fst(args).String())
	return env.Str(dir), nil
}

func fileFileExists(i Interpreter, args []env.Value) (env.Value, error) {
	_, err := os.Stat(slices.Fst(args).String())
	return env.Bool(err == nil), nil
}

func fileIsDir(i Interpreter, args []env.Value) (env.Value, error) {
	fi, err := os.Stat(slices.Fst(args).String())
	if err != nil {
		return env.False(), err
	}
	return env.Bool(fi.IsDir()), nil
}

func fileIsFile(i Interpreter, args []env.Value) (env.Value, error) {
	fi, err := os.Stat(slices.Fst(args).String())
	if err != nil {
		return env.False(), err
	}
	return env.Bool(fi.Mode().IsRegular()), nil
}

func fileIsExecutable(i Interpreter, args []env.Value) (env.Value, error) {
	ok := ioutil.Executable(slices.Fst(args).String(), os.Getuid())
	return env.Bool(ok), nil
}

func fileOwned(i Interpreter, args []env.Value) (env.Value, error) {
	ok := ioutil.Owned(slices.Fst(args).String(), os.Getuid())
	return env.Bool(ok), nil
}

func fileReadable(i Interpreter, args []env.Value) (env.Value, error) {
	ok := ioutil.Readable(slices.Fst(args).String(), os.Getuid())
	return env.Bool(ok), nil
}

func fileWritable(i Interpreter, args []env.Value) (env.Value, error) {
	ok := ioutil.Writable(slices.Fst(args).String(), os.Getuid())
	return env.Bool(ok), nil
}

func fileSize(i Interpreter, args []env.Value) (env.Value, error) {
	fi, err := os.Stat(slices.Fst(args).String())
	if err != nil {
		return nil, err
	}
	return env.Int(fi.Size()), nil
}

func fileModTime(i Interpreter, args []env.Value) (env.Value, error) {
	fi, err := os.Stat(slices.Fst(args).String())
	if err != nil {
		return nil, err
	}
	return env.Int(fi.ModTime().Unix()), nil
}

func fileType(i Interpreter, args []env.Value) (env.Value, error) {
	typ, err := ioutil.Type(slices.Fst(args).String())
	return env.Str(typ), err
}

func fileStat(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func fileJoin(i Interpreter, args []env.Value) (env.Value, error) {
	var list []string
	for _, a := range args {
		vs, err := env.ToStringList(a)
		if err != nil {
			return nil, err
		}
		list = append(list, vs...)
	}
	str := filepath.Join(list...)
	return env.Str(str), nil
}

func fileNormalize(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

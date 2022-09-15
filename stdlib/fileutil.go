package stdlib

import (
	"bufio"
	"os"
	"strings"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/slices"
)

func RunFileForeachLine() Executer {
	return Builtin{
		Name:  "foreachLine",
		Arity: 3,
		Safe:  true,
		Run:   fileutilForeachLine,
	}
}

func RunFileCat() Executer {
	return Builtin{
		Name:  "cat",
		Arity: 1,
		Safe:  true,
		Run:   fileutilCat,
	}
}

func fileutilCat(i Interpreter, args []env.Value) (env.Value, error) {
	bs, err := os.ReadFile(slices.Fst(args).String())
	return env.Str(string(bs)), err
}

func fileutilForeachLine(i Interpreter, args []env.Value) (env.Value, error) {
	r, err := os.Open(slices.Fst(args).String())
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var (
		scan = bufio.NewScanner(r)
		res  env.Value
	)
	for scan.Scan() {
		i.Define(slices.Snd(args).String(), env.Str(scan.Text()))
		res, err = i.Execute(strings.NewReader(slices.Lst(args).String()))
		if err != nil {
			break
		}
	}
	if res == nil {
		res = env.EmptyStr()
	}
	return res, scan.Err()
}

package stdlib

import (
	"fmt"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/slices"
)

func RunList() Executer {
	return Builtin{
		Name:  "list",
		Arity: 1,
		Safe:  true,
		Run:   runList,
	}
}

func RunSplit() Executer {
	return Builtin{
		Name:  "split",
		Arity: 1,
		Safe:  true,
	}
}

func RunLLength() Executer {
	return Builtin{
		Name:  "llength",
		Arity: 1,
		Safe:  true,
		Run:   listLength,
	}
}

func RunLSet() Executer {
	return Builtin{
		Name:  "lset",
		Arity: 1,
		Safe:  true,
		Run:   listSet,
	}
}

func RunLSort() Executer {
	return Builtin{
		Name:  "lsort",
		Arity: 1,
		Safe:  true,
		Run:   listSort,
	}
}

func RunLSearch() Executer {
	return Builtin{
		Name:  "lsearch",
		Arity: 1,
		Safe:  true,
		Run:   listSearch,
	}
}

func RunLReverse() Executer {
	return Builtin{
		Name:  "lreverse",
		Arity: 1,
		Safe:  true,
		Run:   listReverse,
	}
}

func RunLReplace() Executer {
	return Builtin{
		Name:  "lreplace",
		Arity: 1,
		Safe:  true,
		Run:   listReplace,
	}
}

func RunLRepeat() Executer {
	return Builtin{
		Name:     "lrepeat",
		Arity:    2,
		Variadic: true,
		Safe:     true,
		Run:      listRepeat,
	}
}

func RunLIndex() Executer {
	return Builtin{
		Name:  "lindex",
		Arity: 2,
		Safe:  true,
		Run:   listIndex,
	}
}

func RunLMap() Executer {
	return Builtin{
		Name:     "lmap",
		Arity:    3,
		Variadic: true,
		Safe:     true,
		Run:      listMap,
	}
}

func RunLRange() Executer {
	return Builtin{
		Name:  "lrange",
		Arity: 3,
		Safe:  true,
		Run:   listRange,
	}
}

func RunLAssign() Executer {
	return Builtin{
		Name:  "lassign",
		Arity: 1,
		Safe:  true,
		Run:   listAssign,
	}
}

func RunLAppend() Executer {
	return Builtin{
		Name:  "lappend",
		Arity: 1,
		Safe:  true,
		Run:   listAppend,
	}
}

func RunLInsert() Executer {
	return Builtin{
		Name:  "linsert",
		Arity: 1,
		Safe:  true,
		Run:   listInsert,
	}
}

func runList(i Interpreter, args []env.Value) (env.Value, error) {
	return slices.Fst(args).ToList()
}

func listLength(i Interpreter, args []env.Value) (env.Value, error) {
	list, err := slices.Fst(args).ToList()
	if err != nil {
		return nil, err
	}
	n, ok := list.(interface{ Len() int })
	if !ok {
		return env.Int(0), nil
	}
	return env.Int(int64(n.Len())), nil
}

func listInsert(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func listAssign(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func listAppend(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func listMap(i Interpreter, args []env.Value) (env.Value, error) {
	if len(args)%3 != 0 {
		return nil, fmt.Errorf("invalid number of arguments given")
	}
	return nil, nil
}

func listRange(i Interpreter, args []env.Value) (env.Value, error) {
	var (
		fst, err1 = env.ToInt(slices.Snd(args))
		lst, err2 = env.ToInt(slices.Lst(args))
	)
	if err := hasError(err1, err2); err != nil {
		return nil, err
	}
	return env.Range(slices.Fst(args), fst, lst+1)
}

func listIndex(i Interpreter, args []env.Value) (env.Value, error) {
	n, err := env.ToInt(slices.Snd(args))
	if err != nil {
		return nil, err
	}
	return env.At(slices.Fst(args), n)
}

func listRepeat(i Interpreter, args []env.Value) (env.Value, error) {
	n, err := env.ToInt(slices.Fst(args))
	if err != nil {
		return nil, err
	}
	var (
		list []env.Value
		rest = slices.Rest(args)
	)
	for i := 0; i < n; i++ {
		list = append(list, rest...)
	}
	return env.ListFrom(list...), nil
}

func listReplace(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func listReverse(i Interpreter, args []env.Value) (env.Value, error) {
	return env.Reverse(slices.Fst(args))
}

func listSearch(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func listSort(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func listSet(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

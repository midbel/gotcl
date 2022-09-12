package stdlib

import (
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
	}
}

func RunLSort() Executer {
	return Builtin{
		Name:  "lsort",
		Arity: 1,
		Safe:  true,
	}
}

func RunLSearch() Executer {
	return Builtin{
		Name:  "lsearch",
		Arity: 1,
		Safe:  true,
	}
}

func RunLReverse() Executer {
	return Builtin{
		Name:  "lreverse",
		Arity: 1,
		Safe:  true,
	}
}

func RunLReplace() Executer {
	return Builtin{
		Name:  "lreplace",
		Arity: 1,
		Safe:  true,
	}
}

func RunLRepeat() Executer {
	return Builtin{
		Name:  "lrepeat",
		Arity: 1,
		Safe:  true,
	}
}

func RunLIndex() Executer {
	return Builtin{
		Name:  "lindex",
		Arity: 1,
		Safe:  true,
	}
}

func RunLMap() Executer {
	return Builtin{
		Name:  "lmap",
		Arity: 1,
		Safe:  true,
	}
}

func RunLRange() Executer {
	return Builtin{
		Name:  "lrange",
		Arity: 1,
		Safe:  true,
	}
}

func RunLAssign() Executer {
	return Builtin{
		Name:  "lassign",
		Arity: 1,
		Safe:  true,
	}
}

func RunLAppend() Executer {
	return Builtin{
		Name:  "lappend",
		Arity: 1,
		Safe:  true,
	}
}

func RunLInsert() Executer {
	return Builtin{
		Name:  "linsert",
		Arity: 1,
		Safe:  true,
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

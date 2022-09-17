package stdlib

import (
	"github.com/midbel/gotcl/env"
)

func RunPrefixMatch() Executer {
	return Builtin{
		Name: "match",
		Safe: true,
		Run:  prefixMatch,
	}
}

func RunPrefixLongest() Executer {
	return Builtin{
		Name: "longest",
    Arity: 2,
		Safe: true,
		Run:  prefixLongest,
	}
}

func RunPrefixAll() Executer {
	return Builtin{
		Name: "all",
    Arity: 2,
		Safe: true,
		Run:  prefixAll,
	}
}

func prefixMatch(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func prefixLongest(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func prefixAll(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

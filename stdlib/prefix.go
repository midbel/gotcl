package stdlib

import (
	"fmt"
	"strings"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/gotcl/stdlib/strutil"
	"github.com/midbel/slices"
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
		Name:  "longest",
		Arity: 2,
		Safe:  true,
		Run:   prefixLongest,
	}
}

func RunPrefixAll() Executer {
	return Builtin{
		Name:  "all",
		Arity: 2,
		Safe:  true,
		Run:   prefixAll,
	}
}

func prefixMatch(i Interpreter, args []env.Value) (env.Value, error) {
	all, err := getPrefixMatches(slices.Fst(args), slices.Snd(args))
	if err != nil {
		return nil, err
	}
	if len(all) > 1 {
		return nil, fmt.Errorf("too many results returned")
	}
	if len(all) == 0 {
		return nil, fmt.Errorf("empty match")
	}
	return env.Str(all[0]), nil
}

func prefixLongest(i Interpreter, args []env.Value) (env.Value, error) {
	all, err := getPrefixMatches(slices.Fst(args), slices.Snd(args))
	if err != nil {
		return nil, err
	}
	prefix := strutil.LongestCommonPrefix(all)
	return env.Str(prefix), nil
}

func prefixAll(i Interpreter, args []env.Value) (env.Value, error) {
	all, err := getPrefixMatches(slices.Fst(args), slices.Snd(args))
	if err != nil {
		return nil, err
	}
	return env.ListFromStrings(all), nil
}

func getPrefixMatches(list, prefix env.Value) ([]string, error) {
	list, err := list.ToList()
	if err != nil {
		return nil, err
	}
	var match []string
	for _, v := range list.(env.List).Values() {
		if !strings.HasPrefix(v.String(), prefix.String()) {
			continue
		}
		match = append(match, v.String())
	}
	return match, nil
}

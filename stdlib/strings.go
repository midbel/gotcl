package stdlib

import (
	"sort"
	"strconv"
	"strings"

	"github.com/midbel/gotcl/env"
)

func MakeString() Executer {
	e := Ensemble{
		Name: "string",
		List: []Executer{
			Builtin{
				Name:  "tolower",
				Arity: 1,
				Run: func(i Interpreter, args []env.Value) (env.Value, error) {
					return withString(slices.Fst(args), strings.ToLower)
				},
			},
			Builtin{
				Name:  "toupper",
				Arity: 1,
				Run: func(i Interpreter, args []env.Value) (env.Value, error) {
					return withString(slices.Fst(args), strings.ToUpper)
				},
			},
			Builtin{
				Name:  "length",
				Arity: 1,
				Run: func(i Interpreter, args []env.Value) (env.Value, error) {
					return withString(slices.Fst(args), func(s string) string {
						return strconv.Itoa(len(s))
					})
				},
			},
			Builtin{
				Name:  "repeat",
				Arity: 2,
				Run: func(i Interpreter, args []env.Value) (env.Value, error) {
					c, err := env.ToInt(slices.Snd(args))
					if err != nil {
						return nil, err
					}
					return withString(slices.Fst(args), func(s string) string {
						return strings.Repeat(s, c)
					})
				},
			},
		},
	}
	sort.Slice(e.List, func(i, j int) bool {
		return getName(e.List[i]) < getName(e.List[j])
	})
	return e
}

func withString(v env.Value, do func(str string) string) (env.Value, error) {
	str, err := v.ToString()
	if err != nil {
		return nil, err
	}
	return env.Str(do(str.String())), nil
}

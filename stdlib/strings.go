package stdlib

import (
	"strconv"
	"strings"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/slices"
)

func MakeString() Executer {
	e := Ensemble{
		Name: "string",
		List: []Executer{
			Builtin{
				Name:  "tolower",
				Arity: 1,
				Run:   stringToLower,
			},
			Builtin{
				Name:  "toupper",
				Arity: 1,
				Run:   stringToUpper,
			},
			Builtin{
				Name:  "length",
				Arity: 1,
				Run:   stringLength,
			},
			Builtin{
				Name:  "repeat",
				Arity: 2,
				Run:   stringRepeat,
			},
		},
	}
	return sortEnsembleCommands(e)
}

func stringToLower(i Interpreter, args []env.Value) (env.Value, error) {
	return withString(slices.Fst(args), strings.ToLower)
}

func stringToUpper(i Interpreter, args []env.Value) (env.Value, error) {
	return withString(slices.Fst(args), strings.ToUpper)
}

func stringLength(i Interpreter, args []env.Value) (env.Value, error) {
	return withString(slices.Fst(args), func(s string) string {
		return strconv.Itoa(len(s))
	})
}

func stringRepeat(i Interpreter, args []env.Value) (env.Value, error) {
	c, err := env.ToInt(slices.Snd(args))
	if err != nil {
		return nil, err
	}
	return withString(slices.Fst(args), func(s string) string {
		return strings.Repeat(s, c)
	})
}

func withString(v env.Value, do func(str string) string) (env.Value, error) {
	str, err := v.ToString()
	if err != nil {
		return nil, err
	}
	return env.Str(do(str.String())), nil
}

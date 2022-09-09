package stdlib

import (
	"strconv"
	"strings"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/gotcl/glob"
	"github.com/midbel/gotcl/stdlib/strutil"
	"github.com/midbel/slices"
)

func MakeString() Executer {
	e := Ensemble{
		Name: "string",
		List: []Executer{
			Builtin{
				Name:     "cat",
				Variadic: true,
				Run:      stringCat,
			},
			Builtin{
				Name:     "replace",
				Arity:    3,
				Variadic: true,
				Run:      stringReplace,
			},
			Builtin{
				Name:     "trim",
				Arity:    1,
				Variadic: true,
				Run:      stringTrim,
			},
			Builtin{
				Name:     "trimleft",
				Arity:    1,
				Variadic: true,
				Run:      stringTrimLeft,
			},
			Builtin{
				Name:     "trimright",
				Arity:    1,
				Variadic: true,
				Run:      stringTrimRight,
			},
			Builtin{
				Name:  "equal",
				Arity: 2,
				Options: []Option{
					{
						Name:  "nocase",
						Flag:  true,
						Value: env.False(),
						Check: CheckBool,
					},
					{
						Name:  "length",
						Value: env.Zero(),
						Check: CheckNumber,
					},
				},
				Run: stringEqual,
			},
			Builtin{
				Name:  "compare",
				Arity: 2,
				Options: []Option{
					{
						Name:  "nocase",
						Flag:  true,
						Value: env.False(),
						Check: CheckBool,
					},
					{
						Name:  "length",
						Value: env.Zero(),
						Check: CheckNumber,
					},
				},
				Run: stringCompare,
			},
			Builtin{
				Name: "first",
			},
			Builtin{
				Name: "last",
			},
			Builtin{
				Name: "index",
			},
			Builtin{
				Name: "range",
			},
			Builtin{
				Name:  "map",
				Arity: 2,
				Options: []Option{
					{
						Name:  "nocase",
						Flag:  true,
						Value: env.False(),
						Check: CheckBool,
					},
				},
				Run: stringMatch,
			},
			Builtin{
				Name:  "reverse",
				Arity: 1,
				Run:   stringReverse,
			},
			Builtin{
				Name:  "length",
				Arity: 1,
				Run:   stringLength,
			},
			Builtin{
				Name:  "match",
				Arity: 2,
				Options: []Option{
					{
						Name:  "nocase",
						Flag:  true,
						Value: env.False(),
						Check: CheckBool,
					},
				},
				Run: stringMatch,
			},
			Builtin{
				Name:     "totitle",
				Arity:    1,
				Variadic: true,
				Run:      stringToTitle,
			},
			Builtin{
				Name:     "tolower",
				Arity:    1,
				Variadic: true,
				Run:      stringToLower,
			},
			Builtin{
				Name:     "toupper",
				Arity:    1,
				Variadic: true,
				Run:      stringToUpper,
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

func stringEqual(i Interpreter, args []env.Value) (env.Value, error) {
	var (
		length, _ = i.Resolve("length")
		nocase, _ = i.Resolve("nocase")
		first     = cutStr(slices.Fst(args).String(), length)
		last      = cutStr(slices.Snd(args).String(), length)
	)
	if no := env.ToBool(nocase); no {
		first = strings.ToLower(first)
		last = strings.ToLower(last)
	}
	return env.Bool(first == last), nil
}

func stringCompare(i Interpreter, args []env.Value) (env.Value, error) {
	var (
		length, _ = i.Resolve("length")
		nocase, _ = i.Resolve("nocase")
		first     = cutStr(slices.Fst(args).String(), length)
		last      = cutStr(slices.Snd(args).String(), length)
	)
	if no := env.ToBool(nocase); no {
		first = strings.ToLower(first)
		last = strings.ToLower(last)
	}
	res := strings.Compare(first, last)
	return env.Int(int64(res)), nil
}

func stringReverse(i Interpreter, args []env.Value) (env.Value, error) {
	str := strutil.Reverse(slices.Fst(args).String())
	return env.Str(str), nil
}

func stringMatch(i Interpreter, args []env.Value) (env.Value, error) {
	var (
		str       = slices.Fst(args).String()
		pat       = slices.Snd(args).String()
		nocase, _ = i.Resolve("nocase")
	)
	if env.ToBool(nocase) {
		str = strings.ToLower(str)
	}
	return env.Bool(glob.Match(str, pat)), nil
}

func stringMap(i Interpreter, args []env.Value) (env.Value, error) {
	var (
		str       = slices.Fst(args).String()
		list, _   = env.ToStringList(slices.Snd(args))
		nocase, _ = i.Resolve("nocase")
	)
	str, err := strutil.Map(str, list, env.ToBool(nocase))
	if err != nil {
		return nil, err
	}
	return env.Str(str), nil
}

func stringCat(i Interpreter, args []env.Value) (env.Value, error) {
	var str strings.Builder
	for i := range args {
		str.WriteString(args[i].String())
	}
	return env.Str(str.String()), nil
}

func stringReplace(i Interpreter, args []env.Value) (env.Value, error) {
	first, err := env.ToInt(slices.At(args, 1))
	if err != nil {
		return nil, err
	}
	last, err := env.ToInt(slices.At(args, 2))
	if err != nil {
		return nil, err
	}
	var replace string
	if v := slices.At(args, 4); v != nil {
		replace = v.String()
	}
	str, err := strutil.Replace(slices.Fst(args).String(), replace, first, last)
	return env.Str(str), err
}

func stringToLower(i Interpreter, args []env.Value) (env.Value, error) {
	var (
		res env.Value
		str = slices.Fst(args).String()
	)
	fst, lst, err := getRange(str, slices.At(args, 1), slices.At(args, 2))
	if err != nil {
		return nil, err
	}
	str, err = strutil.ToLower(str, fst, lst)
	if err == nil {
		res = env.Str(str)
	}
	return res, nil
}

func stringToUpper(i Interpreter, args []env.Value) (env.Value, error) {
	var (
		res env.Value
		str = slices.Fst(args).String()
	)
	fst, lst, err := getRange(str, slices.At(args, 1), slices.At(args, 2))
	if err != nil {
		return nil, err
	}
	str, err = strutil.ToLower(str, fst, lst)
	if err == nil {
		res = env.Str(str)
	}
	return res, nil
}

func stringToTitle(i Interpreter, args []env.Value) (env.Value, error) {
	var (
		res env.Value
		str = slices.Fst(args).String()
	)
	fst, lst, err := getRange(str, slices.At(args, 1), slices.At(args, 2))
	if err != nil {
		return nil, err
	}
	str, err = strutil.ToLower(str, fst, lst)
	if err == nil {
		res = env.Str(str)
	}
	return res, nil
}

func stringTrim(i Interpreter, args []env.Value) (env.Value, error) {
	var chars string
	if v := slices.Snd(args); v != nil {
		chars = v.String()
	}
	return withString(slices.Fst(args), func(s string) string {
		return strings.Trim(s, chars)
	})
}

func stringTrimLeft(i Interpreter, args []env.Value) (env.Value, error) {
	var chars string
	if v := slices.Snd(args); v != nil {
		chars = v.String()
	}
	return withString(slices.Fst(args), func(s string) string {
		return strings.TrimLeft(s, chars)
	})
}

func stringTrimRight(i Interpreter, args []env.Value) (env.Value, error) {
	var chars string
	if v := slices.Snd(args); v != nil {
		chars = v.String()
	}
	return withString(slices.Fst(args), func(s string) string {
		return strings.TrimRight(s, chars)
	})
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

func getRange(str string, fst, lst env.Value) (int, int, error) {
	var (
		min int
		max = len(str)
		err error
	)
	if fst != nil {
		min, err = env.ToInt(fst)
	}
	if lst != nil {
		max, err = env.ToInt(lst)
	}
	return min, max, err
}

func cutStr(str string, v env.Value) string {
	n, err := env.ToInt(v)
	if err != nil {
		return str
	}
	if n > 0 && len(str) >= n {
		str = str[:n]
	}
	return str
}

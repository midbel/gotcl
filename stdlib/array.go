package stdlib

import (
	"fmt"
	"sort"
	"strings"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/gotcl/word"
	"github.com/midbel/slices"
)

func MakeArray() Executer {
	e := Ensemble{
		Name: "array",
		List: []Executer{
			Builtin{
				Name:  "set",
				Arity: 2,
				Run: func(i Interpreter, args []env.Value) (env.Value, error) {
					arr, err := i.Resolve(slices.Fst(args).String())
					if err != nil {
						arr = env.EmptyArr()
					}
					list, err := scan(slices.Snd(args).String())
					if err != nil {
						return nil, err
					}
					if len(list)%2 != 0 {
						return nil, fmt.Errorf("invalid length")
					}
					s := arr.(env.Array)
					for i := 0; i < len(list); i += 2 {
						s.Set(list[i], env.Str(list[i+1]))
					}
					i.Define(slices.Fst(args).String(), s)
					return nil, nil
				},
			},
			Builtin{
				Name:  "get",
				Arity: 1,
				Run: func(i Interpreter, args []env.Value) (env.Value, error) {
					arr, err := i.Resolve(slices.Fst(args).String())
					if err != nil {
						return nil, err
					}
					arr, err = arr.ToArray()
					if err != nil {
						return nil, err
					}
					var (
						g  = arr.(env.Array)
						vs []env.Value
					)
					for k, v := range g.values {
						vs = append(vs, env.ListFrom(env.Str(k), v))
					}
					return env.ListFrom(vs...), nil
				},
			},
			Builtin{
				Name:  "names",
				Arity: 1,
				Run: func(i Interpreter, args []env.Value) (env.Value, error) {
					arr, err := i.Resolve(slices.Fst(args).String())
					if err != nil {
						return nil, err
					}
					arr, err = arr.ToArray()
					if err != nil {
						return nil, err
					}
					var (
						g  = arr.(env.Array)
						vs []string
					)
					for k := range g.values {
						vs = append(vs, k)
					}
					return env.ListFromStrings(vs), nil
				},
			},
			Builtin{
				Name:  "size",
				Arity: 1,
				Run: func(i Interpreter, args []env.Value) (env.Value, error) {
					return nil, nil
				},
			},
		},
	}
	sort.Slice(e.List, func(i, j int) bool {
		return e.List[i].GetName() < e.List[j].GetName()
	})
	return e
}

func PrintArray() Executer {
	return Builtin{
		Name:  "parray",
		Arity: 1,
		Safe:  true,
		Run: func(i Interpreter, args []env.Value) (env.Value, error) {
			arr, err := i.Resolve(slices.Fst(args).String())
			if err != nil {
				return nil, err
			}
			arr, err = arr.ToArray()
			if err != nil {
				return nil, err
			}
			vs := arr.(env.Array)
			for k, v := range vs.values {
				fmt.Fprintf(i.Out, "%s(%s) = %s", slices.Fst(args), k, v)
				fmt.Fprintln(i.Out)
			}
			return nil, nil
		},
	}
}

func scan(str string) ([]string, error) {
	s, err := word.Scan(strings.NewReader(str))
	if err != nil {
		return nil, err
	}
	var list []string
	for {
		w := s.Scan()
		if w.Type == word.Illegal {
			return nil, fmt.Errorf("illegal token")
		}
		if w.Type == word.EOF {
			break
		}
		if w.Literal != "" {
			list = append(list, w.Literal)
		}
	}
	return list, nil
}

package stdlib

import (
	"fmt"
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
				Run:   arraySet,
			},
			Builtin{
				Name:  "get",
				Arity: 1,
				Run:   arrayGet,
			},
			Builtin{
				Name:  "names",
				Arity: 1,
				Run:   arrayNames,
			},
			Builtin{
				Name:  "size",
				Arity: 1,
				Run:   arraySize,
			},
		},
	}
	return sortEnsembleCommands(e)
}

func PrintArray() Executer {
	return Builtin{
		Name:  "parray",
		Arity: 1,
		Safe:  true,
		Run:   printArray,
	}
}

func printArray(i Interpreter, args []env.Value) (env.Value, error) {
	arr, err := i.Resolve(slices.Fst(args).String())
	if err != nil {
		return nil, err
	}
	arr, err = arr.ToArray()
	if err != nil {
		return nil, err
	}
	ph, ok := i.(PrintHandler)
	if !ok {
		return nil, fmt.Errorf("interpreter can not print array to channel")
	}
	vs := arr.(env.Array)
	for _, n := range vs.Names() {
		msg := fmt.Sprintf("%s(%s) = %s", slices.Fst(args), n, vs.Get(n))
		ph.Println("stdout", msg)
	}
	return nil, nil
}

func arrayNames(i Interpreter, args []env.Value) (env.Value, error) {
	arr, err := i.Resolve(slices.Fst(args).String())
	if err != nil {
		return nil, err
	}
	arr, err = arr.ToArray()
	if err != nil {
		return nil, err
	}
	list := arr.(env.Array).Names()
	return env.ListFromStrings(list), nil
}

func arraySize(i Interpreter, args []env.Value) (env.Value, error) {
	arr, err := slices.Fst(args).ToArray()
	if err != nil {
		return nil, err
	}
	z, ok := arr.(interface{ Len() int })
	if !ok {
		return nil, fmt.Errorf("%s is not an array", slices.Fst(args))
	}
	return env.Int(int64(z.Len())), nil
}

func arrayGet(i Interpreter, args []env.Value) (env.Value, error) {
	arr, err := i.Resolve(slices.Fst(args).String())
	if err != nil {
		return nil, err
	}
	arr, err = arr.ToArray()
	if err != nil {
		return nil, err
	}
	return arr.(env.Array).Pairs(), nil
}

func arraySet(i Interpreter, args []env.Value) (env.Value, error) {
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

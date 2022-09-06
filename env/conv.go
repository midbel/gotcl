package env

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/midbel/gotcl/word"
)

func ToStringList(v Value) ([]string, error) {
	v, err := v.ToList()
	if err != nil {
		return nil, err
	}
	x, ok := v.(List)
	if !ok {
		return nil, nil
	}
	var list []string
	for i := range x.values {
		list = append(list, x.values[i].String())
	}
	return list, nil
}

func ToLevel(v Value) (int, bool, error) {
	var (
		str = v.String()
		abs bool
	)
	if strings.HasPrefix(str, "#") {
		abs = true
		str = strings.TrimPrefix(str, "#")
	}
	lvl, err := strconv.Atoi(str)
	return lvl, abs, err
}

func ToInt(v Value) (int, error) {
	n, err := v.ToNumber()
	if err != nil {
		return 0, err
	}
	x, ok := n.(Number)
	if !ok {
		return 0, nil
	}
	return int(x.value), nil
}

func ToFloat(v Value) (float64, error) {
	n, err := v.ToNumber()
	if err != nil {
		return 0, err
	}
	x, ok := n.(Number)
	if !ok {
		return 0, nil
	}
	return x.value, nil
}

func ToBool(v Value) bool {
	v, err := v.ToBoolean()
	if err != nil {
		return false
	}
	b, ok := v.(Boolean)
	if !ok {
		return ok
	}
	return b.value
}

func split(str string) (Value, error) {
	str = strings.TrimSpace(str)
	scan, err := word.Scan(strings.NewReader(str))
	if err != nil {
		return nil, err
	}
	var list []Value
	for {
		w := scan.Scan()
		if w.Type == word.EOF {
			break
		}
		if w.Type == word.Blank {
			continue
		}
		switch w.Type {
		case word.Literal:
		case word.Block:
			w.Literal = fmt.Sprintf("{%s}", w.Literal)
		case word.Variable:
			w.Literal = fmt.Sprintf("$%s", w.Literal)
		case word.Script:
			w.Literal = fmt.Sprintf("[%s]", w.Literal)
		case word.Quote:
			w.Literal = fmt.Sprintf("\"%s\"", w.Literal)
		default:
			return nil, fmt.Errorf("%s: unsupported token type", w)
		}
		list = append(list, Str(w.Literal))
	}
	if len(list) == 1 {
		return list[0], nil
	}
	return ListFrom(list...), nil
}

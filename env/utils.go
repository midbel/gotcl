package env

import (
	"fmt"

	"github.com/midbel/slices"
)

func Reverse(list Value) (Value, error) {
	ls, err := getListFromValue(list)
	if err != nil {
		return nil, err
	}
	ls.values = slices.Reverse(ls.values)
	return ListFrom(ls.values...), nil
}

func At(list Value, at int) (Value, error) {
	ls, err := getListFromValue(list)
	if err != nil {
		return nil, err
	}
	if at < 0 || at >= len(ls.values) {
		return nil, fmt.Errorf("index out of range")
	}
	return ls.values[at], nil
}

func Range(list Value, first, last int) (Value, error) {
	ls, err := getListFromValue(list)
	if err != nil {
		return nil, err
	}
	if first < 0 || first >= ls.Len() || last < 0 || last >= ls.Len() || first > last {
		return nil, fmt.Errorf("invalid range given: %d - %d", first, last)
	}
	return ListFrom(ls.values[first : last]...), nil
}

func getListFromValue(list Value) (List, error) {
	list, err := list.ToList()
	if err != nil {
		return List{}, err
	}
	ls, ok := list.(List)
	if !ok {
		return ls, fmt.Errorf("a list should be provided, got %T", list)
	}
	return ls, nil
}

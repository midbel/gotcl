package argparse

import (
	"errors"
	"fmt"
	"sort"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/slices"
)

var ErrType = errors.New("invalid type given")

type CheckFunc func(env.Value) error

type Option struct {
	env.Value
	Name     string
	Help     string
	Flag     bool
	Required bool
	Check    CheckFunc
}

func CheckBool(v env.Value) error {
	_, ok := v.(env.Boolean)
	if !ok {
		return ErrType
	}
	return nil
}

func CheckNumber(v env.Value) error {
	_, ok := v.(env.Number)
	if !ok {
		return ErrType
	}
	return nil
}

func CheckString(v env.Value) error {
	_, ok := v.(env.String)
	if !ok {
		return ErrType
	}
	return nil
}

func CombineCheck(cs ...CheckFunc) CheckFunc {
	return func(v env.Value) error {
		for i := range cs {
			if err := cs[i](v); err != nil {
				return err
			}
		}
		return nil
	}
}

func IsValid(list []Option) error {
	ok := slices.Every(list, func(o Option) bool {
		if !o.Required {
			return true
		}
		return o.Required && o.Value != nil
	})
	if !ok {
		return fmt.Errorf("required options are not provided!")
	}
	return nil
}

func IsSet(list []Option, name string) (int, error) {
	x := sort.Search(len(list), func(i int) bool {
		return list[i].Name >= name
	})
	if x < len(list) && list[x].Name == name {
		return x, nil
	}
	return 0, fmt.Errorf("%s: option not supported", name)
}

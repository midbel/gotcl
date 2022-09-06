package env

import (
	"errors"
	"fmt"
)

var ErrUndefined = errors.New("undefined variable")

type Env struct {
	values map[string]Value
}

func EmptyEnv() *Env {
	return &Env{
		values: make(map[string]Value),
	}
}

func (e *Env) Delete(n string) {
	delete(e.values, n)
}

func (e *Env) Define(n string, v Value) {
	e.values[n] = v
}

func (e *Env) Resolve(n string) (Value, error) {
	v, ok := e.values[n]
	if !ok {
		return nil, fmt.Errorf("%s: %w", n, ErrUndefined)
	}
	return v, nil
}

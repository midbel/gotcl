package env

import (
	"errors"
	"fmt"
)

var ErrUndefined = errors.New("variable not defined")

type Environment interface {
	Resolve(string) (string, error)
	Define(string, string) error
	Delete(string) error

	Exists(string) bool
}

type env struct {
	parent Environment
	values map[string]string
}

func EnclosedEnv(parent Environment) Environment {
	return &env{
		parent: parent,
		values: make(map[string]string),
	}
}

func EmptyEnv() Environment {
	return EnclosedEnv(nil)
}

func (e *env) Exists(name string) bool {
	_, ok := e.values[name]
	if !ok && e.parent != nil {
		return e.parent.Exists(name)
	}
	return ok
}

func (e *env) Define(name, value string) error {
	if _, ok := e.values[name]; !ok && e.parent != nil {
		return e.parent.Define(name, value)
	}
	e.values[name] = value
	return nil
}

func (e *env) Delete(name string) error {
	delete(e.values, name)
	return nil
}

func (e *env) Resolve(name string) (string, error) {
	val, ok := e.values[name]
	if ok {
		return val, nil
	}
	if e.parent == nil {
		return "", fmt.Errorf("%s: %w", name, ErrUndefined)
	}
	return e.parent.Resolve(name)
}

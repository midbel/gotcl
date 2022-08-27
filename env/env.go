package env

import (
	"errors"
	"fmt"
)

var (
	ErrUndefined = errors.New("variable not defined")
	ErrForbidden = errors.New("operation not allowed")
)

type Environment interface {
	Resolve(string) (string, error)
	Define(string, string) error
	Delete(string) error

	Exists(string) bool
	IsSet(string) bool
}

type Linker interface {
	LinkVar(Environment, string, string) error
	getVar(string) (*variable, error)
}

type env struct {
	parent Environment
	values map[string]*variable
}

func EnclosedEnv(parent Environment) Environment {
	return &env{
		parent: parent,
		values: make(map[string]*variable),
	}
}

func EmptyEnv() Environment {
	return EnclosedEnv(nil)
}

func (e *env) All() []string {
	list := make([]string, 0, len(e.values))
	for n := range e.values {
		list = append(list, n)
	}
	return list
}

func (e *env) Exists(name string) bool {
	_, ok := e.values[name]
	if !ok && e.parent != nil {
		return e.parent.Exists(name)
	}
	return ok
}

func (e *env) IsSet(name string) bool {
	v, ok := e.values[name]
	if !ok && e.parent != nil {
		e.parent.IsSet(name)
	}
	return v != nil
}

func (e *env) Define(name, value string) error {
	if !e.IsSet(name) {
		e.values[name] = &variable{
			value:    value,
			refcount: 1,
		}
	} else {
		e.values[name].value = value
	}
	return nil
}

func (e *env) Delete(name string) error {
	delete(e.values, name)
	return nil
}

func (e *env) Resolve(name string) (string, error) {
	val, ok := e.values[name]
	if ok {
		return val.value, nil
	}
	if e.parent == nil {
		return "", undefinedVar(name)
	}
	return e.parent.Resolve(name)
}

func (e *env) LinkVar(ev Environment, dst, src string) error {
	other, ok := ev.(Linker)
	if !ok {
		return ErrForbidden
	}
	var err error
	e.values[dst], err = other.getVar(src)
	return err
}

func (e *env) getVar(name string) (*variable, error) {
	if !e.IsSet(name) {
		return nil, fmt.Errorf("%s: %w", name, ErrUndefined)
	}
	return e.values[name], nil
}

type variable struct {
	value    string
	refcount int
}

func undefinedVar(name string) error {
	return fmt.Errorf("%s: %w", name, ErrUndefined)
}

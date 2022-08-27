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
	values map[string]value
}

func EnclosedEnv(parent Environment) Environment {
	return &env{
		parent: parent,
		values: make(map[string]value),
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
		e.values[name] = createVariable(value)
	} else {
		e.values[name].Set("", value)
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
		return val.Get("")
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
	val, err := other.getVar(src)
	if err == nil {
		val.Up()
		e.values[dst] = val
	}
	return err
}

func (e *env) getVar(name string) (value, error) {
	if !e.IsSet(name) {
		return nil, fmt.Errorf("%s: %w", name, ErrUndefined)
	}
	return e.values[name], nil
}

func undefinedVar(name string) error {
	return fmt.Errorf("%s: %w", name, ErrUndefined)
}

type value interface {
	Get(string) (string, error)
	Set(k, v string) error
	Up()
}

type variable struct {
	value    string
	refcount int
}

func createVariable(value string) *variable {
	return &variable{
		value:    value,
		refcount: 1,
	}
}

func (v *variable) Up() {
	v.refcount++
}

func (v *variable) Get(_ string) (string, error) {
	return v.value, nil
}

func (v *variable) Set(_, s string) error {
	v.value = s
	return nil
}

type array struct {
	values   []string
	refcount int
}

func createArray() *array {
	return &array{
		refcount: 1,
	}
}

func (a *array) Up() {
	a.refcount++
}

func (a *array) Get(i string) (string, error) {
	return "", nil
}

type dict struct {
	values   map[string]string
	refcount int
}

func createDict() *dict {
	return &dict{
		values:   make(map[string]string),
		refcount: 1,
	}
}

func (d *dict) Up() {
	d.refcount++
}

func (d *dict) Get(k string) (string, error) {
	return "", nil
}

func (d *dict) Set(k, v string) error {
	return nil
}

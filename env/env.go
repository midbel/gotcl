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
	Link(string) error
	LinkAt(string, string, int) error
	Delete(string) error

	Exists(string) bool
	IsSet(string) bool
}

type env struct {
	parent Environment
	values map[string]*variable
	level  int
}

func EnclosedEnv(parent Environment) Environment {
	var level int
	if e, ok := parent.(*env); ok {
		level = e.level + 1
	}
	return &env{
		parent: parent,
		values: make(map[string]*variable),
		level:  level,
	}
}

func EmptyEnv() Environment {
	return EnclosedEnv(nil)
}

func (e *env) Unlink(name string) error {
	v, ok := e.values[name]
	if !ok {
		return undefinedVar(name)
	}
	v.refcount--
	delete(e.values, name)
	return nil
}

func (e *env) Link(name string) error {
	return e.LinkAt(name, name, -1)
}

func (e *env) LinkAt(name, alias string, level int) error {
	if alias == "" {
		alias = name
	}
	if e.isSet(name, false) || e.isSet(alias, false) {
		return fmt.Errorf("%s (as %s): variable already set", name, alias)
	}
	if e.level == 0 && level > 0 {
		return ErrForbidden
	}
	v, err := e.resolveVar(name, level)
	if err == nil {
		e.values[alias] = v
	}
	return err
}

func (e *env) Exists(name string) bool {
	_, ok := e.values[name]
	if !ok && e.parent != nil {
		return e.parent.Exists(name)
	}
	return ok
}

func (e *env) IsSet(name string) bool {
	return e.isSet(name, true)
}

func (e *env) Define(name, value string) error {
	if _, ok := e.values[name]; !ok && e.parent != nil {
		return e.parent.Define(name, value)
	}
	if e.isSet(name, false) {
		e.values[name].value = value
	} else {
		e.values[name] = &variable{
			value:    value,
			level:    e.level,
			refcount: 1,
		}
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

func (e *env) isSet(name string, parent bool) bool {
	v, ok := e.values[name]
	if !ok && e.parent != nil && parent {
		return e.parent.IsSet(name)
	}
	return v != nil
}

func (e *env) resolveVar(name string, level int) (*variable, error) {
	if level != 0 && e.parent != nil {
		if k, ok := e.parent.(*env); ok {
			return k.resolveVar(name, level-1)
		}
		return nil, ErrForbidden
	}
	v, ok := e.values[name]
	if !ok {
		return nil, undefinedVar(name)
	}
	return v, nil
}

type variable struct {
	value    string
	level    int
	refcount int
}

func undefinedVar(name string) error {
	return fmt.Errorf("%s: %w", name, ErrUndefined)
}

package interp

import (
	"github.com/midbel/gotcl/env"
	"github.com/midbel/slices"
)

type Env struct {
	list []env.Environment
}

func Environ() *Env {
	var e Env
	e.Append()
	return &e
}

func (e *Env) Push(ev env.Environment) {
	e.list = append(e.list, ev)
}

func (e *Env) Append() {
	e.Push(env.EmptyEnv())
}

func (e *Env) Pop() {
	n := len(e.list)
	if n == 0 {
		return
	}
	e.list = e.list[:n-1]
}

func (e *Env) Depth() int {
	return len(e.list)
}

func (e *Env) Resolve(name string) (string, error) {
	return e.Current().Resolve(name)
}

func (e *Env) Define(name, value string) error {
	return e.Current().Define(name, value)
}

func (e *Env) Delete(name string) error {
	return e.Current().Delete(name)
}

func (e *Env) Exists(name string) bool {
	return e.Current().Exists(name)
}

func (e *Env) IsSet(name string) bool {
	return e.Current().IsSet(name)
}

func (e *Env) Current() env.Environment {
	return slices.Lst(e.list)
}

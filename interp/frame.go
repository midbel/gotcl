package interp

import (
	"github.com/midbel/gotcl/env"
	"github.com/midbel/gotcl/stdlib"
)

type Frame struct {
	env      *env.Env
	ns       *Namespace
	deferred []stdlib.Executer

	cmd  string
	args []string
}

func (f *Frame) Names() []string {
	return f.env.Names()
}

func (f *Frame) Define(n string, v env.Value) {
	f.env.Define(n, v)
}

func (f *Frame) Delete(n string) {
	f.env.Delete(n)
}

func (f *Frame) Resolve(n string) (env.Value, error) {
	v, err := f.env.Resolve(n)
	if err == nil {
		return v, err
	}
	return f.ns.Resolve(n)
}

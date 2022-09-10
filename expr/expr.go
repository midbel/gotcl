package expr

import (
	"github.com/midbel/gotcl/env"
	"github.com/midbel/gotcl/expr/types"
	"github.com/midbel/gotcl/word"
)

type Env interface {
	Resolve(string) (env.Value, error)
}

type Expression interface {
	Eval(Env) (types.Value, error)
}

type Choice struct {
	Cdt Expression
	Csq Expression
	Alt Expression
}

func (c Choice) Eval(env Env) (types.Value, error) {
	res, err := c.Cdt.Eval(env)
	if err != nil {
		return nil, err
	}
	val, err := types.AsBool(res)
	if err != nil {
		return nil, err
	}
	if val {
		res, err = c.Csq.Eval(env)
	} else {
		res, err = c.Alt.Eval(env)
	}
	return res, err
}

type Number struct {
	types.Value
}

func (n Number) Eval(_ Env) (types.Value, error) {
	return n.Value, nil
}

type Identifier struct {
	Value string
}

func (i Identifier) Eval(env Env) (types.Value, error) {
	str, err := env.Resolve(i.Value)
	if err != nil {
		return nil, err
	}
	return types.AsValue(str)
}

type Prefix struct {
	Op    rune
	Right Expression
}

func (p Prefix) Eval(env Env) (types.Value, error) {
	v, err := p.Right.Eval(env)
	if err != nil {
		return nil, err
	}
	switch p.Op {
	case word.Not:
		return v.Not()
	case word.Sub:
		return v.Rev()
	case word.Bnot:
		return v.Bnot()
	default:
		return nil, nil
	}
}

type Infix struct {
	Left  Expression
	Right Expression
	Op    rune
}

func (i Infix) Eval(env Env) (types.Value, error) {
	left, err := i.Left.Eval(env)
	if err != nil {
		return nil, err
	}
	right, err := i.Right.Eval(env)
	if err != nil {
		return nil, err
	}
	switch i.Op {
	case word.And:
		return left.And(right)
	case word.Or:
		return left.Or(right)
	case word.Add:
		return left.Add(right)
	case word.Sub:
		return left.Sub(right)
	case word.Mul:
		return left.Mul(right)
	case word.Div:
		return left.Div(right)
	case word.Mod:
		return left.Mod(right)
	case word.Pow:
		return left.Pow(right)
	case word.Eq:
		return left.Eq(right)
	case word.Ne:
		return left.Ne(right)
	case word.Gt:
		return left.Gt(right)
	case word.Ge:
		return left.Ge(right)
	case word.Lt:
		return left.Lt(right)
	case word.Le:
		return left.Le(right)
	case word.Lshift:
		return left.Lshift(right)
	case word.Rshift:
		return left.Rshift(right)
	case word.Band:
		return left.Band(right)
	case word.Bor:
		return left.Bor(right)
	case word.Bxor:
		return left.Bxor(right)
	default:
		return nil, nil
	}
}

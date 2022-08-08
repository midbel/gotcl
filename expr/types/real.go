package types

import (
	"math"
	"strconv"
)

type Real struct {
	value float64
}

func RealValue(n float64) Value {
	return Real{
		value: n,
	}
}

func (r Real) Bool() (Value, error) {
	return BoolValue(r.value != 0), nil
}

func (r Real) Int() (Value, error) {
	return IntValue(int64(r.value)), nil
}

func (r Real) Double() (Value, error) {
	return r, nil
}

func (r Real) String() string {
	return strconv.FormatFloat(r.value, 'g', -1, 64)
}

func (r Real) Not() (Value, error) {
	return BoolValue(r.value == 0), nil
}

func (r Real) Rev() (Value, error) {
	r.value = -r.value
	return r, nil
}

func (r Real) Add(other Value) (Value, error) {
	switch x := other.(type) {
	default:
		return nil, incompatibleType()
	case Real:
		r.value += x.value
	case Integer:
		r.value += float64(x.value)
	}
	return r, nil
}

func (r Real) Sub(other Value) (Value, error) {
	switch x := other.(type) {
	default:
		return nil, incompatibleType()
	case Real:
		r.value -= x.value
	case Integer:
		r.value -= float64(x.value)
	}
	return r, nil
}

func (r Real) Div(other Value) (Value, error) {
	switch x := other.(type) {
	default:
		return nil, incompatibleType()
	case Real:
		if x.value == 0 {
			return nil, ErrZero
		}
		r.value /= x.value
	case Integer:
		if x.value == 0 {
			return nil, ErrZero
		}
		r.value /= float64(x.value)
	}
	return r, nil
}

func (r Real) Mod(other Value) (Value, error) {
	switch x := other.(type) {
	default:
		return nil, incompatibleType()
	case Real:
		if x.value == 0 {
			return nil, ErrZero
		}
		r.value = math.Mod(r.value, x.value)
	case Integer:
		if x.value == 0 {
			return nil, ErrZero
		}
		r.value = math.Mod(r.value, float64(x.value))
	}
	return r, nil
}

func (r Real) Mul(other Value) (Value, error) {
	switch x := other.(type) {
	default:
		return nil, incompatibleType()
	case Real:
		r.value *= x.value
	case Integer:
		r.value *= float64(x.value)
	}
	return r, nil
}

func (r Real) Pow(other Value) (Value, error) {
	switch x := other.(type) {
	default:
		return nil, incompatibleType()
	case Real:
		r.value = math.Pow(r.value, x.value)
	case Integer:
		r.value = math.Pow(r.value, float64(x.value))
	}
	return r, nil
}

func (r Real) And(other Value) (Value, error) {
	var (
		r1, _   = r.Bool()
		r2, err = other.Bool()
	)
	if err != nil {
		return nil, err
	}
	return r1.And(r2)
}

func (r Real) Or(other Value) (Value, error) {
	var (
		r1, _   = r.Bool()
		r2, err = other.Bool()
	)
	if err != nil {
		return nil, err
	}
	return r1.Or(r2)
}

func (r Real) Eq(other Value) (Value, error) {
	x, ok := other.(Real)
	if !ok {
		return nil, incompatibleType()
	}
	return BoolValue(r.value == x.value), nil
}

func (r Real) Ne(other Value) (Value, error) {
	x, err := r.Eq(other)
	if err != nil {
		return nil, err
	}
	return x.Not()
}

func (r Real) Lt(other Value) (Value, error) {
	x, ok := other.(Real)
	if !ok {
		return nil, incompatibleType()
	}
	return BoolValue(r.value < x.value), nil
}

func (r Real) Le(other Value) (Value, error) {
	var (
		r1, err1 = r.Lt(other)
		r2, err2 = r.Eq(other)
	)
	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}
	return r1.Or(r2)
}

func (r Real) Gt(other Value) (Value, error) {
	x, ok := other.(Real)
	if !ok {
		return nil, incompatibleType()
	}
	return BoolValue(r.value > x.value), nil
}

func (r Real) Ge(other Value) (Value, error) {
	var (
		r1, err1 = r.Gt(other)
		r2, err2 = r.Eq(other)
	)
	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}
	return r1.Or(r2)
}

func (r Real) Bnot() (Value, error) {
	return nil, unsupportedOp("~", "real")
}

func (r Real) Lshift(other Value) (Value, error) {
	return nil, unsupportedOp("<<", "real")
}

func (r Real) Rshift(other Value) (Value, error) {
	return nil, unsupportedOp(">>", "real")
}

func (r Real) Band(other Value) (Value, error) {
	return nil, unsupportedOp("&", "real")
}

func (r Real) Bor(other Value) (Value, error) {
	return nil, unsupportedOp("|", "real")
}

func (r Real) Bxor(other Value) (Value, error) {
	return nil, unsupportedOp("^", "real")
}

package types

import (
	"fmt"
	"strconv"
)

type Integer struct {
	value int64
}

func IntValue(n int64) Value {
	return Integer{
		value: n,
	}
}

func (i Integer) Bool() (Value, error) {
	return BoolValue(i.value != 0), nil
}

func (i Integer) Int() (Value, error) {
	return i, nil
}

func (i Integer) Double() (Value, error) {
	return RealValue(float64(i.value)), nil
}

func (i Integer) String() string {
	return strconv.FormatInt(i.value, 10)
}

func (i Integer) Not() (Value, error) {
	return BoolValue(i.value == 0), nil
}

func (i Integer) Rev() (Value, error) {
	i.value = -i.value
	return i, nil
}

func (i Integer) Add(other Value) (Value, error) {
	switch x := other.(type) {
	default:
		return nil, incompatibleType()
	case Integer:
		i.value += x.value
	case Real:
		r, _ := i.Double()
		return r.Add(other)
	}
	return i, nil
}

func (i Integer) Sub(other Value) (Value, error) {
	switch x := other.(type) {
	default:
		return nil, incompatibleType()
	case Integer:
		i.value -= x.value
	case Real:
		r, _ := i.Double()
		return r.Sub(other)
	}
	return i, nil
}

func (i Integer) Div(other Value) (Value, error) {
	switch x := other.(type) {
	default:
		return nil, incompatibleType()
	case Integer:
		if x.value == 0 {
			return nil, ErrZero
		}
		i.value /= x.value
	case Real:
		r, _ := i.Double()
		return r.Div(other)
	}
	return i, nil
}

func (i Integer) Mod(other Value) (Value, error) {
	switch x := other.(type) {
	default:
		return nil, incompatibleType()
	case Integer:
		if x.value == 0 {
			return nil, ErrZero
		}
		i.value %= x.value
	case Real:
		r, _ := i.Double()
		return r.Mod(other)
	}
	return i, nil
}

func (i Integer) Mul(other Value) (Value, error) {
	switch x := other.(type) {
	default:
		return nil, incompatibleType()
	case Integer:
		i.value *= x.value
	case Real:
		r, _ := i.Double()
		return r.Mul(other)
	}
	return i, nil
}

func (i Integer) Pow(other Value) (Value, error) {
	switch x := other.(type) {
	default:
		return nil, incompatibleType()
	case Integer:
		var (
			r1, _ = i.Double()
			r2, _ = x.Double()
			r, _  = r1.Pow(r2)
		)
		return r.Int()
	case Real:
		r, _ := i.Double()
		return r.Pow(other)
	}
	return i, nil
}

func (i Integer) And(other Value) (Value, error) {
	var (
		r1, _   = i.Bool()
		r2, err = other.Bool()
	)
	if err != nil {
		return nil, err
	}
	return r1.And(r2)
}

func (i Integer) Or(other Value) (Value, error) {
	var (
		r1, _   = i.Bool()
		r2, err = other.Bool()
	)
	if err != nil {
		return nil, err
	}
	return r1.Or(r2)
}

func (i Integer) Eq(other Value) (Value, error) {
	x, ok := other.(Integer)
	if !ok {
		return nil, incompatibleType()
	}
	return BoolValue(i.value == x.value), nil
}

func (i Integer) Ne(other Value) (Value, error) {
	r, err := i.Eq(other)
	if err != nil {
		return nil, err
	}
	return r.Not()
}

func (i Integer) Lt(other Value) (Value, error) {
	x, ok := other.(Integer)
	if !ok {
		return nil, incompatibleType()
	}
	return BoolValue(i.value < x.value), nil
}

func (i Integer) Le(other Value) (Value, error) {
	var (
		r1, err1 = i.Lt(other)
		r2, err2 = i.Eq(other)
	)
	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}
	return r1.Or(r2)
}

func (i Integer) Gt(other Value) (Value, error) {
	x, ok := other.(Integer)
	if !ok {
		return nil, incompatibleType()
	}
	return BoolValue(i.value > x.value), nil
}

func (i Integer) Ge(other Value) (Value, error) {
	var (
		r1, err1 = i.Gt(other)
		r2, err2 = i.Eq(other)
	)
	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}
	return r1.Or(r2)
}

func (i Integer) Bnot() (Value, error) {
	i.value = ^i.value
	return i, nil
}

func (i Integer) Lshift(other Value) (Value, error) {
	x, ok := other.(Integer)
	if !ok {
		return nil, incompatibleType()
	}
	if x.value < 0 {
		return nil, fmt.Errorf("negative shift count")
	}
	return Integer{
		value: i.value << x.value,
	}, nil
}

func (i Integer) Rshift(other Value) (Value, error) {
	x, ok := other.(Integer)
	if !ok {
		return nil, incompatibleType()
	}
	if x.value < 0 {
		return nil, fmt.Errorf("negative shift count")
	}
	return Integer{
		value: i.value >> x.value,
	}, nil
}

func (i Integer) Band(other Value) (Value, error) {
	x, ok := other.(Integer)
	if !ok {
		return nil, incompatibleType()
	}
	return Integer{
		value: i.value & x.value,
	}, nil
}

func (i Integer) Bor(other Value) (Value, error) {
	x, ok := other.(Integer)
	if !ok {
		return nil, incompatibleType()
	}
	return Integer{
		value: i.value | x.value,
	}, nil
}

func (i Integer) Bxor(other Value) (Value, error) {
	x, ok := other.(Integer)
	if !ok {
		return nil, incompatibleType()
	}
	return Integer{
		value: i.value ^ x.value,
	}, nil
}

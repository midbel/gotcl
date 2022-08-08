package types

import (
	"errors"
	"fmt"
	"strconv"
)

type Value interface {
	fmt.Stringer

	Int() (Value, error)
	Double() (Value, error)
	Bool() (Value, error)

	Not() (Value, error)
	Rev() (Value, error)

	Add(Value) (Value, error)
	Sub(Value) (Value, error)
	Div(Value) (Value, error)
	Mul(Value) (Value, error)
	Pow(Value) (Value, error)
	Mod(Value) (Value, error)

	And(Value) (Value, error)
	Or(Value) (Value, error)
	Eq(Value) (Value, error)
	Ne(Value) (Value, error)
	Lt(Value) (Value, error)
	Le(Value) (Value, error)
	Gt(Value) (Value, error)
	Ge(Value) (Value, error)

	Bnot() (Value, error)
	Lshift(Value) (Value, error)
	Rshift(Value) (Value, error)
	Band(Value) (Value, error)
	Bor(Value) (Value, error)
	Bxor(Value) (Value, error)
}

func AsFloat(v Value) (float64, error) {
	switch x := v.(type) {
	case Integer:
		return float64(x.value), nil
	case Real:
		return x.value, nil
	case Boolean:
		var res float64
		if x.value {
			res += 1
		}
		return res, nil
	default:
		return 0, incompatibleType()
	}
}

func AsValue(str string) (Value, error) {
	if v, err := strconv.ParseInt(str, 0, 64); err == nil {
		return IntValue(v), nil
	}
	if f, err := strconv.ParseFloat(str, 64); err == nil {
		return RealValue(f), nil
	}

	if b, err := strconv.ParseBool(str); err == nil {
		return BoolValue(b), nil
	}
	return nil, fmt.Errorf("%s can not be converted to Value")
}

func AsInt(v Value) (int64, error) {
	switch x := v.(type) {
	case Integer:
		return x.value, nil
	case Real:
		return int64(x.value), nil
	case Boolean:
		var res int64
		if x.value {
			res++
		}
		return res, nil
	default:
		return 0, incompatibleType()
	}
}

func AsString(v Value) (string, error) {
	return v.String(), nil
}

func AsBool(v Value) (bool, error) {
	switch x := v.(type) {
	case Integer:
		return x.value != 0, nil
	case Real:
		return x.value != 0, nil
	case Boolean:
		return x.value, nil
	default:
		return false, incompatibleType()
	}
}

var (
	ErrOperation = errors.New("unsupported operation")
	ErrType      = errors.New("incompatible type")
	ErrZero      = errors.New("division by zero")
	ErrCast      = errors.New("type can not be casted")
)

func unsupportedOp(op, str string) error {
	return fmt.Errorf("%s: %w on %s type", op, ErrOperation, str)
}

func incompatibleType() error {
	return ErrType
}

func unsupportedCast(src, dst string) error {
	return fmt.Errorf("%s: %w to %s", ErrCast)
}

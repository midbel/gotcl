package types

type Boolean struct {
	value bool
}

func BoolValue(b bool) Value {
	return Boolean{
		value: b,
	}
}

func (b Boolean) Bool() (Value, error) {
	return b, nil
}

func (b Boolean) Int() (Value, error) {
	var i Integer
	if b.value {
		i.value++
	}
	return i, nil
}

func (b Boolean) Double() (Value, error) {
	var r Real
	if b.value {
		r.value += 1
	}
	return r, nil
}

func (b Boolean) String() string {
	if b.value {
		return "1"
	}
	return "0"
}

func (b Boolean) Not() (Value, error) {
	b.value = !b.value
	return b, nil
}

func (b Boolean) Rev() (Value, error) {
	return nil, unsupportedOp("-", "boolean")
}

func (b Boolean) Add(other Value) (Value, error) {
	return nil, unsupportedOp("+", "boolean")
}

func (b Boolean) Sub(other Value) (Value, error) {
	return nil, unsupportedOp("-", "boolean")
}

func (b Boolean) Div(other Value) (Value, error) {
	return nil, unsupportedOp("/", "boolean")
}

func (b Boolean) Mod(other Value) (Value, error) {
	return nil, unsupportedOp("%", "boolean")
}

func (b Boolean) Mul(other Value) (Value, error) {
	return nil, unsupportedOp("*", "boolean")
}

func (b Boolean) Pow(other Value) (Value, error) {
	return nil, unsupportedOp("**", "boolean")
}

func (b Boolean) And(other Value) (Value, error) {
	x, ok := other.(Boolean)
	if !ok {
		return nil, incompatibleType()
	}
	b.value = b.value && x.value
	return b, nil
}

func (b Boolean) Or(other Value) (Value, error) {
	x, ok := other.(Boolean)
	if !ok {
		return nil, incompatibleType()
	}
	b.value = b.value || x.value
	return b, nil
}

func (b Boolean) Eq(other Value) (Value, error) {
	x, ok := other.(Boolean)
	if !ok {
		return nil, incompatibleType()
	}
	b.value = b.value == x.value
	return b, nil
}

func (b Boolean) Ne(other Value) (Value, error) {
	x, ok := other.(Boolean)
	if !ok {
		return nil, incompatibleType()
	}
	b.value = b.value != x.value
	return b, nil
}

func (b Boolean) Lt(other Value) (Value, error) {
	return nil, unsupportedOp("<", "boolean")
}

func (b Boolean) Le(other Value) (Value, error) {
	return nil, unsupportedOp("<=", "boolean")
}

func (b Boolean) Gt(other Value) (Value, error) {
	return nil, unsupportedOp(">", "boolean")
}

func (b Boolean) Ge(other Value) (Value, error) {
	return nil, unsupportedOp(">=", "boolean")
}

func (b Boolean) Bnot() (Value, error) {
	return nil, unsupportedOp("~", "boolean")
}

func (b Boolean) Lshift(other Value) (Value, error) {
	return nil, unsupportedOp("<<", "boolean")
}

func (b Boolean) Rshift(other Value) (Value, error) {
	return nil, unsupportedOp(">>", "boolean")
}

func (b Boolean) Band(other Value) (Value, error) {
	return nil, unsupportedOp("&", "boolean")
}

func (b Boolean) Bor(other Value) (Value, error) {
	return nil, unsupportedOp("|", "boolean")
}

func (b Boolean) Bxor(other Value) (Value, error) {
	return nil, unsupportedOp("^", "boolean")
}

package env

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrCast = errors.New("type can not be casted")

type Value interface {
	fmt.Stringer

	ToList() (Value, error)
	ToArray() (Value, error)
	ToNumber() (Value, error)
	ToString() (Value, error)
	ToBoolean() (Value, error)
}

type Array struct {
	values map[string]Value
}

func ZipArr(keys []string, values []Value) Value {
	return EmptyArr()
}

func EmptyArr() Value {
	return Array{
		values: make(map[string]Value),
	}
}

func (a Array) Get(n string) Value {
	return a.values[n]
}

func (a Array) Set(n string, v Value) {
	a.values[n] = v
}

func (a Array) String() string {
	var str strings.Builder
	for k, v := range a.values {
		str.WriteString(k)
		str.WriteString(" ")
		str.WriteString(v.String())
	}
	return str.String()
}

func (a Array) ToList() (Value, error) {
	return nil, nil
}

func (a Array) ToArray() (Value, error) {
	return a, nil
}

func (a Array) ToNumber() (Value, error) {
	return nil, ErrCast
}

func (a Array) ToString() (Value, error) {
	return Str(a.String()), nil
}

func (a Array) ToBoolean() (Value, error) {
	return Bool(len(a.values) != 0), nil
}

type List struct {
	values []Value
}

func ListFromStrings(vs []string) Value {
	var list []Value
	for i := range vs {
		list = append(list, Str(vs[i]))
	}
	return ListFrom(list...)
}

func ListFrom(vs ...Value) Value {
	if len(vs) == 0 {
		return EmptyList()
	}
	var i List
	i.values = append(i.values, vs...)
	return i
}

func EmptyList() Value {
	var i List
	return i
}

func (i List) String() string {
	var list []string
	for _, v := range i.values {
		list = append(list, v.String())
	}
	return strings.Join(list, " ")
}

func (i List) Len() int {
	return len(i.values)
}

func (i List) ToList() (Value, error) {
	return i, nil
}

func (i List) ToArray() (Value, error) {
	if len(i.values)%2 != 0 {
		return nil, ErrCast
	}
	var (
		ks []string
		vs []Value
	)
	for j := 0; j < len(i.values); j += 2 {
		ks = append(ks, i.values[j].String())
		vs = append(vs, i.values[j+1])
	}
	return ZipArr(ks, vs), nil
}

func (i List) ToNumber() (Value, error) {
	return nil, ErrCast
}

func (i List) ToString() (Value, error) {
	return Str(i.String()), nil
}

func (i List) ToBoolean() (Value, error) {
	return nil, ErrCast
}

type String struct {
	value string
}

func Str(str string) Value {
	return String{value: str}
}

func EmptyStr() Value {
	return Str("")
}

func (s String) String() string {
	return s.value
}

func (s String) ToList() (Value, error) {
	return split(s.value)
}

func (s String) ToArray() (Value, error) {
	list, err := s.ToList()
	if err != nil {
		return nil, err
	}
	return list.ToArray()
}

func (s String) ToNumber() (Value, error) {
	n, err := strconv.ParseFloat(s.value, 64)
	if err != nil {
		return nil, err
	}
	return Float(n), nil
}

func (s String) ToString() (Value, error) {
	return s, nil
}

func (s String) ToBoolean() (Value, error) {
	return Bool(s.value != ""), nil
}

type Boolean struct {
	value bool
}

func False() Value {
	return Bool(false)
}

func True() Value {
	return Bool(true)
}

func Bool(b bool) Value {
	return Boolean{value: b}
}

func (b Boolean) String() string {
	if b.value {
		return "1"
	}
	return "0"
}

func (b Boolean) ToList() (Value, error) {
	return ListFrom(b), nil
}

func (b Boolean) ToArray() (Value, error) {
	return nil, ErrCast
}

func (b Boolean) ToNumber() (Value, error) {
	if !b.value {
		return Zero(), nil
	}
	return Float(1), nil
}

func (b Boolean) ToString() (Value, error) {
	return Str(b.String()), nil
}

func (b Boolean) ToBoolean() (Value, error) {
	return b, nil
}

type Number struct {
	value float64
}

func Float(f float64) Value {
	return Number{value: f}
}

func Int(i int64) Value {
	return Float(float64(i))
}

func Zero() Value {
	return Float(0)
}

func (n Number) String() string {
	return strconv.FormatFloat(n.value, 'g', -1, 64)
}

func (n Number) ToList() (Value, error) {
	return ListFrom(n), nil
}

func (n Number) ToArray() (Value, error) {
	return nil, ErrCast
}

func (n Number) ToNumber() (Value, error) {
	return n, nil
}

func (n Number) ToString() (Value, error) {
	str := strconv.FormatFloat(n.value, 'g', -1, 64)
	return Str(str), nil
}

func (n Number) ToBoolean() (Value, error) {
	return Bool(int(n.value) == 0), nil
}

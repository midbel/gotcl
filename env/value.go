package env

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/midbel/slices"
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

func (a Array) Len() int {
	return len(a.values)
}

func (a Array) Get(n string) Value {
	return a.values[n]
}

func (a Array) Set(n string, v Value) {
	a.values[n] = v
}

func (a Array) Unset(n string) {
	delete(a.values, n)
}

func (a Array) Pairs() Value {
	var list []Value
	for k, v := range a.values {
		list = append(list, ListFrom(Str(k), v))
	}
	return ListFrom(list...)
}

func (a Array) Names() []string {
	var list []string
	for k := range a.values {
		list = append(list, k)
	}
	return list
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

func (i List) At(n int) Value {
	if n < 0 || n >= len(i.values) {
		return EmptyStr()
	}
	return i.values[n]
}

func (i List) Range(fst, lst int) (Value, error) {
	if fst < 0 || fst >= i.Len() || lst < 0 || lst >= i.Len() || fst > lst {
		return nil, fmt.Errorf("invalid range given: %d - %d", fst, lst)
	}
	return ListFrom(i.values[fst:lst]...), nil
}

func (i List) Flat(full bool) Value {
	var flatten func(List, bool) []Value

	flatten = func(i List, full bool) []Value {
		vs := make([]Value, len(i.values))
		for _, v := range i.values {
			if a, ok := v.(List); ok && full {
				xs := flatten(a, full)
				vs = append(vs, xs...)
			} else {
				vs = append(vs, a)
			}
		}
		return vs
	}
	return ListFrom(flatten(i, full)...)
}

func (i List) Reverse() List {
	j := List{}
	j.values = slices.Reverse(i.values)
	return j
}

func (i List) Shuffle() Value {
	vs := make([]Value, len(i.values))
	copy(vs, i.values)
	return ListFrom(slices.Shuffle(i.values)...)
}

func (i List) Equal(other Value) (bool, error) {
	x, err := other.ToList()
	if err != nil {
		return false, err
	}
	j := x.(List)
	if len(i.values) != len(j.values) {
		return false, nil
	}
	for k := range i.values {
		if i.values[k].String() != j.values[k].String() {
			return false, nil
		}
	}
	return true, nil
}

func (i List) Set(v Value, n int) (Value, error) {
	if n < 0 || n >= len(i.values) {
		return nil, fmt.Errorf("index out of range")
	}
	vs := make([]Value, len(i.values))
	copy(vs, i.values)
	vs[n] = v
	return ListFrom(vs...), nil
}

func (i List) Swap(j, k int) List {
	i.values[j], i.values[k] = i.values[k], i.values[j]
	return i
}

func (i List) Shift() (Value, Value) {
	return slices.Fst(i.values), ListFrom(slices.Rest(i.values)...)
}

func (i List) Apply(do func(Value) Value) Value {
	vs := slices.Map(i.values, do)
	return ListFrom(vs...)
}

func (i List) Filter(do func(Value) bool) Value {
	vs := slices.Filter(i.values, do)
	return ListFrom(vs...)
}

func (i List) Replace(v Value, fst, lst int) Value {
	vs := make([]Value, len(i.values))
	if fst < 0 || lst < 0 {
		vs = append([]Value{v}, vs...)
	}
	if fst >= len(vs) || lst >= len(i.values) {
		vs = append(vs, v)
	}
	if lst < fst {
		lst = fst
	}
	vs = append(vs[:fst], append([]Value{v}, vs[lst:]...)...)
	return ListFrom(vs...)
}

func (i List) Insert(n int, v Value) Value {
	vs := make([]Value, len(i.values))
	copy(vs, i.values)
	if n <= 0 {
		vs = append([]Value{v}, vs...)
		return ListFrom(vs...)
	}
	if n >= len(vs) {
		vs = append(vs, v)
		return ListFrom(vs...)
	}
	vs = append(vs[:n], append([]Value{v}, vs[n:]...)...)
	return ListFrom(vs...)
}

func (i List) Append(v Value) {
	i.values = append(i.values, v)
}

func (i List) Values() []Value {
	vs := make([]Value, len(i.values))
	copy(vs, i.values)
	return vs
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

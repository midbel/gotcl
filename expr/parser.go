package expr

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/midbel/gotcl/expr/types"
	"github.com/midbel/gotcl/word"
)

const (
	Lowest int = iota
	Condition
	LogicalOr
	LogicalAnd
	BitOr
	BitXor
	BitAnd
	Equality
	Relational
	Shift
	Additive
	Multiplicative
	Unary
)

var bindings = map[rune]int{
	word.And:     Relational,
	word.Or:      Relational,
	word.Add:     Additive,
	word.Sub:     Additive,
	word.Mul:     Multiplicative,
	word.Div:     Multiplicative,
	word.Mod:     Multiplicative,
	word.Pow:     Multiplicative,
	word.Eq:      Equality,
	word.Ne:      Equality,
	word.Gt:      Relational,
	word.Ge:      Relational,
	word.Lt:      Relational,
	word.Le:      Relational,
	word.Lshift:  Shift,
	word.Rshift:  Shift,
	word.Band:    BitAnd,
	word.Bor:     BitOr,
	word.Bxor:    BitXor,
	word.Ternary: Condition,
}

type Parser struct {
	scan *word.Scanner

	curr word.Word
	peek word.Word

	prefix map[rune]func() (Expression, error)
	infix  map[rune]func(Expression) (Expression, error)
}

func Parse(str string) (*Parser, error) {
	s, err := word.Scan(strings.NewReader(str))
	if err != nil {
		return nil, err
	}
	p := Parser{
		scan:   s,
		prefix: make(map[rune]func() (Expression, error)),
		infix:  make(map[rune]func(Expression) (Expression, error)),
	}
	p.registerPrefix(word.Not, p.parsePrefix)
	p.registerPrefix(word.Sub, p.parsePrefix)
	p.registerPrefix(word.Bnot, p.parsePrefix)
	p.registerPrefix(word.Int, p.parseNumber)
	p.registerPrefix(word.Float, p.parseNumber)
	p.registerPrefix(word.Variable, p.parseVariable)
	p.registerPrefix(word.Paren, p.parseGroup)
	p.registerInfix(word.And, p.parseInfix)
	p.registerInfix(word.Or, p.parseInfix)
	p.registerInfix(word.Add, p.parseInfix)
	p.registerInfix(word.Sub, p.parseInfix)
	p.registerInfix(word.Mul, p.parseInfix)
	p.registerInfix(word.Div, p.parseInfix)
	p.registerInfix(word.Mod, p.parseInfix)
	p.registerInfix(word.Pow, p.parseInfix)
	p.registerInfix(word.Eq, p.parseInfix)
	p.registerInfix(word.Ne, p.parseInfix)
	p.registerInfix(word.Gt, p.parseInfix)
	p.registerInfix(word.Ge, p.parseInfix)
	p.registerInfix(word.Lt, p.parseInfix)
	p.registerInfix(word.Le, p.parseInfix)
	p.registerInfix(word.Lshift, p.parseInfix)
	p.registerInfix(word.Rshift, p.parseInfix)
	p.registerInfix(word.Band, p.parseInfix)
	p.registerInfix(word.Bor, p.parseInfix)
	p.registerInfix(word.Bxor, p.parseInfix)
	p.registerInfix(word.Ternary, p.parseTernary)

	p.next()
	p.next()
	return &p, nil
}

func (p *Parser) Parse() (Expression, error) {
	return p.parseExpression(Lowest)
}

func (p *Parser) parseExpression(binding int) (Expression, error) {
	fn, ok := p.prefix[p.curr.Type]
	if !ok {
		return nil, fmt.Errorf("unsupported prefix operator")
	}
	left, err := fn()
	if err != nil {
		return nil, err
	}
	for p.peek.Type != word.EOF && binding < p.peekPower() {
		p.next()
		fn, ok := p.infix[p.curr.Type]
		if !ok {
			return nil, fmt.Errorf("unsupported infix operator")
		}
		left, err = fn(left)
		if err != nil {
			return nil, err
		}
	}
	return left, nil
}

func (p *Parser) currPower() int {
	pow, ok := bindings[p.curr.Type]
	if !ok {
		pow = Lowest
	}
	return pow
}

func (p *Parser) peekPower() int {
	pow, ok := bindings[p.peek.Type]
	if !ok {
		pow = Lowest
	}
	return pow
}

func (p *Parser) parseGroup() (Expression, error) {
	p.next()
	expr, err := p.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}
	if p.peek.Type != word.Paren {
		return nil, fmt.Errorf("syntax error: missing closing parenthese")
	}
	p.next()
	return expr, nil
}

func (p *Parser) parseNumber() (Expression, error) {
	var val types.Value
	switch p.curr.Type {
	case word.Int:
		v, err := strconv.ParseInt(p.curr.Literal, 0, 64)
		if err != nil {
			return nil, err
		}
		val = types.IntValue(v)
	case word.Float:
		v, err := strconv.ParseFloat(p.curr.Literal, 64)
		if err != nil {
			return nil, err
		}
		val = types.RealValue(v)
	default:
		return nil, fmt.Errorf("unsupported word type: %s", p.curr)
	}
	return Number{Value: val}, nil
}

func (p *Parser) parseVariable() (Expression, error) {
	i := Identifier{
		Value: p.curr.Literal,
	}
	return i, nil
}

func (p *Parser) parsePrefix() (Expression, error) {
	x := Prefix{
		Op: p.curr.Type,
	}
	p.next()
	right, err := p.parseExpression(Unary)
	if err != nil {
		return nil, err
	}
	x.Right = right
	return x, nil
}

func (p *Parser) parseInfix(left Expression) (Expression, error) {
	i := Infix{
		Left: left,
		Op:   p.curr.Type,
	}
	pow := p.currPower()
	p.next()
	right, err := p.parseExpression(pow)
	if err != nil {
		return nil, err
	}
	i.Right = right
	return i, nil
}

func (p *Parser) parseTernary(left Expression) (Expression, error) {
	var (
		expr = Choice{Cdt: left}
		err  error
	)
	p.next()
	expr.Csq, err = p.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}
	p.next()
	if p.curr.Type != word.Alt {
		return nil, fmt.Errorf("unexpected word: %s", p.curr)
	}
	p.next()
	expr.Alt, err = p.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}
	return expr, nil
}

func (p *Parser) registerPrefix(op rune, fn func() (Expression, error)) {
	p.prefix[op] = fn
}

func (p *Parser) registerInfix(op rune, fn func(Expression) (Expression, error)) {
	p.infix[op] = fn
}

func (p *Parser) next() {
	p.curr = p.peek
	p.peek = p.scan.Tokenize()
}

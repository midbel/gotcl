package word

import (
	"fmt"
)

const (
	EOF = -(iota + 1)
	EOL
	Blank
	Literal
	Int
	Float
	Variable
	Script // [...]
	Quote  // "..."
	Comment
	Illegal
	Paren
	Namespace
	Ternary
	Alt
	And
	Or
	Not
	Add
	Sub
	Mul
	Div
	Mod
	Pow
	Eq
	Ne
	Gt
	Ge
	Lt
	Le
	Lshift
	Rshift
	Band
	Bor
	Bxor
	Bnot
)

type Position struct {
	Line   int
	Column int
}

func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}

type Word struct {
	Literal string
	Type    rune
	Position
}

func (w Word) IsEOL() bool {
	return w.Type == Comment || w.Type == EOL || w.Type == EOF
}

func (w Word) String() string {
	return w.Debug()
}

func (w Word) Debug() string {
	switch w.Type {
	case EOF:
		return "<eof>"
	case EOL:
		return "<eol>"
	case Blank:
		return "<blank>"
	case Paren:
		return "<paren>"
	case And:
		return "<and>"
	case Or:
		return "<or>"
	case Not:
		return "<not>"
	case Add:
		return "<add>"
	case Sub:
		return "<subtract>"
	case Mul:
		return "<multiply>"
	case Div:
		return "<divide>"
	case Pow:
		return "<power>"
	case Mod:
		return "<modulo>"
	case Eq:
		return "<eq>"
	case Ne:
		return "<ne>"
	case Gt:
		return "<greatthan>"
	case Ge:
		return "<greateq>"
	case Lt:
		return "<lessthan>"
	case Le:
		return "<lesseq>"
	case Lshift:
		return "<lshift>"
	case Rshift:
		return "<rshift>"
	case Band:
		return "<bin-and>"
	case Bor:
		return "<bin-or>"
	case Bxor:
		return "<bin-xor>"
	case Bnot:
		return "<bin-not>"
	case Namespace:
		return "<namespace>"
	case Ternary:
		return "<ternary>"
	case Alt:
		return "alternative"
	case Int:
		return fmt.Sprintf("integer(%s)", w.Literal)
	case Float:
		return fmt.Sprintf("float(%s)", w.Literal)
	case Script:
		return fmt.Sprintf("script(%s)", w.Literal)
	case Literal:
		return fmt.Sprintf("literal(%s)", w.Literal)
	case Quote:
		return fmt.Sprintf("quote(%s)", w.Literal)
	case Variable:
		return fmt.Sprintf("variable(%s)", w.Literal)
	case Comment:
		return fmt.Sprintf("comment(%s)", w.Literal)
	case Illegal:
		return fmt.Sprintf("illegal(%s)", w.Literal)
	default:
		return "<unknown>"
	}
}

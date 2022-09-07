package interp

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/gotcl/word"
)

var ErrIncomplete = errors.New("incomplete")

type Command struct {
	Name env.Value
	Args []env.Value
}

type Parser struct {
	scan *word.Scanner
	curr word.Word
	peek word.Word
}

func New(r io.Reader) (*Parser, error) {
	scan, err := word.Scan(r)
	if err != nil {
		return nil, err
	}
	p := Parser{
		scan: scan,
	}
	p.next()
	p.next()
	return &p, nil
}

func (p *Parser) Parse(i *Interpreter) (*Command, error) {
	p.skipEmptyLines()
	if p.done() {
		return nil, io.EOF
	}
	var (
		c   Command
		err error
	)
	c.Name, err = p.parse(i)
	if err != nil {
		return nil, err
	}
	for !p.done() && !p.curr.IsEOL() {
		arg, err := p.parse(i)
		if err != nil {
			return nil, err
		}
		c.Args = append(c.Args, arg)
	}
	p.next()
	return &c, nil
}

func (p *Parser) parse(i *Interpreter) (env.Value, error) {
	p.skipBlank()
	var vs []env.Value
	for !p.isEnd() {
		if p.curr.Type == word.Illegal {
			return nil, ErrSyntax
		}
		v, err := substitute(p.curr, i)
		if err != nil {
			return nil, err
		}
		vs = append(vs, v)
		p.next()
	}
	if p.isBlank() {
		p.next()
	}
	return list2str(vs), nil
}

func (p *Parser) next() {
	p.curr = p.peek
	p.peek = p.scan.Scan()
}

func (p *Parser) done() bool {
	return p.curr.Type == word.EOF
}

func (p *Parser) skipEnd() {
	for p.isEnd() {
		p.next()
	}
}

func (p *Parser) skipBlank() {
	for p.isBlank() && !p.done() {
		p.next()
	}
}

func (p *Parser) skipEmptyLines() {
	for p.curr.IsEOL() && !p.done() {
		p.next()
	}
}

func (p *Parser) isEnd() bool {
	return p.curr.IsEOL() || p.isBlank()
}

func (p *Parser) isBlank() bool {
	return p.curr.Type == word.Blank
}

func list2str(list []env.Value) env.Value {
	if len(list) == 1 {
		return list[0]
	}
	var str strings.Builder
	for i := range list {
		str.WriteString(list[i].String())
	}
	return env.Str(str.String())
}

func substitute(curr word.Word, i *Interpreter) (env.Value, error) {
	split := func(str string, i *Interpreter) (env.Value, error) {
		scan, err := word.Scan(strings.NewReader(str))
		if err != nil {
			return nil, err
		}
		var list []env.Value
		for {
			w := scan.Split()
			if w.Type == word.EOF {
				break
			}
			val, err := substitute(w, i)
			if err != nil {
				return nil, err
			}
			list = append(list, val)
		}
		return list2str(list), nil
	}
	var (
		val env.Value
		err error
	)
	switch curr.Type {
	case word.Literal, word.Block:
		val = env.Str(curr.Literal)
	case word.Variable:
		val, err = i.Resolve(curr.Literal)
	case word.Quote:
		val, err = split(curr.Literal, i)
	case word.Script:
		val, err = i.Execute(strings.NewReader(curr.Literal))
	default:
		err = fmt.Errorf("%s: %w", curr, ErrSyntax)
	}
	return val, err
}

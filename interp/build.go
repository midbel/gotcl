package interp

import (
	"fmt"
	"io"
	"strings"

	"github.com/midbel/gotcl/stdlib"
	"github.com/midbel/gotcl/word"
)

type Command struct {
	Cmd  string
	Args []string
}

type Builder struct {
	scan *word.Scanner

	curr word.Word
	peek word.Word
}

func Build(r io.Reader) (*Builder, error) {
	s, err := word.Scan(r)
	if err != nil {
		return nil, err
	}
	b := Builder{
		scan: s,
	}
	b.next()
	b.next()
	return &b, nil
}

func (b *Builder) Next(env stdlib.Interpreter) (*Command, error) {
	if b.done() {
		return nil, io.EOF
	}
	b.skipEnd()
	var (
		c   Command
		err error
	)
	c.Cmd, err = b.prepare(env)
	if err != nil {
		return nil, err
	}
	for !b.done() && !b.curr.IsEOL() {
		str, err := b.prepare(env)
		if err != nil {
			return nil, err
		}
		c.Args = append(c.Args, str)
	}
	b.next()
	return &c, nil
}

func (b *Builder) prepare(env stdlib.Interpreter) (string, error) {
	b.skipBlank()
	var str strings.Builder
	for !b.isEnd() {
		word, err := substitute(b.curr, env)
		if err != nil {
			return "", err
		}
		str.WriteString(word)
		b.next()
	}
	if b.isBlank() {
		b.next()
	}
	return str.String(), nil
}

func (b *Builder) skipEnd() {
	for b.isEnd() {
		b.next()
	}
}

func (b *Builder) skipBlank() {
	for b.isBlank() {
		b.next()
	}
}

func (b *Builder) isEnd() bool {
	return b.curr.IsEOL() || b.isBlank()
}

func (b *Builder) isBlank() bool {
	return b.curr.Type == word.Blank
}

func (b *Builder) done() bool {
	return b.curr.Type == word.EOF
}

func (b *Builder) next() {
	b.curr = b.peek
	b.peek = b.scan.Scan()
}

func scan(str string) ([]string, error) {
	scan, err := word.Scan(strings.NewReader(str))
	if err != nil {
		return nil, err
	}
	var list []string
	for {
		w := scan.Scan()
		if w.Type == word.EOF {
			break
		}
		list = append(list, w.Literal)
	}
	return list, nil
}

func split(str string, env stdlib.Interpreter) (string, error) {
	scan, err := word.Scan(strings.NewReader(str))
	if err != nil {
		return "", err
	}
	var list []string
	for {
		w := scan.Split()
		if w.Type == word.EOF {
			break
		}
		str, err := substitute(w, env)
		if err != nil {
			return "", err
		}
		list = append(list, str)
	}
	return strings.Join(list, ""), nil
}

func substitute(curr word.Word, env stdlib.Interpreter) (string, error) {
	var (
		str = curr.Literal
		err error
	)
	switch curr.Type {
	case word.Literal:
	case word.Variable:
		str, err = env.Resolve(str)
	case word.Quote:
		str, err = split(str, env)
	case word.Script:
		str, err = env.Execute(strings.NewReader(str))
	default:
		err = fmt.Errorf("%s: unsupported word type", curr)
	}
	return str, err
}

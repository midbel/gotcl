package word

import (
	"bytes"
	"io"
	"unicode/utf8"
)

type Scanner struct {
	file       string
	keepBlanks bool

	input []byte
	curr  int
	next  int
	char  rune

	str bytes.Buffer

	pos  Position
	prev Position
}

func (s *Scanner) KeepBlanks(keep bool) {
	s.keepBlanks = keep
}

func Scan(r io.Reader) (*Scanner, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	s := Scanner{
		keepBlanks: true,
		input:      bytes.ReplaceAll(b, []byte{cr, nl}, []byte{nl}),
		pos:        Position{Line: 1},
		prev:       Position{Line: 1},
	}
	if n, ok := r.(interface{ Name() string }); ok {
		s.file = n.Name()
	}
	return &s, nil
}

func (s *Scanner) Tokenize() Word {
	w := s.prepare()
	if w.Type == EOF {
		return w
	}
	if isBlank(s.char) {
		s.skipBlank()
	}
	switch {
	case isVariable(s.char):
		s.scanVariable(&w)
	case isDigit(s.char):
		s.scanNumber(&w)
	case isOperator(s.char):
		s.scanOperator(&w)
	case s.char == lparen || s.char == rparen:
		w.Type = Paren
	default:
	}
	return w
}

func (s *Scanner) Split() Word {
	w := s.prepare()
	if w.Type == EOF {
		return w
	}
	switch {
	case isVariable(s.char):
		s.scanVariable(&w)
	case isScript(s.char):
		s.scanScript(&w)
	default:
		s.scanLiteral(&w, isQuotedDelimiter)
	}
	return w
}

func (s *Scanner) Scan() Word {
	w := s.prepare()
	if w.Type == EOF {
		return w
	}
	if s.char == backslash && s.peek() == nl {
		s.read()
		s.read()
	}
	if isBlank(s.char) {
		s.skipBlank()
		if s.keepBlanks {
			w.Type = Blank
			return w
		}
		s.read()
	}
	switch s.char {
	case dash:
		s.scanComment(&w)
	case lsquare:
		s.scanScript(&w)
	case lcurly:
		s.scanBraces(&w)
	case dquote:
		s.scanQuote(&w)
	case semicolon, nl:
		s.scanEOL(&w)
	case dollar:
		s.scanVariable(&w)
	default:
		s.scanLiteral(&w, isDelimiter)
	}
	return w
}

func (s *Scanner) scanEOL(w *Word) {
	if s.char == semicolon {
		s.read()
	}
	s.skip(func(r rune) bool { return isBlank(r) || isNL(r) })
	w.Type = EOL
}

func (s *Scanner) scanOperator(w *Word) {
	switch s.char {
	case plus:
		w.Type = Add
	case minus:
		w.Type = Sub
	case slash:
		w.Type = Div
	case star:
		w.Type = Mul
		if s.peek() == star {
			s.read()
			w.Type = Pow
		}
	case percent:
		w.Type = Mod
	case equal:
		w.Type = Literal
		if s.peek() == equal {
			w.Type = Eq
		}
	case bang:
		w.Type = Not
		if s.peek() == equal {
			s.read()
			w.Type = Ne
		}
	case ampersand:
		w.Type = Band
		if s.peek() == ampersand {
			s.read()
			w.Type = And
		}
	case pipe:
		w.Type = Bor
		if s.peek() == pipe {
			s.read()
			w.Type = Or
		}
	case langle:
		w.Type = Lt
		if k := s.peek(); k == equal {
			s.read()
			w.Type = Le
		} else if k == langle {
			s.read()
			w.Type = Lshift
		}
	case rangle:
		w.Type = Gt
		if k := s.peek(); k == equal {
			s.read()
			w.Type = Ge
		} else if k == rangle {
			s.read()
			w.Type = Rshift
		}
	case tilde:
		w.Type = Bnot
	case caret:
		w.Type = Bxor
	case question:
		w.Type = Ternary
	case colon:
		w.Type = Alt
		if s.peek() == colon {
			w.Type = Namespace
			s.read()
		}
	default:
		w.Type = Illegal
	}
}

func (s *Scanner) scanNumber(w *Word) {
	defer s.unread()
	if s.char == '0' {
		var accept func(rune) bool
		switch k := s.peek(); k {
		case 'x':
			accept = isHexa
		case 'b':
			accept = isBinary
		case 'o':
			accept = isOctal
		default:
		}
		if accept != nil {
			s.scanInteger(w, accept)
			return
		}
	}
	for isDigit(s.char) {
		s.str.WriteRune(s.char)
		s.read()
	}
	w.Type = Int
	w.Literal = s.str.String()
	if s.char != dot {
		return
	}
	s.str.WriteRune(s.char)
	s.read()
	for isDigit(s.char) {
		s.str.WriteRune(s.char)
		s.read()
	}
	w.Type = Float
	w.Literal = s.str.String()
}

func (s *Scanner) scanInteger(w *Word, accept func(rune) bool) {
	s.str.WriteRune(s.char)
	s.read()
	s.str.WriteRune(s.char)
	s.read()
	for accept(s.char) {
		s.str.WriteRune(s.char)
		s.read()
	}
	w.Type = Int
	w.Literal = s.str.String()
}

func (s *Scanner) scanQuote(w *Word) {
	s.read()
	for s.char != dquote && !s.done() {
		if s.char == lsquare {
			s.str.WriteRune(s.char)
			s.scanScript(w)
		}
		s.str.WriteRune(s.char)
		if s.char == backslash {
			s.read()
			s.str.WriteRune(s.char)
		}
		s.read()
	}
	w.Type = Quote
	if s.char != dquote {
		w.Type = Illegal
	}
	w.Literal = s.str.String()
}

func (s *Scanner) scanVariable(w *Word) {
	defer s.unread()
	s.read()
	var (
		escaped bool
		accept  = func(r rune) bool { return isAlpha(r) || r == colon }
	)
	if escaped = s.char == lcurly; escaped {
		s.read()
		accept = func(c rune) bool { return c != rcurly }
	}
	for accept(s.char) {
		s.str.WriteRune(s.char)
		s.read()
	}
	if escaped && s.char == rcurly {
		s.read()
	}
	w.Type = Variable
	w.Literal = s.str.String()
}

func (s *Scanner) scanUntil(w *Word, starts, ends rune) {
	var scan func(bool)
	scan = func(top bool) {
		s.read()
		for s.char != ends && !s.done() {
			s.str.WriteRune(s.char)
			if s.char == starts {
				scan(false)
			}
			s.read()
		}
		if !top {
			s.str.WriteRune(s.char)
		}
	}
	scan(true)
	if s.char != ends {
		w.Type = Illegal
	}
	w.Literal = s.str.String()
}

func (s *Scanner) scanBraces(w *Word) {
	w.Type = Block
	s.scanUntil(w, lcurly, rcurly)
}

func (s *Scanner) scanScript(w *Word) {
	w.Type = Script
	s.scanUntil(w, lsquare, rsquare)
}

func (s *Scanner) scanLiteral(w *Word, isDone func(rune) bool) {
	defer s.unread()
	for !isDone(s.char) {
		s.str.WriteRune(s.escape())
		s.read()
	}
	w.Type = Literal
	w.Literal = s.str.String()
}

func (s *Scanner) escape() rune {
	if s.char != backslash {
		return s.char
	}
	switch s.peek() {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', dollar, space, backslash, dquote:
		s.read()
		return escapeChar(s.char)
	default:
		return s.char
	}
}

func (s *Scanner) scanComment(w *Word) {
	s.read()
	s.skipBlank()
	s.read()
	for s.char != nl {
		s.str.WriteRune(s.char)
		s.read()
	}
	w.Type = Comment
	w.Literal = s.str.String()
}

func (s *Scanner) done() bool {
	return s.char == null
}

func (s *Scanner) prepare() Word {
	s.str.Reset()
	s.read()
	w := Word{
		Position: s.pos,
	}
	if s.char == null {
		w.Type = EOF
	}
	return w
}

func (s *Scanner) read() {
	s.prev = s.pos
	if s.next >= len(s.input) {
		s.char = null
		s.pos.Column = 1
		s.pos.Line++
		return
	}
	old := s.char
	r, size := utf8.DecodeRune(s.input[s.next:])
	s.curr = s.next
	s.next += size
	s.char = r

	if old == nl {
		s.pos.Line++
		s.pos.Column = 1
	} else {
		s.pos.Column++
	}
}

func (s *Scanner) unread() {
	if s.next == 0 || s.done() {
		return
	}
	s.pos = s.prev
	r, size := utf8.DecodeLastRune(s.input[:s.curr])
	s.char = r
	s.next = s.curr
	s.curr -= size
}

func (s *Scanner) peek() rune {
	r, _ := utf8.DecodeRune(s.input[s.next:])
	if r == utf8.RuneError {
		return null
	}
	return r
}

func (s *Scanner) skip(accept func(rune) bool) {
	if !accept(s.char) {
		return
	}
	defer s.unread()
	for accept(s.char) {
		s.read()
	}
}

func (s *Scanner) skipBlank() {
	s.skip(isBlank)
}

func (s *Scanner) skipNL() {
	s.read()
	s.skip(isNL)
}

const (
	null       = 0
	underscore = '_'
	space      = ' '
	tab        = '\t'
	nl         = '\n'
	cr         = '\r'
	semicolon  = ';'
	backslash  = '\\'
	dollar     = '$'
	dash       = '#'
	dquote     = '"'
	squote     = '\''
	dot        = '.'
	equal      = '='
	bang       = '!'
	langle     = '<'
	rangle     = '>'
	plus       = '+'
	minus      = '-'
	star       = '*'
	slash      = '/'
	percent    = '%'
	pipe       = '|'
	ampersand  = '&'
	caret      = '^'
	question   = '?'
	colon      = ':'
	lparen     = '('
	rparen     = ')'
	lcurly     = '{'
	rcurly     = '}'
	lsquare    = '['
	rsquare    = ']'
	tilde      = '~'
)

func isBlank(c rune) bool {
	return c == space || c == tab
}

func isNL(c rune) bool {
	return c == nl
}

func isEOL(c rune) bool {
	return c == nl || c == semicolon
}

func isLower(c rune) bool {
	return c >= 'a' && c <= 'z'
}

func isUpper(c rune) bool {
	return c >= 'A' && c <= 'Z'
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isHexa(c rune) bool {
	return isDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func isOctal(c rune) bool {
	return c >= '0' && c <= '7'
}

func isBinary(c rune) bool {
	return c == '0' || c == '1'
}

func isLetter(c rune) bool {
	return isLower(c) || isUpper(c) || c == underscore
}

func isAlpha(c rune) bool {
	return isLetter(c) || isDigit(c)
}

func isVariable(c rune) bool {
	return c == dollar
}

func isScript(c rune) bool {
	return c == lsquare
}

func isQuote(c rune) bool {
	return c == dquote
}

func isOperator(c rune) bool {
	return c == plus || c == minus || c == star || c == slash ||
		c == percent || c == equal || c == bang || c == ampersand ||
		c == pipe || c == langle || c == rangle || c == caret ||
		c == tilde || c == question || c == colon
}

func isComment(c rune) bool {
	return c == dash
}

func isParen(c rune) bool {
	return c == lparen
}

func escapeChar(c rune) rune {
	switch c {
	case 'a':
		c = '\a'
	case 'b':
		c = '\b'
	case 'f':
		c = '\f'
	case 'n':
		c = nl
	case 'r':
		c = cr
	case 't':
		c = tab
	case 'v':
		c = '\v'
	}
	return c
}

func isDelimiter(c rune) bool {
	// return isBlank(c) || isParen(c) || isComment(c) || isQuotedDelimiter(c) || isEOL(c)
	return isBlank(c) || isComment(c) || isQuotedDelimiter(c) || isEOL(c)
}

func isQuotedDelimiter(c rune) bool {
	return isQuote(c) || isVariable(c) || isScript(c) || c == null
}

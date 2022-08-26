package glob

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"unicode/utf8"
)

var ErrPattern = errors.New("bad pattern")

func Match(str, pattern string) bool {
	if pattern == "" {
		return true
	}
	return match(str, pattern)
}

func Filter(list []string, pattern string) []string {
	if pattern == "" {
		return list
	}
	var res []string
	for i := range list {
		if ok := match(list[i], pattern); !ok {
			continue
		}
		res = append(res, list[i])
	}
	return res
}

const (
	caret     = '^'
	bang      = '!'
	star      = '*'
	question  = '?'
	lsquare   = '['
	rsquare   = ']'
	backslash = '/'
	dash      = '-'
)

func match(str, pattern string) bool {
	if str == pattern {
		return true
	}
	drain := func(rs *strings.Reader) string {
		var buf bytes.Buffer
		io.Copy(&buf, rs)
		return buf.String()
	}
	var (
		rs  = strings.NewReader(pattern)
		ptr int
	)
	for rs.Len() > 0 {
		r, _, err := rs.ReadRune()
		if err != nil {
			break
		}
		switch r {
		case star:
			want := skipChar(rs, star)
			if want == utf8.RuneError {
				return true
			}
			if want == question || want == lsquare {
				rs.UnreadRune()
				return match(str[ptr:], drain(rs))
			}
			z, ok := matchStar(str[ptr:], want)
			if !ok {
				return ok
			}
			ptr += z
		case question:
			x, z := utf8.DecodeRuneInString(str[ptr:])
			if x == utf8.RuneError {
				return false
			}
			ptr += z
		case lsquare:
			z, ok := matchRange(str[ptr:], rs)
			if !ok {
				return false
			}
			ptr += z
		default:
			if r == backslash {
				r, _, _ = rs.ReadRune()
			}
			x, z := utf8.DecodeRuneInString(str[ptr:])
			if x != r {
				return false
			}
			ptr += z
		}
	}
	return true
}

func matchRange(str string, rs *strings.Reader) (int, bool) {
	x, z := utf8.DecodeRuneInString(str)
	if x == utf8.RuneError {
		return 0, false
	}
	r, _, err := rs.ReadRune()
	if err != nil {
		return 0, false
	}
	var (
		old    rune
		rev    = r == caret || r == bang
		accept = func(p, n rune) bool {
			ok := x >= p && x <= n
			if rev {
				ok = !ok
			}
			return ok
		}
	)
	for rs.Len() > 0 {
		r, _, _ = rs.ReadRune()
		if r == rsquare {
			return z, false
		}
		if r == dash {
			n, _, err := rs.ReadRune()
			if err != nil {
				return 0, false
			}
			if accept(old, n) {
				break
			}
		} else {
			if accept(r, r) {
				break
			}
			old = r
		}
	}
	skipChar(rs, rsquare)
	rs.UnreadRune()
	return z, true
}

func matchStar(str string, char rune) (int, bool) {
	var ptr int
	for ptr < len(str) {
		x, z := utf8.DecodeRuneInString(str[ptr:])
		ptr += z
		if x == char {
			return ptr, true
		}
	}
	return ptr, false
}

func skipChar(rs *strings.Reader, char rune) rune {
	for {
		r, _, err := rs.ReadRune()
		if err != nil {
			return utf8.RuneError
		}
		if r != char {
			return r
		}
	}
}

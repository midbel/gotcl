package strutil

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrRange    = errors.New("invalid interval")
	ErrBoundary = errors.New("invalid boundary")
)

func Range(str string, fst, lst int) (string, error) {
	if err := checkBoundary(str, fst, lst); err != nil {
		return "", err
	}
	return str[fst:lst], nil
}

func Map(str string, list []string, nocase bool) (string, error) {
	if len(list) == 0 || len(str) == 0 {
		return str, nil
	}
	if len(list)%2 != 0 {
		return "", fmt.Errorf("invalid array length")
	}
	for old, pos := 0, 0; pos < len(str); {
		old = pos
		for i := 0; i < len(list); i += 2 {
			pat, rep := list[i], list[i+1]
			txt := str[pos:]
			if nocase {
				pat = strings.ToLower(pat)
				txt = strings.ToLower(txt)
			}
			if strings.HasPrefix(txt, pat) {
				str = str[:pos] + rep + str[pos+len(pat):]
				pos += len(rep)
				break
			}
		}
		if old == pos {
			pos++
		}
	}
	return str, nil
}

func Replace(str, pat string, fst, lst int) (string, error) {
	if err := checkBoundary(str, fst, lst); err != nil {
		return "", err
	}
	if pat == "" {
		return str[:fst] + str[lst+1:], nil
	}
	return strings.ReplaceAll(str, str[fst:lst+1], pat), nil
}

func Reverse(str string) string {
	var (
		in   = []rune(str)
		list = make([]rune, len(in))
	)
	for i, j := len(in)-1, 0; i >= 0; i-- {
		list[j] = in[i]
		j++
	}
	return string(list)
}

func checkBoundary(str string, fst, lst int) error {
	if fst < 0 {
		return ErrBoundary
	}
	if lst > len(str) {
		return fmt.Errorf("%w: %d > %d", ErrBoundary, lst, len(str))
	}
	if lst < fst {
		return rangeError(fst, lst)
	}
	return nil
}

func rangeError(fst, lst int) error {
	return fmt.Errorf("%w: %d < %d", ErrRange, fst, lst)
}

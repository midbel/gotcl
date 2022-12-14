package strutil

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"
)

var (
	ErrRange    = errors.New("invalid interval")
	ErrBoundary = errors.New("invalid boundary")
)

func LongestCommonPrefix(str []string) string {
	if len(str) == 0 {
		return ""
	}
	sort.Slice(str, func(i, j int) bool {
		return len(str[i]) < len(str[j])
	})
	var (
		first = str[0]
		rest  [][]rune
	)
	for _, s := range str[1:] {
		rest = append(rest, []rune(s))
	}
	var n int
	for i, c := range first {
		var ok bool
		for _, str := range rest {
			if ok = str[i] == c; !ok {
				break
			}
		}
		if !ok {
			break
		}
		n += utf8.RuneLen(c)
	}
	if n > 0 {
		return first[:n]
	}
	return ""
}

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

func ToLower(str string, fst, lst int) (string, error) {
	if err := checkBoundary(str, fst, lst); err != nil {
		return "", err
	}
	return strings.ToLower(str[fst:lst]), nil
}

func ToUpper(str string, fst, lst int) (string, error) {
	if err := checkBoundary(str, fst, lst); err != nil {
		return "", err
	}
	return strings.ToUpper(str[fst:lst]), nil
}

func ToTitle(str string, fst, lst int) (string, error) {
	if err := checkBoundary(str, fst, lst); err != nil {
		return "", err
	}
	return strings.ToTitle(str[fst:lst]), nil
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

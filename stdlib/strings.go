package stdlib

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/midbel/slices"
)

func RunString() CommandFunc {
	set := map[string]CommandFunc{
		"cat":       runCat,
		"repeat":    runRepeat,
		"replace":   runReplace,
		"trim":      runTrim,
		"trimleft":  runTrimleft,
		"trimright": runTrimright,
		"equal":     runEqual,
		"compare":   runCompare,
		"first":     runFirst,
		"last":      runLast,
		"index":     runIndex,
		"range":     runRange,
		"map":       runMap,
		"tolower":   runLower,
		"toupper":   runUpper,
		"totitle":   runTitle,
		"reverse":   runReverse,
		"length":    runLength,
		"match":     runMatch,
		// "is":        nil,
	}
	return makeEnsemble("string", set)
}

func runCat(i Interpreter, args []string) (string, error) {
	var str strings.Builder
	for i := range args {
		str.WriteString(args[i])
	}
	return str.String(), nil
}

func runCompare(i Interpreter, args []string) (string, error) {
	var (
		length int
		nocase bool
	)
	args, err := parseArgs("compare", args, func(set *flag.FlagSet) (int, bool) {
		set.BoolVar(&nocase, "nocase", nocase, "nocase")
		set.IntVar(&length, "length", length, "length")
		return 2, true
	})
	if err != nil {
		return "", err
	}
	var (
		fst = cutStr(slices.Fst(args), length)
		lst = cutStr(slices.Lst(args), length)
	)
	if nocase {
		fst = strings.ToLower(fst)
		lst = strings.ToLower(lst)
	}
	res := strings.Compare(fst, lst)
	return strconv.Itoa(res), nil
}

func runEqual(i Interpreter, args []string) (string, error) {
	var (
		length int
		nocase bool
	)
	args, err := parseArgs("equal", args, func(set *flag.FlagSet) (int, bool) {
		set.BoolVar(&nocase, "nocase", nocase, "nocase")
		set.IntVar(&length, "length", length, "length")
		return 2, true
	})
	if err != nil {
		return "", err
	}
	var (
		fst = cutStr(slices.Fst(args), length)
		lst = cutStr(slices.Lst(args), length)
	)
	if nocase {
		fst = strings.ToLower(fst)
		lst = strings.ToLower(lst)
	}
	return strconv.FormatBool(fst == lst), nil
}

func runFirst(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("first", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, false
	})
	if err != nil {
		return "", err
	}
	var index int
	if ix := slices.At(args, 2); ix != "" {
		index, err = strconv.Atoi(ix)
		if err != nil {
			return "", err
		}
		if int(index) > len(slices.Fst(args)) {
			return "", ErrIndex
		}
	}
	x := strings.Index(slices.Fst(args)[index:], slices.Snd(args))
	return strconv.Itoa(x), nil
}

func runLast(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("last", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, false
	})
	if err != nil {
		return "", err
	}
	var index int
	if ix := slices.At(args, 2); ix != "" {
		index, err = strconv.Atoi(ix)
		if err != nil {
			return "", err
		}
		if int(index) > len(slices.Fst(args)) {
			return "", ErrIndex
		}
	}
	x := strings.Index(slices.Fst(args)[index:], slices.Snd(args))
	return strconv.Itoa(x), nil
}

func runIndex(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("index", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, true
	})
	if err != nil {
		return "", err
	}
	index, err := strconv.Atoi(slices.Snd(args))
	if err != nil {
		return "", err
	}
	str := slices.Fst(args)
	if int(index) > len(str) {
		return "", ErrIndex
	}
	return str[index : index+1], nil
}

func runRange(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("range", args, func(_ *flag.FlagSet) (int, bool) {
		return 3, true
	})
	if err != nil {
		return "", err
	}
	first, err := strconv.Atoi(slices.At(args, 1))
	if err != nil {
		return "", err
	}
	last, err := strconv.Atoi(slices.At(args, 2))
	if err != nil {
		return "", err
	}
	res := slices.Fst(args)
	return res[first:last], nil
}

func runMap(i Interpreter, args []string) (string, error) {
	var nocase bool
	args, err := parseArgs("map", args, func(set *flag.FlagSet) (int, bool) {
		set.BoolVar(&nocase, "nocase", nocase, "nocase")
		return 2, true
	})
	if err != nil {
		return "", err
	}
	var (
		str  = slices.Fst(args)
		list = strings.Fields(slices.Snd(args))
	)
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
			if strings.HasPrefix(strings.ToLower(str[pos:]), strings.ToLower(pat)) {
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

func runReplace(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("replace", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, false
	})
	if err != nil {
		return "", err
	}
	first, err := strconv.Atoi(slices.At(args, 1))
	if err != nil {
		return "", err
	}
	last, err := strconv.Atoi(slices.At(args, 2))
	if err != nil {
		return "", err
	}
	res := slices.Fst(args)
	if r := slices.At(args, 4); r != "" {
		res = strings.ReplaceAll(res, res[first:last], r)
	} else {
		res = res[:first] + res[last+1:]
	}
	return res, nil
}

func runRepeat(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("repeat", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, true
	})
	if err != nil {
		return "", err
	}
	count, err := strconv.Atoi(slices.Snd(args))
	if err != nil {
		return "", err
	}
	if count < 0 {
		return "", fmt.Errorf("negative count")
	}
	return strings.Repeat(slices.Fst(args), int(count)), nil
}

func runMatch(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func runLower(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("lower", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return strings.ToLower(slices.Fst(args)), nil
}

func runUpper(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("upper", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return strings.ToUpper(slices.Fst(args)), nil
}

func runTitle(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("title", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return strings.ToTitle(slices.Fst(args)), nil
}

func runTrim(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("trim", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, false
	})
	if err != nil {
		return "", err
	}
	return strings.Trim(slices.Fst(args), slices.At(args, 1)), nil
}

func runTrimleft(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("trimleft", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, false
	})
	if err != nil {
		return "", err
	}
	return strings.TrimLeft(slices.Fst(args), slices.At(args, 1)), nil
}

func runTrimright(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("lower", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, false
	})
	if err != nil {
		return "", err
	}
	return strings.TrimRight(slices.Fst(args), slices.At(args, 1)), nil
}

func runReverse(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("reverse", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	var (
		str  = []rune(slices.Fst(args))
		list = make([]rune, len(str))
	)
	for i, j := len(str)-1, 0; i >= 0; i-- {
		list[j] = str[i]
		j++
	}
	return string(list), nil
}

func runLength(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("length", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	n := len(slices.Fst(args))
	return strconv.Itoa(n), nil
}

func cutStr(str string, n int) string {
	if n > 0 && len(str) >= n {
		str = str[:n]
	}
	return str
}

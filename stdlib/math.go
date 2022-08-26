package stdlib

import (
	"flag"
	"fmt"
	"math"
	"strconv"

	"github.com/midbel/gotcl/stdlib/conv"
	"github.com/midbel/slices"
)

func RunAbs(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunAcos(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunAsin(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunAtan(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunAtan2(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunCos(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunCosh(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunSin(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunSinh(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunTan(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunTanh(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunHypot(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunBool(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunDouble(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunEntier(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunCeil(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunFloor(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunRound(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunFmod(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunInt(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunExp(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunLog(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunLog10(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunMax(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunMin(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunRaise(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunRand(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunSrand(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunIsqrt(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunSqrt(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunWide(i Interpreter, args []string) (string, error) {
	return "", ErrImplemented
}

func RunIncr(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("incr", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, false
	})
	if err != nil {
		return "", err
	}
	by := 1
	if n, err := strconv.Atoi(slices.Snd(args)); err == nil {
		by = n
	}
	return i.Do(slices.Fst(args), func(str string) (string, error) {
		n, err := strconv.Atoi(str)
		if err != nil {
			return "", err
		}
		n += by
		return strconv.Itoa(n), nil
	})
}

func RunDecr(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("decr", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, false
	})
	if err != nil {
		return "", err
	}
	by := 1
	if n, err := strconv.Atoi(slices.Snd(args)); err == nil {
		by = n
	}
	return i.Do(slices.Fst(args), func(str string) (string, error) {
		n, err := strconv.Atoi(str)
		if err != nil {
			return "", err
		}
		n -= by
		return strconv.Itoa(n), nil
	})
}

func RunAdd(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("+", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, false
	})
	if err != nil {
		return "", err
	}
	return applyOp(args, func(left, right float64) (float64, error) {
		return left + right, nil
	})
}

func RunSub(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("-", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, false
	})
	if err != nil {
		return "", err
	}
	if len(args) == 1 {

	}
	return applyOp(args, func(left, right float64) (float64, error) {
		return left - right, nil
	})
}

func RunMul(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("*", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, false
	})
	if err != nil {
		return "", err
	}
	return applyOp(args, func(left, right float64) (float64, error) {
		return left * right, nil
	})
}

func RunDiv(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("/", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, false
	})
	if err != nil {
		return "", err
	}
	return applyOp(args, func(left, right float64) (float64, error) {
		if right == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return left / right, nil
	})
}

func RunMod(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("%", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, false
	})
	if err != nil {
		return "", err
	}
	return applyOp(args, func(left, right float64) (float64, error) {
		if right == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return math.Mod(left, right), nil
	})
}

func RunPow(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("**", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, false
	})
	if err != nil {
		return "", err
	}
	return applyOp(args, func(left, right float64) (float64, error) {
		return math.Pow(left, right), nil
	})
}

func RunEq(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("==", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, false
	})
	if err != nil {
		return "", err
	}
	return compareOp(args, func(left, right float64) (bool, error) {
		return left == right, nil
	})
}

func RunNe(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("!=", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, false
	})
	if err != nil {
		return "", err
	}
	return compareOp(args, func(left, right float64) (bool, error) {
		return left != right, nil
	})
}

func RunLt(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("<", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, false
	})
	if err != nil {
		return "", err
	}
	return compareOp(args, func(left, right float64) (bool, error) {
		return left < right, nil
	})
}

func RunLe(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("<=", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, false
	})
	if err != nil {
		return "", err
	}
	return compareOp(args, func(left, right float64) (bool, error) {
		return left <= right, nil
	})
}

func RunGt(i Interpreter, args []string) (string, error) {
	args, err := parseArgs(">", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, false
	})
	if err != nil {
		return "", err
	}
	return compareOp(args, func(left, right float64) (bool, error) {
		return left > right, nil
	})
}

func RunGe(i Interpreter, args []string) (string, error) {
	args, err := parseArgs(">=", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, false
	})
	if err != nil {
		return "", err
	}
	return compareOp(args, func(left, right float64) (bool, error) {
		return left >= right, nil
	})
}

func RunNot(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("!", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	res, err := strconv.ParseFloat(slices.Fst(args), 64)
	if res == 0 {
		return conv.False(), nil
	}
	return conv.True(), nil
}

type compFunc func(float64, float64) (bool, error)

func compareOp(args []string, cmp compFunc) (string, error) {
	res, err := strconv.ParseFloat(slices.Fst(args), 64)
	if err != nil {
		return "", err
	}
	var ok bool
	for _, val := range slices.Rest(args) {
		val, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return "", err
		}
		ok, err = cmp(res, val)
		if err != nil {
			return "", err
		}
		if !ok {
			break
		}
		res = val
	}
	return conv.Bool(ok), nil
}

type applyFunc func(float64, float64) (float64, error)

func applyOp(args []string, apply applyFunc) (string, error) {
	res, err := strconv.ParseFloat(slices.Fst(args), 64)
	if err != nil {
		return "", err
	}
	for _, val := range slices.Rest(args) {
		val, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return "", err
		}
		res, err = apply(res, val)
		if err != nil {
			return "", err
		}
	}
	var str string
	if math.Floor(res) == res {
		str = strconv.FormatInt(int64(res), 10)
	} else {
		str = strconv.FormatFloat(res, 'g', -1, 64)
	}
	return str, nil
}

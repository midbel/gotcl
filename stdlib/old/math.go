package stdlib

import (
	"flag"
	"fmt"
	"math"
	"strconv"

	"github.com/midbel/gotcl/stdlib/conv"
	"github.com/midbel/slices"
)

type Rander interface {
	Rand() int
	Seed(int)
}

func RunAbs(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("abs", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return withNumber(slices.Fst(args), math.Abs)
}

func RunAcos(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("acos", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return withNumber(slices.Fst(args), math.Acos)
}

func RunAsin(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("asin", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return withNumber(slices.Fst(args), math.Asin)
}

func RunAtan(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("atan", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return withNumber(slices.Fst(args), math.Atan)
}

func RunAtan2(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("atan2", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, true
	})
	if err != nil {
		return "", err
	}
	return withNumber2(slices.Fst(args), slices.Snd(args), math.Atan2)
}

func RunCos(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("cos", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return withNumber(slices.Fst(args), math.Cos)
}

func RunCosh(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("cosh", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return withNumber(slices.Fst(args), math.Cosh)
}

func RunSin(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("sin", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return withNumber(slices.Fst(args), math.Sin)
}

func RunSinh(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("sinh", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return withNumber(slices.Fst(args), math.Sinh)
}

func RunTan(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("tan", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return withNumber(slices.Fst(args), math.Tan)
}

func RunTanh(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("tanh", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return withNumber(slices.Fst(args), math.Tanh)
}

func RunHypot(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("hypot", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, true
	})
	if err != nil {
		return "", err
	}
	return withNumber2(slices.Fst(args), slices.Snd(args), math.Hypot)
}

func RunBool(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("bool", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	v, err := strconv.ParseFloat(slices.Fst(args), 64)
	if err != nil {
		return conv.False(), nil
	}
	return conv.Bool(v != 0), nil
}

func RunDouble(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("double", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	_, err = strconv.ParseFloat(slices.Fst(args), 64)
	return slices.Fst(args), err
}

func RunEntier(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("entier", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	_, err = strconv.ParseInt(slices.Fst(args), 0, 64)
	return slices.Fst(args), err
}

func RunCeil(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("ceil", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return withNumber(slices.Fst(args), math.Ceil)
}

func RunFloor(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("floor", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return withNumber(slices.Fst(args), math.Floor)
}

func RunRound(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("round", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return withNumber(slices.Fst(args), math.Round)
}

func RunFmod(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("fmod", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, true
	})
	if err != nil {
		return "", err
	}
	return withNumber2(slices.Fst(args), slices.Snd(args), math.Mod)
}

func RunInt(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("int", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	v, err := strconv.ParseFloat(slices.Fst(args), 64)
	return strconv.Itoa(int(v)), err
}

func RunExp(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("exp", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return withNumber(slices.Fst(args), math.Exp)
}

func RunLog(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("log", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return withNumber(slices.Fst(args), math.Log)
}

func RunLog10(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("log10", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return withNumber(slices.Fst(args), math.Log10)
}

func RunMax(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("max", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, false
	})
	if err != nil {
		return "", err
	}
	return compareNumber(args, math.Max)
}

func RunMin(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("min", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, false
	})
	if err != nil {
		return "", err
	}
	return compareNumber(args, math.Min)
}

func RunRaise(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("pow", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, true
	})
	if err != nil {
		return "", err
	}
	return withNumber2(slices.Fst(args), slices.Snd(args), math.Pow)
}

func RunRand(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("rand", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, true
	})
	if err != nil {
		return "", err
	}
	r, ok := i.(Rander)
	if !ok {
		return "", fmt.Errorf("interpreter can not generate random number")
	}
	x := r.Rand()
	return strconv.Itoa(x), nil
}

func RunSrand(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("srand", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	r, ok := i.(Rander)
	if !ok {
		return "", fmt.Errorf("interpreter can not be seeded")
	}
	seed, err := strconv.Atoi(slices.Fst(args))
	if err != nil {
		return "", err
	}
	r.Seed(seed)
	return "", nil
}

func RunIsqrt(i Interpreter, args []string) (string, error) {
	_, err := parseArgs("isqrt", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return "", ErrImplemented
}

func RunSqrt(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("sqrt", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return withNumber(slices.Fst(args), math.Sqrt)
}

func RunWide(i Interpreter, args []string) (string, error) {
	_, err := parseArgs("wide", args, func(_ *flag.FlagSet) (int, bool) {
		return 1, true
	})
	if err != nil {
		return "", err
	}
	return "", ErrImplemented
}

func withNumber2(left, right string, do func(float64, float64) float64) (string, error) {
	v0, err := strconv.ParseFloat(left, 64)
	if err != nil {
		return "", err
	}
	v1, err := strconv.ParseFloat(right, 64)
	if err != nil {
		return "", err
	}
	res := do(v0, v1)
	return strconv.FormatFloat(res, 'g', -1, 64), nil
}

func compareNumber(args []string, cmp func(float64, float64) float64) (string, error) {
	var (
		val float64
		tmp float64
		err error
	)
	for i := 0; i < len(args); i++ {
		tmp, err = strconv.ParseFloat(args[i], 64)
		if err != nil {
			break
		}
		val = cmp(tmp, val)
	}
	return strconv.FormatFloat(val, 'g', -1, 64), err
}

func withNumber(str string, do func(float64) float64) (string, error) {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return "", err
	}
	val = do(val)
	return strconv.FormatFloat(val, 'g', -1, 64), nil
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
	if err != nil {
		return "", err
	}
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

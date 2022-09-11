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

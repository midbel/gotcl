package stdlib

import (
	"fmt"
	"math"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/slices"
)

func RunIncr() Executer {
	return Builtin{
		Name:     "incr",
		Arity:    1,
		Variadic: true,
		Safe:     true,
		Run:      runIncr,
	}
}

func RunDecr() Executer {
	return Builtin{
		Name:     "decr",
		Arity:    1,
		Variadic: true,
		Safe:     true,
		Run:      runDecr,
	}
}

func RunAdd() Executer {
	return Builtin{
		Name:     "+",
		Arity:    1,
		Variadic: true,
		Safe:     true,
		Run:      runAdd,
	}
}

func RunSub() Executer {
	return Builtin{
		Name:     "-",
		Arity:    1,
		Variadic: true,
		Safe:     true,
		Run:      runSub,
	}
}

func RunMul() Executer {
	return Builtin{
		Name:     "*",
		Arity:    0,
		Variadic: true,
		Safe:     true,
		Run:      runMul,
	}
}

func RunDiv() Executer {
	return Builtin{
		Name:     "/",
		Arity:    1,
		Variadic: true,
		Safe:     true,
		Run:      runDiv,
	}
}

func RunMod() Executer {
	return Builtin{
		Name:     "%",
		Arity:    1,
		Variadic: true,
		Safe:     true,
		Run:      runMod,
	}
}

func RunPow() Executer {
	return Builtin{
		Name:     "**",
		Arity:    1,
		Variadic: true,
		Safe:     true,
		Run:      runPow,
	}
}

func RunNot() Executer {
	return Builtin{
		Name:  "!",
		Arity: 1,
		Safe:  true,
		Run:   runNot,
	}
}

func RunEq() Executer {
	return Builtin{
		Name:     "==",
		Arity:    1,
		Variadic: true,
		Safe:     true,
		Run:      runEq,
	}
}

func RunNe() Executer {
	return Builtin{
		Name:     "!=",
		Arity:    1,
		Variadic: true,
		Safe:     true,
		Run:      runNe,
	}
}

func RunLt() Executer {
	return Builtin{
		Name:     "<",
		Arity:    1,
		Variadic: true,
		Safe:     true,
		Run:      runLt,
	}
}

func RunLe() Executer {
	return Builtin{
		Name:     "<=",
		Arity:    1,
		Variadic: true,
		Safe:     true,
		Run:      runLe,
	}
}

func RunGt() Executer {
	return Builtin{
		Name:     ">",
		Arity:    1,
		Variadic: true,
		Safe:     true,
		Run:      runGt,
	}
}

func RunGe() Executer {
	return Builtin{
		Name:     ">=",
		Arity:    1,
		Variadic: true,
		Safe:     true,
		Run:      runGe,
	}
}

func RunAbs() Executer {
	return Builtin{}
}

func RunAcos() Executer {
	return Builtin{}
}

func RunAsin() Executer {
	return Builtin{}
}

func RunAtan() Executer {
	return Builtin{}
}

func RunAtan2() Executer {
	return Builtin{}
}

func RunCos() Executer {
	return Builtin{}
}

func RunCosh() Executer {
	return Builtin{}
}

func RunSin() Executer {
	return Builtin{}
}

func RunSinh() Executer {
	return Builtin{}
}

func RunTan() Executer {
	return Builtin{}
}

func RunTanh() Executer {
	return Builtin{}
}

func RunHypot() Executer {
	return Builtin{}
}

func RunBool() Executer {
	return Builtin{}
}

func RunDouble() Executer {
	return Builtin{}
}

func RunEntier() Executer {
	return Builtin{}
}

func RunCeil() Executer {
	return Builtin{}
}

func RunFloor() Executer {
	return Builtin{}
}

func RunRound() Executer {
	return Builtin{}
}

func RunFmod() Executer {
	return Builtin{}
}

func RunInt() Executer {
	return Builtin{}
}

func RunExp() Executer {
	return Builtin{}
}

func RunLog() Executer {
	return Builtin{}
}

func RunLog10() Executer {
	return Builtin{}
}

func RunMax() Executer {
	return Builtin{}
}

func RunMin() Executer {
	return Builtin{}
}

func RunRaise() Executer {
	return Builtin{}
}

func RunRand() Executer {
	return Builtin{}
}

func RunSrand() Executer {
	return Builtin{}
}

func RunIsqrt() Executer {
	return Builtin{}
}

func RunSqrt() Executer {
	return Builtin{}
}

func RunWide() Executer {
	return Builtin{}
}

func runIncr(i Interpreter, args []env.Value) (env.Value, error) {
	var step int
	if v := slices.Snd(args); v != nil {
		x, err := env.ToInt(v)
		if err != nil {
			return nil, err
		}
		step = x
	} else {
		step++
	}
	v, err := i.Resolve(slices.Fst(args).String())
	if err != nil {
		return nil, err
	}
	n, err := env.ToInt(v)
	if err != nil {
		return nil, err
	}
	res := env.Int(int64(n - step))
	i.Define(slices.Fst(args).String(), res)
	return res, nil
}

func runDecr(i Interpreter, args []env.Value) (env.Value, error) {
	var step int
	if v := slices.Snd(args); v != nil {
		x, err := env.ToInt(v)
		if err != nil {
			return nil, err
		}
		step = x
	} else {
		step++
	}
	v, err := i.Resolve(slices.Fst(args).String())
	if err != nil {
		return nil, err
	}
	n, err := env.ToInt(v)
	if err != nil {
		return nil, err
	}
	res := env.Int(int64(n - step))
	i.Define(slices.Fst(args).String(), res)
	return res, nil
}

func runAdd(i Interpreter, args []env.Value) (env.Value, error) {
	return withNumbers(args, func(fst, lst float64) (float64, error) {
		return fst + lst, nil
	})
}

func runSub(i Interpreter, args []env.Value) (env.Value, error) {
	if len(args) == 1 {
		f, err := env.ToFloat(slices.Fst(args))
		if err != nil {
			return nil, err
		}
		return env.Float(-f), nil
	}
	return withNumbers(args, func(fst, lst float64) (float64, error) {
		return fst - lst, nil
	})
}

func runMul(i Interpreter, args []env.Value) (env.Value, error) {
	if len(args) == 0 {
		return env.Int(1), nil
	}
	return withNumbers(args, func(fst, lst float64) (float64, error) {
		return fst * lst, nil
	})
}

func runDiv(i Interpreter, args []env.Value) (env.Value, error) {
	if len(args) == 1 {
		args = slices.Prepend(env.Float(1), args)
	}
	return withNumbers(args, func(fst, lst float64) (float64, error) {
		if lst == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return fst / lst, nil
	})
}

func runMod(i Interpreter, args []env.Value) (env.Value, error) {
	return withNumbers(args, func(fst, lst float64) (float64, error) {
		if lst == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return math.Mod(fst, lst), nil
	})
}

func runPow(i Interpreter, args []env.Value) (env.Value, error) {
	return withNumbers(args, func(fst, lst float64) (float64, error) {
		return math.Pow(fst, lst), nil
	})
}

func runEq(i Interpreter, args []env.Value) (env.Value, error) {
	return withCompare(args, func(fst, lst float64) (bool, error) {
		return fst == lst, nil
	})
}

func runNe(i Interpreter, args []env.Value) (env.Value, error) {
	return withCompare(args, func(fst, lst float64) (bool, error) {
		return fst != lst, nil
	})
}

func runLt(i Interpreter, args []env.Value) (env.Value, error) {
	return withCompare(args, func(fst, lst float64) (bool, error) {
		return fst < lst, nil
	})
}

func runLe(i Interpreter, args []env.Value) (env.Value, error) {
	return withCompare(args, func(fst, lst float64) (bool, error) {
		return fst <= lst, nil
	})
}

func runGt(i Interpreter, args []env.Value) (env.Value, error) {
	return withCompare(args, func(fst, lst float64) (bool, error) {
		return fst > lst, nil
	})
}

func runGe(i Interpreter, args []env.Value) (env.Value, error) {
	return withCompare(args, func(fst, lst float64) (bool, error) {
		return fst > lst, nil
	})
}

func runNot(i Interpreter, args []env.Value) (env.Value, error) {
	f, err := env.ToFloat(slices.Fst(args))
	if err != nil {
		return nil, err
	}
	if f == 0 {
		return env.False(), nil
	}
	return env.True(), nil
}

func withCompare(args []env.Value, do func(float64, float64) (bool, error)) (env.Value, error) {
	r, err := env.ToFloat(slices.Fst(args))
	if err != nil {
		return nil, err
	}
	var ok bool
	for _, v := range slices.Rest(args) {
		c, err := env.ToFloat(v)
		if err != nil {
			return nil, err
		}
		ok, err = do(r, c)
		if !ok || err != nil {
			return env.False(), err
		}
		r = c
	}
	return env.Bool(ok), nil
}

func withNumbers(args []env.Value, do func(float64, float64) (float64, error)) (env.Value, error) {
	res, err := env.ToFloat(slices.Fst(args))
	if err != nil {
		return nil, err
	}
	for _, v := range slices.Rest(args) {
		c, err := env.ToFloat(v)
		if err != nil {
			return nil, err
		}
		res, err = do(res, c)
		if err != nil {
			return nil, err
		}
	}
	return env.Float(res), nil
}

func runAbs(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runAcos(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runAsin(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runAtan(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runAtan2(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runCos(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runCosh(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runSin(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runSinh(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runTan(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runTanh(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runHypot(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runBool(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runDouble(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runEntier(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runCeil(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runFloor(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runRound(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runFmod(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runInt(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runExp(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runLog(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runLog10(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runMax(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runMin(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runRaise(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runRand(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runSrand(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runIsqrt(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runSqrt(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func runWide(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

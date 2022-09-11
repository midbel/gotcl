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
	return Builtin{
		Name:  "abs",
		Arity: 1,
		Safe:  true,
		Run:   runAbs,
	}
}

func RunAcos() Executer {
	return Builtin{
		Name:  "acos",
		Arity: 1,
		Safe:  true,
		Run:   runAcos,
	}
}

func RunAsin() Executer {
	return Builtin{
		Name:  "asin",
		Arity: 1,
		Safe:  true,
		Run:   runAsin,
	}
}

func RunAtan() Executer {
	return Builtin{
		Name:  "atan",
		Arity: 1,
		Safe:  true,
		Run:   runAtan,
	}
}

func RunAtan2() Executer {
	return Builtin{
		Name:  "atan2",
		Arity: 2,
		Safe:  true,
		Run:   runAtan2,
	}
}

func RunCos() Executer {
	return Builtin{
		Name:  "cos",
		Arity: 1,
		Safe:  true,
		Run:   runCos,
	}
}

func RunCosh() Executer {
	return Builtin{
		Name:  "cosh",
		Arity: 1,
		Safe:  true,
		Run:   runCosh,
	}
}

func RunSin() Executer {
	return Builtin{
		Name:  "sin",
		Arity: 1,
		Safe:  true,
		Run:   runSin,
	}
}

func RunSinh() Executer {
	return Builtin{
		Name:  "sinh",
		Arity: 1,
		Safe:  true,
		Run:   runSinh,
	}
}

func RunTan() Executer {
	return Builtin{
		Name:  "tan",
		Arity: 1,
		Safe:  true,
		Run:   runTan,
	}
}

func RunTanh() Executer {
	return Builtin{
		Name:  "tanh",
		Arity: 1,
		Safe:  true,
		Run:   runTanh,
	}
}

func RunHypot() Executer {
	return Builtin{
		Name:  "hypot",
		Arity: 2,
		Safe:  true,
		Run:   runHypot,
	}
}

func RunBool() Executer {
	return Builtin{
		Name:  "bool",
		Arity: 1,
		Safe:  true,
		Run:   runBool,
	}
}

func RunDouble() Executer {
	return Builtin{
		Name:  "double",
		Arity: 1,
		Safe:  true,
		Run:   runDouble,
	}
}

func RunEntier() Executer {
	return Builtin{
		Name:  "entier",
		Arity: 1,
		Safe:  true,
		Run:   runEntier,
	}
}

func RunCeil() Executer {
	return Builtin{
		Name:  "ceil",
		Arity: 1,
		Safe:  true,
		Run:   runCeil,
	}
}

func RunFloor() Executer {
	return Builtin{
		Name:  "floor",
		Arity: 1,
		Safe:  true,
		Run:   runFloor,
	}
}

func RunRound() Executer {
	return Builtin{
		Name:  "round",
		Arity: 1,
		Safe:  true,
		Run:   runRound,
	}
}

func RunFmod() Executer {
	return Builtin{
		Name:  "fmod",
		Arity: 1,
		Safe:  true,
		Run:   runFmod,
	}
}

func RunInt() Executer {
	return Builtin{
		Name:  "int",
		Arity: 1,
		Safe:  true,
		Run:   runInt,
	}
}

func RunExp() Executer {
	return Builtin{
		Name:  "exp",
		Arity: 1,
		Safe:  true,
		Run:   runExp,
	}
}

func RunLog() Executer {
	return Builtin{
		Name:  "log",
		Arity: 1,
		Safe:  true,
		Run:   runLog,
	}
}

func RunLog10() Executer {
	return Builtin{
		Name:  "log10",
		Arity: 1,
		Safe:  true,
		Run:   runLog10,
	}
}

func RunMax() Executer {
	return Builtin{
		Name:     "max",
		Variadic: true,
		Safe:     true,
		Run:      runMax,
	}
}

func RunMin() Executer {
	return Builtin{
		Name:     "min",
		Variadic: true,
		Safe:     true,
		Run:      runMin,
	}
}

func RunRaise() Executer {
	return Builtin{
		Name:  "pow",
		Arity: 2,
		Safe:  true,
		Run:   runRaise,
	}
}

func RunRand() Executer {
	return Builtin{
		Name: "rand",
		Safe: true,
		Run:  runRand,
	}
}

func RunSrand() Executer {
	return Builtin{
		Name:  "srand",
		Arity: 1,
		Safe:  true,
		Run:   runSrand,
	}
}

func RunIsqrt() Executer {
	return Builtin{
		Name:  "isqrt",
		Arity: 1,
		Safe:  true,
		Run:   runIsqrt,
	}
}

func RunSqrt() Executer {
	return Builtin{
		Name:  "sqrt",
		Arity: 1,
		Safe:  true,
		Run:   runSqrt,
	}
}

func RunWide() Executer {
	return Builtin{
		Name:  "wide",
		Arity: 1,
		Safe:  true,
		Run:   runWide,
	}
}

func RunDegree() Executer {
	return Builtin{
		Name:  "deg",
		Arity: 1,
		Safe:  true,
		Run:   runDegree,
	}
}

func RunRadian() Executer {
	return Builtin{
		Name:  "rad",
		Arity: 1,
		Safe:  true,
		Run:   runRadian,
	}
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

func runDegree(i Interpreter, args []env.Value) (env.Value, error) {
	return withFloat(slices.Fst(args), func(f float64) float64 {
		return f * (180 / math.Pi)
	})
}

func runRadian(i Interpreter, args []env.Value) (env.Value, error) {
	return withFloat(slices.Fst(args), func(f float64) float64 {
		return f * (math.Pi / 180)
	})
}

func runAbs(i Interpreter, args []env.Value) (env.Value, error) {
	return withFloat(slices.Fst(args), math.Abs)
}

func runAcos(i Interpreter, args []env.Value) (env.Value, error) {
	return withFloat(slices.Fst(args), math.Acos)
}

func runAsin(i Interpreter, args []env.Value) (env.Value, error) {
	return withFloat(slices.Fst(args), math.Asin)
}

func runAtan(i Interpreter, args []env.Value) (env.Value, error) {
	return withFloat(slices.Fst(args), math.Atan)
}

func runAtan2(i Interpreter, args []env.Value) (env.Value, error) {
	return withFloat2(slices.Fst(args), slices.Snd(args), math.Atan2)
}

func runCos(i Interpreter, args []env.Value) (env.Value, error) {
	return withFloat(slices.Fst(args), math.Cos)
}

func runCosh(i Interpreter, args []env.Value) (env.Value, error) {
	return withFloat(slices.Fst(args), math.Cosh)
}

func runSin(i Interpreter, args []env.Value) (env.Value, error) {
	return withFloat(slices.Fst(args), math.Sin)
}

func runSinh(i Interpreter, args []env.Value) (env.Value, error) {
	return withFloat(slices.Fst(args), math.Sinh)
}

func runTan(i Interpreter, args []env.Value) (env.Value, error) {
	return withFloat(slices.Fst(args), math.Tan)
}

func runTanh(i Interpreter, args []env.Value) (env.Value, error) {
	return withFloat(slices.Fst(args), math.Tan)
}

func runHypot(i Interpreter, args []env.Value) (env.Value, error) {
	return withFloat2(slices.Fst(args), slices.Snd(args), math.Hypot)
}

func runBool(i Interpreter, args []env.Value) (env.Value, error) {
	f, err := env.ToFloat(slices.Fst(args))
	if err != nil {
		return env.False(), nil
	}
	return env.Bool(f != 0), nil
}

func runDouble(i Interpreter, args []env.Value) (env.Value, error) {
	return slices.Fst(args).ToNumber()
}

func runEntier(i Interpreter, args []env.Value) (env.Value, error) {
	return slices.Fst(args).ToNumber()
}

func runCeil(i Interpreter, args []env.Value) (env.Value, error) {
	return withFloat(slices.Fst(args), math.Ceil)
}

func runFloor(i Interpreter, args []env.Value) (env.Value, error) {
	return withFloat(slices.Fst(args), math.Floor)
}

func runRound(i Interpreter, args []env.Value) (env.Value, error) {
	return withFloat(slices.Fst(args), math.Round)
}

func runFmod(i Interpreter, args []env.Value) (env.Value, error) {
	return withFloat2(slices.Fst(args), slices.Snd(args), math.Mod)
}

func runInt(i Interpreter, args []env.Value) (env.Value, error) {
	return slices.Fst(args).ToNumber()
}

func runExp(i Interpreter, args []env.Value) (env.Value, error) {
	return withFloat(slices.Fst(args), math.Exp)
}

func runLog(i Interpreter, args []env.Value) (env.Value, error) {
	return withFloat(slices.Fst(args), math.Log)
}

func runLog10(i Interpreter, args []env.Value) (env.Value, error) {
	return withFloat(slices.Fst(args), math.Log10)
}

func runMax(i Interpreter, args []env.Value) (env.Value, error) {
	return cmpFloat(args, math.Max)
}

func runMin(i Interpreter, args []env.Value) (env.Value, error) {
	return cmpFloat(args, math.Min)
}

func runRaise(i Interpreter, args []env.Value) (env.Value, error) {
	return withFloat2(slices.Fst(args), slices.Snd(args), math.Pow)
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
	return withFloat(slices.Fst(args), math.Sqrt)
}

func runWide(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func withFloat(v env.Value, do func(float64) float64) (env.Value, error) {
	f, err := env.ToFloat(v)
	if err != nil {
		return nil, err
	}
	return env.Float(do(f)), nil
}

func withFloat2(fst, snd env.Value, do func(float64, float64) float64) (env.Value, error) {
	var (
		v1, err1 = env.ToFloat(fst)
		v2, err2 = env.ToFloat(snd)
	)
	if err := hasError(err1, err2); err != nil {
		return nil, err
	}
	return env.Float(do(v1, v2)), nil
}

func cmpFloat(args []env.Value, cmp func(float64, float64) float64) (env.Value, error) {
	var (
		val float64
		tmp float64
		err error
	)
	for i := 0; i < len(args); i++ {
		tmp, err = env.ToFloat(args[i])
		if err != nil {
			break
		}
		val = cmp(val, tmp)
	}
	return env.Float(val), nil
}

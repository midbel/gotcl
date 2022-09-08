package stdlib

import (
	"time"

	"github.com/midbel/gotcl/env"
	"github.com/midbel/gotcl/stdlib/clock"
	"github.com/midbel/slices"
)

func MakeClock() Executer {
	e := Ensemble{
		Name: "clock",
		List: []Executer{
			Builtin{
				Name: "format",
				Run:  clockFormat,
			},
			Builtin{
				Name:  "scan",
				Arity: 2,
				Run:   clockScan,
			},
			Builtin{
				Name: "add",
				Run:  clockAdd,
			},
			Builtin{
				Name: "clicks",
				Run:  clockClicks,
			},
			Builtin{
				Name: "seconds",
				Run:  clockSeconds,
			},
			Builtin{
				Name: "milliseconds",
				Run:  clockMillis,
			},
			Builtin{
				Name: "microseconds",
				Run:  clockMicros,
			},
		},
	}
	return sortEnsembleCommands(e)
}

func clockScan(i Interpreter, args []env.Value) (env.Value, error) {
	unix, err := clock.Scan(slices.Fst(args).String(), slices.Snd(args).String())
	if err != nil {
		return nil, err
	}
	return env.Int(unix), nil
}

func clockFormat(i Interpreter, args []env.Value) (env.Value, error) {
	unix, err := env.ToInt(slices.Fst(args))
	if err != nil {
		return nil, err
	}
	str, err := clock.Format(int64(unix), slices.Snd(args).String())
	return env.Str(str), nil
}

func clockAdd(i Interpreter, args []env.Value) (env.Value, error) {
	return nil, nil
}

func clockClicks(i Interpreter, args []env.Value) (env.Value, error) {
	n := time.Now()
	return env.Int(n.UnixNano()), nil
}

func clockSeconds(i Interpreter, args []env.Value) (env.Value, error) {
	n := time.Now()
	return env.Int(n.Unix()), nil
}

func clockMillis(i Interpreter, args []env.Value) (env.Value, error) {
	var (
		n = time.Now()
		u = n.Unix() * 1000
	)
	return env.Int(u), nil
}

func clockMicros(i Interpreter, args []env.Value) (env.Value, error) {
	var (
		n = time.Now()
		u = n.Unix() * (1000 * 1000)
	)
	return env.Int(u), nil
}

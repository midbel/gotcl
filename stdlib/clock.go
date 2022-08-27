package stdlib

import (
	"flag"
	"strconv"
	"time"

	"github.com/midbel/gotcl/stdlib/clock"
	"github.com/midbel/slices"
)

func RunClock() CommandFunc {
	set := map[string]CommandFunc{
		"format":       runTimeFormat,
		"scan":         runTimeScan,
		"add":          runTimeAdd,
		"clicks":       runTimeClicks,
		"seconds":      runSeconds,
		"milliseconds": runMillisSeconds,
		"microseconds": runMicrosSeconds,
	}
	return makeEnsemble("clock", set)
}

func runTimeAdd(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("add", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, true
	})
	if err != nil {
		return "", err
	}
	unix, err := strconv.Atoi(slices.Fst(args))
	if err != nil {
		return "", err
	}
	offset, err := strconv.Atoi(slices.Snd(args))
	if err != nil {
		return "", err
	}

	now := time.Unix(unix, 0).Add()
	switch unit := slices.Lst(args); unit {
	case "seconds", "":
		now = now.Add(time.Duration(offset) * time.Second)
	case "minutes":
		now = now.Add(time.Duration(offset) * time.Minute)
	case "hours":
		now = now.Add(time.Duration(offset) * time.Hour)
	case "days":
		now = now.AddDate(0, 0, offset)
	case "weeks":
		now = now.AddDate(0, 0, offset*7)
	case "months":
		now = now.AddDate(0, offset, 0)
	case "years":
		now = now.AddDate(offset, 0, 0)
	default:
		return "", fmt.Errorf("%s: unknown time unit", unit)
	}
	return strconv.FormatInt(now.Unix(), 10), ErrImplemented
}

func runTimeFormat(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("format", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, true
	})
	if err != nil {
		return "", err
	}
	unix, err := strconv.ParseInt(slices.Fst(args), 0, 64)
	if err != nil {
		return "", err
	}
	return clock.Format(unix, slices.Snd(args))
}

func runTimeScan(i Interpreter, args []string) (string, error) {
	args, err := parseArgs("scan", args, func(_ *flag.FlagSet) (int, bool) {
		return 2, true
	})
	if err != nil {
		return "", err
	}
	unix, err := clock.Scan(slices.Fst(args), slices.Snd(args))
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(unix, 10), nil
}

func runTimeClicks(i Interpreter, args []string) (string, error) {
	_, err := parseArgs("clicks", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, true
	})
	if err != nil {
		return "", err
	}
	n := time.Now()
	return strconv.FormatInt(n.UnixNano(), 10), nil
}

func runSeconds(i Interpreter, args []string) (string, error) {
	_, err := parseArgs("seconds", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, true
	})
	if err != nil {
		return "", err
	}
	n := time.Now()
	return strconv.FormatInt(n.Unix(), 10), nil
}

func runMillisSeconds(i Interpreter, args []string) (string, error) {
	_, err := parseArgs("milliseconds", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, true
	})
	if err != nil {
		return "", err
	}
	var (
		n = time.Now()
		u = n.Unix() / (1000 * 1000)
	)
	return strconv.FormatInt(u, 10), nil
}

func runMicrosSeconds(i Interpreter, args []string) (string, error) {
	_, err := parseArgs("microseconds", args, func(_ *flag.FlagSet) (int, bool) {
		return 0, true
	})
	if err != nil {
		return "", err
	}
	var (
		n = time.Now()
		u = n.Unix() / 1000 * 1000
	)
	return strconv.FormatInt(u, 10), nil
}

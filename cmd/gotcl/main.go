package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/midbel/gotcl/interp"
	"github.com/midbel/gotcl/stdlib"
)

func main() {
	var (
		echo   = flag.Bool("e", false, "print command to be execute")
		dry    = flag.Bool("n", false, "dry run")
		config = flag.String("i", "", "init file")
	)
	flag.Parse()

	i := interp.New()
	if i, ok := i.(*interp.Interp); ok {
		i.Echo = *echo
	}
	if *dry {
		// TBD
		return
	}
	if *config != "" {
		_, err := executeFile(i, *config)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(5)
		}
	}
	if err := runFile(i, flag.Arg(0)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runFile(i stdlib.Interpreter, file string) error {
	res, err := executeFile(i, file)
	if err != nil {
		return err
	}
	if res != "" {
		fmt.Fprintln(os.Stdout, res)
	}
	if i, ok := i.(*interp.Interp); ok {
		fmt.Println("---")
		fmt.Println("command executed:", i.Count)
	}
	return nil
}

func executeFile(i stdlib.Interpreter, file string) (string, error) {
	r, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer r.Close()
	return i.Execute(r)
}

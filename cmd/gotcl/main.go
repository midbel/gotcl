package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/midbel/gotcl/interp"
	"github.com/midbel/gotcl/stdlib"
)

func main() {
	var (
		config = flag.String("i", "", "init file")
	)
	flag.Parse()

	i := interp.New()
	if *config != "" {
		_, err := executeFile(i, *config)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(5)
		}
	}
	var err error
	if flag.NArg() == 0 {
		err = runREPL(i)
	} else {
		err = runFile(i, flag.Arg(0))
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runREPL(i stdlib.Interpreter) error {
	scan := bufio.NewScanner(os.Stdin)
	io.WriteString(os.Stdin, ">>> ")
	for scan.Scan() {
		line := scan.Text()
		res, err := i.Execute(strings.NewReader(line))
		if err != nil {
			return err
		}
		if res != "" {
			fmt.Println(strings.TrimSpace(res))
		}
		io.WriteString(os.Stdin, ">>> ")
	}
	return scan.Err()
}

func runFile(i stdlib.Interpreter, file string) error {
	res, err := executeFile(i, file)
	if err == nil && res != "" {
		fmt.Fprintln(os.Stdout, res)
	}
	return err
}

func executeFile(i stdlib.Interpreter, file string) (string, error) {
	r, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer r.Close()
	return i.Execute(r)
}

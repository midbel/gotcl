package main

import (
	"bufio"
	"bytes"
	"errors"
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

	i := interp.Interpret()
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

const (
	in  = "\x1b[1;97min [%3d]:\x1b[0m "
	ok  = "\x1b[1;92mout[%3d]:\x1b[0m %s"
	nok = "\x1b[1;91mout[%3d]:\x1b[0m %s"
)

func runREPL(i stdlib.Interpreter) error {
	var (
		buf  bytes.Buffer
		tmp  bytes.Buffer
		scan = bufio.NewScanner(os.Stdin)
		cmd  = 1
	)
	io.WriteString(os.Stdout, fmt.Sprintf(in, cmd))
	for scan.Scan() {
		line := scan.Text()
		if strings.TrimSpace(line) == "" {
			cmd++
			io.WriteString(os.Stdout, fmt.Sprintf(in, cmd))
			continue
		}
		buf.WriteString(line + "\n")
		res, err := i.Execute(io.TeeReader(&buf, &tmp))
		if err != nil {
			if errors.Is(err, stdlib.ErrExit) {
				break
			}
			if errors.Is(err, interp.ErrIncomplete) {
				io.Copy(&buf, &tmp)
				fmt.Fprint(os.Stdout, "... ")
				continue
			}
			fmt.Fprintf(os.Stderr, fmt.Sprintf(nok, cmd, err))
			fmt.Fprintln(os.Stderr)
		} else if res != nil {
			fmt.Fprintf(os.Stdout, ok, cmd, strings.TrimSpace(res.String()))
			fmt.Fprintln(os.Stdout)
		}
		cmd++
		io.WriteString(os.Stdout, fmt.Sprintf(in, cmd))
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
	val, err := i.Execute(r)
	if err != nil {
		return "", err
	}
	return val.String(), nil
}

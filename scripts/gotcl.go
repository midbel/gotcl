package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/midbel/gotcl/interp"
)

func main() {
	flag.Parse()

	r, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer r.Close()

	i := interp.Interpret()
	v, err := i.Execute(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
	if v != nil {
		fmt.Fprintf(os.Stdout, "%[1]s (%[1]T)", v, v)
		fmt.Fprintln(os.Stdout)
	}
}

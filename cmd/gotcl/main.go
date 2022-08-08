package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/midbel/gotcl/interp"
)

func main() {
	var (
		echo = flag.Bool("e", false, "print command to be execute")
		dry  = flag.Bool("n", false, "dry run")
	)
	flag.Parse()

	r, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer r.Close()

	i := interp.New()
	if i, ok := i.(*interp.Interp); ok {
		i.Echo = *echo
	}
	if *dry {
		// TBD
		return
	}
	res, err := i.Execute(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if res != "" {
		fmt.Fprintln(os.Stdout, res)
	}
	if i, ok := i.(*interp.Interp); ok {
		fmt.Println("---")
		fmt.Println("command executed:", i.Count)
	}
}

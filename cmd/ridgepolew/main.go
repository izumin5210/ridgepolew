package main

import (
	"context"
	"fmt"
	"os"

	"github.com/izumin5210/ridgepolew"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	return ridgepolew.NewDefaultRidgepole().Exec(context.Background(), os.Args)
}

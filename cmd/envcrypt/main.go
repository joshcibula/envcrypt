package main

import (
	"fmt"
	"os"

	"github.com/yourorg/envcrypt/cmd/envcrypt/commands"
)

func main() {
	if err := commands.Root().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

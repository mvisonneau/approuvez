package main

import (
	"os"

	"github.com/mvisonneau/approuvez/internal/cli"
)

var version = ""

func main() {
	cli.Run(version, os.Args)
}

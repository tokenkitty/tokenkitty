// Command tokenkitty is the collector that reads local AI coding agent usage.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tokenkitty/tokenkitty/internal/buildinfo"
)

func main() {
	version := flag.Bool("version", false, "print the version and exit")
	flag.Parse()

	if *version {
		fmt.Println(buildinfo.String("tokenkitty"))
		return
	}

	fmt.Fprintln(os.Stderr, "tokenkitty: no commands yet; see --version")
	os.Exit(2)
}

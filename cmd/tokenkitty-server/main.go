// Command tokenkitty-server runs the Tokenkitty API, ingest pipeline and dashboard.
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
		fmt.Println(buildinfo.String("tokenkitty-server"))
		return
	}

	fmt.Fprintln(os.Stderr, "tokenkitty-server: not implemented yet; see --version")
	os.Exit(2)
}

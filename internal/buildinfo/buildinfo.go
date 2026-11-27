// Package buildinfo holds version information stamped in at build time.
package buildinfo

import "fmt"

// Version is set with -ldflags "-X github.com/tokenkitty/tokenkitty/internal/buildinfo.Version=v1.2.3".
var Version = "dev"

// String returns a one-line version string for the named binary.
func String(binary string) string {
	return fmt.Sprintf("%s %s", binary, Version)
}

// Command doctor checks that the tools needed to develop Tokenkitty are installed
// and new enough. Run it with `task doctor`.
package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

type tool struct {
	name    string
	args    []string // version command; nil means only check that it is on PATH
	min     [2]int   // minimum major.minor; zero means any version
	install string
}

var tools = []tool{
	{name: "git", args: []string{"--version"}, install: "winget install Git.Git"},
	{name: "go", args: []string{"version"}, min: [2]int{1, 23}, install: "winget install GoLang.Go"},
	{name: "node", args: []string{"--version"}, min: [2]int{22, 0}, install: "winget install OpenJS.NodeJS.LTS"},
	{name: "npm", args: []string{"--version"}, install: "comes with Node"},
	{name: "task", args: []string{"--version"}, min: [2]int{3, 0}, install: "winget install Task.Task"},
	{name: "golangci-lint", args: []string{"version"}, min: [2]int{2, 0}, install: "go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0"},
	{name: "docker", args: []string{"--version"}, install: "winget install Docker.DockerDesktop"},
	{name: "sqlc", args: []string{"version"}, min: [2]int{1, 0}, install: "go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0"},
	{name: "migrate", install: "go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@v4.20.1"},
	{name: "air", install: "go install github.com/air-verse/air@latest"},
}

var versionRe = regexp.MustCompile(`(\d+)\.(\d+)(?:\.(\d+))?`)

// parseVersion returns the first major.minor found in s.
func parseVersion(s string) (major, minor int, ok bool) {
	m := versionRe.FindStringSubmatch(s)
	if m == nil {
		return 0, 0, false
	}
	major, _ = strconv.Atoi(m[1])
	minor, _ = strconv.Atoi(m[2])
	return major, minor, true
}

// atLeast reports whether major.minor >= floor.
func atLeast(major, minor int, floor [2]int) bool {
	if major != floor[0] {
		return major > floor[0]
	}
	return minor >= floor[1]
}

// check returns a status line and whether the tool passes.
func check(ctx context.Context, t tool) (string, bool) {
	path, err := exec.LookPath(t.name)
	if err != nil {
		return "missing  (" + t.install + ")", false
	}
	if t.args == nil {
		return "ok       " + path, true
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, path, t.args...).CombinedOutput()
	if err != nil {
		return fmt.Sprintf("error    running %s: %v", t.name, err), false
	}
	major, minor, ok := parseVersion(string(out))
	if !ok {
		return "unknown  could not read version", t.min == [2]int{}
	}
	got := fmt.Sprintf("%d.%d", major, minor)
	if t.min != [2]int{} && !atLeast(major, minor, t.min) {
		return fmt.Sprintf("too old  %s, need %d.%d+ (%s)", got, t.min[0], t.min[1], t.install), false
	}
	return "ok       " + got, true
}

func main() {
	ctx := context.Background()
	failed := 0
	for _, t := range tools {
		status, ok := check(ctx, t)
		if !ok {
			failed++
		}
		fmt.Printf("%-14s %s\n", t.name, status)
	}
	if failed > 0 {
		fmt.Printf("\n%d tool(s) need attention\n", failed)
		os.Exit(1)
	}
	fmt.Println("\nall good")
}

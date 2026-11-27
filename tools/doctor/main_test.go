package main

import "testing"

func TestParseVersion(t *testing.T) {
	tests := []struct {
		in           string
		major, minor int
		ok           bool
	}{
		{in: "go version go1.27.4 windows/amd64", major: 1, minor: 27, ok: true},
		{in: "v24.22.1", major: 24, minor: 22, ok: true},
		{in: "3.54.0", major: 3, minor: 54, ok: true},
		{in: "golangci-lint has version 2.15.0 built with go1.27.0", major: 2, minor: 15, ok: true},
		{in: "Docker version 28.3.2, build 578ccf6", major: 28, minor: 3, ok: true},
		{in: "dev", ok: false},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			major, minor, ok := parseVersion(tt.in)
			if ok != tt.ok || major != tt.major || minor != tt.minor {
				t.Errorf("parseVersion(%q) = %d, %d, %v; want %d, %d, %v", tt.in, major, minor, ok, tt.major, tt.minor, tt.ok)
			}
		})
	}
}

func TestAtLeast(t *testing.T) {
	tests := []struct {
		name         string
		major, minor int
		min          [2]int
		want         bool
	}{
		{name: "equal", major: 1, minor: 23, min: [2]int{1, 23}, want: true},
		{name: "newer minor", major: 1, minor: 27, min: [2]int{1, 23}, want: true},
		{name: "older minor", major: 1, minor: 22, min: [2]int{1, 23}, want: false},
		{name: "newer major", major: 3, minor: 0, min: [2]int{2, 9}, want: true},
		{name: "older major", major: 1, minor: 64, min: [2]int{2, 0}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := atLeast(tt.major, tt.minor, tt.min); got != tt.want {
				t.Errorf("atLeast(%d, %d, %v) = %v, want %v", tt.major, tt.minor, tt.min, got, tt.want)
			}
		})
	}
}

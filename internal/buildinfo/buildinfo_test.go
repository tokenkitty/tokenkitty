package buildinfo

import "testing"

func TestString(t *testing.T) {
	tests := []struct {
		name    string
		version string
		binary  string
		want    string
	}{
		{name: "default", version: "dev", binary: "tokenkitty", want: "tokenkitty dev"},
		{name: "stamped", version: "v0.1.0", binary: "tokenkitty-server", want: "tokenkitty-server v0.1.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := Version
			t.Cleanup(func() { Version = old })
			Version = tt.version

			if got := String(tt.binary); got != tt.want {
				t.Errorf("String(%q) = %q, want %q", tt.binary, got, tt.want)
			}
		})
	}
}

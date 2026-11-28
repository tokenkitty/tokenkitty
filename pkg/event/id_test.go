package event

import (
	"testing"
	"time"
)

// The expected IDs were computed by an independent implementation of the same
// scheme. They must never change: a change would duplicate every stored event.
func TestIDsAreStable(t *testing.T) {
	day := time.Date(2026, 11, 27, 15, 4, 5, 0, time.UTC)
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"local", LocalID(SourceClaudeCodeLocal, "sess-1", "msg-1"), "6c2828da26fba62f87a4cc292a12e2663139ffe642e2192b1625a0c23e90734c"},
		{"aggregate", AggregateID(SourceCursorAPI, "user-42", day, "claude-sonnet-4-5"), "86441f5d95376f7d76e17012cb4ab7f9464e8b835f818bda608d99563af39853"},
		{"turn", TurnID(SourceCursorAPI, "evt-9"), "f00c02b319a6bbf344963e8a29c9bb6eda14d9cb0287f8f163fe72066df1d345"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %s, want %s", tt.got, tt.want)
			}
		})
	}
}

func TestIDSeparatorsCannotCollide(t *testing.T) {
	a := LocalID(SourceClaudeCodeLocal, "a|b", "c")
	b := LocalID(SourceClaudeCodeLocal, "a", "b|c")
	if a == b {
		t.Fatalf("different inputs gave the same ID %s", a)
	}
	if a != "9fbed3d565184888016bfdaabc351f3e432ee524e381ded5e32618e533da5dcd" {
		t.Errorf("unexpected ID %s", a)
	}
}

func TestIDsDifferBySourceAndOrigin(t *testing.T) {
	ids := map[string]string{
		"claude local": LocalID(SourceClaudeCodeLocal, "s", "t"),
		"codex local":  LocalID(SourceCodexLocal, "s", "t"),
		"api turn":     TurnID(SourceCursorAPI, "t"),
	}
	seen := map[string]string{}
	for name, id := range ids {
		if other, ok := seen[id]; ok {
			t.Errorf("%s and %s share ID %s", name, other, id)
		}
		seen[id] = name
		if err := validID(id); err != nil {
			t.Errorf("%s: invalid ID %s: %v", name, id, err)
		}
	}
}

func TestValidID(t *testing.T) {
	tests := []struct {
		in string
		ok bool
	}{
		{in: "6c2828da26fba62f87a4cc292a12e2663139ffe642e2192b1625a0c23e90734c", ok: true},
		{in: "6C2828DA26FBA62F87A4CC292A12E2663139FFE642E2192B1625A0C23E90734C", ok: false},
		{in: "abc", ok: false},
		{in: "zz2828da26fba62f87a4cc292a12e2663139ffe642e2192b1625a0c23e90734c", ok: false},
		{in: "", ok: false},
	}
	for _, tt := range tests {
		if err := validID(tt.in); (err == nil) != tt.ok {
			t.Errorf("validID(%q) = %v, want ok=%v", tt.in, err, tt.ok)
		}
	}
}

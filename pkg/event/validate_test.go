package event

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// validLocal returns a valid local event that tests change one field at a time.
func validLocal() UsageEvent {
	cost := Money(1_250_000)
	return UsageEvent{
		EventID:          LocalID(SourceClaudeCodeLocal, "sess-1", "msg-1"),
		Source:           SourceClaudeCodeLocal,
		Tool:             ToolClaudeCode,
		Origin:           OriginLocal,
		Granularity:      GranularityTurn,
		OccurredAt:       time.Date(2026, 11, 27, 8, 15, 2, 0, time.UTC),
		SessionID:        "sess-1",
		Repo:             "github.com/acme/api",
		Model:            "claude-sonnet-4-5",
		Provider:         ProviderAnthropic,
		Tokens:           Tokens{Input: 1200, Output: 340, CacheRead: 50000, CacheWrite: 2000},
		Requests:         1,
		CostReportedUSD:  &cost,
		BillingMode:      BillingSubscription,
		CollectorVersion: "0.1.0",
		SourceVersion:    "2.1.3",
	}
}

func TestValidateAcceptsValidEvents(t *testing.T) {
	local := validLocal()

	api := validLocal()
	api.Source, api.Origin, api.Tool = SourceCursorAPI, OriginAPI, ToolCursor
	api.Granularity, api.Provider = GranularityDaily, ProviderCursor
	api.SessionID, api.Repo, api.CollectorVersion, api.CostReportedUSD = "", "", "", nil

	noRemote := validLocal()
	noRemote.Repo = "local/3f2a9c1b7d0e"

	for name, e := range map[string]UsageEvent{"local": local, "api": api, "local repo hash": noRemote} {
		t.Run(name, func(t *testing.T) {
			if err := e.Validate(); err != nil {
				t.Fatalf("Validate() = %v, want nil", err)
			}
		})
	}
}

func TestValidateRejects(t *testing.T) {
	berlin := time.FixedZone("CET", 3600)
	tests := []struct {
		name   string
		mutate func(e *UsageEvent)
		field  string
	}{
		{"bad event id", func(e *UsageEvent) { e.EventID = "abc" }, "event_id"},
		{"unknown source", func(e *UsageEvent) { e.Source = "chatgpt_web" }, "source"},
		{"origin mismatch", func(e *UsageEvent) { e.Origin = OriginAPI }, "origin"},
		{"tool mismatch", func(e *UsageEvent) { e.Tool = ToolCodex }, "tool"},
		{"tool other on local", func(e *UsageEvent) { e.Tool = ToolOther }, "tool"},
		{"unknown granularity", func(e *UsageEvent) { e.Granularity = "hourly" }, "granularity"},
		{"unknown provider", func(e *UsageEvent) { e.Provider = "acme" }, "provider"},
		{"empty billing mode", func(e *UsageEvent) { e.BillingMode = "" }, "billing_mode"},
		{"missing time", func(e *UsageEvent) { e.OccurredAt = time.Time{} }, "occurred_at"},
		{"non-utc time", func(e *UsageEvent) { e.OccurredAt = e.OccurredAt.In(berlin) }, "occurred_at"},
		{"negative input", func(e *UsageEvent) { e.Tokens.Input = -1 }, "tokens.input"},
		{"negative cache read", func(e *UsageEvent) { e.Tokens.CacheRead = -5 }, "tokens.cache_read"},
		{"negative requests", func(e *UsageEvent) { e.Requests = -1 }, "requests"},
		{"negative cost", func(e *UsageEvent) { c := Money(-1); e.CostReportedUSD = &c }, "cost_reported_usd"},
		{"session id with spaces", func(e *UsageEvent) { e.SessionID = "fix the login bug" }, "session_id"},
		{"session id too long", func(e *UsageEvent) { e.SessionID = strings.Repeat("a", 257) }, "session_id"},
		{"model uppercase", func(e *UsageEvent) { e.Model = "Claude-Sonnet" }, "model"},
		{"unix path as repo", func(e *UsageEvent) { e.Repo = "/home/dev/acme/api" }, "repo"},
		{"windows path as repo", func(e *UsageEvent) { e.Repo = `C:\Users\dev\acme` }, "repo"},
		{"home path as repo", func(e *UsageEvent) { e.Repo = "~/code/api" }, "repo"},
		{"relative path as repo", func(e *UsageEvent) { e.Repo = "./api" }, "repo"},
		{"repo without host", func(e *UsageEvent) { e.Repo = "acme/api/x" }, "repo"},
		{"repo too short", func(e *UsageEvent) { e.Repo = "github.com/acme" }, "repo"},
		{"bad local hash", func(e *UsageEvent) { e.Repo = "local/xyz" }, "repo"},
		{"local without collector version", func(e *UsageEvent) { e.CollectorVersion = "" }, "collector_version"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := validLocal()
			tt.mutate(&e)
			err := e.Validate()
			var fe *FieldError
			if !errors.As(err, &fe) {
				t.Fatalf("Validate() = %v, want a *FieldError", err)
			}
			if fe.Field != tt.field {
				t.Errorf("field = %q, want %q (%v)", fe.Field, tt.field, err)
			}
		})
	}
}

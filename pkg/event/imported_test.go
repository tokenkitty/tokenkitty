package event

import (
	"errors"
	"testing"
	"time"
)

func validImported() ImportedEvent {
	day := time.Date(2026, 11, 27, 0, 0, 0, 0, time.UTC)
	cost := Money(3_400_000)
	return ImportedEvent{
		UsageEvent: UsageEvent{
			EventID:         AggregateID(SourceCursorAPI, "user-42", day, "claude-sonnet-4-5"),
			Source:          SourceCursorAPI,
			Tool:            ToolCursor,
			Origin:          OriginAPI,
			Granularity:     GranularityDaily,
			OccurredAt:      day,
			Model:           "claude-sonnet-4-5",
			Provider:        ProviderCursor,
			Tokens:          Tokens{Input: 90_000, Output: 12_000},
			Requests:        40,
			CostReportedUSD: &cost,
			BillingMode:     BillingSeat,
		},
		Actor: Actor{Email: "dev@example.com", ExternalID: "user-42", Kind: ActorVendorUser},
	}
}

func TestImportedValidate(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(e *ImportedEvent)
		field  string // "" means valid
	}{
		{"valid", func(*ImportedEvent) {}, ""},
		{"no actor at all", func(e *ImportedEvent) { e.Actor = Actor{} }, ""},
		{"api key only", func(e *ImportedEvent) { e.Actor = Actor{ExternalID: "key_123", Kind: ActorAPIKey} }, ""},
		{"local origin", func(e *ImportedEvent) {
			e.Source, e.Origin, e.Tool, e.CollectorVersion = SourceCodexLocal, OriginLocal, ToolCodex, "0.1.0"
		}, "origin"},
		{"bad base event", func(e *ImportedEvent) { e.Tokens.Output = -1 }, "tokens.output"},
		{"email without at", func(e *ImportedEvent) { e.Actor.Email = "dev.example.com" }, "actor.email"},
		{"email uppercase", func(e *ImportedEvent) { e.Actor.Email = "Dev@Example.com" }, "actor.email"},
		{"email with space", func(e *ImportedEvent) { e.Actor.Email = "dev @example.com" }, "actor.email"},
		{"external id without kind", func(e *ImportedEvent) { e.Actor.Kind = "" }, "actor.kind"},
		{"unknown kind", func(e *ImportedEvent) { e.Actor.Kind = "device" }, "actor.kind"},
		{"kind without external id", func(e *ImportedEvent) { e.Actor.ExternalID = "" }, "actor.kind"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := validImported()
			tt.mutate(&e)
			err := e.Validate()
			if tt.field == "" {
				if err != nil {
					t.Fatalf("Validate() = %v, want nil", err)
				}
				return
			}
			var fe *FieldError
			if !errors.As(err, &fe) || fe.Field != tt.field {
				t.Fatalf("Validate() = %v, want a FieldError for %q", err, tt.field)
			}
		})
	}
}

// The ingest decoder must never accept an imported event's actor.
func TestDecodeRejectsImportedEvent(t *testing.T) {
	e := validImported()
	_, err := Decode(eventJSON(t, map[string]any{"actor": e.Actor}))
	if err == nil {
		t.Fatal("Decode accepted an actor")
	}
}

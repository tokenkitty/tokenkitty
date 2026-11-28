package event

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	want := validLocal()
	data, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	got, err := Decode(data)
	if err != nil {
		t.Fatalf("Decode(Marshal(e)) = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("round trip changed the event\n got: %+v\nwant: %+v", got, want)
	}
	if !strings.Contains(string(data), `"cost_reported_usd":"1.250000"`) {
		t.Errorf("cost not written as a 6-decimal string: %s", data)
	}
}

// eventJSON returns a valid event as a JSON object with extra keys merged in.
func eventJSON(t *testing.T, extra map[string]any) []byte {
	t.Helper()
	data, err := json.Marshal(validLocal())
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	for k, v := range extra {
		m[k] = v
	}
	out, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestDecodeRejectsUnknownFields(t *testing.T) {
	tests := []struct {
		name  string
		extra map[string]any
		field string
	}{
		{"prompt", map[string]any{"prompt": "CANARY_PROMPT_7f3a"}, "prompt"},
		{"device id", map[string]any{"device_id": "0192f0c4-0000-7000-8000-000000000000"}, "device_id"},
		{"actor", map[string]any{"actor": map[string]any{"email": "dev@example.com"}}, "actor"},
		{"cost estimated", map[string]any{"cost_estimated_usd": "1.000000"}, "cost_estimated_usd"},
		{"cwd", map[string]any{"cwd": "/home/dev/api"}, "cwd"},
		{"nested in tokens", map[string]any{"tokens": map[string]any{"input": 1, "text": "CANARY"}}, "text"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decode(eventJSON(t, tt.extra))
			if err == nil {
				t.Fatal("Decode accepted an unknown field")
			}
			if !strings.Contains(err.Error(), `unknown field "`+tt.field+`"`) {
				t.Errorf("error %q does not name field %q", err, tt.field)
			}
			if strings.Contains(err.Error(), "CANARY") {
				t.Errorf("error echoes field content: %q", err)
			}
		})
	}
}

func TestDecodeRejectsBadInput(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"cost as number", eventJSON(t, map[string]any{"cost_reported_usd": 1.25})},
		{"tokens as string", eventJSON(t, map[string]any{"tokens": map[string]any{"input": "12"}})},
		{"time not rfc3339", eventJSON(t, map[string]any{"occurred_at": "27/11/2026"})},
		{"trailing data", append(eventJSON(t, nil), []byte(` {}`)...)},
		{"not an object", []byte(`[1,2,3]`)},
		{"empty", []byte(``)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Decode(tt.data); err == nil {
				t.Fatal("Decode accepted bad input")
			}
		})
	}
}

func TestDecodeValidates(t *testing.T) {
	_, err := Decode(eventJSON(t, map[string]any{"origin": "api"}))
	var fe *FieldError
	if !errors.As(err, &fe) || fe.Field != "origin" {
		t.Fatalf("Decode = %v, want a FieldError for origin", err)
	}
}

func TestDecodeAcceptsMissingOptionalFields(t *testing.T) {
	data := `{
		"event_id": "` + LocalID(SourceCodexLocal, "s", "t") + `",
		"source": "codex_local", "tool": "codex", "origin": "local",
		"granularity": "turn", "occurred_at": "2026-11-27T08:15:02Z",
		"session_id": null, "repo": null, "cost_reported_usd": null,
		"provider": "openai", "tokens": {"input": 10},
		"requests": 1, "billing_mode": "unknown", "collector_version": "0.1.0"
	}`
	e, err := Decode([]byte(data))
	if err != nil {
		t.Fatalf("Decode = %v", err)
	}
	if e.Tokens.Output != 0 || e.CostReportedUSD != nil || e.Repo != "" {
		t.Errorf("missing fields not zero: %+v", e)
	}
}

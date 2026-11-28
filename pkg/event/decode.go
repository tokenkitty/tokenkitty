package event

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// Decode parses one event strictly and validates it.
//
// Unknown fields are rejected, including device_id and actor, which the
// server takes from the device token. This is the privacy guard: nothing
// outside the schema can be sent.
func Decode(data []byte) (UsageEvent, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()

	var e UsageEvent
	if err := dec.Decode(&e); err != nil {
		return UsageEvent{}, fmt.Errorf("decode event: %w", err)
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return UsageEvent{}, errors.New("decode event: unexpected data after the event")
	}
	if err := e.Validate(); err != nil {
		return UsageEvent{}, fmt.Errorf("invalid event: %w", err)
	}
	return e, nil
}

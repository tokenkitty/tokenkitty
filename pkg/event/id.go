package event

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"
)

// idScheme is mixed into every hash so the scheme can change later without
// colliding with old IDs.
const idScheme = "v1"

// newID hashes the parts with a length prefix on each one, so no separator
// inside a part (such as "|") can make two different inputs hash the same.
func newID(parts ...string) string {
	h := sha256.New()
	var n [8]byte
	for _, p := range parts {
		binary.BigEndian.PutUint64(n[:], uint64(len(p)))
		h.Write(n[:])
		h.Write([]byte(p))
	}
	return hex.EncodeToString(h.Sum(nil))
}

// LocalID returns the event_id for an event read from a local source.
// sessionID and turnID are the source's own opaque identifiers; reading the
// same turn again always gives the same ID.
func LocalID(source Source, sessionID, turnID string) string {
	return newID(idScheme, string(OriginLocal), string(source), sessionID, turnID)
}

// AggregateID returns the event_id for an aggregated row from a vendor API,
// such as one user's daily total for one model. actorKey is the vendor's user,
// email or API key ID; day is the vendor's calendar date.
func AggregateID(source Source, actorKey string, day time.Time, model string) string {
	return newID(idScheme, string(OriginAPI), string(source), actorKey, day.Format(time.DateOnly), model)
}

// TurnID returns the event_id for an API event that has its own vendor event ID.
func TurnID(source Source, vendorEventID string) string {
	return newID(idScheme, string(OriginAPI), string(source), vendorEventID)
}

// validID reports whether s looks like an event_id: 64 lowercase hex characters.
func validID(s string) error {
	if len(s) != sha256.Size*2 {
		return fmt.Errorf("must be %d hex characters", sha256.Size*2)
	}
	for _, r := range s {
		if (r < '0' || r > '9') && (r < 'a' || r > 'f') {
			return fmt.Errorf("must be lowercase hex")
		}
	}
	return nil
}

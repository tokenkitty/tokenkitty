package event

import (
	"strings"
)

// ActorKind says what an importer's external ID refers to.
type ActorKind string

// Actor kinds for external IDs.
const (
	ActorVendorUser ActorKind = "vendor_user"
	ActorAPIKey     ActorKind = "api_key"
)

const maxActorLen = 256

// Actor is who an imported row belongs to, as the vendor reports it.
// All fields may be empty; the server then uses a per-source catch-all actor.
type Actor struct {
	Email      string    `json:"email,omitempty"`
	ExternalID string    `json:"external_id,omitempty"`
	Kind       ActorKind `json:"kind,omitempty"`
}

// ImportedEvent is an event built on the server by a vendor API importer.
// It never comes over the ingest API, which accepts only UsageEvent.
type ImportedEvent struct {
	UsageEvent
	Actor Actor `json:"actor"`
}

// Validate checks the event and its actor.
func (e *ImportedEvent) Validate() error {
	if err := e.UsageEvent.Validate(); err != nil {
		return err
	}
	if e.Origin != OriginAPI {
		return fieldErr("origin", "imported events must have origin %q", OriginAPI)
	}

	a := e.Actor
	if a.Email != "" {
		if err := checkToken(a.Email, maxActorLen); err != nil {
			return fieldErr("actor.email", "%v", err)
		}
		local, domain, ok := strings.Cut(a.Email, "@")
		if !ok || local == "" || domain == "" || strings.Contains(domain, "@") {
			return fieldErr("actor.email", "is not an email address")
		}
		if a.Email != strings.ToLower(a.Email) {
			return fieldErr("actor.email", "must be lowercase")
		}
	}
	if err := checkToken(a.ExternalID, maxActorLen); err != nil {
		return fieldErr("actor.external_id", "%v", err)
	}
	switch {
	case a.ExternalID != "" && a.Kind != ActorVendorUser && a.Kind != ActorAPIKey:
		return fieldErr("actor.kind", "must be %q or %q when external_id is set", ActorVendorUser, ActorAPIKey)
	case a.ExternalID == "" && a.Kind != "":
		return fieldErr("actor.kind", "requires external_id")
	}
	return nil
}

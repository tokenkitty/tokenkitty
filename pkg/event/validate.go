package event

import (
	"fmt"
	"slices"
	"strings"
	"unicode"
)

// Length limits keep identifiers short so they cannot carry content.
const (
	maxSessionIDLen = 256
	maxRepoLen      = 512
	maxModelLen     = 128
	maxVersionLen   = 64
)

// FieldError says which field of an event is invalid and why.
type FieldError struct {
	Field   string
	Problem string
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Problem)
}

func fieldErr(field, format string, args ...any) error {
	return &FieldError{Field: field, Problem: fmt.Sprintf(format, args...)}
}

// Validate checks the event against the rules in the data model. It returns
// the first problem found as a *FieldError.
func (e *UsageEvent) Validate() error {
	if err := validID(e.EventID); err != nil {
		return fieldErr("event_id", "%v", err)
	}

	info, ok := sourceInfo[e.Source]
	if !ok {
		return fieldErr("source", "unknown source %q", e.Source)
	}
	if e.Origin != info.origin {
		return fieldErr("origin", "source %s requires origin %q, got %q", e.Source, info.origin, e.Origin)
	}
	if !slices.Contains(info.tools, e.Tool) {
		return fieldErr("tool", "source %s does not allow tool %q", e.Source, e.Tool)
	}
	if !granularities[e.Granularity] {
		return fieldErr("granularity", "unknown granularity %q", e.Granularity)
	}
	if !providers[e.Provider] {
		return fieldErr("provider", "unknown provider %q", e.Provider)
	}
	if !billingModes[e.BillingMode] {
		return fieldErr("billing_mode", "unknown billing mode %q", e.BillingMode)
	}

	if e.OccurredAt.IsZero() {
		return fieldErr("occurred_at", "is required")
	}
	if _, offset := e.OccurredAt.Zone(); offset != 0 {
		return fieldErr("occurred_at", "must be UTC")
	}

	for name, v := range map[string]int64{
		"tokens.input":       e.Tokens.Input,
		"tokens.output":      e.Tokens.Output,
		"tokens.cache_read":  e.Tokens.CacheRead,
		"tokens.cache_write": e.Tokens.CacheWrite,
		"tokens.reasoning":   e.Tokens.Reasoning,
		"requests":           e.Requests,
	} {
		if v < 0 {
			return fieldErr(name, "must not be negative")
		}
	}
	if e.CostReportedUSD != nil && *e.CostReportedUSD < 0 {
		return fieldErr("cost_reported_usd", "must not be negative")
	}

	if err := checkToken(e.SessionID, maxSessionIDLen); err != nil {
		return fieldErr("session_id", "%v", err)
	}
	if err := checkModel(e.Model); err != nil {
		return fieldErr("model", "%v", err)
	}
	if err := checkRepo(e.Repo); err != nil {
		return fieldErr("repo", "%v", err)
	}
	if err := checkToken(e.CollectorVersion, maxVersionLen); err != nil {
		return fieldErr("collector_version", "%v", err)
	}
	if err := checkToken(e.SourceVersion, maxVersionLen); err != nil {
		return fieldErr("source_version", "%v", err)
	}
	if e.Origin == OriginLocal && e.CollectorVersion == "" {
		return fieldErr("collector_version", "is required for local events")
	}
	return nil
}

// checkToken allows an empty value, or a short value with no spaces or control characters.
func checkToken(s string, maxLen int) error {
	if len(s) > maxLen {
		return fmt.Errorf("longer than %d characters", maxLen)
	}
	for _, r := range s {
		if unicode.IsSpace(r) || unicode.IsControl(r) {
			return fmt.Errorf("must not contain spaces or control characters")
		}
	}
	return nil
}

// checkModel requires a normalized, lowercase model ID.
func checkModel(s string) error {
	if err := checkToken(s, maxModelLen); err != nil {
		return err
	}
	if s != strings.ToLower(s) {
		return fmt.Errorf("must be lowercase")
	}
	return nil
}

// checkRepo allows "", "host/owner/repo" from a git remote, or "local/<12 hex>"
// for a repo without a remote. Anything that looks like a local path is rejected.
func checkRepo(s string) error {
	if s == "" {
		return nil
	}
	if err := checkToken(s, maxRepoLen); err != nil {
		return err
	}
	if hash, ok := strings.CutPrefix(s, "local/"); ok {
		if len(hash) != 12 || strings.Trim(hash, "0123456789abcdef") != "" {
			return fmt.Errorf("local repo must be local/ followed by 12 lowercase hex characters")
		}
		return nil
	}
	if strings.ContainsAny(s, `\:~`) || strings.HasPrefix(s, "/") || strings.HasPrefix(s, ".") {
		return fmt.Errorf("must be host/owner/repo, not a path")
	}
	parts := strings.Split(s, "/")
	if len(parts) < 3 || slices.Contains(parts, "") || slices.Contains(parts, "..") {
		return fmt.Errorf("must be host/owner/repo")
	}
	if !strings.Contains(parts[0], ".") && parts[0] != "localhost" {
		return fmt.Errorf("must start with a host name")
	}
	return nil
}

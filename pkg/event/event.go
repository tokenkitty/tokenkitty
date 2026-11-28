package event

import "time"

// Source identifies where an event was read from.
type Source string

// Sources. A local source is read by the collector; an api source is imported by the server.
const (
	SourceClaudeCodeLocal Source = "claude_code_local"
	SourceCodexLocal      Source = "codex_local"
	SourceGeminiCLILocal  Source = "gemini_cli_local"
	SourceOpenCodeLocal   Source = "opencode_local"
	SourceAnthropicAPI    Source = "anthropic_api"
	SourceCursorAPI       Source = "cursor_api"
	SourceCopilotAPI      Source = "copilot_api"
	SourceOpenAIAPI       Source = "openai_api"
)

// Tool is the AI coding tool the usage belongs to.
type Tool string

// Tools.
const (
	ToolClaudeCode Tool = "claude_code"
	ToolCodex      Tool = "codex"
	ToolGeminiCLI  Tool = "gemini_cli"
	ToolOpenCode   Tool = "opencode"
	ToolCursor     Tool = "cursor"
	ToolCopilot    Tool = "copilot"
	ToolOther      Tool = "other"
)

// Origin says whether an event came from a developer machine or a vendor API.
type Origin string

// Origins.
const (
	OriginLocal Origin = "local"
	OriginAPI   Origin = "api"
)

// Granularity is how much usage one event covers.
type Granularity string

// Granularities.
const (
	GranularityTurn    Granularity = "turn"
	GranularitySession Granularity = "session"
	GranularityDaily   Granularity = "daily"
)

// Provider is the company serving the model.
type Provider string

// Providers.
const (
	ProviderAnthropic Provider = "anthropic"
	ProviderOpenAI    Provider = "openai"
	ProviderGoogle    Provider = "google"
	ProviderCursor    Provider = "cursor"
	ProviderGitHub    Provider = "github"
	ProviderOther     Provider = "other"
)

// BillingMode says how the usage is paid for.
type BillingMode string

// Billing modes. For subscription usage the estimated cost is an API-equivalent value.
const (
	BillingAPI          BillingMode = "api"
	BillingSubscription BillingMode = "subscription"
	BillingSeat         BillingMode = "seat"
	BillingUnknown      BillingMode = "unknown"
)

// sourceInfo fixes the origin and the allowed tools for each source.
var sourceInfo = map[Source]struct {
	origin Origin
	tools  []Tool
}{
	SourceClaudeCodeLocal: {OriginLocal, []Tool{ToolClaudeCode}},
	SourceCodexLocal:      {OriginLocal, []Tool{ToolCodex}},
	SourceGeminiCLILocal:  {OriginLocal, []Tool{ToolGeminiCLI}},
	SourceOpenCodeLocal:   {OriginLocal, []Tool{ToolOpenCode}},
	// Anthropic usage counts as claude_code only when it comes from Claude Code analytics.
	SourceAnthropicAPI: {OriginAPI, []Tool{ToolClaudeCode, ToolOther}},
	SourceCursorAPI:    {OriginAPI, []Tool{ToolCursor}},
	SourceCopilotAPI:   {OriginAPI, []Tool{ToolCopilot}},
	// OpenAI usage counts as codex only when it can be tied to Codex.
	SourceOpenAIAPI: {OriginAPI, []Tool{ToolCodex, ToolOther}},
}

var (
	granularities = map[Granularity]bool{GranularityTurn: true, GranularitySession: true, GranularityDaily: true}
	providers     = map[Provider]bool{ProviderAnthropic: true, ProviderOpenAI: true, ProviderGoogle: true, ProviderCursor: true, ProviderGitHub: true, ProviderOther: true}
	billingModes  = map[BillingMode]bool{BillingAPI: true, BillingSubscription: true, BillingSeat: true, BillingUnknown: true}
)

// Tokens holds token counts. Missing counts are zero. Cache reads and writes stay
// separate because they are priced differently.
type Tokens struct {
	Input      int64 `json:"input"`
	Output     int64 `json:"output"`
	CacheRead  int64 `json:"cache_read"`
	CacheWrite int64 `json:"cache_write"`
	Reasoning  int64 `json:"reasoning"`
}

// UsageEvent is the normalized event the collector sends to the server.
//
// It has no device and no actor: the server takes both from the device token.
// It has no free-text field, and decoding rejects unknown fields, so content
// such as prompts or file paths cannot be sent.
type UsageEvent struct {
	EventID          string      `json:"event_id"`
	Source           Source      `json:"source"`
	Tool             Tool        `json:"tool"`
	Origin           Origin      `json:"origin"`
	Granularity      Granularity `json:"granularity"`
	OccurredAt       time.Time   `json:"occurred_at"`
	SessionID        string      `json:"session_id,omitempty"`
	Repo             string      `json:"repo,omitempty"`
	Model            string      `json:"model,omitempty"`
	Provider         Provider    `json:"provider"`
	Tokens           Tokens      `json:"tokens"`
	Requests         int64       `json:"requests"`
	CostReportedUSD  *Money      `json:"cost_reported_usd,omitempty"`
	BillingMode      BillingMode `json:"billing_mode"`
	CollectorVersion string      `json:"collector_version,omitempty"`
	SourceVersion    string      `json:"source_version,omitempty"`
}

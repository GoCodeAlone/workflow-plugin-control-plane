package adapters_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/GoCodeAlone/workflow-plugin-control-plane/adapters"
)

func TestProviderHandoffOptionsValidatePublicShape(t *testing.T) {
	opts := adapters.ProviderHandoffOptions{
		ProviderPluginVersion: "v0.4.0",
		CapabilityID:          "payment-intent-create",
		CapabilityVersion:     "v1",
		InputSchemaDigest:     "sha256:1111111111111111111111111111111111111111111111111111111111111111",
		ActionNonce:           "nonce-001",
		IdempotencyKey:        "idem-001",
	}
	if err := opts.Validate(); err != nil {
		t.Fatalf("provider handoff options invalid: %v", err)
	}
	opts.CapabilityID = "https://provider.example.invalid/capability"
	if err := opts.Validate(); err == nil {
		t.Fatal("expected raw capability URL to be rejected")
	}
}

func TestAdapterOptionsUseStableJSONNames(t *testing.T) {
	payload, err := json.Marshal(adapters.ProviderHandoffOptions{
		ProviderPluginVersion: "v0.4.0",
		CapabilityID:          "payment-intent-create",
		CapabilityVersion:     "v1",
		InputSchemaDigest:     "sha256:1111111111111111111111111111111111111111111111111111111111111111",
		ActionNonce:           "nonce-001",
		IdempotencyKey:        "idem-001",
	})
	if err != nil {
		t.Fatalf("marshal provider handoff options: %v", err)
	}
	wire := string(payload)
	for _, want := range []string{`"provider_plugin_version"`, `"capability_id"`, `"input_schema_digest"`, `"action_nonce"`, `"idempotency_key"`} {
		if !strings.Contains(wire, want) {
			t.Fatalf("wire JSON %s missing %s", wire, want)
		}
	}
	for _, forbidden := range []string{"ProviderPluginVersion", "CapabilityID", "InputSchemaDigest"} {
		if strings.Contains(wire, forbidden) {
			t.Fatalf("wire JSON leaked Go field name %q: %s", forbidden, wire)
		}
	}
}

func TestAdminContributionOptionsValidateAndNormalizeMethod(t *testing.T) {
	opts := adapters.AdminContributionOptions{
		SurfaceID:  "product-capture.admin.v1",
		RenderMode: "host.card.v1",
		ActionID:   "submit-validation-task",
		Method:     " post ",
	}
	normalized, err := opts.Normalize()
	if err != nil {
		t.Fatalf("admin contribution options invalid: %v", err)
	}
	if normalized.Method != "POST" {
		t.Fatalf("method = %q want POST", normalized.Method)
	}
	for _, method := range []string{"TRACE", "OPTIONS", "POST-JSON"} {
		opts.Method = method
		if _, err := opts.Normalize(); err == nil {
			t.Fatalf("expected unsupported method %q to be rejected", method)
		}
	}
}

func TestEnvelopeOptionsValidateRetentionShape(t *testing.T) {
	occurredAt := time.Unix(1_700_000_000, 0).UTC()
	opts := adapters.EnvelopeOptions{
		ObservedAt:         occurredAt.Add(time.Second),
		RetentionExpires:   occurredAt.Add(24 * time.Hour),
		CorrelationRekeyID: "rekey-001",
		LegalErasure:       true,
	}
	if err := opts.Validate(occurredAt); err != nil {
		t.Fatalf("envelope options invalid: %v", err)
	}
	opts.CorrelationRekeyID = ""
	if err := opts.Validate(occurredAt); err == nil || !strings.Contains(err.Error(), "correlation_rekey_id") {
		t.Fatalf("expected legal erasure rekey error, got %v", err)
	}
}

package envelopes

import (
	"errors"
	"fmt"

	"github.com/GoCodeAlone/workflow-plugin-control-plane/envelopes/pb"
	"github.com/GoCodeAlone/workflow-plugin-control-plane/internal/validate"
)

const Version = validate.ProtocolVersion

func ValidateEnvelope(e *pb.ControlPlaneEnvelope) error {
	if e == nil {
		return errors.New("envelope is required")
	}
	return validate.Join(
		validate.Protocol(e.ProtocolVersion),
		validate.ID("envelope_id", e.EnvelopeId),
		validate.ID("kind", e.Kind),
		validate.NoRawIdentifier("tenant_handle", e.TenantHandle),
		validate.NoRawIdentifier("actor_handle", e.ActorHandle),
		validate.NoRawIdentifier("resource_handle", e.ResourceHandle),
		validate.ID("correlation_id", e.CorrelationId),
		validate.UnixNano("occurred_at_unix_nano", e.OccurredAtUnixNano),
		validate.UnixNano("observed_at_unix_nano", e.ObservedAtUnixNano),
		validate.NotBefore("observed_at_unix_nano", e.ObservedAtUnixNano, e.OccurredAtUnixNano),
		ValidateRetentionMetadata(e.Retention, e.OccurredAtUnixNano),
	)
}

func ValidateRetentionMetadata(m *pb.RetentionMetadata, occurredAt int64) error {
	if m == nil {
		return errors.New("retention metadata is required")
	}
	var rekeyErr error
	if m.LegalErasure && m.CorrelationRekeyId == "" {
		rekeyErr = fmt.Errorf("correlation_rekey_id is required for legal erasure")
	}
	return validate.Join(
		validate.ID("redaction_state", m.RedactionState),
		validate.UnixNano("retention_expires_at_unix_nano", m.RetentionExpiresAtUnixNano),
		validate.NotBefore("retention_expires_at_unix_nano", m.RetentionExpiresAtUnixNano, occurredAt),
		validate.OptionalID("correlation_rekey_id", m.CorrelationRekeyId),
		rekeyErr,
	)
}

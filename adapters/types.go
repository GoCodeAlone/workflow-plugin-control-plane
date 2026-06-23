package adapters

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/GoCodeAlone/workflow-plugin-control-plane/internal/validate"
)

type ProviderHandoffOptions struct {
	ProviderPluginVersion string `json:"provider_plugin_version"`
	CapabilityID          string `json:"capability_id"`
	CapabilityVersion     string `json:"capability_version"`
	InputSchemaDigest     string `json:"input_schema_digest"`
	ActionNonce           string `json:"action_nonce"`
	IdempotencyKey        string `json:"idempotency_key"`
}

func (o ProviderHandoffOptions) Validate() error {
	return validate.Join(
		requiredVersion("provider_plugin_version", o.ProviderPluginVersion),
		validate.ID("capability_id", o.CapabilityID),
		requiredVersion("capability_version", o.CapabilityVersion),
		validate.Digest("input_schema_digest", o.InputSchemaDigest),
		validate.ID("action_nonce", o.ActionNonce),
		validate.ID("idempotency_key", o.IdempotencyKey),
	)
}

type AdminContributionOptions struct {
	SurfaceID  string `json:"surface_id"`
	RenderMode string `json:"render_mode"`
	ActionID   string `json:"action_id"`
	Method     string `json:"method"`
}

func (o AdminContributionOptions) Normalize() (AdminContributionOptions, error) {
	o.SurfaceID = strings.TrimSpace(o.SurfaceID)
	o.RenderMode = strings.TrimSpace(o.RenderMode)
	o.ActionID = strings.TrimSpace(o.ActionID)
	o.Method = strings.TrimSpace(strings.ToUpper(o.Method))
	if err := validate.Join(
		validate.ID("surface_id", o.SurfaceID),
		validate.ID("render_mode", o.RenderMode),
		validate.ID("action_id", o.ActionID),
		validateAdminContributionMethod(o.Method),
	); err != nil {
		return AdminContributionOptions{}, err
	}
	return o, nil
}

type EnvelopeOptions struct {
	ObservedAt         time.Time `json:"observed_at"`
	RetentionExpires   time.Time `json:"retention_expires"`
	CorrelationRekeyID string    `json:"correlation_rekey_id,omitempty"`
	LegalErasure       bool      `json:"legal_erasure,omitempty"`
	Tombstoned         bool      `json:"tombstoned,omitempty"`
}

func (o EnvelopeOptions) Validate(occurredAt time.Time) error {
	var errs []error
	if occurredAt.IsZero() {
		errs = append(errs, errors.New("occurred_at is required"))
	}
	if o.ObservedAt.IsZero() {
		errs = append(errs, errors.New("observed_at is required"))
	}
	if o.RetentionExpires.IsZero() {
		errs = append(errs, errors.New("retention_expires is required"))
	}
	if !occurredAt.IsZero() && !o.ObservedAt.IsZero() && o.ObservedAt.Before(occurredAt) {
		errs = append(errs, errors.New("observed_at must not be before occurred_at"))
	}
	if !occurredAt.IsZero() && !o.RetentionExpires.IsZero() && o.RetentionExpires.Before(occurredAt) {
		errs = append(errs, errors.New("retention_expires must not be before occurred_at"))
	}
	if o.LegalErasure && strings.TrimSpace(o.CorrelationRekeyID) == "" {
		errs = append(errs, errors.New("correlation_rekey_id is required for legal erasure"))
	}
	errs = append(errs, validate.OptionalID("correlation_rekey_id", o.CorrelationRekeyID))
	return errors.Join(errs...)
}

func validateAdminContributionMethod(method string) error {
	switch method {
	case "GET", "POST", "PUT", "PATCH", "DELETE":
		return nil
	default:
		return fmt.Errorf("method must be one of GET, POST, PUT, PATCH, DELETE")
	}
}

func requiredVersion(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", field)
	}
	return validate.ID(field, value)
}

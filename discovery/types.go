package discovery

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/GoCodeAlone/workflow-plugin-control-plane/internal/validate"
)

const Version = validate.ProtocolVersion

type SignatureEnvelope struct {
	Algorithm string `json:"algorithm,omitempty"`
	KeyID     string `json:"key_id,omitempty"`
	Value     string `json:"value,omitempty"`
	Verified  bool   `json:"-"`
}

type Document struct {
	ProtocolVersion string            `json:"protocol_version"`
	ServerURL       string            `json:"server_url"`
	IssuedAt        time.Time         `json:"issued_at"`
	ExpiresAt       time.Time         `json:"expires_at"`
	Signature       SignatureEnvelope `json:"signature"`
}

// SigningPayload preserves the empty signature object in JSON for compatibility
// with existing workflow-compute discovery signatures.
func (d Document) SigningPayload() Document {
	d.Signature = SignatureEnvelope{}
	return d
}

func (d Document) Validate(now time.Time) error {
	var errs []error
	errs = append(errs, validate.Protocol(d.ProtocolVersion))
	if strings.TrimSpace(d.ServerURL) == "" {
		errs = append(errs, errors.New("server_url is required"))
	} else if err := validateServerURL(d.ServerURL); err != nil {
		errs = append(errs, fmt.Errorf("server_url: %w", err))
	}
	if d.IssuedAt.IsZero() {
		errs = append(errs, errors.New("issued_at is required"))
	}
	if d.ExpiresAt.IsZero() {
		errs = append(errs, errors.New("expires_at is required"))
	}
	if !d.IssuedAt.IsZero() && !d.ExpiresAt.IsZero() && !d.IssuedAt.Before(d.ExpiresAt) {
		errs = append(errs, errors.New("issued_at must be before expires_at"))
	}
	if !d.ExpiresAt.IsZero() && !now.Before(d.ExpiresAt) {
		errs = append(errs, errors.New("discovery document expired"))
	}
	if err := validateSignatureEnvelope(d.Signature, "signature"); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func validateServerURL(raw string) error {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return err
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("must use http or https")
	}
	if parsed.Host == "" {
		return errors.New("host is required")
	}
	if parsed.ForceQuery || parsed.RawQuery != "" || strings.Contains(strings.TrimSpace(raw), "#") {
		return errors.New("must not contain query or fragment")
	}
	return nil
}

func validateSignatureEnvelope(sig SignatureEnvelope, name string) error {
	var errs []error
	if strings.TrimSpace(sig.Algorithm) == "" {
		errs = append(errs, fmt.Errorf("%s.algorithm is required", name))
	}
	if strings.TrimSpace(sig.KeyID) == "" {
		errs = append(errs, fmt.Errorf("%s.key_id is required", name))
	}
	if strings.TrimSpace(sig.Value) == "" {
		errs = append(errs, fmt.Errorf("%s.value is required", name))
	}
	return errors.Join(errs...)
}

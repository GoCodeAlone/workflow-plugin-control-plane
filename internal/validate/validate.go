package validate

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

const ProtocolVersion = "control-plane.v1alpha1"

var (
	idPattern     = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,127}$`)
	refPattern    = regexp.MustCompile(`^[a-z][a-z0-9+.-]*://[A-Za-z0-9._~:/?#\[\]@!$&'()*+,;=%-]+$`)
	digestPattern = regexp.MustCompile(`^sha256:[a-f0-9]{64}$`)
)

var hostScopedSchemes = map[string]struct{}{
	"actionclass": {},
	"auditclass":  {},
	"authclass":   {},
	"provenance":  {},
	"render":      {},
	"resource":    {},
	"risk":        {},
	"route":       {},
	"severity":    {},
	"template":    {},
	"trustroot":   {},
}

func Protocol(version string) error {
	if version != ProtocolVersion {
		return fmt.Errorf("protocol_version must be %s", ProtocolVersion)
	}
	return nil
}

func ID(field, value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("%s is required", field)
	}
	if !idPattern.MatchString(value) {
		return fmt.Errorf("%s must be lowercase DNS-like id", field)
	}
	return nil
}

func OptionalID(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return ID(field, value)
}

func Ref(field, value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("%s is required", field)
	}
	if unsafePayload(value) {
		return fmt.Errorf("%s must be a host-scoped typed ref, not executable or network payload", field)
	}
	if !refPattern.MatchString(value) {
		return fmt.Errorf("%s must be a typed ref", field)
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" {
		return fmt.Errorf("%s must be a typed ref", field)
	}
	if _, ok := hostScopedSchemes[parsed.Scheme]; !ok {
		return fmt.Errorf("%s scheme %q is not a host-scoped control-plane ref", field, parsed.Scheme)
	}
	return nil
}

func OptionalRef(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return Ref(field, value)
}

func Digest(field, value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("%s is required", field)
	}
	if !digestPattern.MatchString(value) {
		return fmt.Errorf("%s must be sha256:<64 lowercase hex>", field)
	}
	return nil
}

func Positive(field string, value uint64) error {
	if value == 0 {
		return fmt.Errorf("%s must be positive", field)
	}
	return nil
}

func UnixNano(field string, value int64) error {
	if value <= 0 {
		return fmt.Errorf("%s must be positive unix nano timestamp", field)
	}
	return nil
}

func NotBefore(field string, later, earlier int64) error {
	if later > 0 && earlier > 0 && later < earlier {
		return fmt.Errorf("%s must not be before prior timestamp", field)
	}
	return nil
}

func NoRawIdentifier(field, value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("%s is required", field)
	}
	if strings.Contains(value, "@") || strings.Contains(value, " ") || strings.Contains(value, "\x00") {
		return fmt.Errorf("%s must be opaque or pseudonymous", field)
	}
	if strings.HasPrefix(value, "user:") || strings.HasPrefix(value, "tenant:") || strings.HasPrefix(value, "customer:") {
		return fmt.Errorf("%s must not expose raw tenant/user/customer identifiers", field)
	}
	return nil
}

func Join(errs ...error) error {
	var kept []error
	for _, err := range errs {
		if err != nil {
			kept = append(kept, err)
		}
	}
	return errors.Join(kept...)
}

func unsafePayload(value string) bool {
	lower := strings.ToLower(strings.TrimSpace(value))
	if lower == "" {
		return false
	}
	if strings.HasPrefix(lower, "http://") ||
		strings.HasPrefix(lower, "https://") ||
		strings.HasPrefix(lower, "webhook://") ||
		strings.HasPrefix(lower, "callback://") ||
		strings.HasPrefix(lower, "command://") ||
		strings.HasPrefix(lower, "shell://") ||
		strings.HasPrefix(lower, "secret://") ||
		strings.HasPrefix(lower, "private-key://") {
		return true
	}
	return strings.ContainsAny(value, "\x00\r\n")
}

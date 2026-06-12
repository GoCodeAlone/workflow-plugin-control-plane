package descriptors

import (
	"errors"
	"fmt"

	"github.com/GoCodeAlone/workflow-plugin-control-plane/descriptors/pb"
	"github.com/GoCodeAlone/workflow-plugin-control-plane/internal/validate"
)

const Version = validate.ProtocolVersion

func ValidateRouteActionDescriptor(d *pb.RouteActionDescriptor) error {
	if d == nil {
		return errors.New("descriptor is required")
	}
	return validate.Join(
		validate.Protocol(d.ProtocolVersion),
		validate.ID("descriptor_id", d.DescriptorId),
		validate.ID("operation_id", d.OperationId),
		validate.Ref("route_handle", d.RouteHandle),
		validate.Ref("auth_class_ref", d.AuthClassRef),
		validate.Ref("audit_class_ref", d.AuditClassRef),
		validate.Digest("input_schema_digest", d.InputSchemaDigest),
		optionalDigest("output_schema_digest", d.OutputSchemaDigest),
		ValidateProviderHandoffRef(d.ProviderHandoff),
		ValidateAdminContribution(d.Admin),
		validate.OptionalRef("presentation_ref", d.PresentationRef),
	)
}

func ValidateProviderHandoffRef(ref *pb.ProviderHandoffRef) error {
	if ref == nil {
		return nil
	}
	return validate.Join(
		validate.ID("provider_plugin_id", ref.ProviderPluginId),
		requiredVersion("provider_plugin_version", ref.ProviderPluginVersion),
		validate.ID("capability_id", ref.CapabilityId),
		requiredVersion("capability_version", ref.CapabilityVersion),
		validate.Digest("provider_handoff.input_schema_digest", ref.InputSchemaDigest),
		validate.ID("action_nonce", ref.ActionNonce),
		validate.ID("idempotency_key", ref.IdempotencyKey),
	)
}

func ValidateAdminContribution(admin *pb.AdminContribution) error {
	if admin == nil {
		return nil
	}
	return validate.Join(
		validate.Ref("admin.resource_ref", admin.ResourceRef),
		validate.Ref("admin.action_class_ref", admin.ActionClassRef),
		validate.Ref("admin.risk_ref", admin.RiskRef),
		validate.Ref("admin.severity_ref", admin.SeverityRef),
		validate.Ref("admin.permission_explanation_template_ref", admin.PermissionExplanationTemplateRef),
		validate.Ref("admin.render_ref", admin.RenderRef),
	)
}

func optionalDigest(field, value string) error {
	if value == "" {
		return nil
	}
	return validate.Digest(field, value)
}

func requiredVersion(field, value string) error {
	if value == "" {
		return fmt.Errorf("%s is required", field)
	}
	return validate.OptionalID(field, value)
}

package registry

import (
	"errors"

	"github.com/GoCodeAlone/workflow-plugin-control-plane/internal/validate"
	"github.com/GoCodeAlone/workflow-plugin-control-plane/registry/pb"
)

const Version = validate.ProtocolVersion

func ValidateDescriptorRegistration(r *pb.DescriptorRegistration) error {
	if r == nil {
		return errors.New("descriptor registration is required")
	}
	return validate.Join(
		validate.Protocol(r.ProtocolVersion),
		validate.Digest("descriptor_digest", r.DescriptorDigest),
		validate.Digest("schema_digest", r.SchemaDigest),
		validate.ID("plugin_id", r.PluginId),
		validate.ID("plugin_version", r.PluginVersion),
		validate.Digest("plugin_package_digest", r.PluginPackageDigest),
		validate.Ref("signing_root_ref", r.SigningRootRef),
		validate.Ref("provenance_source", r.ProvenanceSource),
		validate.Positive("allowlist_epoch", r.AllowlistEpoch),
		validate.ID("downgrade_floor_version", r.DowngradeFloorVersion),
		validate.UnixNano("revocation_fresh_until_unix_nano", r.RevocationFreshUntilUnixNano),
		validate.UnixNano("fetched_at_unix_nano", r.FetchedAtUnixNano),
		validate.UnixNano("validated_at_unix_nano", r.ValidatedAtUnixNano),
		validate.NotBefore("validated_at_unix_nano", r.ValidatedAtUnixNano, r.FetchedAtUnixNano),
		validate.NotBefore("revocation_fresh_until_unix_nano", r.RevocationFreshUntilUnixNano, r.ValidatedAtUnixNano),
		validate.ID("parser_compatibility_version", r.ParserCompatibilityVersion),
		validate.Positive("trust_root_generation", r.TrustRootGeneration),
	)
}

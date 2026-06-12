package registry_test

import (
	"strings"
	"testing"

	"github.com/GoCodeAlone/workflow-plugin-control-plane/registry"
	"github.com/GoCodeAlone/workflow-plugin-control-plane/registry/pb"
)

func TestDescriptorRegistrationValidatesProvenance(t *testing.T) {
	r := validRegistration()
	if err := registry.ValidateDescriptorRegistration(r); err != nil {
		t.Fatalf("registration invalid: %v", err)
	}
}

func TestDescriptorRegistrationRejectsSkewAndStaleProvenance(t *testing.T) {
	cases := map[string]func(*pb.DescriptorRegistration){
		"digest mismatch": func(r *pb.DescriptorRegistration) {
			r.DescriptorDigest = "sha256:bad"
		},
		"network signing root": func(r *pb.DescriptorRegistration) {
			r.SigningRootRef = "https://keys.example.test/root"
		},
		"zero epoch": func(r *pb.DescriptorRegistration) {
			r.AllowlistEpoch = 0
		},
		"stale revocation freshness": func(r *pb.DescriptorRegistration) {
			r.RevocationFreshUntilUnixNano = r.ValidatedAtUnixNano - 1
		},
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			r := validRegistration()
			mutate(r)
			if err := registry.ValidateDescriptorRegistration(r); err == nil {
				t.Fatal("expected registration to fail")
			}
		})
	}
}

func validRegistration() *pb.DescriptorRegistration {
	return &pb.DescriptorRegistration{
		ProtocolVersion:              registry.Version,
		DescriptorDigest:             digest("c"),
		SchemaDigest:                 digest("d"),
		PluginId:                     "workflow-plugin-control-plane",
		PluginVersion:                "v0.1.0",
		PluginPackageDigest:          digest("e"),
		SigningRootRef:               "trustroot://workflow/control-plane",
		ProvenanceSource:             "provenance://github/release",
		AllowlistEpoch:               1,
		DowngradeFloorVersion:        "v0.1.0",
		RevocationFreshUntilUnixNano: 300,
		FetchedAtUnixNano:            100,
		ValidatedAtUnixNano:          200,
		ParserCompatibilityVersion:   "v1",
		TrustRootGeneration:          1,
	}
}

func digest(ch string) string {
	return "sha256:" + strings.Repeat(ch, 64)
}

package descriptors_test

import (
	"strings"
	"testing"

	"github.com/GoCodeAlone/workflow-plugin-control-plane/descriptors"
	"github.com/GoCodeAlone/workflow-plugin-control-plane/descriptors/pb"
)

func TestRouteActionDescriptorValidatesDataOnlyShape(t *testing.T) {
	d := validDescriptor()
	if err := descriptors.ValidateRouteActionDescriptor(d); err != nil {
		t.Fatalf("descriptor invalid: %v", err)
	}
}

func TestRouteActionDescriptorRejectsAuthorityTransferPayloads(t *testing.T) {
	cases := map[string]func(*pb.RouteActionDescriptor){
		"raw route url": func(d *pb.RouteActionDescriptor) {
			d.RouteHandle = "https://example.test/admin"
		},
		"shell presentation": func(d *pb.RouteActionDescriptor) {
			d.PresentationRef = "shell://rm-rf"
		},
		"credential handoff": func(d *pb.RouteActionDescriptor) {
			d.Admin.RenderRef = "secret://control-plane/admin"
		},
		"schema mismatch": func(d *pb.RouteActionDescriptor) {
			d.InputSchemaDigest = "sha256:not-hex"
		},
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			d := validDescriptor()
			mutate(d)
			if err := descriptors.ValidateRouteActionDescriptor(d); err == nil {
				t.Fatal("expected descriptor to fail")
			}
		})
	}
}

func TestRouteActionDescriptorReportsMultipleErrors(t *testing.T) {
	d := validDescriptor()
	d.ProtocolVersion = "wrong"
	d.ProviderHandoff.ProviderPluginId = "Bad Plugin"
	err := descriptors.ValidateRouteActionDescriptor(d)
	if err == nil {
		t.Fatal("expected errors")
	}
	for _, want := range []string{"protocol_version", "provider_plugin_id"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q missing %q", err, want)
		}
	}
}

func validDescriptor() *pb.RouteActionDescriptor {
	return &pb.RouteActionDescriptor{
		ProtocolVersion:   descriptors.Version,
		DescriptorId:      "descriptor.admin.deploy",
		OperationId:       "deploy_app",
		RouteHandle:       "route://admin/deploy-app",
		AuthClassRef:      "authclass://admin/deploy",
		AuditClassRef:     "auditclass://deployment/apply",
		InputSchemaDigest: digest("a"),
		PresentationRef:   "",
		ProviderHandoff: &pb.ProviderHandoffRef{
			ProviderPluginId:      "workflow-plugin-digitalocean",
			ProviderPluginVersion: "v2.0.15",
			CapabilityId:          "app-platform",
			CapabilityVersion:     "v1",
			InputSchemaDigest:     digest("b"),
			ActionNonce:           "nonce-001",
			IdempotencyKey:        "idem-001",
		},
		Admin: &pb.AdminContribution{
			ResourceRef:                      "resource://app-platform/service",
			ActionClassRef:                   "actionclass://deployment/apply",
			RiskRef:                          "risk://destructive",
			SeverityRef:                      "severity://high",
			PermissionExplanationTemplateRef: "template://permission/deployment",
			RenderRef:                        "render://admin/action-form",
		},
	}
}

func digest(ch string) string {
	return "sha256:" + strings.Repeat(ch, 64)
}

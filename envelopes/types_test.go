package envelopes_test

import (
	"strings"
	"testing"

	"github.com/GoCodeAlone/workflow-plugin-control-plane/envelopes"
	"github.com/GoCodeAlone/workflow-plugin-control-plane/envelopes/pb"
)

func TestEnvelopeValidatesOpaqueHandlesAndRetention(t *testing.T) {
	e := validEnvelope()
	if err := envelopes.ValidateEnvelope(e); err != nil {
		t.Fatalf("envelope invalid: %v", err)
	}
}

func TestEnvelopeRejectsRawPIIAndInvalidRetention(t *testing.T) {
	cases := map[string]func(*pb.ControlPlaneEnvelope){
		"raw user email": func(e *pb.ControlPlaneEnvelope) {
			e.ActorHandle = "person@example.test"
		},
		"raw tenant id": func(e *pb.ControlPlaneEnvelope) {
			e.TenantHandle = "tenant:customer-1"
		},
		"observed before occurred": func(e *pb.ControlPlaneEnvelope) {
			e.ObservedAtUnixNano = e.OccurredAtUnixNano - 1
		},
		"legal erasure without rekey": func(e *pb.ControlPlaneEnvelope) {
			e.Retention.LegalErasure = true
			e.Retention.CorrelationRekeyId = ""
		},
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			e := validEnvelope()
			mutate(e)
			if err := envelopes.ValidateEnvelope(e); err == nil {
				t.Fatal("expected envelope to fail")
			}
		})
	}
}

func validEnvelope() *pb.ControlPlaneEnvelope {
	return &pb.ControlPlaneEnvelope{
		ProtocolVersion:    envelopes.Version,
		EnvelopeId:         "evt-001",
		Kind:               "audit",
		TenantHandle:       "opaque-tenant-001",
		ActorHandle:        "opaque-actor-001",
		ResourceHandle:     "opaque-resource-001",
		CorrelationId:      "corr-001",
		OccurredAtUnixNano: 100,
		ObservedAtUnixNano: 200,
		Retention: &pb.RetentionMetadata{
			RedactionState:             "active",
			RetentionExpiresAtUnixNano: 1000,
			CorrelationRekeyId:         "rekey-001",
		},
	}
}

func TestEnvelopeReportsProtocolError(t *testing.T) {
	e := validEnvelope()
	e.ProtocolVersion = "wrong"
	err := envelopes.ValidateEnvelope(e)
	if err == nil || !strings.Contains(err.Error(), "protocol_version") {
		t.Fatalf("expected protocol error, got %v", err)
	}
}

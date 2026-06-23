package discovery_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/GoCodeAlone/workflow-plugin-control-plane/discovery"
)

func TestDocumentValidatesPublicDiscoveryShape(t *testing.T) {
	doc := validDocument()
	if err := doc.Validate(time.Now().UTC()); err != nil {
		t.Fatalf("discovery document invalid: %v", err)
	}
	if doc.SigningPayload().Signature != (discovery.SignatureEnvelope{}) {
		t.Fatalf("signing payload must clear signature: %+v", doc.SigningPayload())
	}
	payload, err := json.Marshal(doc.SigningPayload())
	if err != nil {
		t.Fatalf("marshal signing payload: %v", err)
	}
	payloadJSON := string(payload)
	if !strings.Contains(payloadJSON, `"signature":{}`) {
		t.Fatalf("signing payload must preserve empty signature object for compatibility: %s", payloadJSON)
	}
	if strings.Contains(payloadJSON, "signature-base64") {
		t.Fatalf("signing payload leaked signature value: %s", payloadJSON)
	}
}

func TestDocumentRejectsInvalidData(t *testing.T) {
	cases := map[string]func(*discovery.Document){
		"wrong protocol": func(d *discovery.Document) {
			d.ProtocolVersion = "wrong"
		},
		"raw query": func(d *discovery.Document) {
			d.ServerURL = "https://staging.example.invalid?token=secret"
		},
		"empty query delimiter": func(d *discovery.Document) {
			d.ServerURL = "https://staging.example.invalid?"
		},
		"empty fragment delimiter": func(d *discovery.Document) {
			d.ServerURL = "https://staging.example.invalid#"
		},
		"unsupported scheme": func(d *discovery.Document) {
			d.ServerURL = "ssh://staging.example.invalid"
		},
		"expired": func(d *discovery.Document) {
			d.ExpiresAt = time.Now().UTC().Add(-time.Minute)
		},
		"missing signature": func(d *discovery.Document) {
			d.Signature.Value = ""
		},
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			doc := validDocument()
			mutate(&doc)
			if err := doc.Validate(time.Now().UTC()); err == nil {
				t.Fatal("expected discovery document to fail")
			}
		})
	}
}

func TestSignatureEnvelopeVerifiedIsLocalOnly(t *testing.T) {
	raw := []byte(`{"algorithm":"ed25519","key_id":"server","value":"signature-base64","verified":true}`)
	var sig discovery.SignatureEnvelope
	if err := json.Unmarshal(raw, &sig); err != nil {
		t.Fatalf("unmarshal signature envelope: %v", err)
	}
	if sig.Verified {
		t.Fatal("verified must not be accepted from wire JSON")
	}
	encoded, err := json.Marshal(discovery.SignatureEnvelope{
		Algorithm: "ed25519",
		KeyID:     "server",
		Value:     "signature-base64",
		Verified:  true,
	})
	if err != nil {
		t.Fatalf("marshal signature envelope: %v", err)
	}
	if strings.Contains(string(encoded), "verified") {
		t.Fatalf("verified must not be emitted to wire JSON: %s", encoded)
	}
}

func TestDocumentReportsMultipleErrors(t *testing.T) {
	doc := validDocument()
	doc.ProtocolVersion = "wrong"
	doc.ServerURL = "https://staging.example.invalid#fragment"
	doc.Signature.KeyID = ""
	err := doc.Validate(time.Now().UTC())
	if err == nil {
		t.Fatal("expected discovery document errors")
	}
	for _, want := range []string{"protocol_version", "server_url", "signature.key_id"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q missing %q", err, want)
		}
	}
}

func validDocument() discovery.Document {
	now := time.Now().UTC()
	return discovery.Document{
		ProtocolVersion: discovery.Version,
		ServerURL:       "https://workflow-compute-staging.example.invalid",
		IssuedAt:        now.Add(-time.Minute),
		ExpiresAt:       now.Add(time.Hour),
		Signature: discovery.SignatureEnvelope{
			Algorithm: "ed25519",
			KeyID:     "server",
			Value:     "signature-base64",
		},
	}
}

package internal_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/GoCodeAlone/workflow-plugin-control-plane/internal"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

func TestNewPluginImplementsProvider(t *testing.T) {
	var _ sdk.PluginProvider = internal.NewPlugin()
}

func TestManifestHasContractOnlyIdentity(t *testing.T) {
	manifest := internal.NewPlugin().Manifest()
	if manifest.Name != "workflow-plugin-control-plane" {
		t.Fatalf("name = %q", manifest.Name)
	}
	if manifest.Version == "" {
		t.Fatal("version is required")
	}
	if strings.Contains(strings.ToLower(manifest.Description), "scaffold") {
		t.Fatalf("description contains scaffold text: %q", manifest.Description)
	}
}

func TestPluginJSONAdvertisesNoExecutableAuthority(t *testing.T) {
	data, err := os.ReadFile("../plugin.json")
	if err != nil {
		t.Fatalf("read plugin.json: %v", err)
	}
	var manifest struct {
		Name         string `json:"name"`
		Type         string `json:"type"`
		Private      bool   `json:"private"`
		Capabilities struct {
			ConfigProvider bool     `json:"configProvider"`
			ModuleTypes    []string `json:"moduleTypes"`
			StepTypes      []string `json:"stepTypes"`
			TriggerTypes   []string `json:"triggerTypes"`
		} `json:"capabilities"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("parse plugin.json: %v", err)
	}
	if manifest.Name != "workflow-plugin-control-plane" || manifest.Type != "external" || manifest.Private {
		t.Fatalf("unexpected identity: %+v", manifest)
	}
	if manifest.Capabilities.ConfigProvider ||
		len(manifest.Capabilities.ModuleTypes) != 0 ||
		len(manifest.Capabilities.StepTypes) != 0 ||
		len(manifest.Capabilities.TriggerTypes) != 0 {
		t.Fatalf("control-plane contract plugin must not advertise executable authority: %+v", manifest.Capabilities)
	}
}

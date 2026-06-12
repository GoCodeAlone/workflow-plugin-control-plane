// Package internal implements the workflow-plugin-control-plane plugin.
package internal

import sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"

// Version is set at build time via -ldflags
// "-X github.com/GoCodeAlone/workflow-plugin-control-plane/internal.Version=X.Y.Z".
var Version = "0.0.0"

type ControlPlanePlugin struct{}

func NewPlugin() sdk.PluginProvider {
	return &ControlPlanePlugin{}
}

func (p *ControlPlanePlugin) Manifest() sdk.PluginManifest {
	return sdk.PluginManifest{
		Name:        "workflow-plugin-control-plane",
		Version:     Version,
		Author:      "GoCodeAlone",
		Description: "Product-neutral Workflow control-plane descriptor, envelope, and registry contracts.",
	}
}

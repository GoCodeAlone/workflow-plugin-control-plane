// Command workflow-plugin-control-plane exposes control-plane contract metadata
// as an external Workflow plugin dependency anchor.
package main

import (
	"github.com/GoCodeAlone/workflow-plugin-control-plane/internal"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

func main() {
	sdk.Serve(internal.NewPlugin(),
		sdk.WithBuildVersion(sdk.ResolveBuildVersion(internal.Version)),
	)
}

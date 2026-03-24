package cmd

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// Version is injected at compile time through -ldflags
	// Default to "dev" for development builds (go install, go build without ldflags)
	Version = "dev"
	// BuildTime is injected at compile time through -ldflags
	BuildTime = "unknown"
	// GoVersion contains the Go runtime version
	GoVersion = runtime.Version()
	// GitCommit is injected at compile time through -ldflags
	GitCommit = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display Goca CLI version",
	Long:  "Display the current version of Goca CLI along with build information.",
	Run: func(cmd *cobra.Command, _ []string) {
		short, _ := cmd.Flags().GetBool("short")

		if short {
			ui.Println(Version)
		} else {
			ui.Header(fmt.Sprintf("Goca v%s", Version))
			ui.KeyValue("Build", BuildTime)
			ui.KeyValue("Go Version", GoVersion)
			ui.KeyValue("Git Commit", GitCommit)
		}
	},
}

func init() {
	// When Version is still the default "dev", try to read the module version
	// embedded by the Go toolchain (set when using `go install module@version`).
	if Version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
			Version = strings.TrimPrefix(info.Main.Version, "v")
		}
	}
	versionCmd.Flags().BoolP("short", "s", false, "Display only the version number")
}

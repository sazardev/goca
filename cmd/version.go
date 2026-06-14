package cmd

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// Default to "dev" for development builds (go install, go build without ldflags).
	Version = "dev"
	// BuildTime is injected at compile time through -ldflags.
	BuildTime = "unknown"
	// GoVersion contains the Go runtime version.
	GoVersion = runtime.Version()
	// GitCommit is injected at compile time through -ldflags.
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
			// version is an explicit query command: its output must not be
			// blanked by --quiet, so use the ungated Println path rather than
			// Header/KeyValue (which are gated by verbosity < 1).
			ui.Println(fmt.Sprintf("Goca v%s", Version))
			ui.Println(fmt.Sprintf("Build: %s", BuildTime))
			ui.Println(fmt.Sprintf("Go Version: %s", GoVersion))
			ui.Println(fmt.Sprintf("Git Commit: %s", GitCommit))
		}
	},
}

func init() {
	// When Version is still the default "dev", try to read the module version
	// embedded by the Go toolchain (set when using `go install module@version`).
	// Also fall back to the VCS build settings for BuildTime/GitCommit when the
	// ldflags vars are still their default "unknown" (e.g. `go install`/`go build`).
	if info, ok := debug.ReadBuildInfo(); ok {
		if Version == "dev" && info.Main.Version != "" && info.Main.Version != "(devel)" {
			Version = strings.TrimPrefix(info.Main.Version, "v")
		}

		var vcsRevision, vcsTime string
		var vcsModified bool
		for _, s := range info.Settings {
			switch s.Key {
			case "vcs.revision":
				vcsRevision = s.Value
			case "vcs.time":
				vcsTime = s.Value
			case "vcs.modified":
				vcsModified = s.Value == "true"
			}
		}

		if GitCommit == "unknown" && vcsRevision != "" {
			commit := vcsRevision
			if len(commit) > 12 {
				commit = commit[:12]
			}
			if vcsModified {
				commit += "-dirty"
			}
			GitCommit = commit
		}
		if BuildTime == "unknown" && vcsTime != "" {
			BuildTime = vcsTime
		}
	}
	versionCmd.Flags().BoolP("short", "s", false, "Display only the version number")
}

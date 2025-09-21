package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// Version is injected at compile time through -ldflags
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
			fmt.Println(Version)
		} else {
			fmt.Printf("Goca v%s\n", Version)
			fmt.Printf("Build: %s\n", BuildTime)
			fmt.Printf("Go Version: %s\n", GoVersion)
			fmt.Printf("Git Commit: %s\n", GitCommit)
		}
	},
}

func init() {
	versionCmd.Flags().BoolP("short", "s", false, "Display only the version number")
}

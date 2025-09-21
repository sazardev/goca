package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// Version se inyecta en tiempo de compilación a través de -ldflags
	Version   = "dev"
	BuildTime = "unknown"
	GoVersion = runtime.Version()
	GitCommit = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Muestra la versión de Goca CLI",
	Long:  "Muestra la versión actual de Goca CLI junto con información de compilación.",
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
	versionCmd.Flags().BoolP("short", "s", false, "Muestra solo el número de versión")
}

package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	Version   = "1.0.1"
	BuildTime = "2025-01-19T10:30:00Z"
	GoVersion = runtime.Version()
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Muestra la versión de Goca CLI",
	Long:  "Muestra la versión actual de Goca CLI junto con información de compilación.",
	Run: func(cmd *cobra.Command, args []string) {
		short, _ := cmd.Flags().GetBool("short")

		if short {
			fmt.Println(Version)
		} else {
			fmt.Printf("Goca v%s\n", Version)
			fmt.Printf("Build: %s\n", BuildTime)
			fmt.Printf("Go Version: %s\n", GoVersion)
		}
	},
}

func init() {
	versionCmd.Flags().BoolP("short", "s", false, "Muestra solo el número de versión")
}

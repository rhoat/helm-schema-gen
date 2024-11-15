package version

import (
	"os"

	"github.com/spf13/cobra"
)

// Version identifier populated via the CI/CD process.
//
//nolint:gochecknoglobals // Reason for using global variable (ldflags)
var Version = "HEAD"

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version of the plugin",
		Run: func(*cobra.Command, []string) {
			os.Stdout.Write([]byte(Version))
		},
	}
}

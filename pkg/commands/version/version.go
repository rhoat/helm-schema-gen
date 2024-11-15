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
		RunE: func(*cobra.Command, []string) error {
			_, err := os.Stdout.Write([]byte(Version))
			return err
		},
	}
}

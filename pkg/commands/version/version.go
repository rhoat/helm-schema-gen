package version

import (
	"os"

	"github.com/rhoat/helm-schema-gen/pkg/ctxlogger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// Version identifier populated via the CI/CD process.
//
//nolint:gochecknoglobals // Reason for using global variable (ldflags)
var Version = "HEAD"

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version of the plugin",
		RunE: func(cmd *cobra.Command, _ []string) error {
			logger := ctxlogger.GetLogger(cmd.Context())
			logger.Debug("Called Command to output version")
			_, err := os.Stdout.Write([]byte(Version))
			if err != nil {
				logger.Error("Error writing to standard out, how?", zap.Error(err))
				return err
			}
			return nil
		},
	}
}

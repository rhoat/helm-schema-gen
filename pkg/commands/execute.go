package commands

import (
	"os"

	"github.com/rhoat/helm-schema-gen/pkg/commands/generate"
	"github.com/rhoat/helm-schema-gen/pkg/commands/version"
	"github.com/rhoat/helm-schema-gen/pkg/ctxlogger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger   *zap.Logger
	logLevel string
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:           "helm schema-gen",
	SilenceUsage:  true,
	Args:          generate.Cmd().Args,
	SilenceErrors: true,
	// If no subcommand is provided, use "generate" as the default
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println(`Command "helm schema-gen" is deprecated, use "helm schema-gen generate" instead`)
		return generate.Cmd().RunE(cmd, args)
	},
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		// Parse the log level and initialize the logger
		level, err := zapcore.ParseLevel(logLevel)
		if err != nil {
			cmd.Printf("Invalid log level: %s. Defaulting to 'info'.\n", logLevel)
			level = zapcore.InfoLevel
		}
		loggerCfg := zap.NewProductionConfig()
		loggerCfg.Level = zap.NewAtomicLevelAt(level)
		// Store the logger in the command context
		logger = zap.Must(loggerCfg.Build())
		cmd.SetContext(ctxlogger.SetLogger(cmd.Context(), logger))
	},
	Short: "Helm plugin to generate JSON schema for values YAML",
	Long: `Helm plugin to generate JSON schema for values YAML.

Examples:
  $ helm schema-gen generate values.yaml    # Generate schema JSON
`,
}

// Execute initializes and executes the root command.
func Execute() {
	if err := runCLI(); err != nil {
		logger.Error(err.Error(), zap.Error(err))
		os.Exit(1)
	}
}

// runCLI runs the root command and executes the subcommands.
func runCLI() error {
	// Register subcommands (if not already registered)
	RegisterSubcommands(rootCmd)
	registerPersistentFlags(rootCmd)
	// Execute the root command
	return rootCmd.Execute()
}

// RegisterSubcommands registers all subcommands to the root command.
func RegisterSubcommands(rootCmd *cobra.Command) {
	// Add subcommands dynamically
	rootCmd.AddCommand(generate.Cmd())
	rootCmd.AddCommand(version.Cmd())
}

func registerPersistentFlags(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().StringVarP(
		&logLevel,
		"log-level",
		"l",
		"info",
		"Set the logging level (debug, info, warn, error, dpanic, panic, fatal)",
	)
}

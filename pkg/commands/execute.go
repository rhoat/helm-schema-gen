package commands

import (
	"fmt"
	"os"

	"github.com/rhoat/helm-schema-gen/pkg/commands/generate"
	"github.com/rhoat/helm-schema-gen/pkg/commands/version"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:           "helm schema-gen",
	SilenceUsage:  true,
	Args:          generate.Cmd().Args,
	SilenceErrors: true,
	// If no subcommand is provided, use "generate" as the default
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println(`Command "helm schema-gen" is deprecated, use "helm schema-gen generate" instead`)
		return generate.Cmd().RunE(cmd, args)
	},
	Short: "Helm plugin to generate JSON schema for values YAML",
	Long: `Helm plugin to generate JSON schema for values YAML.

Examples:
  $ helm schema-gen generate values.yaml    # Generate schema JSON
`,
}

// Execute initializes and executes the root command
func Execute() {
	if err := runCLI(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// runCLI runs the root command and executes the subcommands
func runCLI() error {
	// Register subcommands (if not already registered)
	RegisterSubcommands(RootCmd)

	// Execute the root command
	return RootCmd.Execute()
}

// RegisterSubcommands registers all subcommands to the root command
func RegisterSubcommands(rootCmd *cobra.Command) {
	// Add subcommands dynamically
	rootCmd.AddCommand(generate.Cmd())
	rootCmd.AddCommand(version.Cmd())
}

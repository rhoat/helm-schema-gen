package generate

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rhoat/helm-schema-gen/pkg/commands/helper"
	"github.com/rhoat/helm-schema-gen/pkg/ctxlogger"
	"github.com/rhoat/helm-schema-gen/pkg/features"
	jsonschema "github.com/rhoat/helm-schema-gen/pkg/jsonchema-generator"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate <values-yaml-file>",
		Short: "Generate JSON schema from values.yaml",
		Long:  "Generates a JSON schema from the provided Helm values.yaml file.",
		Args: func(_ *cobra.Command, args []string) error {
			return helper.CheckArgsLength(len(args), "values file")
		},
		RunE: generateJSONSchema,
	}
	cmd.PersistentFlags().BoolVar(&features.Schemagen, "schemagen", false, "Allows the user to use +schemagen comments to modify the output")
	cmd.Flags().String("destination", "", "Sets the default output location for the generated schema file")

	return cmd
}

// generateJSONSchema reads a values.yaml file and generates the corresponding JSON schema.
func generateJSONSchema(cmd *cobra.Command, args []string) error {
	logger := ctxlogger.GetLogger(cmd.Context())
	var destination string
	logger.Debug("Looking up destination")
	if cmd.Flags().Lookup("destination") != nil {
		var err error
		destination, err = cmd.Flags().GetString("destination")
		if err != nil {
			return err
		}
		logger.Debug("Destination provided", zap.String("destination", destination))
	}

	valuesFilePath := args[0]
	absPath, err := filepath.Abs(valuesFilePath)
	if err != nil {
		return err
	}
	logger.Debug("values path provided", zap.String("path", valuesFilePath), zap.String("absolutePath", absPath))

	file, err := os.Open(filepath.Clean(absPath))
	if err != nil {
		return fmt.Errorf("error when reading file '%s': %w", valuesFilePath, err)
	}
	defer file.Close()
	logger.Debug("File opened attemtping to generate schema document")
	s, err := jsonschema.Generate(cmd.Context(), file)
	if err != nil {
		return err
	}
	return output(s, destination)
}

func output(jsonDoc *jsonschema.Document, destination string) error {
	var output io.WriteCloser = os.Stdout
	defer output.Close()
	if destination != "" {
		absDest, err := filepath.Abs(destination)
		if err != nil {
			return err
		}
		file, err := os.OpenFile(filepath.Clean(absDest), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		output = file
	}
	_, err := output.Write([]byte(jsonDoc.String()))
	return err
}

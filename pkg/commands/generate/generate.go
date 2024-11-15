package generate

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rhoat/helm-schema-gen/pkg/commands/helper"
	jsonschema "github.com/rhoat/helm-schema-gen/pkg/jsonchema-generator"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate <values-yaml-file>",
		Short: "Generate JSON schema from values.yaml",
		Long:  "Generates a JSON schema from the provided Helm values.yaml file.",
		Args: func(cmd *cobra.Command, args []string) error {
			return helper.CheckArgsLength(len(args), "values file")
		},
		RunE: generateJSONSchema,
	}

	cmd.Flags().String("destination", "", "Sets the default output location for the generated schema file")

	return cmd
}

// generateJSONSchema reads a values.yaml file and generates the corresponding JSON schema.
func generateJSONSchema(cmd *cobra.Command, args []string) error {
	var destination string = ""
	if cmd.Flags().Lookup("destination") != nil {
		var err error
		destination, err = cmd.Flags().GetString("destination")
		if err != nil {
			return err
		}
	}

	valuesFilePath := args[0]
	values := make(map[string]interface{})
	absPath, err := filepath.Abs(valuesFilePath)
	if err != nil {
		return err
	}
	valuesFileData, err := os.ReadFile(filepath.Clean(absPath))
	if err != nil {
		return fmt.Errorf("error when reading file '%s': %v", valuesFilePath, err)
	}
	err = yaml.Unmarshal(valuesFileData, &values)
	if err != nil {
		return err
	}
	s := &jsonschema.Document{}
	s.ReadDeep(&values)

	return output(s, destination)
}

func output(jsonDoc *jsonschema.Document, destination string) error {
	var output io.WriteCloser = os.Stdout
	defer output.Close()
	if destination != "" {
		file, err := os.Create(destination)
		if err != nil {
			return fmt.Errorf("failed to create output file: %v", err)
		}
		output = file
	}
	_, err := output.Write([]byte(jsonDoc.String()))
	return err
}

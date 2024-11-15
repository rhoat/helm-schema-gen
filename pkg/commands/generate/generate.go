package generate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rhoat/helm-schema-gen/pkg/commands/helper"
	jsonschema "github.com/rhoat/helm-schema-gen/pkg/jsonchema-generator"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "generate <values-yaml-file>",
		Short: "Generate JSON schema from values.yaml",
		Long:  "Generates a JSON schema from the provided Helm values.yaml file.",
		Args: func(cmd *cobra.Command, args []string) error {
			return helper.CheckArgsLength(len(args), "values file")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
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
			fmt.Println(s)
			return nil
		},
	}
}

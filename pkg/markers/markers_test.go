package markers_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/rhoat/helm-schema-gen/pkg/markers"
	"gopkg.in/yaml.v3"
)

func TestUncommentYAML(t *testing.T) {
	tests := []struct {
		name      string
		yamlInput string
		expected  *[]markers.SchemaInfo
	}{
		{
			name: "Test with +schemagen marker in key comment",
			yamlInput: `
map_gen_settings:
  # +schemagen:type:number
  width: 0 
`,
			expected: &[]markers.SchemaInfo{{
				Path:       ".map_gen_settings.width",
				SchemaType: "number",
			}},
		},
		{
			name: "Test with +schemagen marker in value comment",
			yamlInput: `
map_gen_settings:
  width: 0

  # +schemagen:type:number
  height: 0 
`,
			expected: &[]markers.SchemaInfo{{
				Path:       ".map_gen_settings.height",
				SchemaType: "number",
			}},
		},
		{
			name: "Test without +schemagen marker",
			yamlInput: `
map_gen_settings:
  width: 0
  height: 0
`,
			expected: &[]markers.SchemaInfo{}, // Expect an empty slice
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var schemaData = []markers.SchemaInfo{}
			// Unmarshal the input YAML into a yaml.Node
			var rootNode yaml.Node
			err := yaml.Unmarshal([]byte(tt.yamlInput), &rootNode)
			if err != nil {
				t.Fatalf("Error unmarshaling YAML: %v", err)
			}

			// Process the YAML nodes to handle comments with +schemagen markers
			rootNode = *markers.UncommentYAML(&rootNode, &schemaData, "")
			if diff := cmp.Diff(tt.expected, &schemaData); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

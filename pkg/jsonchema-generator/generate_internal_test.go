package jsonschema

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rhoat/helm-schema-gen/pkg/features"
)

func TestGenerate(t *testing.T) {
	testCases := []struct {
		desc     string
		mockYaml string
		expected Document
	}{
		{
			desc: "replace minReplicas",
			mockYaml: `
affinity: {}
autoscaling:
  enabled: true

  # +schemagen:type:boolean
  maxReplicas: 5
  minReplicas: 1
  targetCPUUtilizationPercentage: 80
`,
			expected: Document{
				Schema: "http://json-schema.org/schema#",
				property: property{
					Type: "object",
					Properties: map[string]*property{
						"affinity": &property{Type: "object"},
						"autoscaling": &property{
							Type: "object",
							Properties: map[string]*property{
								"enabled": &property{
									Type: "boolean",
								},
								"maxReplicas": &property{
									Type: "boolean",
								},
								"minReplicas": &property{
									Type: "integer",
								},
								"targetCPUUtilizationPercentage": &property{
									Type: "integer",
								},
							},
						},
					},
				}},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			features.Schemagen = true
			result, err := Generate(context.TODO(), strings.NewReader(tC.mockYaml))
			if err != nil {
				t.Errorf("failed to generate %s", err.Error())
			}
			jsonexpected, _ := json.MarshalIndent(tC.expected, "", "  ")
			jsondata, _ := json.MarshalIndent(result, "", "  ")
			t.Log(string(jsondata))
			if diff := cmp.Diff(jsonexpected, jsondata); diff != "" {
				t.Error(diff)
			}
		})
	}
}

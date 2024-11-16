package jsonschema

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestJSONSchemaTypes(t *testing.T) {

	// Table-driven tests for type validation
	tests := []struct {
		name                 string
		path                 string
		changeToType         string
		initialDoc, expected Document
	}{
		{
			name:         "autoscaling.maxReplicas",
			path:         ".autoscaling.maxReplicas",
			changeToType: "number",
			initialDoc: Document{
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
									Type: "integer",
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
									Type: "number",
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

	// Iterate over the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetTypeAtPath(&tt.initialDoc, tt.path, tt.changeToType)
			jsonexpected, _ := json.MarshalIndent(tt.expected, "", "  ")
			jsondata, _ := json.MarshalIndent(tt.initialDoc, "", "  ")
			t.Log(string(jsondata))
			if diff := cmp.Diff(jsonexpected, jsondata); diff != "" {
				t.Errorf("")
			}

		})
	}
}

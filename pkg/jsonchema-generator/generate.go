package jsonschema

import (
	"fmt"
	"io"
	"log"

	"github.com/rhoat/helm-schema-gen/pkg/markers"
	"gopkg.in/yaml.v3"
)

func Generate(yamlValuesFile io.Reader) (*Document, error) {
	values := make(map[string]interface{})
	valuesFileData, err := io.ReadAll(yamlValuesFile)
	if err != nil {
		return nil, fmt.Errorf("error when reading data: %w", err)
	}

	rootNode := yaml.Node{}
	err = yaml.Unmarshal(valuesFileData, &rootNode)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshaling file : %w", err)
	}

	schemaData := []markers.SchemaInfo{}
	rootNode = *markers.UncommentYAML(&rootNode, &schemaData, "")
	err = rootNode.Decode(&values)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshaling file : %w", err)
	}
	log.Printf("SchemaData %+v", schemaData)
	s := &Document{}
	s.ReadDeep(&values)
	// Apply changes to the data based on SchemaInfo
	for _, schemaInfo := range schemaData {
		log.Printf("%+v", schemaInfo)
		// Split the path into components (based on "/")
		if err := SetTypeAtPath(s, schemaInfo.Path, schemaInfo.SchemaType); err != nil {
			log.Printf("Error setting type at path %s: %v", schemaInfo.Path, err)
		}
	}
	return s, nil
}

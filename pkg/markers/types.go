package markers

import (
	"fmt"
	"log"
	"regexp"

	"gopkg.in/yaml.v3"
)

var (
	markerRegex = regexp.MustCompile(`(\+schemagen:.*)`)
	typeMarker  = regexp.MustCompile(`(\+schemagen:type:.*)`)
	getValue    = regexp.MustCompile(`^.*:(\w+)$`)
)

// This will be used to store schema information
type SchemaInfo struct {
	Path       string
	SchemaType string
}

func UncommentYAML(node *yaml.Node, schemaData *[]SchemaInfo, parentPath string) *yaml.Node {
	if schemaData == nil {
		// Initialize an empty slice if schemaData is nil
		schemaData = &[]SchemaInfo{}
	}
	switch node.Kind {
	case yaml.DocumentNode, yaml.SequenceNode:
		// Traverse through the sequence or document node
		for i, child := range node.Content {
			node.Content[i] = UncommentYAML(child, schemaData, parentPath)
		}
		return node
	case yaml.MappingNode:
		log.Printf("content %+v", node.Content)
		// Iterate through key-value pairs in a mapping node
		for i := 0; i < len(node.Content); i = i + 2 {
			key := node.Content[i]
			value := node.Content[i+1]
			fullPath := fmt.Sprintf("%s.%s", parentPath, key.Value)
			log.Printf("fullpath: %s", fullPath)
			log.Printf("HeadComment: %s FootComment: %s LineComment: %s ", key.HeadComment, key.FootComment, key.LineComment)
			// Check for +schemagen markers in both head and foot comments
			if len(key.HeadComment) > 0 {
				// Look for a +schemagen:type marker in the head comment
				matches := typeMarker.FindStringSubmatch(key.HeadComment)
				for _, match := range matches {
					log.Printf("match %s", match)
				}
				if len(matches) > 0 {
					schemaType := getValue.FindStringSubmatch(matches[1])
					if len(schemaType) > 1 {
						// Collect the schema type for this key
						*schemaData = append(*schemaData, SchemaInfo{
							Path:       fullPath,
							SchemaType: schemaType[1],
						})
					}
				}
			}

			// Recursively handle the content of values
			node.Content[i+1] = UncommentYAML(value, schemaData, fullPath)
		}
		return node
	default:
		return node
	}
}

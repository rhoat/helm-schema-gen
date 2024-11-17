package markers

import (
	"context"
	"regexp"

	"github.com/rhoat/helm-schema-gen/pkg/ctxlogger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

var (
	//nolint:unused // Keeping for future expansion of functionality.
	markerRegex = regexp.MustCompile(`(\+schemagen:.*)`)
	typeMarker  = regexp.MustCompile(`(\+schemagen:type:.*)`)
	getValue    = regexp.MustCompile(`^.*:(\w+)$`)
)

// This will be used to store schema information.
type SchemaInfo struct {
	Path       string
	SchemaType string
}

func UncommentYAML(ctx context.Context, node *yaml.Node, schemaData *[]SchemaInfo, parentPath string) *yaml.Node {
	logger := ctxlogger.GetLogger(ctx)
	if schemaData == nil {
		// Initialize an empty slice if schemaData is nil
		schemaData = &[]SchemaInfo{}
	}
	//nolint: exhaustive // These types are not needed yaml.ScalarNode, yaml.AliasNode
	switch node.Kind {
	case yaml.DocumentNode, yaml.SequenceNode:
		// Traverse through the sequence or document node
		for i, child := range node.Content {
			node.Content[i] = UncommentYAML(ctx, child, schemaData, parentPath)
		}
		return node
	case yaml.MappingNode:
		// Iterate through key-value pairs in a mapping node
		for i := 0; i < len(node.Content); i += 2 {
			key := node.Content[i]
			value := node.Content[i+1]
			fullPath := parentPath + "." + key.Value
			logger.Debug("Node", zap.Object("Fields", zapcore.ObjectMarshalerFunc(
				func(oe zapcore.ObjectEncoder) error {
					oe.AddString("Headcomment", key.HeadComment)
					oe.AddString("LineComment", key.LineComment)
					oe.AddString("FootComment", key.FootComment)
					oe.AddString("Key", key.Value)
					return nil
				},
			)))
			// Check for +schemagen markers in both head and foot comments
			if len(key.HeadComment) > 0 {
				logger.Debug("Item has Head Comment")
				// Look for a +schemagen:type marker in the head comment
				matches := typeMarker.FindStringSubmatch(key.HeadComment)
				if len(matches) > 0 {
					logger.Debug("Item has Head +schemagen:type marker")

					schemaType := getValue.FindStringSubmatch(matches[1])
					if len(schemaType) > 1 {
						logger.Debug("appending to list", zap.String("path", fullPath), zap.String("type", schemaType[1]))
						// Collect the schema type for this key
						*schemaData = append(*schemaData, SchemaInfo{
							Path:       fullPath,
							SchemaType: schemaType[1],
						})
					}
				}
			}

			// Recursively handle the content of values
			node.Content[i+1] = UncommentYAML(ctx, value, schemaData, fullPath)
		}
		return node
	default:
		return node
	}
}

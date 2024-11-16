package jsonschema

import (
	"context"
	"fmt"
	"io"

	"github.com/rhoat/helm-schema-gen/pkg/ctxlogger"
	"github.com/rhoat/helm-schema-gen/pkg/features"
	"github.com/rhoat/helm-schema-gen/pkg/markers"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

func Generate(ctx context.Context, yamlValuesFile io.Reader) (*Document, error) {
	logger := ctxlogger.GetLogger(ctx)
	logger.Debug("Reading yaml from io reader")
	valuesFileData, err := io.ReadAll(yamlValuesFile)
	if err != nil {
		return nil, fmt.Errorf("error when reading data: %w", err)
	}
	logger.Debug("unmarshaling the yaml")
	rootNode := yaml.Node{}
	err = yaml.Unmarshal(valuesFileData, &rootNode)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshaling file : %w", err)
	}
	values := make(map[string]interface{})
	err = rootNode.Decode(&values)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshaling file : %w", err)
	}
	s := &Document{}
	s.ReadDeep(&values)
	if features.Schemagen {
		logger.Debug("attempting to find markers")
		schemaData := []markers.SchemaInfo{}
		rootNode = *markers.UncommentYAML(ctx, &rootNode, &schemaData, "")
		logger.Debug("Found Schema Data", zap.Array("schemaData", zapcore.ArrayMarshalerFunc(
			func(ae zapcore.ArrayEncoder) error {
				for _, v := range schemaData {
					ae.AppendString(v.Path)
					ae.AppendString(v.SchemaType)
				}
				return nil
			},
		)))
		logger.Debug("placing values into document")
		// Apply changes to the data based on SchemaInfo
		for _, schemaInfo := range schemaData {
			logger.Debug(
				"attempting to modify document",
				zap.String("Path", schemaInfo.Path),
				zap.String("Type", schemaInfo.SchemaType),
			)
			if err = SetTypeAtPath(s, schemaInfo.Path, schemaInfo.SchemaType); err != nil {
				logger.Error(
					"Setting type at path",
					zap.String("Path", schemaInfo.Path),
					zap.String("Type", schemaInfo.SchemaType),
					zap.Error(err),
				)
			}
		}
	}
	return s, nil
}

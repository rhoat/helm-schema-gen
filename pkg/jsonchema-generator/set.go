package jsonschema

import (
	"errors"
	"fmt"
	"strings"
)

func SetTypeAtPath(doc *Document, path string, newType string) error {
	// Split the path into keys
	keys := strings.Split(path, ".")
	if len(keys) == 0 {
		return errors.New("invalid path")
	}
	keys = keys[1:]

	current := &doc.Properties
	for i, key := range keys {
		// Check if the key exists
		prop, exists := (*current)[key]
		if !exists {
			return fmt.Errorf("path not found: %s", key)
		}

		// If it's the last key, update the type
		if i == len(keys)-1 {
			prop.Type = newType
			(*current)[key] = prop
			return nil
		}

		// Move deeper into the structure
		current = &prop.Properties
	}
	return nil
}

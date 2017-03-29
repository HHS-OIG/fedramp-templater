package implementation

import (
	"github.com/opencontrol/fedramp-templater/common/source"
	"gopkg.in/fatih/set.v0"
	"log"
	"reflect"
	"strings"
)

// Key is a unique value that represents the different types of possible implementation status value.
type Key uint8

// Origination prefixes.
const (
	NoStatus Key = iota
	ImplementedImplementation
	PartialImplementation
	PlannedImplementation
	NotApplicableImplementation
)

// SrcMapping is a data structure that represents the text for a particular implementation in a particular source.
type SrcMapping map[source.Source]string

// IsDocMappingASubstrOf is wrapper that checks if the input string contains the SSP mapping.
// This is useful because the input string value may have extra characters so we can't do a == (equal to) check.
func (o SrcMapping) IsDocMappingASubstrOf(value string) bool {
	return strings.Contains(value, o[source.SSP])
}

// IsYAMLMappingEqualTo is a wrapper that checks if the input string equals to the YAML mapping.
func (o SrcMapping) IsYAMLMappingEqualTo(value string) bool {
	return value == o[source.YAML]
}

// GetSourceMappings returns a mapping of each implementation to their respective sources.
func GetSourceMappings() map[Key]SrcMapping {
	return map[Key]SrcMapping{
		ImplementedImplementation: {
			source.YAML: "complete",
			source.SSP:  "Implemented",
		},
		PartialImplementation: {
			source.YAML: "partial",
			source.SSP:  "Partially implemented",
		},
		PlannedImplementation: {
			source.YAML: "planned",
			source.SSP:  "Planned",
		},
		NotApplicableImplementation: {
			source.YAML: "none",
			source.SSP:  "Not applicable",
		},
		// For some reason the OpenControl Schema does not specify alternative implementation
		//AlternativeImplementation: {
		//	source.YAML: "<not-defined>",
		//	source.SSP:  "Alternative implementation",
		//},
	}
}

// ConvertSetToKeys will convert the set, which each value is of type interface{} to the Key.
func ConvertSetToKeys(s set.Interface) []Key {
	keys := []Key{}
	for _, item := range s.List() {
		key, isType := item.(Key)
		if isType {
			keys = append(keys, key)
		} else {
			log.Printf("Unable to use value as ImplementationStatus 'Key' Type: %v. Value: %v.\n",
				reflect.TypeOf(item), item)
		}
	}
	return keys
}

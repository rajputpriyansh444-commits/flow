package jq

import (
	"github.com/rajputpriyansh444-commits/flow/internal/executor/registry"
	"github.com/google/jsonschema-go/jsonschema"
)

var configSchema = &jsonschema.Schema{
	Type: "object",
	Properties: map[string]*jsonschema.Schema{
		"raw":   {Type: "boolean", Description: "Output raw strings without JSON encoding (like jq -r)"},
		"input": {Type: "string", Description: "File path to read JSON input from. Mutually exclusive with script."},
		"args": {
			Type:                 "object",
			Description:          "Variables to bind in the jq filter. Each key is available as $<key>; values keep their YAML type (strings may use ${...} references). Supplying args, including {}, makes the filter literal jq source.",
			AdditionalProperties: &jsonschema.Schema{},
		},
	},
}

func init() {
	registry.RegisterExecutorConfigSchema("jq", configSchema)
}

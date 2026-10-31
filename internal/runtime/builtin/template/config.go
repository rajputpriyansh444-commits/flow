package template

import (
	"github.com/rajputpriyansh444-commits/flow/internal/executor/registry"
	"github.com/google/jsonschema-go/jsonschema"
)

var configSchema = &jsonschema.Schema{
	Type:                 "object",
	AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
	Properties: map[string]*jsonschema.Schema{
		"data":         {Type: "object", Description: "Template data variables accessible as {{ .key }} in the template"},
		"output":       {Type: "string", Description: "File path to write the rendered output to. If empty, output is written to stdout."},
		"template_ref": {Type: "string", Description: "Complete scoped Flow value reference that resolves to template text."},
	},
}

func init() {
	registry.RegisterExecutorConfigSchema("template", configSchema)
}

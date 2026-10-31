// Package allproviders imports all LLM providers to register them.
// Import this package if you want all providers to be available:
//
//	import _ "github.com/rajputpriyansh444-commits/flow/internal/llm/allproviders"
package allproviders

import (
	_ "github.com/rajputpriyansh444-commits/flow/internal/llm/providers/anthropic"
	_ "github.com/rajputpriyansh444-commits/flow/internal/llm/providers/gemini"
	_ "github.com/rajputpriyansh444-commits/flow/internal/llm/providers/local"
	_ "github.com/rajputpriyansh444-commits/flow/internal/llm/providers/openai"
	_ "github.com/rajputpriyansh444-commits/flow/internal/llm/providers/openaicodex"
	_ "github.com/rajputpriyansh444-commits/flow/internal/llm/providers/openrouter"
	_ "github.com/rajputpriyansh444-commits/flow/internal/llm/providers/zai"
)

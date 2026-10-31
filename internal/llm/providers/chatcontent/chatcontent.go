// Package chatcontent builds message content for OpenAI-compatible chat
// completion APIs.
package chatcontent

import "github.com/rajputpriyansh444-commits/flow/internal/llm"

// Content returns the text as a plain string when there are no images, and
// otherwise a list of content parts with the images ahead of the text.
func Content(text string, images []llm.Image) any {
	if len(images) == 0 {
		return text
	}
	parts := make([]any, 0, len(images)+1)
	for _, image := range images {
		parts = append(parts, map[string]any{
			"type":      "image_url",
			"image_url": map[string]string{"url": image.DataURL()},
		})
	}
	if text != "" {
		parts = append(parts, map[string]any{"type": "text", "text": text})
	}
	return parts
}

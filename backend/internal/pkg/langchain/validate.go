package langchain

import (
	"context"
	"fmt"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/telemetry"
	"go.opentelemetry.io/otel/codes"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

// NormalizeLLMResponse provides the generated text from an LLM response and asks the LLM
// to ensure it is restructured to valid JSON.
func NormalizeLLMResponse(ctx context.Context, input string, llm *ollama.LLM) (string, error) {
	ctx, span := telemetry.StartSpan(ctx, telemetry.WithSpanName("Normalize LLM Response"))
	defer span.End()
	grounding := `
	Please pretty print the following into valid json. 
	Do not include any markdown, new lines, or extra whitespace. Please specifically remove any 
	backticks and markdown indicating this is a json object.
	Please remove any you find in the provided text: 
	
	%s
	so that the output is only a valid json object and nothing else.
	`

	recommendation, err := llms.GenerateFromSinglePrompt(ctx, llm, fmt.Sprintf(grounding, input))
	if err != nil {
		span.RecordError(err)
		return "", err
	}
	span.SetStatus(codes.Ok, "normalization complete")
	return recommendation, nil
}

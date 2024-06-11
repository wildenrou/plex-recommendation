package langchain

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func GenerateRecommendation(ctx context.Context, recentlyViewed, fullCollection string, llm *ollama.LLM) (string, error) {
	log.Println("generating recommendation...")
	grounding := `Please recommend me up to 3 different movies to watch based on my recent watch
	history provided here: %+v. Please do not suggest any titles that do not exist in the following 
	collection: %+v. Do not recommend me any titles that have a content rating exceeding the highest
	content rating in my recent watch history. I have a child in the family, so take into consideration
	something a toddler would also enjoy watching. Please provide your recommendation as a json array of
	objects, whose members have this shape:
	{
		"title": title,
		"summary": summary,
		"content_rating": content_rating
	}
	Please do not recomment more than 3 titles. Please do ensure your response is valid json before
	returning it to me.
	`

	recommendation, err := llms.GenerateFromSinglePrompt(ctx, llm, fmt.Sprintf(grounding, recentlyViewed, fullCollection))
	if err != nil {
		return "", err
	}

	log.Println("generated")
	return recommendation, nil

}

func GenerateSimpleRecommendation(ctx context.Context, llm *ollama.LLM) (string, error) {
	log.Println("generating recommendation...")
	grounding := "Hello! How are you today?"

	recommendation, err := llms.GenerateFromSinglePrompt(ctx, llm, grounding)
	if err != nil {
		return "", err
	}

	log.Println("generated")
	return recommendation, nil
}

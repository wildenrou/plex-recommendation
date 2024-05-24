package langchain

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
)

func GenerateRecommendation(ctx context.Context, recentlyViewed, fullCollection string, llm *llms.Model) (string, error) {
	log.Println("generating recommendation...")
	grounding := `Please recommend me up to five different movies to watch based on my recent watch
	history provided here: %+v. Please do not suggest any titles that do not exist in the following 
	collection: %+v. Do not recommend me any titles that have a content rating exceeding the highest
	content rating in my recent watch history. I have a child in the family, so take into consideration
	something a toddler would also enjoy watching. Provide your recommendations in the following format:
	Title (ContentRating): Summary
	`

	recommendation, err := llms.GenerateFromSinglePrompt(ctx, *llm, fmt.Sprintf(grounding, recentlyViewed, fullCollection))
	if err != nil {
		return "", err
	}

	log.Println("generated\n")
	return recommendation, nil

}

func GenerateSimpleRecommendation(ctx context.Context, llm *llms.Model) (string, error) {
	log.Println("generating recommendation...")
	grounding := "Hello! How are you today?"

	recommendation, err := llms.GenerateFromSinglePrompt(ctx, *llm, grounding)
	if err != nil {
		return "", err
	}

	log.Println("generated")
	return recommendation, nil
}

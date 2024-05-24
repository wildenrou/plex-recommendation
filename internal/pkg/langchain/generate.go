package langchain

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/plex"
)

func GenerateRecommendation(ctx context.Context, videos []plex.VideoShort, llm *llms.Model) (string, error) {
	log.Println("generating recommendation...")
	grounding := `Given the following videos from my watch history, 
	recommend me one or two videos to watch next. Ensure that the rating on your 
	recommendations do not exceed that of my latest viewings 
	(EG, if I have recently watching only G and PG rated films, do not 
	recommend a PG-13 rated film to me.). Provide me the title, summary, 
	and letter rating of the film in your recommentation. 
	
	%+v\n`

	recommendation, err := llms.GenerateFromSinglePrompt(ctx, *llm, fmt.Sprintf(grounding, videos))
	if err != nil {
		return "", err
	}

	log.Println("generated")
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

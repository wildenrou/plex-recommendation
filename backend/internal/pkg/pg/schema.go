package pg

import (
	"gorm.io/gorm"
)

type RecommendationCache struct {
	gorm.Model
	// InputTitles are the titles provided to the LLM
	// that it used as grounding to generate the output
	// cached in GeneratedOutput
	InputTitles string
	// GeneratedOutput is the response from the LLM
	// when ased for a recommendation using the
	// provided InputTitles
	GeneratedOutput string
}

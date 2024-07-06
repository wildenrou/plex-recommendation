package httpinternal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/wgeorgecook/plex-recommendation/internal/pkg/plex"
)

func formatHttpError(err error) []byte {
	return []byte(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
}

func recommendationHandler(w http.ResponseWriter, r *http.Request) {
	section := r.PathValue("movieSection")
	var limit int
	limitQuery, ok := r.URL.Query()["limit"]
	if ok {
		limit, _ = strconv.Atoi(limitQuery[0])
	}

	recommendation, err := getRecommendation(r.Context(), section, limit)
	if err != nil {
		w.Write(formatHttpError(err))
		return
	}

	var respStruct []*plex.VideoShort
	if err := json.Unmarshal([]byte(recommendation), &respStruct); err != nil {
		w.Write(formatHttpError(err))
		return
	}
	respBytes, err := json.Marshal(&respStruct)
	if err != nil {
		w.Write(formatHttpError(err))
		return
	}

	_, err = w.Write(respBytes)
	if err != nil {
		log.Println("could not write back to client: ", err.Error())
	}
}

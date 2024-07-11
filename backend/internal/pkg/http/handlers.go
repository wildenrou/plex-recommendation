package httpinternal

import (
	"encoding/json"
	"fmt"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"log"
	"net/http"
	"strconv"

	"github.com/wgeorgecook/plex-recommendation/internal/pkg/plex"
)

func formatHttpError(err error) []byte {
	return []byte(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
}

const recommendationPathway = "/recommendation/{movieSection}"

func recommendationHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := telemetry.StartSpan(r.Context(),
		telemetry.WithSpanName("Get Recommendation HTTP Handler"),
		telemetry.WithSpanPackage("httpinternal"),
	)
	defer span.End()
	section := r.PathValue("movieSection")
	span.SetAttributes(attribute.String("movieSection", section))
	var limit int
	limitQuery, ok := r.URL.Query()["limit"]
	if ok {
		limit, _ = strconv.Atoi(limitQuery[0])
	}
	span.SetAttributes(attribute.Int("limit", limit))

	recommendation, err := getRecommendation(ctx, section, limit)
	if err != nil {
		w.Write(formatHttpError(err))
		span.SetStatus(codes.Error, err.Error())
		return
	}
	span.AddEvent("recommendation generated")
	var respStruct []*plex.VideoShort
	if err := json.Unmarshal([]byte(recommendation), &respStruct); err != nil {
		w.Write(formatHttpError(err))
		span.SetStatus(codes.Error, err.Error())
		return
	}
	respBytes, err := json.Marshal(&respStruct)
	if err != nil {
		w.Write(formatHttpError(err))
		span.SetStatus(codes.Error, err.Error())
		return
	}
	span.AddEvent("marshal complete")
	_, err = w.Write(respBytes)
	if err != nil {
		log.Println("could not write back to client: ", err.Error())
	}
	span.AddEvent("write complete")
	span.SetStatus(codes.Ok, "recommendation successfully retrieved")
}

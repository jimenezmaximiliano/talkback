package talkback

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// LogError represents a function that can log errors.
type LogError func(ctx context.Context, err error)

// RespondWithJSON responds to an HTTP request with the given JSON body.
func RespondWithJSON(
	ctx context.Context,
	logError LogError,
	responseWriter http.ResponseWriter,
	statusCode int,
	payload any,
) {
	body, err := json.Marshal(payload)
	if err != nil {
		logError(ctx, errors.Wrap(err, "failed to marshal JSON response"))
		responseWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	if string(body) == "null" {
		logError(ctx, errors.New("marshalled body was null"))
		responseWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)

	if _, err := responseWriter.Write(body); err != nil {
		logError(ctx, errors.Wrap(err, "failed to write JSON body to response"))
		responseWriter.WriteHeader(http.StatusInternalServerError)

		return
	}
}

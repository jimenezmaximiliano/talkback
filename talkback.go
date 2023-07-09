package talkback

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// LogError represents a function that can log errors.
type LogError func(err error)

// Talkback is a service that provides helper functions to respond to http requests.
type Talkback struct {
	logError LogError
}

// NewTalkback is a constructor for the Talkback service.
func NewTalkback(logError LogError) Talkback {
	return Talkback{
		logError: logError,
	}
}

// RespondSuccessWithJSON responds to an HTTP request with a 200 and the given JSON body.
func (service Talkback) RespondSuccessWithJSON(responseWriter http.ResponseWriter, JSONBody []byte) {
	setContentTypeToJSON(responseWriter)

	if _, err := responseWriter.Write(JSONBody); err != nil {
		service.logInternalErrorAndRespond(responseWriter, errors.Wrap(err, "failed to respond with JSON body"))

		return
	}
}

// RespondCreatedWithJSON responds to an HTTP request with a 201 and the given JSON body.
func (service Talkback) RespondCreatedWithJSON(responseWriter http.ResponseWriter, JSONBody []byte) {
	setContentTypeToJSON(responseWriter)
	responseWriter.WriteHeader(http.StatusCreated)

	if _, err := responseWriter.Write(JSONBody); err != nil {
		service.logInternalErrorAndRespond(responseWriter, errors.Wrap(err, "failed to respond with JSON body"))

		return
	}
}

// RespondWithBadRequestJSONMessage responds to an HTTP request with a 400 and a JSON payload with the error message provided.
// In case of an error while responding, it logs an internal error and responds with a 500.
func (service Talkback) RespondWithBadRequestJSONMessage(responseWriter http.ResponseWriter, errorMessage string) {
	setContentTypeToJSON(responseWriter)
	responseWriter.WriteHeader(http.StatusBadRequest)

	JSONPayload, err := json.Marshal(struct {
		Error string `json:"error"`
	}{
		Error: errorMessage,
	})
	if err != nil {
		service.logInternalErrorAndRespond(responseWriter, errors.Wrap(err, "failed to marshal JSON response"))

		return
	}

	if _, err := responseWriter.Write(JSONPayload); err != nil {
		service.logInternalErrorAndRespond(responseWriter, errors.Wrap(err, "failed to respond with JSON body"))

		return
	}
}

func setContentTypeToJSON(responseWriter http.ResponseWriter) {
	responseWriter.Header().Set("Content-Type", "application/json")
}

// logInternalErrorAndRespond logs an error as Error and responds with a 500.
func (service Talkback) logInternalErrorAndRespond(responseWriter http.ResponseWriter, err error) {
	service.logError(err)
	responseWriter.WriteHeader(http.StatusInternalServerError)
}

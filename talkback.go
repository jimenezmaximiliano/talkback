package talkback

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// LogError represents a function that logs an error.
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

// RedirectTo responds to an HTTP request with a temporary redirection to the given URL.
func (service Talkback) RedirectTo(responseWriter http.ResponseWriter, request *http.Request, url string) {
	http.Redirect(responseWriter, request, url, http.StatusTemporaryRedirect)
}

// LogInternalErrorAndRespond logs an error as Error and responds with a 500.
func (service Talkback) LogInternalErrorAndRespond(responseWriter http.ResponseWriter, err error) {
	service.logError(err)
	responseWriter.WriteHeader(http.StatusInternalServerError)
}

// RespondSuccessWithJSON responds to an HTTP request with a 200 and the given JSON body.
func (service Talkback) RespondSuccessWithJSON(responseWriter http.ResponseWriter, JSONBody []byte) {
	responseWriter.Header().Set("Content-Type", "application/json")

	if _, err := responseWriter.Write(JSONBody); err != nil {
		service.LogInternalErrorAndRespond(responseWriter, errors.Wrap(err, "failed to respond with JSON body"))
		return
	}
}

// RespondSuccess responds to an HTTP request with a 200.
func (service Talkback) RespondSuccess(responseWriter http.ResponseWriter) {
	responseWriter.WriteHeader(http.StatusOK)
}

// RespondUnauthorized responds to an HTTP request with a 401.
func (service Talkback) RespondUnauthorized(responseWriter http.ResponseWriter) {
	responseWriter.WriteHeader(http.StatusUnauthorized)
}

// RespondCreatedWithJSON responds to an HTTP request with a 201 and the given JSON body.
func (service Talkback) RespondCreatedWithJSON(responseWriter http.ResponseWriter, JSONBody []byte) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusCreated)
	if _, err := responseWriter.Write(JSONBody); err != nil {
		service.LogInternalErrorAndRespond(responseWriter, errors.Wrap(err, "failed to respond with JSON body"))
		return
	}
}

// RespondWithBadRequest responds to an HTTP request with a 400 and a JSON payload with the error message provided.
// In case of an error while responding, it logs an internal error and responds with a 500.
func (service Talkback) RespondWithBadRequest(responseWriter http.ResponseWriter, errorMessage string) {
	responseWriter.WriteHeader(http.StatusBadRequest)
	responseWriter.Header().Set("Content-Type", "application/json")
	JSONPayload, err := json.Marshal(struct {
		Error string
	}{
		Error: errorMessage,
	})
	if err != nil {
		service.LogInternalErrorAndRespond(responseWriter, errors.Wrap(err, "failed to marshal JSON response"))
		return
	}
	if _, err := responseWriter.Write(JSONPayload); err != nil {
		service.LogInternalErrorAndRespond(responseWriter, errors.Wrap(err, "failed to respond with JSON body"))
		return
	}
}

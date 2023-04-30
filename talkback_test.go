package talkback_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jimenezmaximiliano/talkback"
	"github.com/jimenezmaximiliano/talkback/mocks"
)

var jsonBody = []byte("{ \"property\": \"value\" }")

func TestRedirection(test *testing.T) {
	test.Parallel()

	errorLoggerFunc := mocks.NewLogError(test)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "https://obladi.com", nil)
	talk := talkback.NewTalkback(errorLoggerFunc.Execute)

	talk.RedirectTo(responseRecorder, request, "https://oblada.com")

	response := responseRecorder.Result()
	url, err := response.Location()
	require.NoError(test, err)

	assert.Equal(test, http.StatusTemporaryRedirect, response.StatusCode)
	assert.Equal(test, "https://oblada.com", url.String())
}

func TestRespondingWithACreatedJSON(test *testing.T) {
	test.Parallel()

	errorLoggerFunc := mocks.NewLogError(test)
	responseRecorder := httptest.NewRecorder()
	talk := talkback.NewTalkback(errorLoggerFunc.Execute)

	talk.RespondCreatedWithJSON(responseRecorder, jsonBody)

	response := responseRecorder.Result()
	body, err := io.ReadAll(response.Body)
	require.NoError(test, err)

	assert.Equal(test, http.StatusCreated, response.StatusCode)
	assert.Equal(test, "application/json", response.Header.Get("Content-Type"))
	assert.Equal(test, jsonBody, body)
}

func TestRespondingWithOkJSON(test *testing.T) {
	test.Parallel()

	errorLoggerFunc := mocks.NewLogError(test)
	responseRecorder := httptest.NewRecorder()
	talk := talkback.NewTalkback(errorLoggerFunc.Execute)

	talk.RespondSuccessWithJSON(responseRecorder, jsonBody)

	response := responseRecorder.Result()
	body, err := io.ReadAll(response.Body)
	require.NoError(test, err)

	assert.Equal(test, http.StatusOK, response.StatusCode)
	assert.Equal(test, "application/json", response.Header.Get("Content-Type"))
	assert.Equal(test, jsonBody, body)
}

func TestRespondingWithBadRequestJSONMessage(test *testing.T) {
	test.Parallel()

	errorLoggerFunc := mocks.NewLogError(test)
	responseRecorder := httptest.NewRecorder()
	talk := talkback.NewTalkback(errorLoggerFunc.Execute)

	talk.RespondWithBadRequestJSONMessage(responseRecorder, "error message")

	response := responseRecorder.Result()
	body, err := io.ReadAll(response.Body)
	require.NoError(test, err)

	assert.Equal(test, http.StatusBadRequest, response.StatusCode)
	assert.Equal(test, "application/json", response.Header.Get("Content-Type"))
	assert.Equal(test, "{\"error\":\"error message\"}", string(body))
}

func TestRespondingWithInternalServerErrorAndLogging(test *testing.T) {
	test.Parallel()

	errorLoggerFunc := mocks.NewLogError(test)
	responseRecorder := httptest.NewRecorder()
	talk := talkback.NewTalkback(errorLoggerFunc.Execute)
	internalErr := errors.New("oops")

	errorLoggerFunc.On("Execute", internalErr).Return()

	talk.LogInternalErrorAndRespond(responseRecorder, internalErr)

	response := responseRecorder.Result()

	assert.Equal(test, http.StatusInternalServerError, response.StatusCode)
}

func TestRespondingWithSuccess(test *testing.T) {
	test.Parallel()

	errorLoggerFunc := mocks.NewLogError(test)
	responseRecorder := httptest.NewRecorder()
	talk := talkback.NewTalkback(errorLoggerFunc.Execute)

	talk.RespondSuccess(responseRecorder)
	response := responseRecorder.Result()

	assert.Equal(test, http.StatusOK, response.StatusCode)
}

func TestRespondingUnauthorized(test *testing.T) {
	test.Parallel()

	errorLoggerFunc := mocks.NewLogError(test)
	responseRecorder := httptest.NewRecorder()
	talk := talkback.NewTalkback(errorLoggerFunc.Execute)

	talk.RespondUnauthorized(responseRecorder)
	response := responseRecorder.Result()

	assert.Equal(test, http.StatusUnauthorized, response.StatusCode)
}

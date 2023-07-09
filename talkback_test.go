package talkback_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jimenezmaximiliano/talkback"
	"github.com/jimenezmaximiliano/talkback/mocks"
)

var jsonBody = []byte("{ \"property\": \"value\" }")

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

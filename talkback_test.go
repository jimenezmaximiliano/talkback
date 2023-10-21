package talkback_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/jimenezmaximiliano/talkback"
	"github.com/jimenezmaximiliano/talkback/mocks"
)

var payload = struct {
	Message string
}{
	Message: "talkback",
}

func TestRespondingWithJSON(test *testing.T) {
	test.Parallel()

	errorLoggerFunc := mocks.NewLogError(test)
	responseRecorder := httptest.NewRecorder()

	talkback.RespondWithJSON(context.Background(), errorLoggerFunc.Execute, responseRecorder, http.StatusOK, payload)

	response := responseRecorder.Result()
	body, err := io.ReadAll(response.Body)
	require.NoError(test, err)

	assert.Equal(test, http.StatusOK, response.StatusCode)
	assert.Equal(test, "application/json", response.Header.Get("Content-Type"))
	assert.JSONEq(test, `{"Message":"talkback"}`, string(body))
}

func TestRespondingWithJSONFails(test *testing.T) {
	test.Parallel()

	errorLoggerFunc := mocks.NewLogError(test)
	responseRecorder := httptest.NewRecorder()
	errorLoggerFunc.
		On("Execute", mock.Anything, mock.Anything).
		Return().
		Once()

	talkback.RespondWithJSON(context.Background(), errorLoggerFunc.Execute, responseRecorder, http.StatusOK, nil)

	response := responseRecorder.Result()

	assert.Equal(test, http.StatusInternalServerError, response.StatusCode)
}

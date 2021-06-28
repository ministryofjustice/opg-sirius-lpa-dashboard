package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockFeedbackClient struct {
	feedback struct {
		count       int
		lastCtx     sirius.Context
		lastMessage string
		err         error
	}
}

func (m *mockFeedbackClient) Feedback(ctx sirius.Context, message string) error {
	m.feedback.count += 1
	m.feedback.lastCtx = ctx
	m.feedback.lastMessage = message

	return m.feedback.err
}

func TestGetFeedback(t *testing.T) {
	assert := assert.New(t)

	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r.Header.Add("Referer", "http://example.com/previous")

	err := feedback(nil, template)(w, r)
	assert.Nil(err)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(feedbackVars{
		XSRFToken: getContext(r).XSRFToken,
		Redirect:  "http://example.com/previous",
	}, template.lastVars)
}

func TestPostFeedback(t *testing.T) {
	assert := assert.New(t)

	client := &mockFeedbackClient{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("redirect=a&feedback=b"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := feedback(client, nil)(w, r)
	assert.Equal(RedirectError("a"), err)

	assert.Equal(1, client.feedback.count)
	assert.Equal(getContext(r), client.feedback.lastCtx)
	assert.Equal("b", client.feedback.lastMessage)
}

func TestPostFeedbackError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockFeedbackClient{}
	client.feedback.err = expectedError

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("redirect=a&feedback=b"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := feedback(client, nil)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.feedback.count)
	assert.Equal(getContext(r), client.feedback.lastCtx)
	assert.Equal("b", client.feedback.lastMessage)
}

func TestBadMethodFeedback(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/path", nil)

	err := feedback(nil, nil)(w, r)
	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)
}

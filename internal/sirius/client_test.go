package sirius

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func newPact() *dsl.Pact {
	return &dsl.Pact{
		Consumer:          "sirius-lpa-dashboard",
		Provider:          "sirius",
		Host:              "localhost",
		PactFileWriteMode: "merge",
		LogDir:            "../../logs",
		PactDir:           "../../pacts",
	}
}

func teapotServer() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		}),
	)
}

func getContext(cookies []*http.Cookie) Context {
	return Context{
		Context:   context.Background(),
		Cookies:   cookies,
		XSRFToken: "abcde",
	}
}

func TestClientError(t *testing.T) {
	assert.Equal(t, "message", ClientError("message").Error())
}

func TestStatusError(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/some/url", nil)

	resp := &http.Response{
		StatusCode: http.StatusTeapot,
		Request:    req,
	}

	err := newStatusError(resp)

	assert.Equal(t, "POST /some/url returned 418", err.Error())
	assert.Equal(t, "unexpected response from Sirius", err.Title())
	assert.Equal(t, err, err.Data())
}

package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockLogger struct {
	count       int
	lastRequest *http.Request
	lastError   error
}

func (m *mockLogger) Request(r *http.Request, err error) {
	m.count += 1
	m.lastRequest = r
	m.lastError = err
}

type mockTemplate struct {
	count    int
	lastName string
	lastVars interface{}
}

func (m *mockTemplate) ExecuteTemplate(w io.Writer, name string, vars interface{}) error {
	m.count += 1
	m.lastName = name
	m.lastVars = vars
	return nil
}

func TestNew(t *testing.T) {
	assert.Implements(t, (*http.Handler)(nil), New(nil, nil, nil, "", "", "", ""))
}

func TestSecurityHeaders(t *testing.T) {
	assert := assert.New(t)

	handler := securityHeaders(http.NotFoundHandler())

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler.ServeHTTP(w, r)

	resp := w.Result()

	assert.Equal("default-src 'self'", resp.Header.Get("Content-Security-Policy"))
	assert.Equal("same-origin", resp.Header.Get("Referrer-Policy"))
	assert.Equal("max-age=31536000; includeSubDomains; preload", resp.Header.Get("Strict-Transport-Security"))
	assert.Equal("nosniff", resp.Header.Get("X-Content-Type-Options"))
	assert.Equal("SAMEORIGIN", resp.Header.Get("X-Frame-Options"))
	assert.Equal("1; mode=block", resp.Header.Get("X-XSS-Protection"))
}

func TestErrorHandler(t *testing.T) {
	assert := assert.New(t)

	tmplError := &mockTemplate{}

	wrap := errorHandler(nil, tmplError, "/prefix", "http://sirius")
	handler := wrap(func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusTeapot)
		return nil
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler.ServeHTTP(w, r)

	resp := w.Result()

	assert.Equal(http.StatusTeapot, resp.StatusCode)
	assert.Equal(0, tmplError.count)
}

func TestErrorHandlerUnauthorized(t *testing.T) {
	assert := assert.New(t)

	tmplError := &mockTemplate{}

	wrap := errorHandler(nil, tmplError, "/prefix", "http://sirius")
	handler := wrap(func(w http.ResponseWriter, r *http.Request) error {
		return sirius.ErrUnauthorized
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusFound, resp.StatusCode)
	assert.Equal("http://sirius/auth", resp.Header.Get("Location"))

	assert.Equal(0, tmplError.count)
}

func TestErrorHandlerRedirect(t *testing.T) {
	assert := assert.New(t)

	tmplError := &mockTemplate{}

	wrap := errorHandler(nil, tmplError, "/prefix", "http://sirius")
	handler := wrap(func(w http.ResponseWriter, r *http.Request) error {
		return RedirectError("/here")
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusFound, resp.StatusCode)
	assert.Equal("/prefix/here", resp.Header.Get("Location"))

	assert.Equal(0, tmplError.count)
}

func TestErrorHandlerStatus(t *testing.T) {
	assert := assert.New(t)

	logger := &mockLogger{}
	tmplError := &mockTemplate{}

	wrap := errorHandler(logger, tmplError, "/prefix", "http://sirius")
	handler := wrap(func(w http.ResponseWriter, r *http.Request) error {
		return StatusError(http.StatusTeapot)
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusInternalServerError, resp.StatusCode)

	assert.Equal(1, tmplError.count)
	assert.Equal(errorVars{SiriusURL: "http://sirius", Code: http.StatusInternalServerError, Error: "418 I'm a teapot"}, tmplError.lastVars)

	assert.Equal(1, logger.count)
	assert.Equal(r, logger.lastRequest)
	assert.Equal(StatusError(http.StatusTeapot), logger.lastError)
}

func TestErrorHandlerStatusKnown(t *testing.T) {
	for name, code := range map[string]int{
		"Forbidden": http.StatusForbidden,
		"Not Found": http.StatusNotFound,
	} {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			logger := &mockLogger{}
			tmplError := &mockTemplate{}

			wrap := errorHandler(logger, tmplError, "/prefix", "http://sirius")
			handler := wrap(func(w http.ResponseWriter, r *http.Request) error {
				return StatusError(code)
			})

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "/path", nil)

			handler.ServeHTTP(w, r)

			resp := w.Result()
			assert.Equal(code, resp.StatusCode)

			assert.Equal(1, tmplError.count)
			assert.Equal(errorVars{SiriusURL: "http://sirius", Code: code, Error: fmt.Sprintf("%d %s", code, name)}, tmplError.lastVars)

			assert.Equal(1, logger.count)
			assert.Equal(r, logger.lastRequest)
			assert.Equal(StatusError(code), logger.lastError)
		})
	}
}

func TestGetContext(t *testing.T) {
	assert := assert.New(t)

	r, _ := http.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "XSRF-TOKEN", Value: "z3tVRZ00yx4dHz3KWYv3boLWHZ4/RsCsVAKbvo2SBNc%3D"})
	r.AddCookie(&http.Cookie{Name: "another", Value: "one"})

	ctx := getContext(r)
	assert.Equal(r.Context(), ctx.Context)
	assert.Equal(r.Cookies(), ctx.Cookies)
	assert.Equal("z3tVRZ00yx4dHz3KWYv3boLWHZ4/RsCsVAKbvo2SBNc=", ctx.XSRFToken)
}

func TestGetContextBadXSRFToken(t *testing.T) {
	assert := assert.New(t)

	r, _ := http.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "XSRF-TOKEN", Value: "%"})
	r.AddCookie(&http.Cookie{Name: "another", Value: "one"})

	ctx := getContext(r)
	assert.Equal(r.Context(), ctx.Context)
	assert.Equal(r.Cookies(), ctx.Cookies)
	assert.Equal("", ctx.XSRFToken)
}

func TestGetContextMissingXSRFToken(t *testing.T) {
	assert := assert.New(t)

	r, _ := http.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "another", Value: "one"})

	ctx := getContext(r)
	assert.Equal(r.Context(), ctx.Context)
	assert.Equal(r.Cookies(), ctx.Cookies)
	assert.Equal("", ctx.XSRFToken)
}

func TestGetContextForPostRequest(t *testing.T) {
	assert := assert.New(t)

	r, _ := http.NewRequest("POST", "/", strings.NewReader("xsrfToken=the-real-one"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{Name: "XSRF-TOKEN", Value: "z3tVRZ00yx4dHz3KWYv3boLWHZ4/RsCsVAKbvo2SBNc%3D"})
	r.AddCookie(&http.Cookie{Name: "another", Value: "one"})

	ctx := getContext(r)
	assert.Equal(r.Context(), ctx.Context)
	assert.Equal(r.Cookies(), ctx.Cookies)
	assert.Equal("the-real-one", ctx.XSRFToken)
}

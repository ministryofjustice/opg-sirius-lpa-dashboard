package server

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type Logger interface {
	Request(*http.Request, error)
}

type Client interface {
	AllCasesClient
	CentralCasesClient
	MarkWorkedClient
	PendingCasesClient
	RequestNextCasesClient
	TasksClient
}

type Template interface {
	ExecuteTemplate(io.Writer, string, interface{}) error
}

func New(logger Logger, client Client, templates map[string]*template.Template, prefix, siriusURL, siriusPublicURL, webDir string) http.Handler {
	wrap := errorHandler(logger, templates["error.gotmpl"], prefix, siriusPublicURL)

	mux := http.NewServeMux()

	mux.Handle("/", wrap(landingPage(client)))

	mux.Handle("/pending-cases",
		wrap(
			pendingCases(client, templates["pending-cases.gotmpl"])))

	mux.Handle("/tasks",
		wrap(
			tasks(client, templates["tasks.gotmpl"])))

	mux.Handle("/all-cases",
		wrap(
			allCases(client, templates["all-cases.gotmpl"])))

	mux.Handle("/teams/central",
		wrap(
			centralCases(client, templates["central-cases.gotmpl"])))

	mux.Handle("/request-next-cases",
		wrap(
			requestNextCases(client)))

	mux.Handle("/mark-worked",
		wrap(
			markWorked(client)))

	mux.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {})

	static := http.FileServer(http.Dir(webDir + "/static"))
	mux.Handle("/assets/", static)
	mux.Handle("/javascript/", static)
	mux.Handle("/stylesheets/", static)

	return http.StripPrefix(prefix, mux)
}

type RedirectError string

func (e RedirectError) Error() string {
	return "redirect to " + string(e)
}

func (e RedirectError) To() string {
	return string(e)
}

type StatusError int

func (e StatusError) Error() string {
	code := e.Code()

	return fmt.Sprintf("%d %s", code, http.StatusText(code))
}

func (e StatusError) Code() int {
	return int(e)
}

type Handler func(w http.ResponseWriter, r *http.Request) error

type errorVars struct {
	SiriusURL string
	Path      string

	Code  int
	Error string
}

func errorHandler(logger Logger, tmplError Template, prefix, siriusURL string) func(next Handler) http.Handler {
	return func(next Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := next(w, r)

			if err != nil {
				if err == sirius.ErrUnauthorized {
					http.Redirect(w, r, siriusURL+"/auth", http.StatusFound)
					return
				}

				if redirect, ok := err.(RedirectError); ok {
					http.Redirect(w, r, prefix+redirect.To(), http.StatusFound)
					return
				}

				logger.Request(r, err)

				code := http.StatusInternalServerError
				if status, ok := err.(StatusError); ok {
					if status.Code() == http.StatusForbidden || status.Code() == http.StatusNotFound {
						code = status.Code()
					}
				}

				w.WriteHeader(code)
				err = tmplError.ExecuteTemplate(w, "page", errorVars{
					SiriusURL: siriusURL,
					Path:      "",
					Code:      code,
					Error:     err.Error(),
				})

				if err != nil {
					logger.Request(r, err)
					http.Error(w, "Could not generate error template", http.StatusInternalServerError)
				}
			}
		})
	}
}

func getContext(r *http.Request) sirius.Context {
	token := ""

	if r.Method == http.MethodGet {
		if cookie, err := r.Cookie("XSRF-TOKEN"); err == nil {
			token, _ = url.QueryUnescape(cookie.Value)
		}
	} else {
		token = r.FormValue("xsrfToken")
	}

	return sirius.Context{
		Context:   r.Context(),
		Cookies:   r.Cookies(),
		XSRFToken: token,
	}
}

func getPage(r *http.Request) int {
	page := r.FormValue("page")
	if page == "" {
		return 1
	}

	v, err := strconv.Atoi(page)
	if err != nil {
		return 1
	}

	return v
}

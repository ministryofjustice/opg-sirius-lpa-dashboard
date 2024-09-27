package server

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Client interface {
	AllCasesClient
	TasksDashboardClient
	CentralCasesClient
	FeedbackClient
	MarkWorkedClient
	PendingCasesClient
	ReassignClient
	RedirectClient
	RequestNextCasesClient
	RequestNextTaskClient
	TasksClient
	TeamWorkInProgressClient
	UserAllCasesClient
	UserPendingCasesClient
	UserTasksClient
}

type Template interface {
	ExecuteTemplate(io.Writer, string, interface{}) error
}

func New(logger *slog.Logger, client Client, templates map[string]*template.Template, prefix, siriusURL, siriusPublicURL, webDir string) http.Handler {
	wrap := errorHandler(templates["error.gotmpl"], prefix, siriusPublicURL)

	mux := http.NewServeMux()

	mux.Handle("/", wrap(redirect(client)))

	mux.Handle("/pending-cases",
		wrap(
			pendingCases(client, templates["pending-cases.gotmpl"])))

	mux.Handle("/tasks-dashboard",
		wrap(
			tasksDashboard(client, templates["tasks-dashboard.gotmpl"])))

	mux.Handle("/tasks",
		wrap(
			tasks(client, templates["tasks.gotmpl"])))

	mux.Handle("/all-cases",
		wrap(
			allCases(client, templates["all-cases.gotmpl"])))

	mux.Handle("/teams/central",
		wrap(
			centralCases(client, templates["central-cases.gotmpl"])))

	mux.Handle("/teams/work-in-progress/",
		wrap(
			teamWorkInProgress(client, templates["team-work-in-progress.gotmpl"])))

	mux.Handle("/users/pending-cases/",
		wrap(
			userPendingCases(client, templates["user-pending-cases.gotmpl"])))

	mux.Handle("/users/tasks/",
		wrap(
			userTasks(client, templates["user-tasks.gotmpl"])))

	mux.Handle("/users/all-cases/",
		wrap(
			userAllCases(client, templates["user-all-cases.gotmpl"])))

	mux.Handle("/reassign",
		wrap(
			reassign(client, templates["reassign.gotmpl"])))

	mux.Handle("/request-next-cases",
		wrap(
			requestNextCases(client)))

	mux.Handle("/request-next-task",
		wrap(
			requestNextTask(client)))

	mux.Handle("/mark-worked",
		wrap(
			markWorked(client)))

	mux.Handle("/feedback",
		wrap(
			feedback(client, templates["feedback.gotmpl"])))

	mux.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {})

	static := http.FileServer(http.Dir(webDir + "/static"))
	mux.Handle("/assets/", static)
	mux.Handle("/javascript/", static)
	mux.Handle("/stylesheets/", static)

	middleware := telemetry.Middleware(logger)

	return otelhttp.NewHandler(http.StripPrefix(prefix, securityheaders.Use(middleware(mux))), "lpa-dashboard")
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

func errorHandler(tmplError Template, prefix, siriusURL string) func(next Handler) http.Handler {
	return func(next Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := next(w, r)

			if err != nil {
				if errors.Is(err, context.Canceled) {
					w.WriteHeader(499)
					return
				}

				if err == sirius.ErrUnauthorized {
					http.Redirect(w, r, fmt.Sprintf("%s/auth?redirect=%s", siriusURL, url.QueryEscape(prefix+r.URL.Path)), http.StatusFound)
					return
				}

				if redirect, ok := err.(RedirectError); ok {
					http.Redirect(w, r, prefix+redirect.To(), http.StatusFound)
					return
				}

				code := http.StatusInternalServerError
				message := err.Error()
				if status, ok := err.(StatusError); ok {
					if status.Code() == http.StatusForbidden || status.Code() == http.StatusNotFound {
						code = status.Code()
					}
				}

				if statusError, ok := err.(*sirius.StatusError); ok {
					code = statusError.Code
					message = statusError.Title()
				}

				logger := telemetry.LoggerFromContext(r.Context())
				if code == http.StatusInternalServerError {
					logger.Error(err.Error())
				}

				w.WriteHeader(code)
				err = tmplError.ExecuteTemplate(w, "page", errorVars{
					SiriusURL: siriusURL,
					Path:      "",
					Code:      code,
					Error:     message,
				})

				if err != nil {
					logger.Error("could not generate error template", slog.Any("err", err.Error()))
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

package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type FeedbackClient interface {
	Feedback(sirius.Context, string) error
}

type feedbackVars struct {
	XSRFToken string
	Redirect  string
}

func feedback(client FeedbackClient, tmpl Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		switch r.Method {
		case http.MethodGet:
			ctx := getContext(r)

			return tmpl.ExecuteTemplate(w, "page", feedbackVars{
				XSRFToken: ctx.XSRFToken,
				Redirect:  r.Header.Get("Referer"),
			})

		case http.MethodPost:
			ctx := getContext(r)

			if err := client.Feedback(ctx, r.FormValue("feedback")); err != nil {
				return err
			}

			return RedirectError(r.FormValue("redirect"))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}

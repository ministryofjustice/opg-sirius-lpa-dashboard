package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type CardPaymentsClient interface {
	TasksByAssignee(sirius.Context, int, sirius.Criteria) ([]sirius.Task, *sirius.Pagination, error)
	MyDetails(sirius.Context) (sirius.MyDetails, error)
}

type cardPaymentsVars struct {
	Tasks     []sirius.Task
	XSRFToken string
}

func cardPayments(client CardPaymentsClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		myDetails, err := client.MyDetails(ctx)
		if err != nil {
			return err
		}

		criteria := sirius.Criteria{}.
			Filter("status", "Not started").
			Sort("dueDate", sirius.Ascending).
			Sort("name", sirius.Descending)

		tasks, _, err := client.TasksByAssignee(ctx, myDetails.ID, criteria)
		if err != nil {
			return err
		}

		vars := cardPaymentsVars{
			Tasks:     tasks,
			XSRFToken: ctx.XSRFToken,
		}

		return tmpl(w, vars)
	}
}

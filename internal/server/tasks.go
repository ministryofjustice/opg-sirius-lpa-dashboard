package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type TasksClient interface {
	CasesWithOpenTasksByAssignee(sirius.Context, int, int) ([]sirius.Case, *sirius.Pagination, error)
	HasWorkableCase(sirius.Context, int) (bool, error)
	MyDetails(sirius.Context) (sirius.MyDetails, error)
}

type tasksVars struct {
	Cases           []sirius.Case
	Pagination      *Pagination
	HasWorkableCase bool
	XSRFToken       string
}

func tasks(client TasksClient, tmpl Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		myDetails, err := client.MyDetails(ctx)
		if err != nil {
			return err
		}

		cases, pagination, err := client.CasesWithOpenTasksByAssignee(ctx, myDetails.ID, getPage(r))
		if err != nil {
			return err
		}

		hasWorkableCase, err := client.HasWorkableCase(ctx, myDetails.ID)
		if err != nil {
			return err
		}

		return tmpl.ExecuteTemplate(w, "page", tasksVars{
			Cases:           cases,
			Pagination:      newPagination(pagination),
			HasWorkableCase: hasWorkableCase,
			XSRFToken:       ctx.XSRFToken,
		})
	}
}

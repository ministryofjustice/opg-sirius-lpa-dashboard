package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type DashboardClient interface {
	CasesByAssignee(sirius.Context, int, string, int) ([]sirius.Case, *sirius.Pagination, error)
	MyDetails(sirius.Context) (sirius.MyDetails, error)
}

type dashboardVars struct {
	Cases           []sirius.Case
	Pagination      *sirius.Pagination
	HasWorkableCase bool
	XSRFToken       string
}

func dashboard(client DashboardClient, tmpl Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		myDetails, err := client.MyDetails(ctx)
		if err != nil {
			return err
		}

		myCases, pagination, err := client.CasesByAssignee(ctx, myDetails.ID, "Pending", getPage(r))
		if err != nil {
			return err
		}

		vars := dashboardVars{
			Cases:           myCases,
			Pagination:      pagination,
			HasWorkableCase: pagination.TotalItems > 0,
			XSRFToken:       ctx.XSRFToken,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

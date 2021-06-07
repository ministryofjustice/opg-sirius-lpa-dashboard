package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type PendingCasesClient interface {
	CasesByAssignee(sirius.Context, int, sirius.Criteria) ([]sirius.Case, *sirius.Pagination, error)
	HasWorkableCase(sirius.Context, int) (bool, error)
	MyDetails(sirius.Context) (sirius.MyDetails, error)
}

type pendingCasesVars struct {
	Cases           []sirius.Case
	Pagination      *Pagination
	HasWorkableCase bool
	XSRFToken       string
}

func pendingCases(client PendingCasesClient, tmpl Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		myDetails, err := client.MyDetails(ctx)
		if err != nil {
			return err
		}

		criteria := sirius.Criteria{}.Filter("status", "Pending").Page(getPage(r)).Sort("receiptDate", sirius.Ascending)
		myCases, pagination, err := client.CasesByAssignee(ctx, myDetails.ID, criteria)

		if err != nil {
			return err
		}

		hasWorkableCase, err := client.HasWorkableCase(ctx, myDetails.ID)
		if err != nil {
			return err
		}

		vars := pendingCasesVars{
			Cases:           myCases,
			Pagination:      newPagination(pagination),
			HasWorkableCase: hasWorkableCase,
			XSRFToken:       ctx.XSRFToken,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

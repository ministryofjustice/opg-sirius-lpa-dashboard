package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type AllCasesClient interface {
	CasesByAssignee(sirius.Context, int, sirius.Criteria) ([]sirius.Case, *sirius.Pagination, error)
	HasWorkableCase(sirius.Context, int) (bool, error)
	MyDetails(sirius.Context) (sirius.MyDetails, error)
}

type allCasesVars struct {
	Cases           []sirius.Case
	Pagination      *Pagination
	HasWorkableCase bool
	CanRequestCase  bool
	XSRFToken       string
}

func allCases(client AllCasesClient, tmpl Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		myDetails, err := client.MyDetails(ctx)
		if err != nil {
			return err
		}

		myCases, pagination, err := client.CasesByAssignee(ctx, myDetails.ID, sirius.Criteria{}.Page(getPage(r)).Sort("receiptDate", sirius.Ascending))
		if err != nil {
			return err
		}

		hasWorkableCase, err := client.HasWorkableCase(ctx, myDetails.ID)
		if err != nil {
			return err
		}

		vars := allCasesVars{
			Cases:           myCases,
			Pagination:      newPagination(pagination),
			HasWorkableCase: hasWorkableCase,
			CanRequestCase:  myDetails.HasRole("Self Allocation User"),
			XSRFToken:       ctx.XSRFToken,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type CasesClient interface {
	CasesByAssignee(sirius.Context, int, string, int) ([]sirius.Case, *sirius.Pagination, error)
	MyDetails(sirius.Context) (sirius.MyDetails, error)
}

type allCasesVars struct {
	Cases      []sirius.Case
	Pagination *sirius.Pagination
}

func cases(client CasesClient, tmpl Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		myDetails, err := client.MyDetails(ctx)
		if err != nil {
			return err
		}

		myCases, pagination, err := client.CasesByAssignee(ctx, myDetails.ID, "", getPage(r))
		if err != nil {
			return err
		}

		vars := allCasesVars{
			Cases:      myCases,
			Pagination: pagination,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

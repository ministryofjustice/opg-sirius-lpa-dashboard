package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type UserPendingCasesClient interface {
	CasesByAssignee(sirius.Context, int, sirius.Criteria) ([]sirius.Case, *sirius.Pagination, error)
	MyDetails(sirius.Context) (sirius.MyDetails, error)
	User(sirius.Context, int) (sirius.Assignee, error)
}

type userPendingCasesVars struct {
	Assignee   sirius.Assignee
	Cases      []sirius.Case
	Pagination *Pagination
	XSRFToken  string
}

func userPendingCases(client UserPendingCasesClient, tmpl Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		myDetails, err := client.MyDetails(ctx)
		if err != nil {
			return err
		}

		if !myDetails.IsManager() {
			return StatusError(http.StatusForbidden)
		}

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/users/pending-cases/"))
		if err != nil {
			return StatusError(http.StatusNotFound)
		}

		assignee, err := client.User(ctx, id)

		if err != nil {
			return err
		}

		criteria := sirius.Criteria{}.Filter("status", "Pending").Page(getPage(r)).Sort("receiptDate", sirius.Ascending)
		myCases, pagination, err := client.CasesByAssignee(ctx, id, criteria)

		if err != nil {
			return err
		}

		vars := userPendingCasesVars{
			Assignee:   assignee,
			Cases:      myCases,
			Pagination: newPagination(pagination),
			XSRFToken:  ctx.XSRFToken,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

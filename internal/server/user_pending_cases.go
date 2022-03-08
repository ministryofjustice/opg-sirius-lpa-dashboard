package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type UserPendingCasesClient interface {
	CasesByAssignee(sirius.Context, int, sirius.Criteria) ([]sirius.Case, *sirius.Pagination, error)
	MyDetails(sirius.Context) (sirius.MyDetails, error)
	User(sirius.Context, int) (sirius.Assignee, error)
}

type userPendingCasesVars struct {
	Assignee   sirius.Assignee
	Team       sirius.Team
	Cases      []sirius.Case
	Pagination *Pagination
	XSRFToken  string
}

func userPendingCases(client UserPendingCasesClient, tmpl template.Template) Handler {
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
		cases, pagination, err := client.CasesByAssignee(ctx, id, criteria)

		if err != nil {
			return err
		}

		var team sirius.Team
		if len(assignee.Teams) > 0 {
			team = assignee.Teams[0]
		}

		vars := userPendingCasesVars{
			Assignee:   assignee,
			Team:       team,
			Cases:      cases,
			Pagination: newPagination(pagination),
			XSRFToken:  ctx.XSRFToken,
		}

		return tmpl(w, vars)
	}
}

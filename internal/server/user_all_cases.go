package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type UserAllCasesClient interface {
	CasesByAssignee(sirius.Context, int, sirius.Criteria) ([]sirius.Case, *sirius.Pagination, error)
	MyDetails(sirius.Context) (sirius.MyDetails, error)
	User(sirius.Context, int) (sirius.Assignee, error)
}

type userAllCasesVars struct {
	Assignee   sirius.Assignee
	Team       sirius.Team
	Cases      []sirius.Case
	Pagination *Pagination
	XSRFToken  string
}

func userAllCases(client UserAllCasesClient, tmpl template.Template) Handler {
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

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/users/all-cases/"))
		if err != nil {
			return StatusError(http.StatusNotFound)
		}

		assignee, err := client.User(ctx, id)

		if err != nil {
			return err
		}

		cases, pagination, err := client.CasesByAssignee(ctx, id, sirius.Criteria{}.Page(getPage(r)).Sort("receiptDate", sirius.Ascending))

		if err != nil {
			return err
		}

		var team sirius.Team
		if len(assignee.Teams) > 0 {
			team = assignee.Teams[0]
		}

		vars := userAllCasesVars{
			Assignee:   assignee,
			Team:       team,
			Cases:      cases,
			Pagination: newPagination(pagination),
			XSRFToken:  ctx.XSRFToken,
		}

		return tmpl(w, vars)
	}
}

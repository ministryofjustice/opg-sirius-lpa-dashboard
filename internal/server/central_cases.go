package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type CentralCasesClient interface {
	CasesByAssignee(sirius.Context, int, sirius.Criteria) ([]sirius.Case, *sirius.Pagination, error)
	MyDetails(sirius.Context) (sirius.MyDetails, error)
	UserByEmail(sirius.Context, string) (sirius.User, error)
}

type centralCasesVars struct {
	Cases          []sirius.Case
	OldestCaseDate sirius.SiriusDate
	Pagination     *Pagination
	TeamID         int
	TeamName       string
	IsCaseWorker   bool
}

func centralCases(client CentralCasesClient, tmpl Template) Handler {
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

		centralPotUser, err := client.UserByEmail(ctx, sirius.PotUserEmail)
		if err != nil {
			return err
		}

		criteria := sirius.Criteria{}.Filter("status", "Pending").Page(getPage(r)).Sort("receiptDate", sirius.Ascending)
		teamCases, pagination, err := client.CasesByAssignee(ctx, centralPotUser.ID, criteria)

		if err != nil {
			return err
		}

		criteria = sirius.Criteria{}.Filter("status", "Pending").Sort("receiptDate", sirius.Ascending).Limit(1).Page(1)
		oldestCases, _, err := client.CasesByAssignee(ctx, centralPotUser.ID, criteria)

		if err != nil {
			return err
		}

		vars := centralCasesVars{
			Cases:        teamCases,
			Pagination:   newPagination(pagination),
			IsCaseWorker: myDetails.HasRole("Self Allocation User"),
		}

		if len(oldestCases) > 0 {
			vars.OldestCaseDate = oldestCases[0].ReceiptDate
		}

		if len(myDetails.Teams) > 0 {
			vars.TeamID = myDetails.Teams[0].ID
			vars.TeamName = myDetails.Teams[0].DisplayName
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

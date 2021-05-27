package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type CentralCasesClient interface {
	CasesByAssignee(sirius.Context, int, sirius.CasesByAssigneeCriteria) ([]sirius.Case, *sirius.Pagination, error)
	MyDetails(sirius.Context) (sirius.MyDetails, error)
	UserByEmail(sirius.Context, string) (sirius.User, error)
}

type centralCasesVars struct {
	Cases          []sirius.Case
	OldestCaseDate sirius.SiriusDate
	Pagination     *sirius.Pagination
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
			return StatusError(http.StatusUnauthorized)
		}

		centralPotUser, err := client.UserByEmail(ctx, "manager@opgtest.com")
		if err != nil {
			return err
		}

		teamCases, pagination, err := client.CasesByAssignee(ctx, centralPotUser.ID, sirius.CasesByAssigneeCriteria{
			Filter: sirius.CasesByAssigneeFilter{
				Status: "Pending",
			},
			Page: getPage(r),
		})

		if err != nil {
			return err
		}

		oldestCases, _, err := client.CasesByAssignee(ctx, centralPotUser.ID, sirius.CasesByAssigneeCriteria{
			Filter: sirius.CasesByAssigneeFilter{
				Status: "Pending",
			},
			Sort: sirius.CasesByAssigneeSort{
				Field: "receiptDate",
				Order: sirius.Ascending,
			},
			Limit: 1,
		})

		if err != nil {
			return err
		}

		vars := centralCasesVars{
			Cases:      teamCases,
			Pagination: pagination,
		}

		if len(oldestCases) > 0 {
			vars.OldestCaseDate = oldestCases[0].ReceiptDate
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

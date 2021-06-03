package server

import (
	"net/http"
	"time"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type TeamWorkInProgressClient interface {
	CasesByTeam(sirius.Context, int, sirius.Criteria) (*sirius.CasesByTeam, error)
	MyDetails(sirius.Context) (sirius.MyDetails, error)
}

type teamWorkInProgressVars struct {
	Cases          []sirius.Case
	OldestCaseDate sirius.SiriusDate
	Pagination     *sirius.Pagination
	TeamName       string
	Today          time.Time
	Stats          sirius.CasesByTeamMetadata
}

func teamWorkInProgress(client TeamWorkInProgressClient, tmpl Template) Handler {
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

		if len(myDetails.Teams) == 0 {
			return StatusError(http.StatusBadRequest)
		}

		result, err := client.CasesByTeam(ctx, myDetails.Teams[0].ID, sirius.Criteria{}.Page(getPage(r)))
		if err != nil {
			return err
		}

		vars := teamWorkInProgressVars{
			Cases:      result.Cases,
			Stats:      result.Stats,
			Pagination: result.Pagination,
			TeamName:   myDetails.Teams[0].DisplayName,
			Today:      time.Now(),
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

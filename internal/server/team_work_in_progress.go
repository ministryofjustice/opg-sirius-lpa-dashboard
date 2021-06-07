package server

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type TeamWorkInProgressClient interface {
	CasesByTeam(sirius.Context, int, sirius.Criteria) (*sirius.CasesByTeam, error)
	MyDetails(sirius.Context) (sirius.MyDetails, error)
	Teams(sirius.Context) ([]sirius.Team, error)
}

type teamWorkInProgressVars struct {
	Cases          []sirius.Case
	OldestCaseDate sirius.SiriusDate
	Pagination     *sirius.Pagination
	Today          time.Time
	Stats          sirius.CasesByTeamMetadata
	Team           sirius.Team
	Teams          []sirius.Team
	Filters        teamWorkInProgressFilters
}

type teamWorkInProgressFilters struct {
	Set        bool
	Allocation []int
	Status     []string
	DateFrom   time.Time
	DateTo     time.Time
	LpaType    string
}

func newTeamWorkInProgressFilters(form url.Values) teamWorkInProgressFilters {
	filters := teamWorkInProgressFilters{}

	if allocation, ok := form["allocation"]; ok {
		for _, v := range allocation {
			if i, err := strconv.Atoi(v); err == nil {
				filters.Allocation = append(filters.Allocation, i)
				filters.Set = true
			}
		}
	}

	if status, ok := form["status"]; ok {
		for _, v := range status {
			if v == "Pending" || v == "Pending, worked" {
				filters.Status = append(filters.Status, v)
				filters.Set = true
			}
		}
	}

	if v, err := time.Parse("2006-01-02", form.Get("date-from")); err == nil {
		filters.DateFrom = v
		filters.Set = true
	}

	if v, err := time.Parse("2006-01-02", form.Get("date-to")); err == nil {
		filters.DateTo = v
		filters.Set = true
	}

	if v := form.Get("lpa-type"); v == "pfa" || v == "hw" || v == "both" {
		filters.LpaType = v
		filters.Set = true
	}

	return filters
}

func teamWorkInProgress(client TeamWorkInProgressClient, tmpl Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/teams/work-in-progress/"))
		if err != nil {
			return StatusError(http.StatusNotFound)
		}

		ctx := getContext(r)

		myDetails, err := client.MyDetails(ctx)
		if err != nil {
			return err
		}

		if !myDetails.IsManager() {
			return StatusError(http.StatusForbidden)
		}

		teams, err := client.Teams(ctx)
		if err != nil {
			return err
		}

		currentTeam, ok := findTeam(id, teams)
		if !ok {
			return StatusError(http.StatusNotFound)
		}

		result, err := client.CasesByTeam(ctx, id, sirius.Criteria{}.Page(getPage(r)))
		if err != nil {
			return err
		}

		vars := teamWorkInProgressVars{
			Cases:      result.Cases,
			Stats:      result.Stats,
			Pagination: result.Pagination,
			Today:      time.Now(),
			Team:       currentTeam,
			Teams:      teams,
			Filters:    newTeamWorkInProgressFilters(r.Form),
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

func findTeam(id int, teams []sirius.Team) (sirius.Team, bool) {
	for _, team := range teams {
		if id == team.ID {
			return team, true
		}
	}

	return sirius.Team{}, false
}

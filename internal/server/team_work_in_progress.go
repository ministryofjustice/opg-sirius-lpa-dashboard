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
	Pagination     *Pagination
	Today          time.Time
	Stats          sirius.CasesByTeamMetadata
	Team           sirius.Team
	Teams          []sirius.Team
	Filters        teamWorkInProgressFilters
	IsCaseWorker   bool
}

type teamWorkInProgressFilters struct {
	Set        bool
	Allocation []int
	Status     []string
	DateFrom   time.Time
	DateTo     time.Time
	LpaType    string
}

func (f teamWorkInProgressFilters) Encode() string {
	if !f.Set {
		return ""
	}

	form := url.Values{}
	for _, v := range f.Allocation {
		form.Add("allocation", strconv.Itoa(v))
	}
	for _, v := range f.Status {
		form.Add("status", v)
	}
	if !f.DateFrom.IsZero() {
		form.Add("date-from", f.DateFrom.Format("2006-01-02"))
	}
	if !f.DateTo.IsZero() {
		form.Add("date-to", f.DateTo.Format("2006-01-02"))
	}
	if f.LpaType != "" {
		form.Add("lpa-type", f.LpaType)
	}

	return form.Encode()
}

func (f teamWorkInProgressFilters) Criteria() sirius.Criteria {
	c := sirius.Criteria{}
	if !f.Set {
		return c
	}

	for _, v := range f.Allocation {
		c = c.Filter("allocation", strconv.Itoa(v))
	}
	for _, v := range f.Status {
		c = c.Filter("status", v)
	}
	if !f.DateFrom.IsZero() {
		c = c.Filter("date-from", f.DateFrom.Format("2006-01-02"))
	}
	if !f.DateTo.IsZero() {
		c = c.Filter("date-to", f.DateTo.Format("2006-01-02"))
	}
	if f.LpaType != "" {
		c = c.Filter("lpa-type", f.LpaType)
	}

	return c
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
			if v == "pending" || v == "pending-worked" {
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

		var caseworkTeams []sirius.Team
		for _, team := range teams {
			if strings.HasPrefix(team.DisplayName, "Casework Team") {
				caseworkTeams = append(caseworkTeams, team)
			}
		}

		currentTeam, ok := findTeam(id, teams)
		if !ok {
			return StatusError(http.StatusNotFound)
		}

		page := getPage(r)
		filters := newTeamWorkInProgressFilters(r.Form)

		result, err := client.CasesByTeam(ctx, id, filters.Criteria().Page(page))
		if err != nil {
			return err
		}

		vars := teamWorkInProgressVars{
			Cases:        result.Cases,
			Stats:        result.Stats,
			Pagination:   newPaginationWithQuery(result.Pagination, filters.Encode()),
			Today:        time.Now(),
			Team:         currentTeam,
			Teams:        caseworkTeams,
			Filters:      filters,
			IsCaseWorker: myDetails.HasRole("Self Allocation User"),
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

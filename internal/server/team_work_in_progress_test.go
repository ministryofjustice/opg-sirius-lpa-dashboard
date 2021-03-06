package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockTeamWorkInProgressClient struct {
	casesByTeam struct {
		count        int
		lastCtx      sirius.Context
		lastId       int
		lastCriteria sirius.Criteria
		data         *sirius.CasesByTeam
		err          error
	}
	myDetails struct {
		count   int
		lastCtx sirius.Context
		data    sirius.MyDetails
		err     error
	}
	teams struct {
		count   int
		lastCtx sirius.Context
		data    []sirius.Team
		err     error
	}
}

func (m *mockTeamWorkInProgressClient) CasesByTeam(ctx sirius.Context, id int, criteria sirius.Criteria) (*sirius.CasesByTeam, error) {
	m.casesByTeam.count += 1
	m.casesByTeam.lastCtx = ctx
	m.casesByTeam.lastId = id
	m.casesByTeam.lastCriteria = criteria

	return m.casesByTeam.data, m.casesByTeam.err
}

func (m *mockTeamWorkInProgressClient) MyDetails(ctx sirius.Context) (sirius.MyDetails, error) {
	m.myDetails.count += 1
	m.myDetails.lastCtx = ctx

	return m.myDetails.data, m.myDetails.err
}

func (m *mockTeamWorkInProgressClient) Teams(ctx sirius.Context) ([]sirius.Team, error) {
	m.teams.count += 1
	m.teams.lastCtx = ctx

	return m.teams.data, m.teams.err
}

func TestGetTeamWorkInProgress(t *testing.T) {
	assert := assert.New(t)

	client := &mockTeamWorkInProgressClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Manager"},
		Teams: []sirius.MyDetailsTeam{{ID: 123, DisplayName: "team"}},
	}
	client.casesByTeam.data = &sirius.CasesByTeam{
		Cases: []sirius.Case{{
			ID: 78,
			Donor: sirius.Donor{
				ID: 79,
			},
		}},
		Pagination: &sirius.Pagination{
			TotalItems: 1,
		},
		Stats: sirius.CasesByTeamMetadata{
			WorkedTotal: 1,
		},
	}
	client.teams.data = []sirius.Team{
		{
			ID:          1,
			DisplayName: "my team",
		},
		{
			ID:          2,
			DisplayName: "Casework Team 2",
		},
		{
			ID:          3,
			DisplayName: "Nottingham casework team 2",
		},
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/teams/work-in-progress/1", nil)

	err := teamWorkInProgress(client, template)(w, r)
	assert.Nil(err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.teams.count)
	assert.Equal(getContext(r), client.teams.lastCtx)

	assert.Equal(1, client.casesByTeam.count)
	assert.Equal(getContext(r), client.casesByTeam.lastCtx)
	assert.Equal(1, client.casesByTeam.lastId)
	assert.Equal(sirius.Criteria{}.Page(1), client.casesByTeam.lastCriteria)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)

	vars := template.lastVars.(teamWorkInProgressVars)
	assert.WithinDuration(time.Now(), vars.Today, time.Second)
	vars.Today = time.Time{}

	assert.Equal(teamWorkInProgressVars{
		Cases:      client.casesByTeam.data.Cases,
		Team:       client.teams.data[0],
		Pagination: newPaginationWithQuery(client.casesByTeam.data.Pagination, ""),
		Stats:      client.casesByTeam.data.Stats,
		Teams:      []sirius.Team{client.teams.data[1], client.teams.data[2]},
	}, vars)
}

func TestGetTeamWorkInProgressBadPath(t *testing.T) {
	testCases := map[string]string{
		"not a number": "/teams/work-in-progress/what",
		"no value":     "/teams/work-in-progress/",
	}

	for name, url := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", url, nil)

			err := teamWorkInProgress(nil, nil)(w, r)
			assert.Equal(StatusError(http.StatusNotFound), err)
		})
	}
}

func TestGetTeamWorkInProgressPage(t *testing.T) {
	assert := assert.New(t)

	client := &mockTeamWorkInProgressClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Case Manager", "Manager", "System Admin"},
		Teams: []sirius.MyDetailsTeam{{ID: 123, DisplayName: "team"}},
	}
	client.casesByTeam.data = &sirius.CasesByTeam{
		Cases: []sirius.Case{{
			ID: 78,
			Donor: sirius.Donor{
				ID: 79,
			},
		}},
	}
	client.teams.data = []sirius.Team{{
		ID:          1,
		DisplayName: "Casework Team 1",
	}}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/teams/work-in-progress/1?page=4", nil)

	err := teamWorkInProgress(client, template)(w, r)
	assert.Nil(err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.teams.count)
	assert.Equal(getContext(r), client.teams.lastCtx)

	assert.Equal(1, client.casesByTeam.count)
	assert.Equal(getContext(r), client.casesByTeam.lastCtx)
	assert.Equal(1, client.casesByTeam.lastId)
	assert.Equal(sirius.Criteria{}.Page(4), client.casesByTeam.lastCriteria)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)

	vars := template.lastVars.(teamWorkInProgressVars)
	assert.WithinDuration(time.Now(), vars.Today, time.Second)
	vars.Today = time.Time{}

	assert.Equal(teamWorkInProgressVars{
		Cases: client.casesByTeam.data.Cases,
		Team:  client.teams.data[0],
		Teams: client.teams.data,
	}, vars)
}

func TestGetTeamWorkInProgressFiltered(t *testing.T) {
	assert := assert.New(t)

	client := &mockTeamWorkInProgressClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Case Manager", "Manager", "System Admin"},
		Teams: []sirius.MyDetailsTeam{{ID: 123, DisplayName: "team"}},
	}
	client.casesByTeam.data = &sirius.CasesByTeam{
		Cases: []sirius.Case{{
			ID: 78,
			Donor: sirius.Donor{
				ID: 79,
			},
		}},
	}
	client.teams.data = []sirius.Team{{
		ID:          1,
		DisplayName: "Casework Team 1",
	}}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/teams/work-in-progress/1?allocation=123", nil)

	err := teamWorkInProgress(client, template)(w, r)
	assert.Nil(err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.teams.count)
	assert.Equal(getContext(r), client.teams.lastCtx)

	assert.Equal(1, client.casesByTeam.count)
	assert.Equal(getContext(r), client.casesByTeam.lastCtx)
	assert.Equal(1, client.casesByTeam.lastId)
	assert.Equal(sirius.Criteria{}.Filter("allocation", "123").Page(1), client.casesByTeam.lastCriteria)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)

	vars := template.lastVars.(teamWorkInProgressVars)
	assert.WithinDuration(time.Now(), vars.Today, time.Second)
	vars.Today = time.Time{}

	assert.Equal(teamWorkInProgressVars{
		Cases: client.casesByTeam.data.Cases,
		Team:  client.teams.data[0],
		Teams: client.teams.data,
		Filters: teamWorkInProgressFilters{
			Set:        true,
			Allocation: []int{123},
		},
	}, vars)
}

func TestGetTeamWorkInProgressForbidden(t *testing.T) {
	assert := assert.New(t)

	client := &mockTeamWorkInProgressClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Case Manager", "System Admin"},
		Teams: []sirius.MyDetailsTeam{{DisplayName: "team"}},
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/teams/work-in-progress/1", nil)

	err := teamWorkInProgress(client, template)(w, r)
	assert.Equal(StatusError(http.StatusForbidden), err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(0, client.casesByTeam.count)
}

func TestGetTeamWorkInProgressTeamDoesNotExist(t *testing.T) {
	assert := assert.New(t)

	client := &mockTeamWorkInProgressClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Manager"},
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/teams/work-in-progress/12", nil)

	err := teamWorkInProgress(client, template)(w, r)
	assert.Equal(StatusError(http.StatusNotFound), err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.teams.count)
	assert.Equal(getContext(r), client.teams.lastCtx)

	assert.Equal(0, client.casesByTeam.count)
}

func TestGetTeamWorkInProgressMyDetailsError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockTeamWorkInProgressClient{}
	client.myDetails.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/teams/work-in-progress/1", nil)

	err := teamWorkInProgress(client, template)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(0, client.teams.count)
	assert.Equal(0, client.casesByTeam.count)
}

func TestGetTeamWorkInProgressTeamsError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockTeamWorkInProgressClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Manager"},
	}
	client.teams.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/teams/work-in-progress/1", nil)

	err := teamWorkInProgress(client, template)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.teams.count)
	assert.Equal(getContext(r), client.teams.lastCtx)

	assert.Equal(0, client.casesByTeam.count)
}

func TestGetTeamWorkInProgressQueryError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockTeamWorkInProgressClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Manager"},
		Teams: []sirius.MyDetailsTeam{{ID: 123, DisplayName: "team"}},
	}
	client.teams.data = []sirius.Team{{
		ID:          1,
		DisplayName: "my team",
	}}
	client.casesByTeam.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/teams/work-in-progress/1", nil)

	err := teamWorkInProgress(client, template)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.teams.count)
	assert.Equal(getContext(r), client.teams.lastCtx)

	assert.Equal(1, client.casesByTeam.count)
	assert.Equal(getContext(r), client.casesByTeam.lastCtx)
	assert.Equal(1, client.casesByTeam.lastId)
	assert.Equal(sirius.Criteria{}.Page(1), client.casesByTeam.lastCriteria)
}

func TestBadMethodTeamWorkInProgress(t *testing.T) {
	assert := assert.New(t)

	client := &mockTeamWorkInProgressClient{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/path", nil)

	err := teamWorkInProgress(client, template)(w, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)

	assert.Equal(0, client.casesByTeam.count)
	assert.Equal(0, template.count)
}

func TestTeamWorkInProgressFilters(t *testing.T) {
	testCases := map[string]struct {
		Input    string
		Encoded  string
		Criteria sirius.Criteria
	}{
		"empty": {
			Input:    "",
			Encoded:  "",
			Criteria: sirius.Criteria{},
		},
		"all": {
			Input:   "allocation=123&allocation=456&date-from=2021-01-02&date-to=2021-01-03&lpa-type=both&status=pending&status=pending-worked",
			Encoded: "allocation=123&allocation=456&date-from=2021-01-02&date-to=2021-01-03&lpa-type=both&status=pending&status=pending-worked",
			Criteria: sirius.Criteria{}.
				Filter("allocation", "123").
				Filter("allocation", "456").
				Filter("status", "pending").
				Filter("status", "pending-worked").
				Filter("date-from", "2021-01-02").
				Filter("date-to", "2021-01-03").
				Filter("lpa-type", "both"),
		},
		"date-range-bad": {
			Input:    "date-from=what&date-to=huh",
			Encoded:  "",
			Criteria: sirius.Criteria{},
		},
		"lpa-type-hw": {
			Input:    "lpa-type=hw",
			Encoded:  "lpa-type=hw",
			Criteria: sirius.Criteria{}.Filter("lpa-type", "hw"),
		},
		"lpa-type-pfa": {
			Input:    "lpa-type=pfa",
			Encoded:  "lpa-type=pfa",
			Criteria: sirius.Criteria{}.Filter("lpa-type", "pfa"),
		},
		"lpa-type-unknown": {
			Input:    "lpa-type=what",
			Encoded:  "",
			Criteria: sirius.Criteria{},
		},
		"status-bad": {
			Input:    "status=what",
			Encoded:  "",
			Criteria: sirius.Criteria{},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			query, _ := url.ParseQuery(tc.Input)
			filters := newTeamWorkInProgressFilters(query)

			assert.Equal(t, tc.Encoded, filters.Encode())
			assert.Equal(t, tc.Criteria, filters.Criteria())
		})
	}
}

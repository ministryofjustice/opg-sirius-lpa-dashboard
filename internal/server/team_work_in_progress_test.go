package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
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
		data         []sirius.Case
		pagination   *sirius.Pagination
		err          error
	}
	myDetails struct {
		count   int
		lastCtx sirius.Context
		data    sirius.MyDetails
		err     error
	}
}

func (m *mockTeamWorkInProgressClient) CasesByTeam(ctx sirius.Context, id int, criteria sirius.Criteria) ([]sirius.Case, *sirius.Pagination, error) {
	m.casesByTeam.count += 1
	m.casesByTeam.lastCtx = ctx
	m.casesByTeam.lastId = id
	m.casesByTeam.lastCriteria = criteria

	return m.casesByTeam.data, m.casesByTeam.pagination, m.casesByTeam.err
}

func (m *mockTeamWorkInProgressClient) MyDetails(ctx sirius.Context) (sirius.MyDetails, error) {
	m.myDetails.count += 1
	m.myDetails.lastCtx = ctx

	return m.myDetails.data, m.myDetails.err
}

func TestGetTeamWorkInProgress(t *testing.T) {
	assert := assert.New(t)

	client := &mockTeamWorkInProgressClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Manager"},
		Teams: []sirius.MyDetailsTeam{{ID: 123, DisplayName: "team"}},
	}
	client.casesByTeam.data = []sirius.Case{{
		ID: 78,
		Donor: sirius.Donor{
			ID: 79,
		},
	}}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := teamWorkInProgress(client, template)(w, r)
	assert.Nil(err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.casesByTeam.count)
	assert.Equal(getContext(r), client.casesByTeam.lastCtx)
	assert.Equal(123, client.casesByTeam.lastId)
	assert.Equal(sirius.Criteria{}.Page(1), client.casesByTeam.lastCriteria)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)

	vars := template.lastVars.(teamWorkInProgressVars)
	assert.WithinDuration(time.Now(), vars.Today, time.Second)
	vars.Today = time.Time{}

	assert.Equal(teamWorkInProgressVars{
		Cases:    client.casesByTeam.data,
		TeamName: "team",
	}, vars)
}

func TestGetTeamWorkInProgressPage(t *testing.T) {
	assert := assert.New(t)

	client := &mockTeamWorkInProgressClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Case Manager", "Manager", "System Admin"},
		Teams: []sirius.MyDetailsTeam{{ID: 123, DisplayName: "team"}},
	}
	client.casesByTeam.data = []sirius.Case{{
		ID: 78,
		Donor: sirius.Donor{
			ID: 79,
		},
	}}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path?page=4", nil)

	err := teamWorkInProgress(client, template)(w, r)
	assert.Nil(err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.casesByTeam.count)
	assert.Equal(getContext(r), client.casesByTeam.lastCtx)
	assert.Equal(123, client.casesByTeam.lastId)
	assert.Equal(sirius.Criteria{}.Page(4), client.casesByTeam.lastCriteria)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)

	vars := template.lastVars.(teamWorkInProgressVars)
	assert.WithinDuration(time.Now(), vars.Today, time.Second)
	vars.Today = time.Time{}

	assert.Equal(teamWorkInProgressVars{
		Cases:    client.casesByTeam.data,
		TeamName: "team",
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
	r, _ := http.NewRequest("GET", "/path", nil)

	err := teamWorkInProgress(client, template)(w, r)
	assert.Equal(StatusError(http.StatusForbidden), err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(0, client.casesByTeam.count)
}

func TestGetTeamWorkInProgressNotInTeam(t *testing.T) {
	assert := assert.New(t)

	client := &mockTeamWorkInProgressClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Manager"},
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := teamWorkInProgress(client, template)(w, r)
	assert.Equal(StatusError(http.StatusBadRequest), err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(0, client.casesByTeam.count)
}

func TestGetTeamWorkInProgressMyDetailsError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockTeamWorkInProgressClient{}
	client.myDetails.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := teamWorkInProgress(client, template)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

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
	client.casesByTeam.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := teamWorkInProgress(client, template)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.casesByTeam.count)
	assert.Equal(getContext(r), client.casesByTeam.lastCtx)
	assert.Equal(123, client.casesByTeam.lastId)
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

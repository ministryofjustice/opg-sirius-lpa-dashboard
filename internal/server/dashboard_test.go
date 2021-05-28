package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockDashboardClient struct {
	myDetails struct {
		count   int
		lastCtx sirius.Context
		data    sirius.MyDetails
		err     error
	}
	casesByAssignee struct {
		count        int
		lastCtx      sirius.Context
		lastId       int
		lastCriteria sirius.Criteria
		data         []sirius.Case
		pagination   *sirius.Pagination
		err          error
	}
}

func (m *mockDashboardClient) MyDetails(ctx sirius.Context) (sirius.MyDetails, error) {
	m.myDetails.count += 1
	m.myDetails.lastCtx = ctx

	return m.myDetails.data, m.myDetails.err
}

func (m *mockDashboardClient) CasesByAssignee(ctx sirius.Context, id int, criteria sirius.Criteria) ([]sirius.Case, *sirius.Pagination, error) {
	m.casesByAssignee.count += 1
	m.casesByAssignee.lastCtx = ctx
	m.casesByAssignee.lastId = id
	m.casesByAssignee.lastCriteria = criteria

	return m.casesByAssignee.data, m.casesByAssignee.pagination, m.casesByAssignee.err
}

func TestGetDashboard(t *testing.T) {
	assert := assert.New(t)

	client := &mockDashboardClient{}
	client.myDetails.data = sirius.MyDetails{
		ID: 14,
	}
	client.casesByAssignee.data = []sirius.Case{{
		ID: 78,
		Donor: sirius.Donor{
			ID: 79,
		},
	}}
	client.casesByAssignee.pagination = &sirius.Pagination{
		TotalItems: 20,
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := dashboard(client, template)(w, r)
	assert.Nil(err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.casesByAssignee.count)
	assert.Equal(getContext(r), client.casesByAssignee.lastCtx)
	assert.Equal(14, client.casesByAssignee.lastId)
	assert.Equal(sirius.Criteria{}.Filter("status", "Pending").Page(1), client.casesByAssignee.lastCriteria)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(dashboardVars{
		Cases:           client.casesByAssignee.data,
		Pagination:      client.casesByAssignee.pagination,
		HasWorkableCase: true,
	}, template.lastVars)
}

func TestGetDashboardPage(t *testing.T) {
	assert := assert.New(t)

	client := &mockDashboardClient{}
	client.myDetails.data = sirius.MyDetails{
		ID: 14,
	}
	client.casesByAssignee.data = []sirius.Case{{
		ID: 78,
		Donor: sirius.Donor{
			ID: 79,
		},
	}}
	client.casesByAssignee.pagination = &sirius.Pagination{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path?page=4", nil)

	err := dashboard(client, template)(w, r)
	assert.Nil(err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.casesByAssignee.count)
	assert.Equal(getContext(r), client.casesByAssignee.lastCtx)
	assert.Equal(14, client.casesByAssignee.lastId)
	assert.Equal(sirius.Criteria{}.Filter("status", "Pending").Page(4), client.casesByAssignee.lastCriteria)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(dashboardVars{
		Cases:      client.casesByAssignee.data,
		Pagination: client.casesByAssignee.pagination,
	}, template.lastVars)
}

func TestGetDashboardMyDetailsError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockDashboardClient{}
	client.myDetails.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := dashboard(client, template)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(0, client.casesByAssignee.count)
}

func TestGetDashboardQueryError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockDashboardClient{}
	client.myDetails.data = sirius.MyDetails{
		ID: 14,
	}
	client.casesByAssignee.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := dashboard(client, template)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.casesByAssignee.count)
	assert.Equal(getContext(r), client.casesByAssignee.lastCtx)
	assert.Equal(14, client.casesByAssignee.lastId)
	assert.Equal(sirius.Criteria{}.Filter("status", "Pending").Page(1), client.casesByAssignee.lastCriteria)
}

func TestBadMethodDashboard(t *testing.T) {
	assert := assert.New(t)

	client := &mockDashboardClient{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/path", nil)

	err := dashboard(client, template)(w, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)

	assert.Equal(0, client.myDetails.count)
	assert.Equal(0, client.casesByAssignee.count)
	assert.Equal(0, template.count)
}

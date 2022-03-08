package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockUserAllCasesClient struct {
	myDetails struct {
		count   int
		lastCtx sirius.Context
		data    sirius.MyDetails
		err     error
	}
	user struct {
		count   int
		lastCtx sirius.Context
		lastId  int
		data    sirius.Assignee
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

func (m *mockUserAllCasesClient) MyDetails(ctx sirius.Context) (sirius.MyDetails, error) {
	m.myDetails.count += 1
	m.myDetails.lastCtx = ctx

	return m.myDetails.data, m.myDetails.err
}

func (m *mockUserAllCasesClient) User(ctx sirius.Context, id int) (sirius.Assignee, error) {
	m.user.count += 1
	m.user.lastCtx = ctx
	m.user.lastId = id

	return m.user.data, m.user.err
}

func (m *mockUserAllCasesClient) CasesByAssignee(ctx sirius.Context, id int, criteria sirius.Criteria) ([]sirius.Case, *sirius.Pagination, error) {
	m.casesByAssignee.count += 1
	m.casesByAssignee.lastCtx = ctx
	m.casesByAssignee.lastId = id
	m.casesByAssignee.lastCriteria = criteria

	return m.casesByAssignee.data, m.casesByAssignee.pagination, m.casesByAssignee.err
}

func TestGetUserAllCases(t *testing.T) {
	assert := assert.New(t)

	client := &mockUserAllCasesClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Manager"},
	}
	client.user.data = sirius.Assignee{
		ID:          74,
		DisplayName: "Elfriede Giesing",
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
	r, _ := http.NewRequest("GET", "/users/all-cases/74", nil)

	err := userAllCases(client, template.Func)(w, r)
	assert.Nil(err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.user.count)
	assert.Equal(getContext(r), client.user.lastCtx)
	assert.Equal(74, client.user.lastId)

	assert.Equal(1, client.casesByAssignee.count)
	assert.Equal(getContext(r), client.casesByAssignee.lastCtx)
	assert.Equal(74, client.casesByAssignee.lastId)
	assert.Equal(sirius.Criteria{}.Page(1).Sort("receiptDate", sirius.Ascending), client.casesByAssignee.lastCriteria)

	assert.Equal(1, template.count)
	assert.Equal(userAllCasesVars{
		Assignee:   client.user.data,
		Cases:      client.casesByAssignee.data,
		Pagination: newPagination(client.casesByAssignee.pagination),
	}, template.lastVars)
}

func TestGetUserAllCasesPage(t *testing.T) {
	assert := assert.New(t)

	client := &mockUserAllCasesClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Manager"},
	}
	client.user.data = sirius.Assignee{
		ID:          74,
		DisplayName: "Elfriede Giesing",
		Teams: []sirius.Team{{
			ID:          281,
			DisplayName: "Casework Team 6",
		}},
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
	r, _ := http.NewRequest("GET", "/users/all-cases/74?page=4", nil)

	err := userAllCases(client, template.Func)(w, r)
	assert.Nil(err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.user.count)
	assert.Equal(getContext(r), client.user.lastCtx)
	assert.Equal(74, client.user.lastId)

	assert.Equal(1, client.casesByAssignee.count)
	assert.Equal(getContext(r), client.casesByAssignee.lastCtx)
	assert.Equal(74, client.casesByAssignee.lastId)
	assert.Equal(sirius.Criteria{}.Page(4).Sort("receiptDate", sirius.Ascending), client.casesByAssignee.lastCriteria)

	assert.Equal(1, template.count)
	assert.Equal(userAllCasesVars{
		Assignee:   client.user.data,
		Team:       client.user.data.Teams[0],
		Cases:      client.casesByAssignee.data,
		Pagination: newPagination(client.casesByAssignee.pagination),
	}, template.lastVars)
}

func TestGetUserAllCasesForbidden(t *testing.T) {
	assert := assert.New(t)

	client := &mockUserAllCasesClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Case Worker"},
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users/all-cases/74", nil)

	err := userAllCases(client, template.Func)(w, r)
	assert.Equal(StatusError(http.StatusForbidden), err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(0, client.user.count)
	assert.Equal(0, client.casesByAssignee.count)
}

func TestGetUserAllCasesMyDetailsError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockUserAllCasesClient{}
	client.myDetails.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users/all-cases/74", nil)

	err := userAllCases(client, template.Func)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(0, client.user.count)
	assert.Equal(0, client.casesByAssignee.count)
}

func TestGetUserAllCasesGetUserError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockUserAllCasesClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Manager"},
	}
	client.user.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users/all-cases/74", nil)

	err := userAllCases(client, template.Func)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.user.count)
	assert.Equal(getContext(r), client.user.lastCtx)
	assert.Equal(74, client.user.lastId)

	assert.Equal(0, client.casesByAssignee.count)
}

func TestGetUserAllCasesQueryError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockUserAllCasesClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Manager"},
	}
	client.user.data = sirius.Assignee{
		ID:          74,
		DisplayName: "Elfriede Giesing",
	}
	client.casesByAssignee.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users/all-cases/74", nil)

	err := userAllCases(client, template.Func)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.user.count)
	assert.Equal(getContext(r), client.user.lastCtx)
	assert.Equal(74, client.user.lastId)

	assert.Equal(1, client.casesByAssignee.count)
	assert.Equal(getContext(r), client.casesByAssignee.lastCtx)
	assert.Equal(74, client.casesByAssignee.lastId)
	assert.Equal(sirius.Criteria{}.Page(1).Sort("receiptDate", sirius.Ascending), client.casesByAssignee.lastCriteria)
}

func TestBadMethodUserAllCases(t *testing.T) {
	assert := assert.New(t)

	client := &mockUserAllCasesClient{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/users/all-cases/74", nil)

	err := userAllCases(client, template.Func)(w, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)

	assert.Equal(0, client.user.count)
	assert.Equal(0, client.casesByAssignee.count)
	assert.Equal(0, template.count)
}

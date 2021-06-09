package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockUserPendingCasesClient struct {
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

func (m *mockUserPendingCasesClient) MyDetails(ctx sirius.Context) (sirius.MyDetails, error) {
	m.myDetails.count += 1
	m.myDetails.lastCtx = ctx

	return m.myDetails.data, m.myDetails.err
}

func (m *mockUserPendingCasesClient) User(ctx sirius.Context, id int) (sirius.Assignee, error) {
	m.user.count += 1
	m.user.lastCtx = ctx
	m.user.lastId = id

	return m.user.data, m.user.err
}

func (m *mockUserPendingCasesClient) CasesByAssignee(ctx sirius.Context, id int, criteria sirius.Criteria) ([]sirius.Case, *sirius.Pagination, error) {
	m.casesByAssignee.count += 1
	m.casesByAssignee.lastCtx = ctx
	m.casesByAssignee.lastId = id
	m.casesByAssignee.lastCriteria = criteria

	return m.casesByAssignee.data, m.casesByAssignee.pagination, m.casesByAssignee.err
}

func TestGetUserPendingCases(t *testing.T) {
	assert := assert.New(t)

	client := &mockUserPendingCasesClient{}
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
	r, _ := http.NewRequest("GET", "/users/pending-cases/74", nil)

	err := userPendingCases(client, template)(w, r)
	assert.Nil(err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.user.count)
	assert.Equal(getContext(r), client.user.lastCtx)
	assert.Equal(74, client.user.lastId)

	assert.Equal(1, client.casesByAssignee.count)
	assert.Equal(getContext(r), client.casesByAssignee.lastCtx)
	assert.Equal(74, client.casesByAssignee.lastId)
	assert.Equal(sirius.Criteria{}.Filter("status", "Pending").Page(1).Sort("receiptDate", sirius.Ascending), client.casesByAssignee.lastCriteria)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(userPendingCasesVars{
		Assignee:   client.user.data,
		Cases:      client.casesByAssignee.data,
		Pagination: newPagination(client.casesByAssignee.pagination),
	}, template.lastVars)
}

func TestGetUserPendingCasesPage(t *testing.T) {
	assert := assert.New(t)

	client := &mockUserPendingCasesClient{}
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
	client.casesByAssignee.pagination = &sirius.Pagination{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users/pending-cases/74?page=4", nil)

	err := userPendingCases(client, template)(w, r)
	assert.Nil(err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.user.count)
	assert.Equal(getContext(r), client.user.lastCtx)
	assert.Equal(74, client.user.lastId)

	assert.Equal(1, client.casesByAssignee.count)
	assert.Equal(getContext(r), client.casesByAssignee.lastCtx)
	assert.Equal(74, client.casesByAssignee.lastId)
	assert.Equal(sirius.Criteria{}.Filter("status", "Pending").Page(4).Sort("receiptDate", sirius.Ascending), client.casesByAssignee.lastCriteria)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(userPendingCasesVars{
		Assignee:   client.user.data,
		Cases:      client.casesByAssignee.data,
		Pagination: newPagination(client.casesByAssignee.pagination),
	}, template.lastVars)
}

func TestGetUserPendingCasesForbidden(t *testing.T) {
	assert := assert.New(t)

	client := &mockUserPendingCasesClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Case Worker"},
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users/pending-cases/74", nil)

	err := userPendingCases(client, template)(w, r)
	assert.Equal(StatusError(http.StatusForbidden), err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(0, client.user.count)
	assert.Equal(0, client.casesByAssignee.count)
}

func TestGetUserPendingCasesMyDetailsError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockUserPendingCasesClient{}
	client.myDetails.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users/pending-cases/74", nil)

	err := userPendingCases(client, template)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(0, client.user.count)
	assert.Equal(0, client.casesByAssignee.count)
}

func TestGetUserPendingCasesGetUserError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockUserPendingCasesClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Manager"},
	}
	client.user.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users/pending-cases/74", nil)

	err := userPendingCases(client, template)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.user.count)
	assert.Equal(getContext(r), client.user.lastCtx)
	assert.Equal(74, client.user.lastId)

	assert.Equal(0, client.casesByAssignee.count)
}

func TestGetUserPendingCasesQueryError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockUserPendingCasesClient{}
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
	r, _ := http.NewRequest("GET", "/users/pending-cases/74", nil)

	err := userPendingCases(client, template)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.user.count)
	assert.Equal(getContext(r), client.user.lastCtx)
	assert.Equal(74, client.user.lastId)

	assert.Equal(1, client.casesByAssignee.count)
	assert.Equal(getContext(r), client.casesByAssignee.lastCtx)
	assert.Equal(74, client.casesByAssignee.lastId)
	assert.Equal(sirius.Criteria{}.Filter("status", "Pending").Page(1).Sort("receiptDate", sirius.Ascending), client.casesByAssignee.lastCriteria)
}

func TestBadMethodUserPendingCases(t *testing.T) {
	assert := assert.New(t)

	client := &mockUserPendingCasesClient{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/users/pending-cases/74", nil)

	err := userPendingCases(client, template)(w, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)

	assert.Equal(0, client.user.count)
	assert.Equal(0, client.casesByAssignee.count)
	assert.Equal(0, template.count)
}

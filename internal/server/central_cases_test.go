package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockCentralCasesClient struct {
	casesByAssignee struct {
		count      int
		lastCtx    sirius.Context
		lastId     int
		criteria   []sirius.Criteria
		data       []sirius.Case
		pagination *sirius.Pagination
		err        error
	}
	myDetails struct {
		count   int
		lastCtx sirius.Context
		data    sirius.MyDetails
		err     error
	}
	userByEmail struct {
		count     int
		lastCtx   sirius.Context
		lastEmail string
		data      sirius.User
		err       error
	}
}

func (m *mockCentralCasesClient) CasesByAssignee(ctx sirius.Context, id int, criteria sirius.Criteria) ([]sirius.Case, *sirius.Pagination, error) {
	m.casesByAssignee.count += 1
	m.casesByAssignee.lastCtx = ctx
	m.casesByAssignee.lastId = id
	m.casesByAssignee.criteria = append(m.casesByAssignee.criteria, criteria)

	return m.casesByAssignee.data, m.casesByAssignee.pagination, m.casesByAssignee.err
}

func (m *mockCentralCasesClient) MyDetails(ctx sirius.Context) (sirius.MyDetails, error) {
	m.myDetails.count += 1
	m.myDetails.lastCtx = ctx

	return m.myDetails.data, m.myDetails.err
}

func (m *mockCentralCasesClient) UserByEmail(ctx sirius.Context, email string) (sirius.User, error) {
	m.userByEmail.count += 1
	m.userByEmail.lastCtx = ctx
	m.userByEmail.lastEmail = email

	return m.userByEmail.data, m.userByEmail.err
}

func TestGetCentralCases(t *testing.T) {
	assert := assert.New(t)

	client := &mockCentralCasesClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Manager"},
		Teams: []sirius.MyDetailsTeam{{ID: 123, DisplayName: "team"}},
	}
	client.userByEmail.data = sirius.User{
		ID: 14,
	}
	client.casesByAssignee.data = []sirius.Case{{
		ID: 78,
		Donor: sirius.Donor{
			ID: 79,
		},
	}}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := centralCases(client, template.Func)(w, r)
	assert.Nil(err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.userByEmail.count)
	assert.Equal(getContext(r), client.userByEmail.lastCtx)
	assert.Equal("manager@opgtest.com", client.userByEmail.lastEmail)

	assert.Equal(2, client.casesByAssignee.count)
	assert.Equal(getContext(r), client.casesByAssignee.lastCtx)
	assert.Equal(14, client.casesByAssignee.lastId)
	assert.Equal(sirius.Criteria{}.Filter("status", "Pending").Page(1).Sort("receiptDate", sirius.Ascending), client.casesByAssignee.criteria[0])
	assert.Equal(sirius.Criteria{}.Filter("status", "Pending").Page(1).Sort("receiptDate", sirius.Ascending).Limit(1), client.casesByAssignee.criteria[1])

	assert.Equal(1, template.count)
	assert.Equal(centralCasesVars{
		Cases:    client.casesByAssignee.data,
		TeamID:   123,
		TeamName: "team",
	}, template.lastVars)
}

func TestGetCentralCasesPage(t *testing.T) {
	assert := assert.New(t)

	client := &mockCentralCasesClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Case Manager", "Manager", "System Admin"},
	}
	client.userByEmail.data = sirius.User{
		ID: 14,
	}
	client.casesByAssignee.data = []sirius.Case{{
		ID: 78,
		Donor: sirius.Donor{
			ID: 79,
		},
	}}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path?page=4", nil)

	err := centralCases(client, template.Func)(w, r)
	assert.Nil(err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.userByEmail.count)
	assert.Equal(getContext(r), client.userByEmail.lastCtx)
	assert.Equal("manager@opgtest.com", client.userByEmail.lastEmail)

	assert.Equal(2, client.casesByAssignee.count)
	assert.Equal(getContext(r), client.casesByAssignee.lastCtx)
	assert.Equal(14, client.casesByAssignee.lastId)
	assert.Equal(sirius.Criteria{}.Filter("status", "Pending").Page(4).Sort("receiptDate", sirius.Ascending), client.casesByAssignee.criteria[0])
	assert.Equal(sirius.Criteria{}.Filter("status", "Pending").Page(1).Sort("receiptDate", sirius.Ascending).Limit(1), client.casesByAssignee.criteria[1])

	assert.Equal(1, template.count)
	assert.Equal(centralCasesVars{
		Cases: client.casesByAssignee.data,
	}, template.lastVars)
}

func TestGetCentralCasesForbidden(t *testing.T) {
	assert := assert.New(t)

	client := &mockCentralCasesClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Case Manager", "System Admin"},
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := centralCases(client, template.Func)(w, r)
	assert.Equal(StatusError(http.StatusForbidden), err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(0, client.userByEmail.count)
	assert.Equal(0, client.casesByAssignee.count)
}

func TestGetCentralCasesMyDetailsError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockCentralCasesClient{}
	client.myDetails.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := centralCases(client, template.Func)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(0, client.userByEmail.count)
	assert.Equal(0, client.casesByAssignee.count)
}

func TestGetCentralCasesUserError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockCentralCasesClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Manager"},
	}
	client.userByEmail.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := centralCases(client, template.Func)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.userByEmail.count)
	assert.Equal(getContext(r), client.userByEmail.lastCtx)
	assert.Equal("manager@opgtest.com", client.userByEmail.lastEmail)

	assert.Equal(0, client.casesByAssignee.count)
}

func TestGetCentralCasesQueryError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockCentralCasesClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Manager"},
	}
	client.userByEmail.data = sirius.User{
		ID: 14,
	}
	client.casesByAssignee.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := centralCases(client, template.Func)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.userByEmail.count)
	assert.Equal(getContext(r), client.userByEmail.lastCtx)
	assert.Equal("manager@opgtest.com", client.userByEmail.lastEmail)

	assert.Equal(1, client.casesByAssignee.count)
	assert.Equal(getContext(r), client.casesByAssignee.lastCtx)
	assert.Equal(14, client.casesByAssignee.lastId)
	assert.Equal(sirius.Criteria{}.Filter("status", "Pending").Page(1).Sort("receiptDate", sirius.Ascending), client.casesByAssignee.criteria[0])
}

func TestBadMethodCentralCases(t *testing.T) {
	assert := assert.New(t)

	client := &mockCentralCasesClient{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/path", nil)

	err := centralCases(client, template.Func)(w, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)

	assert.Equal(0, client.userByEmail.count)
	assert.Equal(0, client.casesByAssignee.count)
	assert.Equal(0, template.count)
}

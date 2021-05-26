package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockCasesClient struct {
	myDetails struct {
		count   int
		lastCtx sirius.Context
		data    sirius.MyDetails
		err     error
	}
	casesByAssignee struct {
		count      int
		lastCtx    sirius.Context
		lastId     int
		lastStatus string
		lastPage   int
		data       []sirius.Case
		pagination *sirius.Pagination
		err        error
	}
	hasWorkableCase struct {
		count   int
		lastCtx sirius.Context
		lastId  int
		data    bool
		err     error
	}
}

func (m *mockCasesClient) MyDetails(ctx sirius.Context) (sirius.MyDetails, error) {
	m.myDetails.count += 1
	m.myDetails.lastCtx = ctx

	return m.myDetails.data, m.myDetails.err
}

func (m *mockCasesClient) CasesByAssignee(ctx sirius.Context, id int, status string, page int) ([]sirius.Case, *sirius.Pagination, error) {
	m.casesByAssignee.count += 1
	m.casesByAssignee.lastCtx = ctx
	m.casesByAssignee.lastId = id
	m.casesByAssignee.lastStatus = status
	m.casesByAssignee.lastPage = page

	return m.casesByAssignee.data, m.casesByAssignee.pagination, m.casesByAssignee.err
}

func (m *mockCasesClient) HasWorkableCase(ctx sirius.Context, id int) (bool, error) {
	m.hasWorkableCase.count += 1
	m.hasWorkableCase.lastCtx = ctx
	m.hasWorkableCase.lastId = id

	return m.hasWorkableCase.data, m.hasWorkableCase.err
}

func TestGetCases(t *testing.T) {
	assert := assert.New(t)

	client := &mockCasesClient{}
	client.myDetails.data = sirius.MyDetails{
		ID: 14,
	}
	client.casesByAssignee.data = []sirius.Case{{
		ID: 78,
		Donor: sirius.Donor{
			ID: 79,
		},
	}}
	client.hasWorkableCase.data = true
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := cases(client, template)(w, r)
	assert.Nil(err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.casesByAssignee.count)
	assert.Equal(getContext(r), client.casesByAssignee.lastCtx)
	assert.Equal(14, client.casesByAssignee.lastId)
	assert.Equal("", client.casesByAssignee.lastStatus)
	assert.Equal(1, client.casesByAssignee.lastPage)

	assert.Equal(1, client.hasWorkableCase.count)
	assert.Equal(getContext(r), client.hasWorkableCase.lastCtx)
	assert.Equal(14, client.hasWorkableCase.lastId)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(allCasesVars{
		Cases:           client.casesByAssignee.data,
		HasWorkableCase: true,
	}, template.lastVars)
}

func TestGetCasesPage(t *testing.T) {
	assert := assert.New(t)

	client := &mockCasesClient{}
	client.myDetails.data = sirius.MyDetails{
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

	err := cases(client, template)(w, r)
	assert.Nil(err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.casesByAssignee.count)
	assert.Equal(getContext(r), client.casesByAssignee.lastCtx)
	assert.Equal(14, client.casesByAssignee.lastId)
	assert.Equal("", client.casesByAssignee.lastStatus)
	assert.Equal(4, client.casesByAssignee.lastPage)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(allCasesVars{
		Cases: client.casesByAssignee.data,
	}, template.lastVars)
}

func TestGetCasesMyDetailsError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockCasesClient{}
	client.myDetails.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := cases(client, template)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(0, client.casesByAssignee.count)
}

func TestGetCasesQueryError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockCasesClient{}
	client.myDetails.data = sirius.MyDetails{
		ID: 14,
	}
	client.casesByAssignee.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := cases(client, template)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.casesByAssignee.count)
	assert.Equal(getContext(r), client.casesByAssignee.lastCtx)
	assert.Equal(14, client.casesByAssignee.lastId)
	assert.Equal("", client.casesByAssignee.lastStatus)
}

func TestBadMethodCases(t *testing.T) {
	assert := assert.New(t)

	client := &mockCasesClient{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/path", nil)

	err := cases(client, template)(w, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)

	assert.Equal(0, client.myDetails.count)
	assert.Equal(0, client.casesByAssignee.count)
	assert.Equal(0, template.count)
}

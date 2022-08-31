package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockTasksDashboardClient struct {
	myDetails struct {
		count   int
		lastCtx sirius.Context
		data    sirius.MyDetails
		err     error
	}
	tasksByAssignee struct {
		count        int
		lastCtx      sirius.Context
		lastId       int
		lastCriteria sirius.Criteria
		data         []sirius.Task
		pagination   *sirius.Pagination
		err          error
	}
}

func (m *mockTasksDashboardClient) MyDetails(ctx sirius.Context) (sirius.MyDetails, error) {
	m.myDetails.count += 1
	m.myDetails.lastCtx = ctx

	return m.myDetails.data, m.myDetails.err
}

func (m *mockTasksDashboardClient) TasksByAssignee(ctx sirius.Context, id int, criteria sirius.Criteria) ([]sirius.Task, *sirius.Pagination, error) {
	m.tasksByAssignee.count += 1
	m.tasksByAssignee.lastCtx = ctx
	m.tasksByAssignee.lastId = id
	m.tasksByAssignee.lastCriteria = criteria

	return m.tasksByAssignee.data, m.tasksByAssignee.pagination, m.tasksByAssignee.err
}

func TestGetTasksDashboard(t *testing.T) {
	assert := assert.New(t)

	client := &mockTasksDashboardClient{}
	client.myDetails.data = sirius.MyDetails{
		ID: 14,
	}
	client.tasksByAssignee.data = []sirius.Task{{
		ID: 78,
	}}
	client.tasksByAssignee.pagination = &sirius.Pagination{
		TotalItems: 20,
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := tasksDashboard(client, template)(w, r)
	assert.Nil(err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.tasksByAssignee.count)
	assert.Equal(getContext(r), client.tasksByAssignee.lastCtx)
	assert.Equal(14, client.tasksByAssignee.lastId)
	assert.Equal(sirius.Criteria{}.Filter("status", "Not started").Sort("dueDate", sirius.Ascending).Sort("name", sirius.Descending), client.tasksByAssignee.lastCriteria)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(tasksDashboardVars{
		Tasks:     client.tasksByAssignee.data,
		Title:     "Tasks Dashboard",
		XSRFToken: getContext(r).XSRFToken,
	}, template.lastVars)
}

func TestGetTasksDashboardTitle(t *testing.T) {
	testCases := map[string]struct {
		expectedTitle string
	}{
		"Data Quality Team": {
			expectedTitle: "Data Quality Dashboard",
		},
		"System Administrators": {
			expectedTitle: "System Administrators Dashboard",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			client := &mockTasksDashboardClient{}
			client.myDetails.data = sirius.MyDetails{
				ID: 14,
				Teams: []sirius.MyDetailsTeam{
					{DisplayName: name},
				},
			}
			template := &mockTemplate{}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "/path", nil)

			err := tasksDashboard(client, template)(w, r)
			assert.Nil(err)

			vars, _ := template.lastVars.(tasksDashboardVars)
			assert.Equal(tc.expectedTitle, vars.Title)
		})
	}
}

func TestGetTasksDashboardMyDetailsError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockTasksDashboardClient{}
	client.myDetails.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := tasksDashboard(client, template)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(0, client.tasksByAssignee.count)
}

func TestGetTasksDashboardQueryError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockTasksDashboardClient{}
	client.myDetails.data = sirius.MyDetails{
		ID: 14,
	}
	client.tasksByAssignee.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := tasksDashboard(client, template)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(1, client.tasksByAssignee.count)
}

func TestBadMethodTasksDashboard(t *testing.T) {
	assert := assert.New(t)

	client := &mockTasksDashboardClient{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	err := tasksDashboard(client, template)(w, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)

	assert.Equal(0, client.myDetails.count)
	assert.Equal(0, client.tasksByAssignee.count)
	assert.Equal(0, template.count)
}

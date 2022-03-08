package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockReassignClient struct {
	myDetails struct {
		count   int
		lastCtx sirius.Context
		data    sirius.MyDetails
		err     error
	}
	user struct {
		count int
		ctx   []sirius.Context
		id    []int
		data  []sirius.Assignee
		err   []error
	}
	userByEmail struct {
		count     int
		lastCtx   sirius.Context
		lastEmail string
		data      sirius.User
		err       error
	}
	team struct {
		count   int
		lastCtx sirius.Context
		lastId  int
		data    sirius.Team
		err     error
	}
	assign struct {
		count        int
		lastCtx      sirius.Context
		lastCases    []int
		lastAssignee int
		err          error
	}
}

func (m *mockReassignClient) MyDetails(ctx sirius.Context) (sirius.MyDetails, error) {
	m.myDetails.count += 1
	m.myDetails.lastCtx = ctx

	return m.myDetails.data, m.myDetails.err
}

func (m *mockReassignClient) User(ctx sirius.Context, id int) (sirius.Assignee, error) {
	i := m.user.count

	m.user.count += 1
	m.user.ctx = append(m.user.ctx, ctx)
	m.user.id = append(m.user.id, id)

	return m.user.data[i], m.user.err[i]
}

func (m *mockReassignClient) UserByEmail(ctx sirius.Context, email string) (sirius.User, error) {
	m.userByEmail.count += 1
	m.userByEmail.lastCtx = ctx
	m.userByEmail.lastEmail = email

	return m.userByEmail.data, m.userByEmail.err
}

func (m *mockReassignClient) Team(ctx sirius.Context, id int) (sirius.Team, error) {
	m.team.count += 1
	m.team.lastCtx = ctx
	m.team.lastId = id

	return m.team.data, m.team.err
}

func (m *mockReassignClient) Assign(ctx sirius.Context, cases []int, assignee int) error {
	m.assign.count += 1
	m.assign.lastCtx = ctx
	m.assign.lastCases = cases
	m.assign.lastAssignee = assignee

	return m.assign.err
}

func TestGetReassign(t *testing.T) {
	assert := assert.New(t)

	client := &mockReassignClient{}
	client.myDetails.data = sirius.MyDetails{
		ID:    14,
		Roles: []string{"Manager"},
	}
	client.user.data = []sirius.Assignee{{
		ID:          47,
		DisplayName: "some person",
		Teams: []sirius.Team{{
			ID: 439,
		}},
	}}
	client.user.err = []error{nil}
	client.team.data = sirius.Team{
		ID: 439,
		Members: []sirius.TeamMember{
			{
				ID:          440,
				DisplayName: "person 1",
			},
		},
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path?selected=1&selected=4&assignee=47", nil)

	err := reassign(client, template.Func)(w, r)
	assert.Nil(err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.user.count)
	assert.Equal(getContext(r), client.user.ctx[0])
	assert.Equal(47, client.user.id[0])

	assert.Equal(1, client.team.count)
	assert.Equal(getContext(r), client.team.lastCtx)
	assert.Equal(439, client.team.lastId)

	assert.Equal(1, template.count)
	assert.Equal(reassignVars{
		XSRFToken:   getContext(r).XSRFToken,
		Selected:    []int{1, 4},
		Assignee:    client.user.data[0],
		TeamMembers: client.team.data.Members,
	}, template.lastVars)
}

func TestGetReassignNotManager(t *testing.T) {
	assert := assert.New(t)

	client := &mockReassignClient{}
	client.myDetails.data = sirius.MyDetails{
		ID:    14,
		Roles: []string{},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path?selected=1&selected=4&assignee=47", nil)

	err := reassign(client, nil)(w, r)
	assert.Equal(StatusError(http.StatusForbidden), err)
}

func TestGetReassignBadRequest(t *testing.T) {
	testCases := map[string]string{
		"bad-assignee": "/path?selected=1&selected=4&assignee=what",
		"bad-selected": "/path?selected=1&selected=what&assignee=47",
	}

	for name, path := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			client := &mockReassignClient{}
			client.myDetails.data = sirius.MyDetails{
				ID:    14,
				Roles: []string{"Manager"},
			}
			client.user.data = []sirius.Assignee{{
				ID:          47,
				DisplayName: "some person",
				Teams: []sirius.Team{{
					ID: 439,
				}},
			}}
			client.user.err = []error{nil}
			client.team.data = sirius.Team{
				ID: 439,
				Members: []sirius.TeamMember{
					{
						ID:          440,
						DisplayName: "person 1",
					},
				},
			}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", path, nil)

			err := reassign(client, nil)(w, r)
			assert.Equal(StatusError(http.StatusBadRequest), err)
		})
	}
}

func TestGetReassignMyDetailsError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockReassignClient{}
	client.myDetails.err = expectedError

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path?selected=1&selected=4&assignee=47", nil)

	err := reassign(client, nil)(w, r)
	assert.Equal(expectedError, err)
}

func TestGetReassignUserError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockReassignClient{}
	client.myDetails.data = sirius.MyDetails{
		ID:    14,
		Roles: []string{"Manager"},
	}
	client.user.data = []sirius.Assignee{{}}
	client.user.err = []error{expectedError}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path?selected=1&selected=4&assignee=47", nil)

	err := reassign(client, nil)(w, r)
	assert.Equal(expectedError, err)
}

func TestGetReassignTeamError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockReassignClient{}
	client.myDetails.data = sirius.MyDetails{
		ID:    14,
		Roles: []string{"Manager"},
	}
	client.user.data = []sirius.Assignee{{
		ID:          47,
		DisplayName: "some person",
		Teams: []sirius.Team{{
			ID: 439,
		}},
	}}
	client.user.err = []error{nil}
	client.team.err = expectedError

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path?selected=1&selected=4&assignee=47", nil)

	err := reassign(client, nil)(w, r)
	assert.Equal(expectedError, err)
}

func TestPostReassignToCentralPot(t *testing.T) {
	assert := assert.New(t)

	client := &mockReassignClient{}
	client.myDetails.data = sirius.MyDetails{
		ID:    14,
		Roles: []string{"Manager"},
	}
	client.user.data = []sirius.Assignee{{
		ID:          47,
		DisplayName: "some person",
		Teams: []sirius.Team{{
			ID: 439,
		}},
	}}
	client.user.err = []error{nil}
	client.team.data = sirius.Team{
		ID: 439,
		Members: []sirius.TeamMember{
			{
				ID:          440,
				DisplayName: "person 1",
			},
		},
	}
	client.userByEmail.data = sirius.User{
		ID: 50,
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("selected=1&selected=4&assignee=47&reassign=central-pot"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := reassign(client, template.Func)(w, r)
	assert.Nil(err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.user.count)
	assert.Equal(getContext(r), client.user.ctx[0])
	assert.Equal(47, client.user.id[0])

	assert.Equal(1, client.team.count)
	assert.Equal(getContext(r), client.team.lastCtx)
	assert.Equal(439, client.team.lastId)

	assert.Equal(1, client.userByEmail.count)
	assert.Equal(getContext(r), client.userByEmail.lastCtx)
	assert.Equal("manager@opgtest.com", client.userByEmail.lastEmail)

	assert.Equal(1, client.assign.count)
	assert.Equal(getContext(r), client.assign.lastCtx)
	assert.Equal([]int{1, 4}, client.assign.lastCases)
	assert.Equal(50, client.assign.lastAssignee)

	assert.Equal(1, template.count)
	assert.Equal(reassignVars{
		XSRFToken:   getContext(r).XSRFToken,
		Selected:    []int{1, 4},
		Assignee:    client.user.data[0],
		TeamMembers: client.team.data.Members,
		Success:     true,
		AssignedTo: sirius.Assignee{
			ID:          50,
			DisplayName: "Central Pot",
		},
	}, template.lastVars)
}

func TestPostReassignToCentralPotError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockReassignClient{}
	client.myDetails.data = sirius.MyDetails{
		ID:    14,
		Roles: []string{"Manager"},
	}
	client.user.data = []sirius.Assignee{{
		ID:          47,
		DisplayName: "some person",
		Teams: []sirius.Team{{
			ID: 439,
		}},
	}}
	client.user.err = []error{nil}
	client.team.data = sirius.Team{
		ID: 439,
		Members: []sirius.TeamMember{
			{
				ID:          440,
				DisplayName: "person 1",
			},
		},
	}
	client.userByEmail.err = expectedError

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("selected=1&selected=4&assignee=47&reassign=central-pot"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := reassign(client, nil)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(1, client.user.count)
	assert.Equal(1, client.team.count)
	assert.Equal(1, client.userByEmail.count)
	assert.Equal(0, client.assign.count)
}

func TestPostReassignToUser(t *testing.T) {
	assert := assert.New(t)

	client := &mockReassignClient{}
	client.myDetails.data = sirius.MyDetails{
		ID:    14,
		Roles: []string{"Manager"},
	}
	client.user.data = []sirius.Assignee{
		{
			ID:          47,
			DisplayName: "some person",
			Teams: []sirius.Team{{
				ID: 439,
			}},
		},
		{
			ID:          99,
			DisplayName: "Assigned to user",
		},
	}
	client.user.err = []error{nil, nil}
	client.team.data = sirius.Team{
		ID: 439,
		Members: []sirius.TeamMember{
			{
				ID:          440,
				DisplayName: "person 1",
			},
		},
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("selected=1&selected=4&assignee=47&reassign=user&caseworker=99"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := reassign(client, template.Func)(w, r)
	assert.Nil(err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(2, client.user.count)
	assert.Equal(getContext(r), client.user.ctx[0])
	assert.Equal(47, client.user.id[0])
	assert.Equal(getContext(r), client.user.ctx[1])
	assert.Equal(99, client.user.id[1])

	assert.Equal(1, client.team.count)
	assert.Equal(getContext(r), client.team.lastCtx)
	assert.Equal(439, client.team.lastId)

	assert.Equal(1, client.assign.count)
	assert.Equal(getContext(r), client.assign.lastCtx)
	assert.Equal([]int{1, 4}, client.assign.lastCases)
	assert.Equal(99, client.assign.lastAssignee)

	assert.Equal(1, template.count)
	assert.Equal(reassignVars{
		XSRFToken:   getContext(r).XSRFToken,
		Selected:    []int{1, 4},
		Assignee:    client.user.data[0],
		TeamMembers: client.team.data.Members,
		Success:     true,
		AssignedTo:  client.user.data[1],
	}, template.lastVars)
}

func TestPostReassignToUserError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockReassignClient{}
	client.myDetails.data = sirius.MyDetails{
		ID:    14,
		Roles: []string{"Manager"},
	}
	client.user.data = []sirius.Assignee{
		{
			ID:          47,
			DisplayName: "some person",
			Teams: []sirius.Team{{
				ID: 439,
			}},
		},
		{},
	}
	client.user.err = []error{nil, expectedError}
	client.team.data = sirius.Team{
		ID: 439,
		Members: []sirius.TeamMember{
			{
				ID:          440,
				DisplayName: "person 1",
			},
		},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("selected=1&selected=4&assignee=47&reassign=user&caseworker=99"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := reassign(client, nil)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(2, client.user.count)
	assert.Equal(getContext(r), client.user.ctx[0])
	assert.Equal(47, client.user.id[0])
	assert.Equal(getContext(r), client.user.ctx[1])
	assert.Equal(99, client.user.id[1])

	assert.Equal(1, client.team.count)
	assert.Equal(getContext(r), client.team.lastCtx)
	assert.Equal(439, client.team.lastId)

	assert.Equal(0, client.assign.count)
}

func TestPostReassignBadRequest(t *testing.T) {
	testCases := map[string]string{
		"bad-assignee":   "selected=1&selected=4&assignee=what&reassign=central-pot",
		"bad-selected":   "selected=1&selected=what&assignee=47&reassign=central-pot",
		"bad-reassign":   "selected=1&selected=4&assignee=47&reassign=what&caseworker=5",
		"bad-caseworker": "selected=1&selected=4&assignee=47&reassign=user&caseworker=what",
	}

	for name, path := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			client := &mockReassignClient{}
			client.myDetails.data = sirius.MyDetails{
				ID:    14,
				Roles: []string{"Manager"},
			}
			client.user.data = []sirius.Assignee{{
				ID:          47,
				DisplayName: "some person",
				Teams: []sirius.Team{{
					ID: 439,
				}},
			}}
			client.user.err = []error{nil}
			client.team.data = sirius.Team{
				ID: 439,
				Members: []sirius.TeamMember{
					{
						ID:          440,
						DisplayName: "person 1",
					},
				},
			}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "/path", strings.NewReader(path))
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			err := reassign(client, nil)(w, r)
			assert.Equal(StatusError(http.StatusBadRequest), err)
		})
	}
}

func TestPostReassignError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockReassignClient{}
	client.myDetails.data = sirius.MyDetails{
		ID:    14,
		Roles: []string{"Manager"},
	}
	client.user.data = []sirius.Assignee{
		{
			ID:          47,
			DisplayName: "some person",
			Teams: []sirius.Team{{
				ID: 439,
			}},
		},
		{
			ID:          99,
			DisplayName: "Assigned to user",
		},
	}
	client.user.err = []error{nil, nil}
	client.team.data = sirius.Team{
		ID: 439,
		Members: []sirius.TeamMember{
			{
				ID:          440,
				DisplayName: "person 1",
			},
		},
	}
	client.assign.err = expectedError

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("selected=1&selected=4&assignee=47&reassign=user&caseworker=99"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := reassign(client, nil)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(2, client.user.count)
	assert.Equal(1, client.team.count)
	assert.Equal(1, client.assign.count)
}

func TestBadMethodReassign(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/path", nil)

	err := reassign(nil, nil)(w, r)
	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)
}

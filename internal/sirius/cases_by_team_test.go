package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestCasesByTeam(t *testing.T) {
	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name           string
		criteria       Criteria
		setup          func()
		expectedResult *CasesByTeam
		expectedError  error
	}{
		{
			name:     "OK",
			criteria: Criteria{}.Page(1),
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to get my team's cases").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/teams/66/cases"),
						Query: dsl.MapMatcher{
							"page": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Like(map[string]interface{}{
							"total": dsl.Like(1),
							"limit": dsl.Like(25),
							"pages": dsl.Like(map[string]interface{}{
								"current": dsl.Like(1),
								"total":   dsl.Like(1),
							}),
							"metadata": dsl.Like(map[string]interface{}{
								"workedTotal": dsl.Like(1),
								"worked": dsl.EachLike(map[string]interface{}{
									"assignee": dsl.Like(map[string]interface{}{
										"id":          dsl.Like(17),
										"displayName": dsl.Like("John Smith"),
									}),
									"total": dsl.Like(1),
								}, 1),
								"tasksCompleted": dsl.EachLike(map[string]interface{}{
									"assignee": dsl.Like(map[string]interface{}{
										"id":          dsl.Like(17),
										"displayName": dsl.Like("John Smith"),
									}),
									"total": dsl.Like(3),
								}, 1),
							}),
							"cases": dsl.EachLike(map[string]interface{}{
								"id":  dsl.Like(36),
								"uId": dsl.Term("7000-8548-8461", `\d{4}-\d{4}-\d{4}`),
								"donor": dsl.Like(map[string]interface{}{
									"id":        dsl.Like(23),
									"uId":       dsl.Term("7000-5382-4438", `\d{4}-\d{4}-\d{4}`),
									"firstname": dsl.Like("Adrian"),
									"surname":   dsl.Like("Kurkjian"),
								}),
								"assignee": dsl.Like(map[string]interface{}{
									"id":          dsl.Like(17),
									"displayName": dsl.Like("John Smith"),
								}),
								"caseSubtype": dsl.Term("pfa", "hw|pfa"),
								"receiptDate": dsl.Term("12/05/2021", `\d{1,2}/\d{1,2}/\d{4}`),
								"status":      dsl.Like("Perfect"),
							}, 1),
						}),
					})
			},
			expectedResult: &CasesByTeam{
				Cases: []Case{{
					ID:  36,
					Uid: "7000-8548-8461",
					Donor: Donor{
						ID:        23,
						Uid:       "7000-5382-4438",
						Firstname: "Adrian",
						Surname:   "Kurkjian",
					},
					Assignee: Assignee{
						ID:          17,
						DisplayName: "John Smith",
					},
					SubType:     "pfa",
					ReceiptDate: SiriusDate{time.Date(2021, 5, 12, 0, 0, 0, 0, time.UTC)},
					Status:      "Perfect",
				}},
				Pagination: &Pagination{
					TotalItems:  1,
					CurrentPage: 1,
					TotalPages:  1,
					PageSize:    25,
				},
				Stats: CasesByTeamMetadata{
					WorkedTotal: 1,
					Worked: []CasesByTeamMetadataMember{{
						Assignee: Assignee{
							ID:          17,
							DisplayName: "John Smith",
						},
						Total: 1,
					}},
					TasksCompleted: []CasesByTeamMetadataMember{{
						Assignee: Assignee{
							ID:          17,
							DisplayName: "John Smith",
						},
						Total: 3,
					}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				result, err := client.CasesByTeam(Context{Context: context.Background()}, 66, tc.criteria)
				assert.Equal(t, tc.expectedResult, result)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestCasesByTeamIgnored(t *testing.T) {
	pact := newIgnoredPact()
	defer pact.Teardown()

	testCases := []struct {
		name           string
		criteria       Criteria
		setup          func()
		expectedResult *CasesByTeam
		expectedError  error
	}{
		{
			name:     "OK",
			criteria: Criteria{}.Page(1),
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to get team 67's cases").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/teams/67/cases"),
						Query: dsl.MapMatcher{
							"page": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Like(map[string]interface{}{
							"total": dsl.Like(1),
							"limit": dsl.Like(25),
							"pages": dsl.Like(map[string]interface{}{
								"current": dsl.Like(1),
								"total":   dsl.Like(1),
							}),
							"metadata": dsl.Like(map[string]interface{}{
								"workedTotal": dsl.Like(1),
								"worked": dsl.EachLike(map[string]interface{}{
									"assignee": dsl.Like(map[string]interface{}{
										"id":          dsl.Like(17),
										"displayName": dsl.Like("John Smith"),
									}),
									"total": dsl.Like(1),
								}, 1),
								"tasksCompleted": dsl.EachLike(map[string]interface{}{
									"assignee": dsl.Like(map[string]interface{}{
										"id":          dsl.Like(17),
										"displayName": dsl.Like("John Smith"),
									}),
									"total": dsl.Like(3),
								}, 1),
							}),
							"cases": dsl.EachLike(map[string]interface{}{
								"id":  dsl.Like(36),
								"uId": dsl.Term("7000-8548-8461", `\d{4}-\d{4}-\d{4}`),
								"donor": dsl.Like(map[string]interface{}{
									"id":        dsl.Like(23),
									"uId":       dsl.Term("7000-5382-4438", `\d{4}-\d{4}-\d{4}`),
									"firstname": dsl.Like("Adrian"),
									"surname":   dsl.Like("Kurkjian"),
								}),
								"assignee": dsl.Like(map[string]interface{}{
									"id":          dsl.Like(17),
									"displayName": dsl.Like("John Smith"),
								}),
								"caseSubtype": dsl.Term("pfa", "hw|pfa"),
								"receiptDate": dsl.Term("12/05/2021", `\d{1,2}/\d{1,2}/\d{4}`),
								"status":      dsl.Like("Perfect"),
							}, 1),
						}),
					})
			},
			expectedResult: &CasesByTeam{
				Cases: []Case{{
					ID:  36,
					Uid: "7000-8548-8461",
					Donor: Donor{
						ID:        23,
						Uid:       "7000-5382-4438",
						Firstname: "Adrian",
						Surname:   "Kurkjian",
					},
					Assignee: Assignee{
						ID:          17,
						DisplayName: "John Smith",
					},
					SubType:     "pfa",
					ReceiptDate: SiriusDate{time.Date(2021, 5, 12, 0, 0, 0, 0, time.UTC)},
					Status:      "Perfect",
				}},
				Pagination: &Pagination{
					TotalItems:  1,
					CurrentPage: 1,
					TotalPages:  1,
					PageSize:    25,
				},
				Stats: CasesByTeamMetadata{
					WorkedTotal: 1,
					Worked: []CasesByTeamMetadataMember{{
						Assignee: Assignee{
							ID:          17,
							DisplayName: "John Smith",
						},
						Total: 1,
					}},
					TasksCompleted: []CasesByTeamMetadataMember{{
						Assignee: Assignee{
							ID:          17,
							DisplayName: "John Smith",
						},
						Total: 3,
					}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				result, err := client.CasesByTeam(Context{Context: context.Background()}, 67, tc.criteria)
				assert.Equal(t, tc.expectedResult, result)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestCasesByTeamWithAllocationIgnored(t *testing.T) {
	pact := newIgnoredPact()
	defer pact.Teardown()

	testCases := []struct {
		name           string
		criteria       Criteria
		setup          func()
		expectedResult *CasesByTeam
		expectedError  error
	}{
		{
			name:     "OK with allocation",
			criteria: Criteria{}.Filter("allocation", "47").Page(1),
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to get my team's cases filtered").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/teams/66/cases"),
						Query: dsl.MapMatcher{
							"page":   dsl.String("1"),
							"filter": dsl.String("allocation:47"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Like(map[string]interface{}{
							"total": dsl.Like(1),
							"limit": dsl.Like(25),
							"pages": dsl.Like(map[string]interface{}{
								"current": dsl.Like(1),
								"total":   dsl.Like(1),
							}),
							"metadata": dsl.Like(map[string]interface{}{
								"workedTotal": dsl.Like(1),
								"worked": dsl.EachLike(map[string]interface{}{
									"assignee": dsl.Like(map[string]interface{}{
										"id":          dsl.Like(17),
										"displayName": dsl.Like("John Smith"),
									}),
									"total": dsl.Like(1),
								}, 1),
								"tasksCompleted": dsl.EachLike(map[string]interface{}{
									"assignee": dsl.Like(map[string]interface{}{
										"id":          dsl.Like(17),
										"displayName": dsl.Like("John Smith"),
									}),
									"total": dsl.Like(3),
								}, 1),
							}),
							"cases": dsl.EachLike(map[string]interface{}{
								"id":  dsl.Like(36),
								"uId": dsl.Term("7000-8548-8461", `\d{4}-\d{4}-\d{4}`),
								"donor": dsl.Like(map[string]interface{}{
									"id":        dsl.Like(23),
									"uId":       dsl.Term("7000-5382-4438", `\d{4}-\d{4}-\d{4}`),
									"firstname": dsl.Like("Someone"),
									"surname":   dsl.Like("Else"),
								}),
								"assignee": dsl.Like(map[string]interface{}{
									"id":          dsl.Like(17),
									"displayName": dsl.Like("John Smith"),
								}),
								"caseSubtype": dsl.Term("pfa", "hw|pfa"),
								"receiptDate": dsl.Term("12/05/2021", `\d{1,2}/\d{1,2}/\d{4}`),
								"status":      dsl.Like("Perfect"),
							}, 1),
						}),
					})
			},
			expectedResult: &CasesByTeam{
				Cases: []Case{{
					ID:  36,
					Uid: "7000-8548-8461",
					Donor: Donor{
						ID:        23,
						Uid:       "7000-5382-4438",
						Firstname: "Someone",
						Surname:   "Else",
					},
					Assignee: Assignee{
						ID:          17,
						DisplayName: "John Smith",
					},
					SubType:     "pfa",
					ReceiptDate: SiriusDate{time.Date(2021, 5, 12, 0, 0, 0, 0, time.UTC)},
					Status:      "Perfect",
				}},
				Pagination: &Pagination{
					TotalItems:  1,
					CurrentPage: 1,
					TotalPages:  1,
					PageSize:    25,
				},
				Stats: CasesByTeamMetadata{
					WorkedTotal: 1,
					Worked: []CasesByTeamMetadataMember{{
						Assignee: Assignee{
							ID:          17,
							DisplayName: "John Smith",
						},
						Total: 1,
					}},
					TasksCompleted: []CasesByTeamMetadataMember{{
						Assignee: Assignee{
							ID:          17,
							DisplayName: "John Smith",
						},
						Total: 3,
					}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				result, err := client.CasesByTeam(Context{Context: context.Background()}, 66, tc.criteria)
				assert.Equal(t, tc.expectedResult, result)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestCasesByTeamStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.CasesByTeam(Context{Context: context.Background()}, 66, Criteria{}.Page(2))
	assert.Equal(t, &StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/teams/66/cases?page=2",
		Method: http.MethodGet,
	}, err)
}

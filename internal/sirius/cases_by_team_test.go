package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestCasesByTeam(t *testing.T) {
	pact, err := newPact()

	assert.NoError(t, err)

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
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/api/v1/teams/66/cases"),
						Query: matchers.MapMatcher{
							"page": matchers.String("1"),
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.Like(map[string]interface{}{
							"total": matchers.Like(1),
							"limit": matchers.Like(25),
							"pages": matchers.Like(map[string]interface{}{
								"current": matchers.Like(1),
								"total":   matchers.Like(1),
							}),
							"metadata": matchers.Like(map[string]interface{}{
								"workedTotal": matchers.Like(1),
								"worked": matchers.EachLike(map[string]interface{}{
									"assignee": matchers.Like(map[string]interface{}{
										"id":          matchers.Like(17),
										"displayName": matchers.Like("John Smith"),
									}),
									"total": matchers.Like(1),
								}, 1),
								"tasksCompleted": matchers.EachLike(map[string]interface{}{
									"assignee": matchers.Like(map[string]interface{}{
										"id":          matchers.Like(17),
										"displayName": matchers.Like("John Smith"),
									}),
									"total": matchers.Like(3),
								}, 1),
							}),
							"cases": matchers.EachLike(map[string]interface{}{
								"id":  matchers.Like(36),
								"uId": matchers.Term("7000-8548-8461", `\d{4}-\d{4}-\d{4}`),
								"donor": matchers.Like(map[string]interface{}{
									"id":        matchers.Like(23),
									"uId":       matchers.Term("7000-5382-4438", `\d{4}-\d{4}-\d{4}`),
									"firstname": matchers.Like("Adrian"),
									"surname":   matchers.Like("Kurkjian"),
								}),
								"assignee": matchers.Like(map[string]interface{}{
									"id":          matchers.Like(17),
									"displayName": matchers.Like("John Smith"),
								}),
								"caseSubtype": matchers.Term("pfa", "hw|pfa"),
								"receiptDate": matchers.Term("12/05/2021", `\d{1,2}/\d{1,2}/\d{4}`),
								"status":      matchers.Like("Perfect"),
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

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				result, err := client.CasesByTeam(Context{Context: context.Background()}, 66, tc.criteria)
				assert.Equal(t, tc.expectedResult, result)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestCasesByTeamIgnored(t *testing.T) {
	pact, err := newIgnoredPact()

	assert.NoError(t, err)

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
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/api/v1/teams/67/cases"),
						Query: matchers.MapMatcher{
							"page": matchers.String("1"),
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.Like(map[string]interface{}{
							"total": matchers.Like(1),
							"limit": matchers.Like(25),
							"pages": matchers.Like(map[string]interface{}{
								"current": matchers.Like(1),
								"total":   matchers.Like(1),
							}),
							"metadata": matchers.Like(map[string]interface{}{
								"workedTotal": matchers.Like(1),
								"worked": matchers.EachLike(map[string]interface{}{
									"assignee": matchers.Like(map[string]interface{}{
										"id":          matchers.Like(17),
										"displayName": matchers.Like("John Smith"),
									}),
									"total": matchers.Like(1),
								}, 1),
								"tasksCompleted": matchers.EachLike(map[string]interface{}{
									"assignee": matchers.Like(map[string]interface{}{
										"id":          matchers.Like(17),
										"displayName": matchers.Like("John Smith"),
									}),
									"total": matchers.Like(3),
								}, 1),
							}),
							"cases": matchers.EachLike(map[string]interface{}{
								"id":  matchers.Like(36),
								"uId": matchers.Term("7000-8548-8461", `\d{4}-\d{4}-\d{4}`),
								"donor": matchers.Like(map[string]interface{}{
									"id":        matchers.Like(23),
									"uId":       matchers.Term("7000-5382-4438", `\d{4}-\d{4}-\d{4}`),
									"firstname": matchers.Like("Adrian"),
									"surname":   matchers.Like("Kurkjian"),
								}),
								"assignee": matchers.Like(map[string]interface{}{
									"id":          matchers.Like(17),
									"displayName": matchers.Like("John Smith"),
								}),
								"caseSubtype": matchers.Term("pfa", "hw|pfa"),
								"receiptDate": matchers.Term("12/05/2021", `\d{1,2}/\d{1,2}/\d{4}`),
								"status":      matchers.Like("Perfect"),
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

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				result, err := client.CasesByTeam(Context{Context: context.Background()}, 67, tc.criteria)
				assert.Equal(t, tc.expectedResult, result)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestCasesByTeamWithAllocationIgnored(t *testing.T) {
	pact, err := newIgnoredPact()

	assert.NoError(t, err)

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
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/api/v1/teams/66/cases"),
						Query: matchers.MapMatcher{
							"page":   matchers.String("1"),
							"filter": matchers.String("allocation:47"),
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.Like(map[string]interface{}{
							"total": matchers.Like(1),
							"limit": matchers.Like(25),
							"pages": matchers.Like(map[string]interface{}{
								"current": matchers.Like(1),
								"total":   matchers.Like(1),
							}),
							"metadata": matchers.Like(map[string]interface{}{
								"workedTotal": matchers.Like(1),
								"worked": matchers.EachLike(map[string]interface{}{
									"assignee": matchers.Like(map[string]interface{}{
										"id":          matchers.Like(17),
										"displayName": matchers.Like("John Smith"),
									}),
									"total": matchers.Like(1),
								}, 1),
								"tasksCompleted": matchers.EachLike(map[string]interface{}{
									"assignee": matchers.Like(map[string]interface{}{
										"id":          matchers.Like(17),
										"displayName": matchers.Like("John Smith"),
									}),
									"total": matchers.Like(3),
								}, 1),
							}),
							"cases": matchers.EachLike(map[string]interface{}{
								"id":  matchers.Like(36),
								"uId": matchers.Term("7000-8548-8461", `\d{4}-\d{4}-\d{4}`),
								"donor": matchers.Like(map[string]interface{}{
									"id":        matchers.Like(23),
									"uId":       matchers.Term("7000-5382-4438", `\d{4}-\d{4}-\d{4}`),
									"firstname": matchers.Like("Someone"),
									"surname":   matchers.Like("Else"),
								}),
								"assignee": matchers.Like(map[string]interface{}{
									"id":          matchers.Like(17),
									"displayName": matchers.Like("John Smith"),
								}),
								"caseSubtype": matchers.Term("pfa", "hw|pfa"),
								"receiptDate": matchers.Term("12/05/2021", `\d{1,2}/\d{1,2}/\d{4}`),
								"status":      matchers.Like("Perfect"),
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

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

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

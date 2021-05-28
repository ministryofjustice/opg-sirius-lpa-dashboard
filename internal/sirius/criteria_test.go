package sirius

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCriteria(t *testing.T) {
	testCases := map[string]struct {
		Input Criteria
		Query url.Values
	}{
		"defaultValues": {
			Input: Criteria{},
			Query: url.Values{},
		},
		"filter": {
			Input: Criteria{}.Filter("type", "LPA"),
			Query: url.Values{
				"filter": {"type:LPA"},
			},
		},
		"multiple-filters": {
			Input: Criteria{}.Filter("type", "LPA").Filter("status", "Registered"),
			Query: url.Values{
				"filter": {"type:LPA,status:Registered"},
			},
		},
		"limit": {
			Input: Criteria{}.Limit(16),
			Query: url.Values{
				"limit": {"16"},
			},
		},
		"page": {
			Input: Criteria{}.Page(4),
			Query: url.Values{
				"page": {"4"},
			},
		},
		"sort": {
			Input: Criteria{}.Sort("type", Ascending),
			Query: url.Values{
				"sort": {"type:asc"},
			},
		},
		"multiple-sorts": {
			Input: Criteria{}.Sort("type", Ascending).Sort("assignee", Descending),
			Query: url.Values{
				"sort": {"type:asc,assignee:desc"},
			},
		},
		"everything": {

			Input: Criteria{}.Filter("type", "LPA").Page(4).Sort("type", Ascending).Filter("status", "Registered").Limit(16).Sort("assignee", Descending),
			Query: url.Values{
				"filter": {"type:LPA,status:Registered"},
				"limit":  {"16"},
				"page":   {"4"},
				"sort":   {"type:asc,assignee:desc"},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			assert.Equal(tc.Query.Encode(), tc.Input.String())
		})
	}
}

package server

import (
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
	"github.com/stretchr/testify/assert"
)

func TestPagination(t *testing.T) {
	testCases := map[string]struct {
		TotalItems, CurrentPage, TotalPages, PageSize int
		Start, End, PreviousPage, NextPage            int
		HasPrevious, HasNext                          bool
		Pages                                         []int
	}{
		"empty": {
			TotalItems:  0,
			CurrentPage: 1,
			TotalPages:  0,
			PageSize:    25,
			Start:       1,
			End:         0,
			HasPrevious: false,
			HasNext:     false,
			Pages:       []int{},
		},
		"one-item": {
			TotalItems:  1,
			CurrentPage: 1,
			TotalPages:  1,
			PageSize:    25,
			Start:       1,
			End:         1,
			HasPrevious: false,
			HasNext:     false,
			Pages:       []int{1},
		},
		"one-page": {
			TotalItems:  25,
			CurrentPage: 1,
			TotalPages:  1,
			PageSize:    25,
			Start:       1,
			End:         25,
			HasPrevious: false,
			HasNext:     false,
			Pages:       []int{1},
		},
		"many-pages": {
			TotalItems:   76,
			CurrentPage:  2,
			TotalPages:   4,
			PageSize:     25,
			Start:        26,
			End:          50,
			HasPrevious:  true,
			PreviousPage: 1,
			HasNext:      true,
			NextPage:     3,
			Pages:        []int{1, 2, 3, 4},
		},
		"first-of-many-pages": {
			TotalItems:  76,
			CurrentPage: 1,
			TotalPages:  4,
			PageSize:    25,
			Start:       1,
			End:         25,
			HasPrevious: false,
			HasNext:     true,
			NextPage:    2,
			Pages:       []int{1, 2, 3, 4},
		},
		"last-of-many-pages": {
			TotalItems:   76,
			CurrentPage:  4,
			TotalPages:   4,
			PageSize:     25,
			Start:        76,
			End:          76,
			HasPrevious:  true,
			PreviousPage: 3,
			HasNext:      false,
			Pages:        []int{1, 2, 3, 4},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			pagination := newPagination(&sirius.Pagination{
				TotalItems:  tc.TotalItems,
				CurrentPage: tc.CurrentPage,
				TotalPages:  tc.TotalPages,
				PageSize:    tc.PageSize,
			})

			assert.Equal("?", pagination.Query)
			assert.Equal(tc.Start, pagination.Start())
			assert.Equal(tc.End, pagination.End())
			assert.Equal(tc.HasPrevious, pagination.HasPrevious())
			if tc.HasPrevious {
				assert.Equal(tc.PreviousPage, pagination.PreviousPage())
			}
			assert.Equal(tc.HasNext, pagination.HasNext())
			if tc.HasNext {
				assert.Equal(tc.NextPage, pagination.NextPage())
			}
			assert.Equal(tc.Pages, pagination.Pages())
		})
	}
}

func TestPaginationWithQuery(t *testing.T) {
	assert := assert.New(t)

	pagination := newPaginationWithQuery(&sirius.Pagination{}, "this-is-here")

	assert.Equal("?this-is-here&", pagination.Query)
}

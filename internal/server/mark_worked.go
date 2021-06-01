package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type MarkWorkedClient interface {
	MarkWorked(sirius.Context, int) error
}

func markWorked(client MarkWorkedClient) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}

		if err := r.ParseForm(); err != nil {
			return err
		}

		ctx := getContext(r)

		for _, workedID := range r.PostForm["worked"] {
			id, err := strconv.Atoi(workedID)
			if err != nil {
				return err
			}

			if err := client.MarkWorked(ctx, id); err != nil {
				return err
			}
		}

		return RedirectError("/pending-cases")
	}
}

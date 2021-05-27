package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type RequestNextCasesClient interface {
	RequestNextCases(sirius.Context) error
}

func requestNextCases(client RequestNextCasesClient) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}

		if err := client.RequestNextCases(getContext(r)); err != nil {
			return err
		}

		return RedirectError("/pending-cases")
	}
}

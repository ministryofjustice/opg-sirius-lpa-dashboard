package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type LandingPageClient interface {
	MyDetails(sirius.Context) (sirius.MyDetails, error)
}

func landingPage(client LandingPageClient) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		myDetails, err := client.MyDetails(ctx)
		if err != nil {
			return err
		}

		if myDetails.IsManager() {
			return RedirectError("/teams/central")
		} else {
			return RedirectError("/pending-cases")
		}
	}
}

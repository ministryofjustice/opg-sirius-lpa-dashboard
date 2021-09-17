package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type HomeRedirectClient interface {
	MyDetails(sirius.Context) (sirius.MyDetails, error)
}

func homeRedirect(client HomeRedirectClient) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		myDetails, err := client.MyDetails(ctx)
		if err != nil {
			return err
		}

		if myDetails.IsCardPaymentUser() {
			return RedirectError("/tasks")
		} else {
			return RedirectError("/pending-cases")
		}

	}
}

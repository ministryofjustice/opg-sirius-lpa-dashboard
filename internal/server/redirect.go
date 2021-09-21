package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type RedirectClient interface {
	MyDetails(sirius.Context) (sirius.MyDetails, error)
}

func redirect(client RedirectClient) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		myDetails, err := client.MyDetails(ctx)
		if err != nil {
			return err
		}

		if myDetails.IsCardPaymentUser() {
			return RedirectError("/card-payments")
		}

		return RedirectError("/pending-cases")
	}
}

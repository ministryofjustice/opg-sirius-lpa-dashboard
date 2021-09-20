package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type RequestNextPaymentTaskClient interface {
	RequestNextPaymentTask(sirius.Context) error
}

func requestNextPaymentTask(client RequestNextPaymentTaskClient) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}

		if err := client.RequestNextPaymentTask(getContext(r)); err != nil {
			return err
		}

		return RedirectError("/card-payments")
	}
}

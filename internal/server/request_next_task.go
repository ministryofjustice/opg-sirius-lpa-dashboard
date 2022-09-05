package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type RequestNextTaskClient interface {
	RequestNextTask(sirius.Context) error
}

func requestNextTask(client RequestNextTaskClient) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}

		if err := client.RequestNextTask(getContext(r)); err != nil {
			return err
		}

		return RedirectError("/tasks-dashboard")
	}
}

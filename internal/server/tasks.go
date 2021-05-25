package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type TasksClient interface {
	CasesWithOpenTasksByAssignee(sirius.Context, int, int) ([]sirius.Case, *sirius.Pagination, error)
	MyDetails(sirius.Context) (sirius.MyDetails, error)
}

type tasksVars struct {
	Cases      []sirius.Case
	Pagination *sirius.Pagination
}

func tasks(client TasksClient, tmpl Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		myDetails, err := client.MyDetails(ctx)
		if err != nil {
			return err
		}

		cases, pagination, err := client.CasesWithOpenTasksByAssignee(ctx, myDetails.ID, getPage(r))
		if err != nil {
			return err
		}

		return tmpl.ExecuteTemplate(w, "page", tasksVars{
			Cases:      cases,
			Pagination: pagination,
		})
	}
}

func getPage(r *http.Request) int {
	page := r.FormValue("page")
	if page == "" {
		return 1
	}

	v, err := strconv.Atoi(page)
	if err != nil {
		return 1
	}

	return v
}
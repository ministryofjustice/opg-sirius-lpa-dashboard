package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type UserTasksClient interface {
	CasesWithOpenTasksByAssignee(sirius.Context, int, int) ([]sirius.Case, *sirius.Pagination, error)
	MyDetails(sirius.Context) (sirius.MyDetails, error)
	User(sirius.Context, int) (sirius.Assignee, error)
}

type userTasksVars struct {
	Assignee   sirius.Assignee
	Cases      []sirius.Case
	Pagination *Pagination
	XSRFToken  string
}

func userTasks(client UserTasksClient, tmpl Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		myDetails, err := client.MyDetails(ctx)
		if err != nil {
			return err
		}

		if !myDetails.IsManager() {
			return StatusError(http.StatusForbidden)
		}

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/users/tasks/"))
		if err != nil {
			return StatusError(http.StatusNotFound)
		}

		assignee, err := client.User(ctx, id)

		if err != nil {
			return err
		}

		cases, pagination, err := client.CasesWithOpenTasksByAssignee(ctx, id, getPage(r))
		if err != nil {
			return err
		}

		if err != nil {
			return err
		}

		vars := userTasksVars{
			Assignee:   assignee,
			Cases:      cases,
			Pagination: newPagination(pagination),
			XSRFToken:  ctx.XSRFToken,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

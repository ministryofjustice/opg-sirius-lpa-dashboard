package server

import (
	"net/http"
	"strings"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type TasksDashboardClient interface {
	TasksByAssignee(sirius.Context, int, sirius.Criteria) ([]sirius.Task, *sirius.Pagination, error)
	MyDetails(sirius.Context) (sirius.MyDetails, error)
}

type tasksDashboardVars struct {
	Tasks     []sirius.Task
	Title     string
	XSRFToken string
}

func tasksDashboard(client TasksDashboardClient, tmpl Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		myDetails, err := client.MyDetails(ctx)
		if err != nil {
			return err
		}

		criteria := sirius.Criteria{}.
			Filter("status", "Not started").
			Sort("dueDate", sirius.Ascending).
			Sort("name", sirius.Descending)

		tasks, _, err := client.TasksByAssignee(ctx, myDetails.ID, criteria)
		if err != nil {
			return err
		}

		vars := tasksDashboardVars{
			Tasks:     tasks,
			Title:     "Tasks Dashboard",
			XSRFToken: ctx.XSRFToken,
		}

		if len(myDetails.Teams) > 0 {
			teamName := strings.Trim(strings.ReplaceAll(myDetails.Teams[0].DisplayName, "Team", ""), " ")
			vars.Title = teamName + " Dashboard"
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

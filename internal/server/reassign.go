package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
)

type ReassignClient interface {
	MyDetails(sirius.Context) (sirius.MyDetails, error)
	User(sirius.Context, int) (sirius.Assignee, error)
	UserByEmail(sirius.Context, string) (sirius.User, error)
	Team(sirius.Context, int) (sirius.Team, error)
	Assign(sirius.Context, []int, int) error
}

type reassignVars struct {
	XSRFToken   string
	Selected    []int
	Assignee    sirius.Assignee
	TeamMembers []sirius.TeamMember
	Success     bool
	AssignedTo  sirius.Assignee
}

func reassign(client ReassignClient, tmpl Template) Handler {
	getAssignee := func(ctx sirius.Context, id string) (sirius.Assignee, error) {
		assigneeID, err := strconv.Atoi(id)
		if err != nil {
			return sirius.Assignee{}, err
		}

		return client.User(ctx, assigneeID)
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
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

		assignee, err := getAssignee(ctx, r.FormValue("assignee"))
		if err != nil {
			return err
		}

		var selected []int
		for _, v := range r.Form["selected"] {
			i, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			selected = append(selected, i)
		}

		team, err := client.Team(ctx, assignee.Teams[0].ID)
		if err != nil {
			return err
		}

		vars := reassignVars{
			XSRFToken:   ctx.XSRFToken,
			Selected:    selected,
			Assignee:    assignee,
			TeamMembers: team.Members,
		}

		if r.Method == http.MethodPost {
			var reassignTo sirius.Assignee
			if r.FormValue("reassign") == "central-pot" {
				centralPot, err := client.UserByEmail(ctx, "manager@opgtest.com")
				if err != nil {
					return err
				}
				reassignTo = sirius.Assignee{
					ID:          centralPot.ID,
					DisplayName: "Central Pot",
				}
			} else {
				reassignTo, err = getAssignee(ctx, r.FormValue("caseworker"))
				if err != nil {
					return err
				}
			}

			if err := client.Assign(ctx, selected, reassignTo.ID); err != nil {
				return err
			}

			vars.Success = true
			vars.AssignedTo = reassignTo
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

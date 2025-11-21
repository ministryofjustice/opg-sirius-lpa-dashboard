describe("Tasks dashboard", () => {
  it("shows your tasks", () => {
    cy.addMock("/lpa-api/v1/users/current", "GET", {
      status: 200,
      body: {
        displayName: "Manager",
        id: 106,
        teams: [{ displayName: "my team" }],
      },
    });

    cy.addTaskFilterMock(
      {
        assigneeId: 106,
        filter: "status:Not started",
        sort: "dueDate:asc,name:desc",
      },
      [
        {
          caseItems: [
            {
              caseSubtype: "pfa",
              donor: {
                firstname: "Adrian",
                id: 23,
                surname: "Kurkjian",
                uId: "7000-5382-4438",
              },
              id: 1,
              uId: "7000-8548-8461",
            },
          ],
          dueDate: "19/05/2021",
          id: 36,
          name: "something",
          status: "Not started",
        },
      ],
    );

    cy.visit("/tasks-dashboard");

    cy.title().should("contain", "my team Dashboard");
    cy.get("h1").should("contain", "my team Dashboard");

    const $row = cy.get("table > tbody > tr");
    $row.should("contain", "Adrian Kurkjian");
    $row.should("contain", "PF");
    $row.should("contain", "19 May 2021");
    $row
      .contains("7000-8548-8461")
      .should("have.attr", "href")
      .should("contain", "/person/23/1");
  });
});

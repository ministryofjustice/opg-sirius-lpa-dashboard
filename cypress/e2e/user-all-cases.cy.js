describe("All of another user's cases", () => {
  it("shows all the user's cases", () => {
    cy.addMock("/api/v1/users/current", "GET", {
      status: 200,
      body: {
        displayName: "A Manager",
        id: 114,
        roles: ["Manager"],
      },
    });

    cy.addMock("/api/v1/users/47", "GET", {
      status: 200,
      body: {
        displayName: "John Paulson",
        id: 47,
        teams: [
          {
            displayName: "Cool Team",
            id: 12,
          },
        ],
      },
    });

    cy.addCaseFilterMock(
      {
        assigneeId: 47,
        filter: "caseType:lpa,active:true",
      },
      [
        {
          caseSubtype: "pfa",
          donor: {
            id: 23,
            firstname: "Adrian",
            surname: "Kurkjian",
            uId: "7000-5382-4438",
          },
          id: 36,
          receiptDate: "12/05/2021",
          status: "Perfect",
          uId: "7000-8548-8461",
        },
      ],
    );

    cy.visit("/users/all-cases/47");

    cy.title().should("contain", "John");
    cy.get("h1").should("contain", "John");

    cy.get(".govuk-back-link").should("contain", "Cool Team");

    cy.get(".govuk-tabs__list-item--selected").should("contain", "All cases");

    const $row = cy.get("table > tbody > tr");
    $row.should("contain", "Adrian Kurkjian");
    $row.should("contain", "PF");
    $row.should("contain", "12 May 2021");
    $row
      .get("a")
      .contains("7000-8548-8461")
      .should("have.attr", "href")
      .should("contain", "/person/23/36");
    $row.get(".govuk-tag").should("contain", "Perfect");
  });
});

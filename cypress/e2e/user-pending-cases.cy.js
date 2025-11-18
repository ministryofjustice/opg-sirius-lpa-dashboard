describe("Another user's pending cases", () => {
  it("shows the user's pending cases", () => {
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
        filter: "status:Pending,caseType:lpa,active:true",
      },
      [
        {
          caseSubtype: "hw",
          donor: {
            firstname: "Wilma",
            id: 17,
            surname: "Ruthman",
            uId: "7000-5382-4438",
          },
          id: 58,
          receiptDate: "14/05/2021",
          status: "Pending",
          taskCount: 1,
          uId: "7000-2830-9492",
        },
      ],
    );

    cy.visit("/users/pending-cases/47");

    cy.title().should("contain", "John");
    cy.get("h1").should("contain", "John");

    cy.get(".govuk-back-link").should("contain", "Cool Team");

    cy.get(".govuk-tabs__list-item--selected").should(
      "contain",
      "Pending cases",
    );

    const $row = cy.get("table > tbody > tr");
    $row.should("contain", "Wilma Ruthman");
    $row.should("contain", "HW");
    $row.should("contain", "14 May 2021");
    $row
      .get("a")
      .contains("7000-2830-9492")
      .should("have.attr", "href")
      .should("contain", "/person/17/58");
    $row.get(".govuk-tag").should("contain", "Pending");
  });
});

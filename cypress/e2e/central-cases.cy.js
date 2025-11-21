describe("Central pot pending cases", () => {
  it("should not be available to non-managers", () => {
    cy.visit("/teams/central", {
      failOnStatusCode: false,
    });

    cy.title().should("contain", "Forbidden");
    cy.get("h1").should("contain", "Forbidden");
  });

  it("shows your cases", () => {
    cy.addMock("/api/v1/users/current", "GET", {
      status: 200,
      body: {
        displayName: "Central Manager",
        id: 293,
        roles: ["Manager"],
      },
    });

    cy.addMock("/api/v1/users?email=opgcasework@publicguardian.gov.uk", "GET", {
      status: 200,
      body: {
        id: 14,
      },
    });

    // Latest cases
    cy.addCaseFilterMock(
      {
        assigneeId: 14,
        filter: "status:Pending,caseType:lpa,active:true",
        sort: "receiptDate:asc",
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
          uId: "7000-2830-9492",
        },
      ],
    );

    // Oldest cases
    cy.addCaseFilterMock(
      {
        assigneeId: 14,
        filter: "status:Pending,caseType:lpa,active:true",
        limit: 1,
        page: 1,
        sort: "receiptDate:asc",
      },
      [
        {
          caseSubtype: "hw",
          donor: {
            firstname: "Mario",
            id: 363,
            surname: "Evanosky",
            uId: "7000-5382-4435",
          },
          id: 453,
          receiptDate: "28/11/2017",
          status: "Pending",
          uId: "7000-2830-9429",
        },
      ],
    );

    cy.visit("/teams/central");

    cy.title().should("contain", "LPA Allocations");
    cy.get("h1").should("contain", "LPA allocations");

    cy.get(".govuk-tabs__list-item--selected").should(
      "contain",
      "Central pot - unallocated cases",
    );
    cy.get(".moj-ticket-panel .govuk-heading-xl")
      .invoke("text")
      .should("equal", "1");

    cy.get(".moj-ticket-panel").should(
      "contain",
      "Oldest case date: 28 Nov 2017",
    );

    const $row = cy.get("table > tbody > tr");
    $row.should("contain", "Wilma Ruthman");
    $row.should("contain", "HW");
    $row.should("contain", "14 May 2021");
    $row
      .contains("7000-2830-9492")
      .should("have.attr", "href")
      .should("contain", "/person/17/58");
  });
});

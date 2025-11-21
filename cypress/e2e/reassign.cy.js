describe("Reassign", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/users/current", "GET", {
      status: 200,
      body: {
        displayName: "Manager",
        id: 105,
        roles: ["Manager"],
      },
    });

    cy.addMock("/lpa-api/v1/users/47", "GET", {
      status: 200,
      body: {
        id: 47,
        displayName: "John",
        teams: [
          {
            id: 3,
            name: "Team A",
          },
        ],
      },
    });

    cy.addCaseFilterMock(
      {
        assigneeId: 47,
        filter: "status:Pending,caseType:lpa,active:true",
        sort: "receiptDate:asc",
      },
      [
        {
          id: 58,
          uId: "7000-2830-9492",
          assignee: {
            id: 17,
          },
        },
      ],
    );

    cy.visit("/users/pending-cases/47");
  });

  it("allows reassigning cases", () => {
    cy.get("tbody label[for=58]").click();

    cy.addCaseFilterMock({
      assigneeId: 105,
      filter: "status:Pending,caseType:lpa,active:true",
      sort: "workedDate:desc,receiptDate:asc",
    });

    cy.addCaseFilterMock({
      assigneeId: 105,
      filter: "status:Pending,worked:false,caseType:lpa,active:true",
    });

    cy.addMock("/lpa-api/v1/teams/3", "GET", {
      status: 200,
      body: {
        id: 47,
        teams: [
          {
            id: 3,
            name: "Team A",
          },
        ],
      },
    });

    cy.contains("Reassign or return selected case(s)").click();

    cy.addMock(
      "/lpa-api/v1/users?email=opgcasework@publicguardian.gov.uk",
      "GET",
      {
        status: 200,
        body: {
          id: 14,
        },
      },
    );

    cy.addMock("/lpa-api/v1/users/14/cases/58", "PUT", {
      status: 200,
    });

    cy.contains("Return to central pot").click();
    cy.contains("Submit").click();

    cy.contains("The case has been reassigned from John to Central Pot.");
    cy.contains("Continue").click();

    cy.get(".govuk-tabs__list-item--selected").should(
      "contain",
      "Pending cases",
    );
  });
});

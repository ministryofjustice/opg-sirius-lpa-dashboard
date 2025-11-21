describe("Team work in progress", () => {
  beforeEach(() => {
    cy.addMock("/api/v1/users/current", "GET", {
      status: 200,
      body: {
        displayName: "Manager",
        id: 107,
        roles: ["Manager"],
      },
    });

    cy.addMock("/api/v1/teams", "GET", {
      status: 200,
      body: [
        {
          displayName: "Casework Team",
          id: 66,
          members: [
            {
              displayName: "John",
              id: 47,
            },
          ],
        },
        {
          displayName: "Nottingham casework team",
          id: 67,
        },
      ],
    });

    cy.addMock("/api/v1/teams/66/cases?page=1", "GET", {
      status: 200,
      body: {
        cases: [
          {
            assignee: {
              displayName: "John Smith",
              id: 17,
            },
            caseSubtype: "pfa",
            donor: {
              firstname: "Adrian",
              id: 23,
              surname: "Kurkjian",
              uId: "7000-5382-4438",
            },
            id: 36,
            receiptDate: "12/05/2021",
            status: "Perfect",
            uId: "7000-8548-8461",
          },
        ],
        metadata: {
          tasksCompleted: [
            {
              assignee: {
                displayName: "John Smith",
                id: 17,
              },
              total: 3,
            },
          ],
          worked: [
            {
              assignee: {
                displayName: "John Smith",
                id: 17,
              },
              total: 1,
            },
          ],
          workedTotal: 1,
        },
      },
    });

    cy.visit("/teams/work-in-progress/66");
  });

  it("shows cases for my team", () => {
    cy.title().should("contain", "LPA Allocations");
    cy.get("h1").should("contain", "LPA allocations");

    cy.get(".govuk-tabs__list-item--selected").should(
      "contain",
      "Casework Team - work in progress",
    );
    cy.get(".moj-ticket-panel .govuk-heading-xl")
      .invoke("text")
      .should("equal", "1");

    cy.get(".moj-ticket-panel").should("contain", "Casework Team");

    cy.get("table > tbody > tr").within(() => {
      cy.contains("Adrian Kurkjian");
      cy.contains("PF");
      cy.contains("12 May 2021");
      cy.contains("John Smith");
      cy.contains("Perfect");
      cy.contains("7000-8548-8461")
        .should("have.attr", "href")
        .should("contain", "/person/23/36");
    });
  });

  it("shows worked cases for each team member", () => {
    cy.get(".app-name-grid:visible").within(() => {
      cy.contains("John Smith");
      cy.contains("1").should("be.visible");
    });
  });

  it("shows tasks completed for each team member", () => {
    cy.get("#data-select").select("Tasks completed today");

    cy.get(".app-name-grid:visible").within(() => {
      cy.contains("John Smith");
      cy.contains("3").should("be.visible");
    });
  });

  it("can be filtered", () => {
    cy.contains("Apply filters").should("not.be.visible");

    cy.contains("Show filters").click();
    cy.contains("label", "John").click();

    cy.addCaseFilterMock({
      assigneeId: 107,
      filter: "status:Pending,worked:false,caseType:lpa,active:true",
    });

    cy.addCaseFilterMock({
      assigneeId: 107,
      filter: "status:Pending,caseType:lpa,active:true",
      sort: "workedDate:desc,receiptDate:asc",
    });

    cy.addMock("/api/v1/teams/66/cases?filter=allocation%3A47&page=1", "GET", {
      status: 200,
      body: {
        cases: [
          {
            assignee: {
              displayName: "John Smith",
              id: 17,
            },
            caseSubtype: "pfa",
            donor: {
              firstname: "Someone",
              id: 23,
              surname: "Else",
              uId: "7000-5382-4438",
            },
            id: 36,
            receiptDate: "12/05/2021",
            status: "Perfect",
            uId: "7000-8548-8461",
          },
        ],
        limit: 25,
        metadata: {
          tasksCompleted: [
            {
              assignee: {
                displayName: "John Smith",
                id: 17,
              },
              total: 3,
            },
          ],
          worked: [
            {
              assignee: {
                displayName: "John Smith",
                id: 17,
              },
              total: 1,
            },
          ],
          workedTotal: 1,
        },
        pages: {
          current: 1,
          total: 1,
        },
        total: 1,
      },
    });

    cy.contains("Apply filters").click();
    cy.url().should("contain", "allocation=47");

    cy.get("table > tbody > tr").within(() => {
      cy.contains("Someone Else");
    });

    cy.contains("Reset").click();
    cy.url().should("not.contain", "allocation=123");
    cy.contains("Apply filters").should("not.be.visible");
  });

  it("enables navigation to other teams via dropdown", () => {
    cy.addCaseFilterMock({
      assigneeId: 107,
      filter: "status:Pending,worked:false,caseType:lpa,active:true",
    });

    cy.addCaseFilterMock({
      assigneeId: 107,
      filter: "status:Pending,caseType:lpa,active:true",
      sort: "workedDate:desc,receiptDate:asc",
    });

    cy.addMock("/api/v1/teams/67/cases?page=1", "GET", {
      status: 200,
      body: {
        cases: [
          {
            id: 853,
            receiptDate: "12/05/2021",
            status: "Perfect",
            uId: "7000-8548-2721",
          },
        ],
      },
    });

    cy.get("[data-select-navigate]").select("Nottingham casework team");
    cy.url().should("contain", "teams/work-in-progress/67");
    cy.get(".govuk-tabs__list-item--selected").should(
      "contain",
      "Nottingham casework team - work in progress",
    );
    cy.get(".moj-ticket-panel").should("contain", "Nottingham casework team");
  });
});

describe("Case navigation", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/users/current", "GET", {
      status: 200,
      body: {
        displayName: "Central Manager",
        id: 105,
        roles: ["Manager", "Self Allocation User"],
        teams: [
          {
            displayName: "my team",
            id: 66,
          },
        ],
      },
    });

    cy.addCaseFilterMock({
      assigneeId: 105,
      filter: "status:Pending,worked:false,caseType:lpa,active:true",
    });

    cy.addCaseFilterMock({
      assigneeId: 105,
      filter: "status:Pending,caseType:lpa,active:true",
    });

    cy.visit("/");
  });

  it("can direct to pending-cases via tab", () => {
    cy.url().should("include", "/pending-cases");
    cy.title().should("contain", "Your cases");
    cy.get("h1").should("contain", "Your cases");
  });

  it("can direct to your cases tasks tab", () => {
    cy.contains("Tasks").click();
    cy.url().should("include", "/tasks");
  });

  it("can direct to your cases all cases tab", () => {
    cy.contains("All cases").click();
    cy.url().should("include", "/all-cases");
  });

  it("can direct to feedback page", () => {
    cy.contains("Feedback").click();
    cy.url().should("include", "/feedback");
  });

  it("can direct to LPA allocations from your cases", () => {
    cy.addCaseFilterMock({
      assigneeId: 105,
      filter: "caseType:lpa,active:true",
    });

    cy.contains("LPA allocations").click();
    cy.url().should("include", "/teams/central");
  });

  it("can direct to my team tab on LPA allocations", () => {
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

    cy.addCaseFilterMock({
      assigneeId: 14,
      filter: "status:Pending,caseType:lpa,active:true",
    });

    cy.contains("LPA allocations").click();
    cy.contains("my team - work in progress").click();
    cy.url().should("include", "/work-in-progress");
  });

  it("can direct to central pot tab on LPA allocations", () => {
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

    cy.addCaseFilterMock({
      assigneeId: 14,
      filter: "status:Pending,caseType:lpa,active:true",
    });

    cy.contains("LPA allocations").click();
    cy.contains("Central pot - unallocated cases").click();
    cy.url().should("include", "/teams/central");
  });

  it("can direct to your cases page from LPA allocations", () => {
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

    cy.addCaseFilterMock({
      assigneeId: 14,
      filter: "status:Pending,caseType:lpa,active:true",
    });

    cy.contains("LPA allocations").click();
    cy.contains("Your cases").click();
    cy.url().should("include", "/pending-cases");
    cy.get("h1").should("contain", "Your cases");
  });

  it("can direct to LPA allocations team tab from another user's cases", () => {
    cy.addMock("/lpa-api/v1/users/47", "GET", {
      status: 200,
      body: {
        teams: [
          {
            displayName: "Cool Team",
            id: 88,
          },
        ],
      },
    });

    cy.addMock("/lpa-api/v1/assignees/47/cases-with-open-tasks?page=1", "GET", {
      status: 200,
      body: {},
    });

    cy.visit("/users/tasks/47");

    cy.contains("Cool Team").click();
    cy.url().should("include", "/work-in-progress");
  });
});

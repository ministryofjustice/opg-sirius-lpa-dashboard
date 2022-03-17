describe("Case navigation", () => {
  beforeEach(() => {
    cy.setCookie("Other", "other");
    cy.setCookie("XSRF-TOKEN", "abcde");
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
    cy.contains("LPA allocations").click();
    cy.url().should("include", "/teams/central");
  });

  it("can direct to my team tab on LPA allocations", () => {
    cy.contains("LPA allocations").click();
    cy.contains("my team - work in progress").click();
    cy.url().should("include", "/work-in-progress");
  });

  it("can direct to central pot tab on LPA allocations", () => {
    cy.contains("LPA allocations").click();
    cy.contains("Central pot - unallocated cases").click();
    cy.url().should("include", "/teams/central");
  });

  it("can direct to your cases page from LPA allocations", () => {
    cy.contains("LPA allocations").click();
    cy.contains("Your cases").click();
    cy.url().should("include", "/pending-cases");
    cy.get("h1").should("contain", "Your cases");
  });

  it("can direct to LPA allocations team tab from another users cases", () => {
    cy.visit("/users/tasks/47");
    cy.contains("Cool Team").click();
    cy.url().should("include", "/work-in-progress");
  });
});

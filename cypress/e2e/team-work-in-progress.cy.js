describe("Team work in progress", () => {
  beforeEach(() => {
    cy.visit("/teams/work-in-progress/66");
  });

  it("shows cases for my team", () => {
    cy.title().should("contain", "LPA Allocations");
    cy.get("h1").should("contain", "LPA allocations");

    cy.get(".govuk-tabs__list-item--selected").should(
      "contain",
      "Casework Team - work in progress"
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
    cy.contains("Apply filters").click();
    cy.url().should("contain", "allocation=47");

    cy.get("table > tbody > tr").within(() => {
      cy.contains("Someone Else");
    });

    cy.contains("Reset").click();
    cy.url().should("not.contain", "allocation=123");
    cy.contains("Apply filters").should("not.be.visible");
  });
});

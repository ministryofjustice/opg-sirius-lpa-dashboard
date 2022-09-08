describe("All of another user's cases", () => {
  beforeEach(() => {
    cy.visit("/users/all-cases/47");
  });

  it("shows all the user's cases", () => {
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

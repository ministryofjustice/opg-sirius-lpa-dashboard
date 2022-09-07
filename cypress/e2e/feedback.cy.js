describe("Feedback", () => {
  beforeEach(() => {
    cy.visit("/all-cases");
  });

  it("returns to previous page on send", () => {
    cy.contains("Feedback").click();

    cy.get("textarea").type("Hey");
    cy.contains("Submit").click();

    cy.url().should("include", "/all-cases");
  });
});

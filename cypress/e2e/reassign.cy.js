describe("Reassign", () => {
  beforeEach(() => {
    cy.visit("/users/pending-cases/47");
  });

  it("allows reassigning cases", () => {
    cy.get("tbody label[for=58]").click();
    cy.contains("Reassign or return selected case(s)").click();

    cy.contains("Return to central pot").click();
    cy.contains("Submit").click();

    cy.contains("The case has been reassigned from John to Central Pot.");
    cy.contains("Continue").click();

    cy.get(".govuk-tabs__list-item--selected").should(
      "contain",
      "Pending cases"
    );
  });
});

describe("Pending cases", () => {
  beforeEach(() => {
    cy.setCookie("Other", "other");
    cy.setCookie("XSRF-TOKEN", "abcde");
    cy.visit("/pending-cases");
  });

  it("shows your cases", () => {
    cy.title().should("contain", "Your cases");
    cy.get("h1").should("contain", "Your cases");

    const $row = cy.get("table > tbody > tr");
    $row.should("contain", "Wilma Ruthman");
    $row.should("contain", "HW");
    $row.should("contain", "14 May 2021");
    $row
      .contains("7000-2830-9492")
      .should("have.attr", "href")
      .should("contain", "/person/17/58");
  });

  it("enables the 'Progress worked cases' button on selection", () => {
    cy.contains("button", "Progress worked cases").should(
      "have.attr",
      "disabled"
    );
    cy.get("input[type=checkbox]").check();
    cy.contains("button", "Progress worked cases").should(
      "not.have.attr",
      "disabled"
    );
  });
});

describe("All cases", () => {
  beforeEach(() => {
    cy.setCookie("Other", "other");
    cy.setCookie("XSRF-TOKEN", "abcde");
    cy.visit("/all-cases");
  });

  it("shows your cases", () => {
    cy.title().should("contain", "Your cases");
    cy.get("h1").should("contain", "Your cases");

    const $row = cy.get("table > tbody > tr");
    $row.should("contain", "Adrian Kurkjian");
    $row.should("contain", "PF");
    $row.should("contain", "12 May 2021");
    $row.should("contain", "Perfect");
    $row
      .contains("7000-8548-8461")
      .should("have.attr", "href")
      .should("contain", "/person/23/36");
  });
});

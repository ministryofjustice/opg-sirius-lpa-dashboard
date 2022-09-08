describe("Tasks", () => {
  beforeEach(() => {
    cy.visit("/tasks");
  });

  it("shows your cases", () => {
    cy.title().should("contain", "Your cases");
    cy.get("h1").should("contain", "Your cases");

    const $row = cy.get("table > tbody > tr");
    $row.should("contain", "Wilma Ruthman");
    $row.should("contain", "HW");
    $row.should("contain", "1 task");
    $row.should("contain", "Pending");
    $row
      .contains("7000-2830-9492")
      .should("have.attr", "href")
      .should("contain", "/person/17/58");
  });
});

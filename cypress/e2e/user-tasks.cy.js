describe("Another user's tasks", () => {
  beforeEach(() => {
    cy.visit("/users/tasks/47");
  });

  it("shows the user's tasks", () => {
    cy.title().should("contain", "John");
    cy.get("h1").should("contain", "John");

    cy.get(".govuk-back-link").should("contain", "Cool Team");

    cy.get(".govuk-tabs__list-item--selected").should("contain", "Tasks");

    const $row = cy.get("table > tbody > tr");
    $row.should("contain", "Wilma Ruthman");
    $row.should("contain", "HW");
    $row.should("contain", "1 task");
    $row
      .get("a")
      .contains("7000-2830-9492")
      .should("have.attr", "href")
      .should("contain", "/person/17/58");
    $row.get(".govuk-tag").should("contain", "Pending");
  });
});

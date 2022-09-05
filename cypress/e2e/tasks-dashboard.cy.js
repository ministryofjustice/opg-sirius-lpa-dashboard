describe("Tasks dashboard", () => {
  beforeEach(() => {
    cy.setCookie("Other", "other");
    cy.setCookie("XSRF-TOKEN", "abcde");
    cy.visit("/tasks-dashboard");
  });

  it("shows your tasks", () => {
    cy.title().should("contain", "my team Dashboard");
    cy.get("h1").should("contain", "my team Dashboard");

    const $row = cy.get("table > tbody > tr");
    $row.should("contain", "Adrian Kurkjian");
    $row.should("contain", "PF");
    $row.should("contain", "19 May 2021");
    $row
      .contains("7000-8548-8461")
      .should("have.attr", "href")
      .should("contain", "/person/23/1");
  });
});

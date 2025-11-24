describe("Feedback", () => {
  it("returns to previous page on send", () => {
    // My assigned cases
    cy.addCaseFilterMock(
      {
        assigneeId: 104,
        filter: "caseType:lpa,active:true",
        sort: "receiptDate:asc",
      },
      [],
    );

    // Unworked cases assigned to me
    cy.addCaseFilterMock(
      {
        assigneeId: 104,
        filter: "status:Pending,worked:false,caseType:lpa,active:true",
        limit: 1,
      },
      [],
    );

    cy.visit("/all-cases");

    cy.contains("Feedback").click();

    cy.get("textarea").type("Hey");

    cy.addMock("/lpa-api/v1/feedback/poas", "POST", {
      status: 200,
    });

    cy.contains("Submit").click();

    cy.url().should("include", "/all-cases");
  });
});

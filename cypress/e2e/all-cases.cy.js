describe("All cases", () => {
  it("shows your cases", () => {
    // My assigned cases
    cy.addCaseFilterMock(
      {
        assigneeId: 104,
        filter: "caseType:lpa,active:true",
        sort: "receiptDate:asc",
      },
      [
        {
          caseSubtype: "pfa",
          donor: {
            firstname: "Adrian",
            id: 23,
            surname: "Kurkjian",
            uId: "7000-5382-4438",
          },
          id: 36,
          receiptDate: "12/05/2021",
          status: "Perfect",
          uId: "7000-8548-8461",
        },
      ],
    );

    // Unworked cases assigned to me
    cy.addCaseFilterMock(
      {
        assigneeId: 104,
        filter: "status:Pending,worked:false,caseType:lpa,active:true",
        limit: 1,
      },
      [
        {
          caseSubtype: "hw",
          donor: {
            firstname: "Mario",
            id: 363,
            surname: "Evanosky",
            uId: "7000-5382-4435",
          },
          id: 453,
          receiptDate: "28/11/2017",
          status: "Pending",
          uId: "7000-2830-9429",
        },
      ],
    );

    cy.visit("/all-cases");

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

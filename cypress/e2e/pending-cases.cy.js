describe("Pending cases", () => {
  beforeEach(() => {
    cy.addCaseFilterMock(
      {
        assigneeId: 104,
        filter: "status:Pending,caseType:lpa,active:true",
        sort: "workedDate:desc,receiptDate:asc",
      },
      [
        {
          caseSubtype: "hw",
          donor: {
            firstname: "Wilma",
            id: 17,
            surname: "Ruthman",
            uId: "7000-5382-4438",
          },
          id: 58,
          receiptDate: "14/05/2021",
          status: "Pending",
          uId: "7000-2830-9492",
        },
      ],
    );

    cy.addCaseFilterMock({
      assigneeId: 104,
      filter: "status:Pending,worked:false,caseType:lpa,active:true",
    });

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
      "disabled",
    );
    cy.get("label[for=58]").click();
    cy.contains("button", "Progress worked cases").should(
      "not.have.attr",
      "disabled",
    );
  });
});

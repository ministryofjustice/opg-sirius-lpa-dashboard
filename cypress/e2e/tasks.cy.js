describe("Tasks", () => {
  it("shows your cases", () => {
    cy.addMock("/api/v1/assignees/104/cases-with-open-tasks?page=1", "GET", {
      status: 200,
      body: {
        cases: [
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
            taskCount: 1,
            uId: "7000-2830-9492",
          },
        ],
        limit: 25,
        pages: {
          current: 1,
          total: 1,
        },
        total: 1,
      },
    });

    cy.addCaseFilterMock({
      assigneeId: 104,
      filter: "status:Pending,worked:false,caseType:lpa,active:true",
      limit: 1,
    });

    cy.visit("/tasks");

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

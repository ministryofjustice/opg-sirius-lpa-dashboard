
describe("All cases", () => {
  beforeEach(() => {
      cy.setCookie("Other", "other");
      cy.setCookie("XSRF-TOKEN", "abcde");
      cy.visit("/");
  });

  it("will redirect to pending-cases", () => {
    cy.url().should("include", "/pending-cases");
    cy.title().should('contain', 'Your cases');
    cy.get('h1').should('contain', 'Your cases');
  });

  it("can direct to tasks hyperlink via tab", () => {
    cy.contains('Tasks').click();;
    cy.url().should("include", "/tasks");
  });

  it("can direct to all cases hyperlink via tab", () => {
    cy.contains('All cases').click();;
    cy.url().should("include", "/all-cases");
  });

  it("can direct to sirius using task hyperlink", () => {
    cy.contains('7000-2830-9492').click();
    cy.get('.app-auth-error').should('contain', 'You have been logged out due to session inactivity');
  });

})
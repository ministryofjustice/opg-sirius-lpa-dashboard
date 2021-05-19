describe("The tests run", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/", { failOnStatusCode: false });
    });

    it("goes to a page", () => {
        cy.contains('Your cases');
    });
});

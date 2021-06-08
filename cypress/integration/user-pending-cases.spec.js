describe("Pending cases", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/users/pending-cases/47");
    });

    it("shows your cases", () => {
        cy.title().should('contain', 'John');
        cy.get('h1').should('contain', 'John');

        const $row = cy.get('table > tbody > tr');
        $row.should('contain', 'Wilma Ruthman');
        $row.should('contain', 'HW');
        $row.should('contain', '14 May 2021');
        $row.get('a').contains('7000-2830-9492').should('have.attr', 'href').should('contain', '/person/17/58');
    });
});

describe("Central pot pending cases", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/teams/central");
    });

    it("shows your cases", () => {
        cy.title().should('contain', 'LPA Allocations');
        cy.get('h1').should('contain', 'LPA allocations');

        cy.get('.govuk-tabs__list-item--selected').should('contain', 'Central pot - pending cases');
        cy.get('.moj-ticket-panel .govuk-heading-xl').invoke('text').should('equal', '1')

        cy.get('.moj-ticket-panel').should('contain', 'Oldest case date: 28 Nov 2017');

        const $row = cy.get('table > tbody > tr');
        $row.should('contain', 'Wilma Ruthman');
        $row.should('contain', 'HW');
        $row.should('contain', '14 May 2021');
        $row.contains('7000-2830-9492').should('have.attr', 'href').should('contain', '/person/17/58');
    });

    it('cross-links to caseworker view', () => {
        cy.contains('Your cases').click();
        cy.url().should("include", "/pending-cases");
    });
});

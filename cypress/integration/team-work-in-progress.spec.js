describe("Team work in progress", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/teams/work-in-progress");
    });

    it("shows cases for my team", () => {
        cy.title().should('contain', 'LPA Allocations');
        cy.get('h1').should('contain', 'LPA allocations');

        cy.get('.govuk-tabs__list-item--selected').should('contain', 'my team - work in progress');
        cy.get('.moj-ticket-panel .govuk-heading-xl').invoke('text').should('equal', '1')

        cy.get('.moj-ticket-panel').should('contain', 'my team');

        const $row = cy.get('table > tbody > tr');
        $row.should('contain', 'Adrian Kurkjian');
        $row.should('contain', 'PF');
        $row.should('contain', '12 May 2021');
        $row.should('contain', 'John Smith');
        $row.should('contain', 'Perfect');
        $row.contains('7000-8548-8461').should('have.attr', 'href').should('contain', '/person/23/36');
    });

    it('shows statistics for each team member', () => {
        const $row = cy.get('.app-name-grid');

        $row.should('contain', 'John Smith');
        $row.should('contain', '1');
    });
});

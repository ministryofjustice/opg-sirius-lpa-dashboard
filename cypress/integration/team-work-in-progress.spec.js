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

        cy.get('table > tbody > tr').within(() => {
            cy.contains('Adrian Kurkjian');
            cy.contains('PF');
            cy.contains('12 May 2021');
            cy.contains('John Smith');
            cy.contains('Perfect');
            cy.contains('7000-8548-8461').should('have.attr', 'href').should('contain', '/person/23/36');
        });
    });

    it('shows worked cases for each team member', () => {
        cy.get('.app-name-grid').within(() => {
            cy.contains('John Smith');
            cy.contains('1').should('be.visible');
        });
    });

    it('shows tasks completed for each team member', () => {
        cy.get('select').select('Tasks completed today');

        cy.get('.app-name-grid').within(() => {
            cy.contains('John Smith');
            cy.contains('3').should('be.visible');
        });
    });
});

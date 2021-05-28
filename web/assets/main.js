import MOJFrontend from '@ministryofjustice/frontend/moj/all.js';
import GOVUKFrontend from 'govuk-frontend/govuk/all.js';
import './main.scss';

function initEnableWhenSelection() {
    const button = document.querySelector('button[data-enable-when-selection]');
    if (button) {
        const checkboxes = Array.from(document.querySelectorAll('input[type=checkbox]'));
        button.disabled = !checkboxes.some(x => x.checked);

        checkboxes.forEach(checkbox => {
            checkbox.onchange = () => {
                button.disabled = !checkboxes.some(x => x.checked);
            };
        });
    }
}

// we aren't using the JS tabs, but they try to initialise this will stop them breaking
GOVUKFrontend.Tabs.prototype.setup = () => { };

GOVUKFrontend.initAll();
MOJFrontend.initAll();
initEnableWhenSelection();

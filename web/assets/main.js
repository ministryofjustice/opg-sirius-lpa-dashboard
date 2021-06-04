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

function initSelectShow() {
    const select = document.querySelector('select[data-select-show]');
    if (select) {
        function update() {
            const selectIds = Array.from(document.querySelectorAll('[data-select-id]'));
            selectIds.forEach(x => x.classList.add('govuk-!-display-none'));

            const show = document.querySelector(`[data-select-id='${select.value}']`);
            show.classList.remove('govuk-!-display-none');
        }

        select.onchange = update;
        update();
    }
}

function initSelectNavigate() {
    const select = document.querySelector('select[data-select-navigate]');
    if (select) {
        select.onchange = () => {
            window.location.href = select.value;
        };
    }
}

// we aren't using the JS tabs, but they try to initialise this will stop them breaking
GOVUKFrontend.Tabs.prototype.setup = () => { };

GOVUKFrontend.initAll();
MOJFrontend.initAll();
initEnableWhenSelection();
initSelectShow();
initSelectNavigate();

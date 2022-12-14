import MOJFrontend from "@ministryofjustice/frontend/moj/all.js";
import * as GOVUKFrontend from "govuk-frontend";
import $ from "jquery";
import 'opg-sirius-search-ui';

function initEnableWhenSelection() {
  const button = document.querySelector("button[data-enable-when-selection]");
  if (button) {
    const checkboxes = Array.from(
      document.querySelectorAll("table input[type=checkbox]")
    );
    const bodyCheckboxes = Array.from(
      document.querySelectorAll("tbody input[type=checkbox]")
    );
    button.disabled = !bodyCheckboxes.some((x) => x.checked);

    checkboxes.forEach((checkbox) => {
      checkbox.onchange = () => {
        button.disabled = !bodyCheckboxes.some((x) => x.checked);
      };
    });
  }
}

function initSelectShow() {
  const select = document.querySelector("select[data-select-show]");
  if (select) {
    function update() {
      const selectIds = Array.from(
        document.querySelectorAll("[data-select-id]")
      );
      selectIds.forEach((x) => x.classList.add("govuk-!-display-none"));

      const show = document.querySelector(`[data-select-id='${select.value}']`);
      show.classList.remove("govuk-!-display-none");
    }

    select.onchange = update;
    update();
  }
}

function initSelectNavigate() {
  const select = document.querySelector("select[data-select-navigate]");
  if (select) {
    select.onchange = () => {
      window.location.href = select.value;
    };
  }
}

function initFilterToggle() {
  const button = document.querySelector("button[data-filter-toggle]");
  const filters = document.querySelector(".moj-filter-layout__filter");

  if (button && filters) {
    button.onclick = () => {
      if (button.innerText === "Hide filters") {
        button.innerText = "Show filters";
        filters.classList.add("govuk-!-display-none");
      } else {
        button.innerText = "Hide filters";
        filters.classList.remove("govuk-!-display-none");
      }
    };
  }
}

function initFilterHeadings() {
  const buttons = document.querySelectorAll(".app-c-option-select__button");

  for (const button of buttons) {
    button.onclick = () => {
      const content = document.getElementById(
        button.getAttribute("aria-controls")
      );

      content.classList.toggle("govuk-!-display-none");

      if (content.classList.contains("govuk-!-display-none")) {
        button.setAttribute("aria-expanded", "false");
      } else {
        button.setAttribute("aria-expanded", "true");
      }
    };
  }
}

document.body.className = document.body.className
  ? document.body.className + " js-enabled"
  : "js-enabled";

// Expose jQuery on window so MOJFrontend can use it
window.$ = $;

// we aren't using the JS tabs, but they try to initialise this will stop them breaking
GOVUKFrontend.Tabs.prototype.setup = () => {};

GOVUKFrontend.initAll();
MOJFrontend.initAll();
initEnableWhenSelection();
initSelectShow();
initSelectNavigate();
initFilterToggle();
initFilterHeadings();

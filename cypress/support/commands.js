import { addMock, reset } from "../mocks/wiremock";

Cypress.Commands.add("addMock", async (url, method, response) => {
  await addMock(url, method, response);
});

Cypress.Commands.add("resetMocks", async () => {
  await reset();
});

Cypress.Commands.add(
  "addCaseFilterMock",
  async (query = {}, cases = [], body = {}) => {
    const { assigneeId, ...rest } = query;

    await addMock(
      {
        path: `/api/v1/assignees/${assigneeId}/cases`,
        query: {
          page: 1,
          ...rest,
        },
      },
      "GET",
      {
        status: 200,
        body: {
          limit: query["limit"] || 25,
          pages: {
            current: 1,
            total: 1,
          },
          total: 1,
          cases,
          ...body,
        },
      },
    );
  },
);

Cypress.Commands.add(
  "addTaskFilterMock",
  async (query = {}, tasks = [], body = {}) => {
    const { assigneeId, ...rest } = query;

    await addMock(
      {
        path: `/api/v1/assignees/${assigneeId}/tasks`,
        query: {
          ...rest,
        },
      },
      "GET",
      {
        status: 200,
        body: {
          limit: query["limit"] || 25,
          pages: {
            current: 1,
            total: 1,
          },
          total: 1,
          tasks,
          ...body,
        },
      },
    );
  },
);

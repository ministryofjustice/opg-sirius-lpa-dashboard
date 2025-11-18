async function addMock(url, method, response) {
  if (typeof response.body !== "string") {
    response.body = JSON.stringify(response.body);
  }

  const request = {
    method,
  };

  if (typeof url === "string") {
    request.url = url;
  } else {
    request.urlPath = url.path;
    request.queryParameters = Object.entries(url.query).reduce(
      (acc, [key, value]) => ({
        ...acc,
        [key]: { equalTo: decodeURIComponent(value) },
      }),
      {},
    );
  }

  await fetch(`${Cypress.env("MOCK_SERVER_URI")}/__admin/mappings`, {
    method: "POST",
    body: JSON.stringify({
      request,
      response,
    }),
  });
}

async function reset() {
  await fetch(`${Cypress.env("MOCK_SERVER_URI")}/__admin/mappings/reset`, {
    method: "POST",
  });
}

export { addMock, reset };

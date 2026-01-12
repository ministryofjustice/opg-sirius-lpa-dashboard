# Sirius LPA dashboard

[![PkgGoDev](https://pkg.go.dev/badge/github.com/ministryofjustice/opg-sirius-lpa-dashboard)](https://pkg.go.dev/github.com/ministryofjustice/opg-sirius-lpa-dashboard)

## Quick start

### Major dependencies

- [Go](https://golang.org/) (>= 1.14)
- [Pact](https://github.com/pact-foundation/pact-ruby-standalone) (>= 1.88.3)
- [docker-compose](https://docs.docker.com/compose/install/) (>= 1.27.4)
- [Node](https://nodejs.org/en/) (>= 14.15.1)

### Running the application

```
docker compose up -d lpa-dashboard
```

This will run the application at http://localhost:8888/, and against the mock.

Alternatively the application can be run without the use of Docker

```
npm ci && npm run build
SIRIUS_PUBLIC_URL=http://localhost:8080 SIRIUS_URL=http://localhost:8080 PORT=8888 go run main.go
```

If you get an error message from running the `npm ci` this could be due to you having a package.json and you will
need to run `npm install` first.

If you want to run your local changes in the context of local sirius then build the local image and start up sirius.

```
make build
# cd to sirius repo
make dev-up
```

### Testing

```
make unit-test
```

This will run the Go unit tests and generate the local pact contracts for the pact stub for further testing.

```
make cypress
```

This will run Cypress aginst the lpa-dashboard and the pact-stub

## Development

Linting
This will include a check on formatting so it is recommended to setup your editor to use `go fmt`.
You can run linting locally with

```
make lint
```

If you want to build the application and run all test suites that get run in CI just run:

```
make
```

## Environment variables

| Name                | Description                         |
| ------------------- | ----------------------------------- |
| `PORT`              | Port to run on                      |
| `WEB_DIR`           | Path to the 'web' directory         |
| `SIRIUS_URL`        | Base URL to call Sirius             |
| `SIRIUS_PUBLIC_URL` | Base URL to redirect to Sirius      |
| `PREFIX`            | Path to prefix to each page's route |

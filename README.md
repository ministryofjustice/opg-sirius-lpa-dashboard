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
docker-compose -f docker/docker-compose.yml up -d --build
```

This will run the application at http://localhost:8888/, and assumes that Sirius
is running at http://localhost:8080/.

Alternatively the application can be run without the use of Docker

```
yarn && yarn build
SIRIUS_PUBLIC_URL=http://localhost:8080 SIRIUS_URL=http://localhost:8080 PORT=8888 go run main.go
```

### Testing

```
go test ./...
```

This will run the Go unit tests. It relies on `pact` being available on your
`$PATH`. The tests will produce a `./pacts` directory which is then used to
provide a stub service for the Cypress tests.

```
docker-compose -f docker/docker-compose.cypress.yml up -d --build
yarn && yarn cypress
```

Will start the application in a way that uses the stub service, then opens
Cypress in the current project.

## Development

On CI we lint using [golangci-lint](https://golangci-lint.run/). It may be
useful to install locally to check changes. This will include a check on
formatting so it is recommended to setup your editor to use `go fmt`.

## Environment variables

| Name                | Description                         |
| ------------------- | ----------------------------------- |
| `PORT`              | Port to run on                      |
| `WEB_DIR`           | Path to the 'web' directory         |
| `SIRIUS_URL`        | Base URL to call Sirius             |
| `SIRIUS_PUBLIC_URL` | Base URL to redirect to Sirius      |
| `PREFIX`            | Path to prefix to each page's route |

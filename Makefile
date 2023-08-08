.PHONY: cypress

all: lint unit-test build-all scan pa11y lighthouse cypress down

lint:
	docker compose -f docker/docker-compose.ci.yml run --rm go-lint

test-results:
	mkdir -p -m 0777 test-results .gocache pacts logs cypress/screenshots

setup-directories: test-results

unit-test: setup-directories
	docker compose -f docker/docker-compose.ci.yml run --rm test-runner

build:
	docker compose -f docker/docker-compose.ci.yml build lpa-dashboard

build-ci:
	docker compose -f docker/docker-compose.ci.yml build --parallel lpa-dashboard pact-stub

build-all:
	docker compose -f docker/docker-compose.ci.yml build --parallel lpa-dashboard pact-stub puppeteer

scan:
	trivy image lpa-dashboard:latest

pa11y: setup-directories
	docker compose -f docker/docker-compose.ci.yml run --rm --entrypoint="pa11y-ci" puppeteer

lighthouse: setup-directories
	docker compose -f docker/docker-compose.ci.yml run --rm --entrypoint="lhci autorun" puppeteer

cypress: setup-directories
	docker compose -f docker/docker-compose.ci.yml run --rm cypress

down:
	docker compose -f docker/docker-compose.ci.yml down

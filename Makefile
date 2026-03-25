.PHONY: cypress

all: lint unit-test build-all pa11y lighthouse cypress down

lint:
	docker compose run --rm go-lint

test-results:
	mkdir -p -m 0777 test-results .gocache pacts logs cypress/screenshots

setup-directories: test-results

unit-test: setup-directories
	docker compose run --rm test-runner

build:
	docker compose build lpa-dashboard

build-ci:
	docker compose build --parallel lpa-dashboard pact-stub

build-all:
	docker compose build --parallel lpa-dashboard pact-stub puppeteer

pa11y: setup-directories
	docker compose run --rm --entrypoint="pa11y-ci" puppeteer

lighthouse: setup-directories
	docker compose run --rm --entrypoint="lhci autorun" puppeteer

cypress: setup-directories
	docker compose run --rm cypress

down:
	docker compose down

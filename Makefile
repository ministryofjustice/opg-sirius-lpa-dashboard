.PHONY: cypress

all: lint unit-test build-all scan pa11y lighthouse cypress down

lint:
	docker compose run --rm go-lint

test-results:
	mkdir -p -m 0777 test-results .gocache pacts logs cypress/screenshots .trivy-cache

setup-directories: test-results

unit-test: setup-directories
	docker compose run --rm test-runner

build:
	docker compose build lpa-dashboard

build-ci:
	docker compose build --parallel lpa-dashboard pact-stub

build-all:
	docker compose build --parallel lpa-dashboard pact-stub puppeteer

scan: setup-directories
	docker compose run --rm trivy image --format table --exit-code 0 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-lpa-dashboard:latest
	docker compose run --rm trivy image --format sarif --output /test-results/trivy.sarif --exit-code 1 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-lpa-dashboard:latest

pa11y: setup-directories
	docker compose run --rm --entrypoint="pa11y-ci" puppeteer

lighthouse: setup-directories
	docker compose run --rm --entrypoint="lhci autorun" puppeteer

cypress: setup-directories
	docker compose run --rm cypress

down:
	docker compose down

FROM golang:1.25-bookworm@sha256:916d09f50df6af501cda192d349f57c66a6d738ca80241a8d18e4f1e1d2ed8a5

RUN go install gotest.tools/gotestsum@latest
RUN go install github.com/pact-foundation/pact-go/v2@v2.2.0
## This is needed to install a dependency not included in the main package
RUN pact-go -l DEBUG install

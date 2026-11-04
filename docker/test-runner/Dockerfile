FROM golang:1.27.3-bookworm@sha256:d9b6a1beb9474adc2a1ce842573639b44b2a0dd49b6669f51d3821baf9ec37c9

RUN go install gotest.tools/gotestsum@latest
RUN go install github.com/pact-foundation/pact-go/v2@v2.4.3
## This is needed to install a dependency not included in the main package
RUN pact-go -l DEBUG install

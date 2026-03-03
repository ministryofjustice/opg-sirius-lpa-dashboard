FROM golang:1.26-bookworm@sha256:2a0ba12e116687098780d3ce700f9ce3cb340783779646aafbabed748fa6677c

RUN go install gotest.tools/gotestsum@latest
RUN go install github.com/pact-foundation/pact-go/v2@v2.2.0
## This is needed to install a dependency not included in the main package
RUN pact-go -l DEBUG install
